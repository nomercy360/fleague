package syncer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/user/project/internal/contract"
	"log"
	"net/http"
	"time"

	"github.com/user/project/internal/db"
)

type storager interface {
	SaveMatch(ctx context.Context, match db.Match) error
	GetTeamByName(ctx context.Context, name string) (db.Team, error)
	GetTeamByID(ctx context.Context, id string) (db.Team, error)
	GetCompletedMatchesWithoutCompletedPredictions(ctx context.Context) ([]db.Match, error)
	GetPredictionsForMatch(ctx context.Context, matchID string) ([]db.Prediction, error)
	UpdatePredictionResult(ctx context.Context, matchID, userID string, points int) error
	GetActiveSeasons(ctx context.Context) ([]db.Season, error)
	GetActiveSeason(ctx context.Context, seasonType string) (db.Season, error)
	UpdateUserLeaderboardPoints(ctx context.Context, userID, seasonID string, points int) error
	UpdateUserPoints(ctx context.Context, userID string, points int, isCorrect bool) error
	SaveTeam(ctx context.Context, team db.Team) error
	GetUserByID(id string) (db.User, error)
	UpdateUserStreak(ctx context.Context, userID string, currentStreak, longestStreak int) error
	MarkSeasonInactive(ctx context.Context, seasonID string) error
	CreateSeason(ctx context.Context, season db.Season) error
	CountSeasons(ctx context.Context, seasonType string) (int, error)
	GetMatchesForTeam(ctx context.Context, teamID string, hoursAhead int) ([]db.Match, error)
	GetAllUsers(ctx context.Context) ([]db.User, error)
	GetWeeklyRecap(ctx context.Context, userID string) (db.WeeklyRecap, error)
	HasNotificationBeenSent(ctx context.Context, userID, notificationType, relatedID string) (bool, error)
	LogNotification(ctx context.Context, userID, notificationType, relatedID string) error
	GetAllUsersWithFavoriteTeam(ctx context.Context) ([]db.User, error)
	GetActiveMatches(ctx context.Context, userID string) ([]db.Match, error)
	CreateUser(user db.User) error
	GetLastMatchesByTeamID(ctx context.Context, teamID string, limit int) ([]db.Match, error)
	SavePrediction(ctx context.Context, prediction db.Prediction) error
}
type Config struct {
	APIBaseURL      string
	APIKey          string
	WebAppURL       string
	OpenAIKey       string
	ImagePreviewURL string
}
type Syncer struct {
	storage  storager
	notifier noifier
	cfg      Config
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
		HomeWin *float64 `json:"homeWin"`
		Draw    *float64 `json:"draw"`
		AwayWin *float64 `json:"awayWin"`
	} `json:"odds"`
	Referees []struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Nationality string `json:"nationality"`
	} `json:"referees"`
}

type noifier interface {
	SendTextNotification(params contract.SendNotificationParams) error
	SendPhotoNotification(params contract.SendNotificationParams) error
}

type APIResponse struct {
	Matches []APIMatch `json:"matches"`
}

// NewSyncer creates a new instance of the syncer
func NewSyncer(storage storager, notifier noifier, cfg Config) *Syncer {
	return &Syncer{
		storage:  storage,
		notifier: notifier,
		cfg:      cfg,
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", s.cfg.APIBaseURL, endpoint), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-Auth-Token", s.cfg.APIKey)

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
	}

	return nil
}

func (s *Syncer) SyncMatches(ctx context.Context) error {
	competitions := []string{"PL", "PD", "FL1", "SA", "BL1", "CL"} // Competition codes
	lastRequestTime := time.Now().Add(-time.Minute)

	for _, competition := range competitions {
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
				HomeOdds:   match.Odds.HomeWin,
				DrawOdds:   match.Odds.Draw,
				AwayOdds:   match.Odds.AwayWin,
			})

			if err != nil {
				log.Printf("Failed to save match %d in competition %s: %v", match.Id, competition, err)
			}
		}
	}

	return nil
}
