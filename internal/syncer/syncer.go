package syncer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/user/project/internal/db"
)

type storager interface {
	SaveMatch(ctx context.Context, match db.Match) error
	GetTeamByName(ctx context.Context, name string) (db.Team, error)
	GetCompletedMatchesWithoutCompletedPredictions(ctx context.Context) ([]db.Match, error)
	GetPredictionsForMatch(ctx context.Context, matchID string) ([]db.Prediction, error)
	UpdatePredictionResult(ctx context.Context, matchID, userID string, points int) error
	GetActiveSeason(ctx context.Context) (db.Season, error)
	UpdateUserLeaderboardPoints(ctx context.Context, userID, seasonID string, points int) error
	UpdateUserPoints(ctx context.Context, userID string, points int) error
	SaveTeam(ctx context.Context, team db.Team) error
}

type Syncer struct {
	storage    storager
	apiBaseURL string
	apiKey     string
}

type APIMatch struct {
	Area struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
		Flag string `json:"flag"`
	} `json:"area"`
	Competition struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Code   string `json:"code"`
		Type   string `json:"type"`
		Emblem string `json:"emblem"`
	} `json:"competition"`
	Season struct {
		Id              int     `json:"id"`
		StartDate       string  `json:"startDate"`
		EndDate         string  `json:"endDate"`
		CurrentMatchday int     `json:"currentMatchday"`
		Winner          *string `json:"winner"`
	} `json:"season"`
	Id          int       `json:"id"`
	UtcDate     time.Time `json:"utcDate"`
	Status      string    `json:"status"`
	Matchday    int       `json:"matchday"`
	Stage       string    `json:"stage"`
	Group       *string   `json:"group"`
	LastUpdated time.Time `json:"lastUpdated"`
	HomeTeam    struct {
		Id        *int    `json:"id"`
		Name      *string `json:"name"`
		ShortName *string `json:"shortName"`
		Tla       *string `json:"tla"`
		Crest     *string `json:"crest"`
	} `json:"homeTeam"`
	AwayTeam struct {
		Id        *int    `json:"id"`
		Name      *string `json:"name"`
		ShortName *string `json:"shortName"`
		Tla       *string `json:"tla"`
		Crest     *string `json:"crest"`
	} `json:"awayTeam"`
	Score struct {
		Winner   *string `json:"winner"`
		Duration string  `json:"duration"`
		FullTime struct {
			Home *int `json:"home"`
			Away *int `json:"away"`
		} `json:"fullTime"`
		HalfTime struct {
			Home *int `json:"home"`
			Away *int `json:"away"`
		} `json:"halfTime"`
	} `json:"score"`
	Odds struct {
		Msg string `json:"msg"`
	} `json:"odds"`
	Referees []struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Nationality string `json:"nationality"`
	} `json:"referees"`
}

type APIResponse struct {
	Matches []APIMatch `json:"matches"`
}

// NewSyncer creates a new instance of the syncer
func NewSyncer(storage storager, apiBaseURL, apiKey string) *Syncer {
	return &Syncer{
		storage:    storage,
		apiBaseURL: apiBaseURL,
		apiKey:     apiKey,
	}
}

func statusMapper(status string) string {
	// in db, we have only scheduled, ongoing, completed
	switch status {
	case "SCHEDULED", "TIMED":
		return db.MatchStatusScheduled
	case "IN_PLAY", "PAUSED":
		return db.MatchStatusOngoing
	case "FINISHED":
		return db.MatchStatusCompleted
	default:
		return "unknown"
	}
}

func executeWithRateLimit(ctx context.Context, client *http.Client, req *http.Request, limit int, lastRequestTime *time.Time) (*http.Response, error) {
	requestInterval := time.Minute / time.Duration(limit)
	elapsed := time.Since(*lastRequestTime)
	if elapsed < requestInterval {
		time.Sleep(requestInterval - elapsed)
	}
	*lastRequestTime = time.Now()

	var retryCount int
	for {
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to execute request: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			retryCount++
			if retryCount > 3 { // Limit retries
				log.Printf("Too many retries. Skipping request.")
				return nil, fmt.Errorf("too many retries for request")
			}

			resetTime := time.Now().Add(10 * time.Second) // Default retry delay
			if val := resp.Header.Get("X-RequestCounter-Reset"); val != "" {
				// value is a number of seconds to wait
				if waitDuration, err := time.ParseDuration(val + "s"); err == nil {
					resetTime = time.Now().Add(waitDuration)
				}
			}

			waitDuration := time.Until(resetTime)
			log.Printf("Rate limit reached. Retrying after %v...", waitDuration)
			time.Sleep(waitDuration)
			continue
		}

		return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}
}

func (s *Syncer) fetchAPIData(ctx context.Context, endpoint string, lastRequestTime *time.Time, result interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", s.apiBaseURL, endpoint), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-Auth-Token", s.apiKey)

	resp, err := executeWithRateLimit(ctx, client, req, 10, lastRequestTime)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode API response: %w", err)
	}

	return nil
}

func (s *Syncer) SyncTeams(ctx context.Context) error {
	competitions := []string{"PL", "PD", "FL1", "SA", "BL1", "CL"} // England, Spain, France, Italy, Germany, Champions League
	lastRequestTime := time.Now().Add(-time.Minute)

	for _, competition := range competitions {
		log.Printf("Starting team sync for competition: %s", competition)

		var apiResp struct {
			Teams []struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				ShortName string `json:"shortName"`
				Tla       string `json:"tla"`
				Crest     string `json:"crest"`
				Area      struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
					Code string `json:"code"`
				}
			} `json:"teams"`
		}

		if err := s.fetchAPIData(ctx, fmt.Sprintf("/competitions/%s/teams", competition), &lastRequestTime, &apiResp); err != nil {
			log.Printf("Failed to fetch teams for competition %s: %v", competition, err)
			continue
		}

		for _, team := range apiResp.Teams {
			err := s.storage.SaveTeam(ctx, db.Team{
				ID:           fmt.Sprintf("%d", team.ID),
				Name:         team.Name,
				ShortName:    team.ShortName,
				Abbreviation: team.Tla,
				CrestURL:     team.Crest,
				Country:      team.Area.Code,
			})
			if err != nil {
				log.Printf("Failed to save team %d (%s): %v", team.ID, team.Name, err)
			}
		}

		log.Printf("Completed team sync for competition: %s", competition)
	}

	return nil
}

// SyncMatches fetches matches from the API and saves them to the database
func (s *Syncer) SyncMatches(ctx context.Context) error {
	if err := s.SyncTeams(ctx); err != nil {
		return fmt.Errorf("failed to sync teams: %w", err)
	}

	competitions := []string{"PL", "PD", "FL1", "SA", "BL1", "CL"} // Competition codes
	lastRequestTime := time.Now().Add(-time.Minute)

	for _, competition := range competitions {
		log.Printf("Starting sync for competition: %s", competition)

		var apiResp APIResponse
		if err := s.fetchAPIData(ctx, fmt.Sprintf("/competitions/%s/matches", competition), &lastRequestTime, &apiResp); err != nil {
			log.Printf("Failed to fetch matches for competition %s: %v", competition, err)
			continue
		}

		for _, match := range apiResp.Matches {
			if match.HomeTeam.Name == nil || match.AwayTeam.Name == nil {
				log.Printf("Skipping match with missing team names in competition %s", competition)
				continue
			}

			homeTeam, err := s.storage.GetTeamByName(ctx, *match.HomeTeam.Name)
			if err != nil {
				log.Printf("Failed to retrieve home team ID in competition %s: %v", competition, err)
				continue
			}

			awayTeam, err := s.storage.GetTeamByName(ctx, *match.AwayTeam.Name)
			if err != nil {
				log.Printf("Failed to retrieve away team ID in competition %s: %v", competition, err)
				continue
			}

			matchDate := match.UtcDate // Assume valid date parsing here
			err = s.storage.SaveMatch(ctx, db.Match{
				ID:         fmt.Sprintf("%d", match.Id),
				Tournament: match.Competition.Name,
				HomeTeamID: homeTeam.ID,
				AwayTeamID: awayTeam.ID,
				MatchDate:  matchDate,
				Status:     statusMapper(match.Status),
				HomeScore:  match.Score.FullTime.Home,
				AwayScore:  match.Score.FullTime.Away,
			})
			if err != nil {
				log.Printf("Failed to save match %d in competition %s: %v", match.Id, competition, err)
			}
		}

		log.Printf("Completed sync for competition: %s", competition)
	}

	return nil
}

func (s *Syncer) ProcessPredictions(ctx context.Context) error {
	matches, err := s.storage.GetCompletedMatchesWithoutCompletedPredictions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get completed matches: %w", err)
	}

	season, err := s.storage.GetActiveSeason(ctx)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return fmt.Errorf("failed to get active season: %w", err)
	} else if errors.Is(err, db.ErrNotFound) {
		return fmt.Errorf("no active season found")
	}

	for _, match := range matches {
		predictions, err := s.storage.GetPredictionsForMatch(ctx, match.ID)
		if err != nil {
			log.Printf("Failed to fetch predictions for match %d: %v", match.ID, err)
			continue
		}

		for _, prediction := range predictions {
			if match.AwayScore == nil || match.HomeScore == nil {
				log.Printf("Skipping prediction for match %d with missing scores", match.ID)
				continue
			}

			points := calculatePoints(match, prediction)
			if err := s.storage.UpdatePredictionResult(ctx, prediction.MatchID, prediction.UserID, points); err != nil {
				log.Printf("Failed to update prediction result for match %d, user %d: %v", prediction.MatchID, prediction.UserID, err)
			}

			// update user points in leaderboard for current season
			if err := s.storage.UpdateUserLeaderboardPoints(ctx, prediction.UserID, season.ID, points); err != nil {
				log.Printf("Failed to update leaderboard for user %d: %v", prediction.UserID, err)
			}

			// update user stats
			if err := s.storage.UpdateUserPoints(ctx, prediction.UserID, points); err != nil {
				log.Printf("Failed to update user points for user %d: %v", prediction.UserID, err)
			}
		}
	}

	return nil
}

func calculatePoints(match db.Match, prediction db.Prediction) int {
	awayScore := *match.AwayScore
	homeScore := *match.HomeScore

	// If prediction was by exact score
	if prediction.PredictedHomeScore != nil && prediction.PredictedAwayScore != nil {
		predictedHomeScore := *prediction.PredictedHomeScore
		predictedAwayScore := *prediction.PredictedAwayScore

		if homeScore == predictedHomeScore && awayScore == predictedAwayScore {
			return 5
		}
	}

	// If prediction was by outcome
	if prediction.PredictedOutcome != nil {
		outcome := *prediction.PredictedOutcome

		if outcome == db.MatchOutcomeDraw && homeScore == awayScore {
			return 3
		}
		if outcome == db.MatchOutcomeHome && homeScore > awayScore {
			return 3
		}
		if outcome == db.MatchOutcomeAway && awayScore > homeScore {
			return 3
		}
	}

	return 0
}
