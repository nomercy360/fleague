package contract

import (
	"fmt"
	"github.com/user/project/internal/db"
	"time"
)

type Error struct {
	Message string `json:"message"`
}

type UserAuthResponse struct {
	ID                 int       `json:"id"`
	FirstName          *string   `json:"first_name"`
	LastName           *string   `json:"last_name"`
	Username           string    `json:"username"`
	ChatID             int64     `json:"chat_id"`
	LanguageCode       *string   `json:"language_code"`
	CreatedAt          time.Time `json:"created_at"`
	Token              string    `json:"token"`
	TotalPoints        int       `json:"total_points"`
	TotalPredictions   int       `json:"total_predictions"`
	CorrectPredictions int       `json:"correct_predictions"`
}

type TeamResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ShortName    string `json:"short_name"`
	CrestURL     string `json:"crest_url"`
	Country      string `json:"country"`
	Abbreviation string `json:"abbreviation"`
}

type PredictionResponse struct {
	ID                 int           `json:"id"`
	UserID             int           `json:"user_id"`
	MatchID            int           `json:"match_id"`
	PredictedOutcome   *string       `json:"predicted_outcome"`
	PredictedHomeScore *int          `json:"predicted_home_score"`
	PredictedAwayScore *int          `json:"predicted_away_score"`
	PointsAwarded      int           `json:"points_awarded"`
	CreatedAt          time.Time     `json:"created_at"`
	CompletedAt        *time.Time    `json:"completed_at"`
	Match              MatchResponse `json:"match"`
}

type PredictionRequest struct {
	MatchID            int     `json:"match_id"`
	PredictedOutcome   *string `json:"predicted_outcome"`
	PredictedHomeScore *int    `json:"predicted_home_score"`
	PredictedAwayScore *int    `json:"predicted_away_score"`
}

func (p PredictionRequest) Validate() error {
	if p.PredictedOutcome == nil && p.PredictedHomeScore == nil && p.PredictedAwayScore == nil {
		return fmt.Errorf("at least one of the fields must be filled")
	}
	if p.PredictedOutcome != nil && (p.PredictedHomeScore != nil || p.PredictedAwayScore != nil) {
		return fmt.Errorf("predicted outcome cannot be set with predicted score")
	}
	if p.PredictedOutcome != nil && (*p.PredictedOutcome != db.MatchOutcomeHome && *p.PredictedOutcome != db.MatchOutcomeAway && *p.PredictedOutcome != db.MatchOutcomeDraw) {
		return fmt.Errorf("predicted outcome must be one of home, away, or draw")
	}
	if p.PredictedHomeScore != nil && *p.PredictedHomeScore < 0 {
		return fmt.Errorf("predicted home score must be a positive number")
	}
	if p.PredictedAwayScore != nil && *p.PredictedAwayScore < 0 {
		return fmt.Errorf("predicted away score must be a positive number")
	}
	return nil
}

type MatchResponse struct {
	ID         int            `json:"id"`
	Tournament string         `json:"tournament"`
	HomeTeam   TeamResponse   `json:"home_team"`
	AwayTeam   TeamResponse   `json:"away_team"`
	MatchDate  time.Time      `json:"match_date"`
	Status     string         `json:"status"`
	AwayScore  *int           `json:"away_score"`
	HomeScore  *int           `json:"home_score"`
	Prediction *db.Prediction `json:"prediction"`
}
