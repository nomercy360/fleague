package contract

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/user/project/internal/db"
	"time"
)

type Error struct {
	Message string `json:"message"`
}

type UserAuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID                 string           `json:"id"`
	FirstName          *string          `json:"first_name"`
	LastName           *string          `json:"last_name"`
	Username           string           `json:"username"`
	ChatID             int64            `json:"chat_id"`
	LanguageCode       *string          `json:"language_code"`
	CreatedAt          time.Time        `json:"created_at"`
	TotalPoints        int              `json:"total_points"`
	TotalPredictions   int              `json:"total_predictions"`
	CorrectPredictions int              `json:"correct_predictions"`
	AvatarURL          *string          `json:"avatar_url"`
	ReferredBy         *string          `json:"referred_by"`
	GlobalRank         int              `json:"global_rank"`
	FavoriteTeam       *db.FavoriteTeam `json:"favorite_team"`
	CurrentWinStreak   int              `json:"current_win_streak"`
	LongestWinStreak   int              `json:"longest_win_streak"`
}

type PredictionResponse struct {
	UserID             string        `json:"user_id"`
	MatchID            string        `json:"match_id"`
	PredictedOutcome   *string       `json:"predicted_outcome"`
	PredictedHomeScore *int          `json:"predicted_home_score"`
	PredictedAwayScore *int          `json:"predicted_away_score"`
	PointsAwarded      int           `json:"points_awarded"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
	CompletedAt        *time.Time    `json:"completed_at"`
	Match              MatchResponse `json:"match"`
}

type PredictionRequest struct {
	MatchID            string  `json:"match_id"`
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
	ID         string         `json:"id"`
	Tournament string         `json:"tournament"`
	HomeTeam   db.Team        `json:"home_team"`
	AwayTeam   db.Team        `json:"away_team"`
	MatchDate  time.Time      `json:"match_date"`
	Status     string         `json:"status"`
	AwayScore  *int           `json:"away_score"`
	HomeScore  *int           `json:"home_score"`
	Prediction *db.Prediction `json:"prediction"`
}

type UserProfile struct {
	ID                 string           `json:"id"`
	FirstName          *string          `json:"first_name"`
	LastName           *string          `json:"last_name"`
	Username           string           `json:"username"`
	AvatarURL          *string          `json:"avatar_url"`
	TotalPoints        int              `json:"total_points"`
	TotalPredictions   int              `json:"total_predictions"`
	CorrectPredictions int              `json:"correct_predictions"`
	GlobalRank         int              `json:"global_rank"`
	FavoriteTeam       *db.FavoriteTeam `json:"favorite_team"`
	CurrentWinStreak   int              `json:"current_win_streak"`
	LongestWinStreak   int              `json:"longest_win_streak"`
}

type LeaderboardEntry struct {
	UserID   string      `json:"user_id"`
	Points   int         `json:"points"`
	SeasonID string      `json:"season_id"`
	User     UserProfile `json:"user"`
}

type UserInfoResponse struct {
	User        UserProfile          `json:"user"`
	Predictions []PredictionResponse `json:"predictions"`
}

type SeasonResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UID    string `json:"uid"`
	ChatID int64  `json:"chat_id"`
}

type AuthTelegramRequest struct {
	Query      string  `json:"query"`
	ReferrerID *string `json:"referrer_id"`
}

func (a AuthTelegramRequest) Validate() error {
	if a.Query == "" {
		return fmt.Errorf("query cannot be empty")
	}

	// if not nil, check not empty
	if a.ReferrerID != nil && *a.ReferrerID == "" {
		return fmt.Errorf("referrer id cannot be empty")
	}

	return nil
}

type UpdateUserRequest struct {
	FirstName      *string `json:"first_name"`
	LastName       *string `json:"last_name"`
	FavoriteTeamID *string `json:"favorite_team_id"`
}

func (u UpdateUserRequest) Validate() error {
	if u.FirstName != nil && *u.FirstName == "" {
		return fmt.Errorf("first name cannot be empty")
	}
	if u.LastName != nil && *u.LastName == "" {
		return fmt.Errorf("last name cannot be empty")
	}
	if u.FavoriteTeamID != nil && *u.FavoriteTeamID == "" {
		return fmt.Errorf("favorite team id cannot be empty")
	}
	return nil
}

type SendNotificationParams struct {
	ChatID     int64
	Message    string
	BotWebApp  string
	WebAppURL  string
	Image      []byte
	ButtonText string
}
