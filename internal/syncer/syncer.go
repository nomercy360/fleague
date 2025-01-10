package syncer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	telegram "github.com/go-telegram/bot"
	"github.com/user/project/internal/contract"
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
	UpdateUserPoints(ctx context.Context, userID string, points int, isCorrect bool) error
	SaveTeam(ctx context.Context, team db.Team) error
	GetUserByID(id string) (db.User, error)
	UpdateUserStreak(ctx context.Context, userID string, currentStreak, longestStreak int) error
	MarkSeasonInactive(ctx context.Context, seasonID string) error
	CreateSeason(ctx context.Context, season db.Season) error
	CountSeasons(ctx context.Context) (int, error)
}

type Syncer struct {
	storage    storager
	notifier   noifier
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

type noifier interface {
	SendTextNotification(params contract.SendNotificationParams) error
}

type APIResponse struct {
	Matches []APIMatch `json:"matches"`
}

// NewSyncer creates a new instance of the syncer
func NewSyncer(storage storager, notifier noifier, apiBaseURL, apiKey string) *Syncer {
	return &Syncer{
		storage:    storage,
		notifier:   notifier,
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
	}

	return nil
}

// SyncMatches fetches matches from the API and saves them to the database
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
			})
			if err != nil {
				log.Printf("Failed to save match %d in competition %s: %v", match.Id, competition, err)
			}
		}
	}

	return nil
}

func (s *Syncer) ProcessPredictions(ctx context.Context) error {
	// Retrieve all completed matches that have predictions pending evaluation
	matches, err := s.storage.GetCompletedMatchesWithoutCompletedPredictions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get completed matches: %w", err)
	}

	// Fetch the active season
	season, err := s.storage.GetActiveSeason(ctx)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return fmt.Errorf("failed to get active season: %w", err)
	} else if errors.Is(err, db.ErrNotFound) {
		return fmt.Errorf("no active season found")
	}

	for _, match := range matches {
		// Retrieve all predictions for the current match
		predictions, err := s.storage.GetPredictionsForMatch(ctx, match.ID)
		if err != nil {
			log.Printf("Failed to fetch predictions for match %s: %v", match.ID, err)
			continue
		}

		for _, prediction := range predictions {
			// Ensure match scores are available
			if match.AwayScore == nil || match.HomeScore == nil {
				log.Printf("Skipping prediction for match %s with missing scores", match.ID)
				continue
			}

			// Calculate base points based on prediction correctness
			basePoints := calculateBasePoints(match, prediction)
			isExactCorrect := basePoints == 7
			isOutcomeCorrect := basePoints == 3

			// Determine if the prediction was correct
			isCorrect := isExactCorrect || isOutcomeCorrect

			// Fetch the user to update streaks
			user, err := s.storage.GetUserByID(prediction.UserID)
			if err != nil {
				log.Printf("Failed to fetch user %s: %v", prediction.UserID, err)
				continue
			}

			// Initialize bonus points
			bonusPoints := 0

			if isCorrect {
				// Increment the user's current streak
				user.CurrentWinStreak += 1

				// Update the longest streak if necessary
				if user.CurrentWinStreak > user.LongestWinStreak {
					user.LongestWinStreak = user.CurrentWinStreak
				}

				// Calculate bonus based on the new streak
				bonusPoints = calculateBonus(user.CurrentWinStreak)
			} else {
				// Reset the user's current streak
				user.CurrentWinStreak = 0
			}

			// Calculate total points (base + bonus)
			totalPoints := basePoints + bonusPoints

			// Update the prediction result with total points
			err = s.storage.UpdatePredictionResult(ctx, prediction.MatchID, prediction.UserID, totalPoints)
			if err != nil {
				log.Printf("Failed to update prediction result for match %s, user %s: %v", prediction.MatchID, prediction.UserID, err)
				continue
			}

			err = s.storage.UpdateUserLeaderboardPoints(ctx, prediction.UserID, season.ID, totalPoints)
			if err != nil {
				log.Printf("Failed to update leaderboard for user %s: %v", prediction.UserID, err)
				continue
			}

			err = s.storage.UpdateUserPoints(ctx, prediction.UserID, totalPoints, isCorrect)
			if err != nil {
				log.Printf("Failed to update user points for user %s: %v", prediction.UserID, err)
				continue
			}

			err = s.storage.UpdateUserStreak(ctx, user.ID, user.CurrentWinStreak, user.LongestWinStreak)
			if err != nil {
				log.Printf("Failed to update streak for user %s: %v", user.ID, err)
				continue
			}

			// go s.notifyUser(ctx, user, user.CurrentWinStreak, bonusPoints)
		}
	}

	return nil
}

func calculateBonus(currentStreak int) int {
	switch {
	case currentStreak >= 11:
		return 10
	case currentStreak >= 7:
		return 5
	case currentStreak >= 4:
		return 2
	default:
		return 0
	}
}

func calculateBasePoints(match db.Match, prediction db.Prediction) int {
	awayScore := *match.AwayScore
	homeScore := *match.HomeScore

	// Exact score prediction
	if prediction.PredictedHomeScore != nil && prediction.PredictedAwayScore != nil {
		predictedHomeScore := *prediction.PredictedHomeScore
		predictedAwayScore := *prediction.PredictedAwayScore

		if homeScore == predictedHomeScore && awayScore == predictedAwayScore {
			return 7
		}
	}

	// Outcome prediction
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

func (s *Syncer) notifyUser(ctx context.Context, user db.User, streak int, bonusPoints int) {
	if streak < 4 && bonusPoints == 0 {
		return
	}

	message := fmt.Sprintf("🎉 Congratulations! You've achieved a streak of %d correct predictions and earned an extra %d points!", streak, bonusPoints)
	err := s.notifier.SendTextNotification(contract.SendNotificationParams{
		ChatID:  user.ChatID,
		Message: telegram.EscapeMarkdown(message),
	})

	if err != nil {
		log.Printf("Failed to send notification to user %s: %v", user.ID, err)
	}
}
