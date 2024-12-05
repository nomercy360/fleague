package service

import (
	"context"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
)

// storager interface for database operations
type storager interface {
	Health() (db.HealthStats, error)
	GetLeaderboard(ctx context.Context, leagueID int) ([]db.LeaderboardEntry, error)
	AddPrediction(ctx context.Context, prediction db.Prediction) error
	GetActiveMatches(ctx context.Context, leagueID *int) ([]db.Match, error)
	AddUserToLeague(ctx context.Context, leagueID int, userID string) error
	GetUserByChatID(chatID int64) (*db.User, error)
	CreateUser(user db.User) error
}

// Service struct for handling business logic
type Service struct {
	storage  storager
	botToken string
}

// New creates a new Service instance
func New(storage storager, botToken string) *Service {
	return &Service{
		storage:  storage,
		botToken: botToken,
	}
}

// Health checks the database health
func (s Service) Health() (db.HealthStats, error) {
	return s.storage.Health()
}

// GetLeaderboard fetches the leaderboard for a specific league
func (s Service) GetLeaderboard(ctx context.Context, leagueID int) ([]db.LeaderboardEntry, error) {
	return s.storage.GetLeaderboard(ctx, leagueID)
}

// AddPrediction adds a prediction for a user
func (s Service) AddPrediction(ctx context.Context, prediction db.Prediction) error {
	return s.storage.AddPrediction(ctx, prediction)
}

// GetActiveMatches fetches active matches for a league or all matches
func (s Service) GetActiveMatches(ctx context.Context, leagueID *int) ([]contract.MatchResponse, error) {
	res, err := s.storage.GetActiveMatches(ctx, leagueID)

	if err != nil {
		return nil, err
	}

	var matches []contract.MatchResponse
	for _, match := range res {
		matches = append(matches, contract.MatchResponse{
			ID:         match.ID,
			Tournament: match.Tournament,
			HomeTeam:   match.HomeTeam,
			AwayTeam:   match.AwayTeam,
			MatchDate:  match.MatchDate,
			Status:     match.Status,
		})
	}

	return matches, nil
}

// GetUserProfile fetches a user's profile
func (s Service) GetUserProfile(ctx context.Context, userID string) (db.User, error) {
	user, err := s.storage.GetUserByChatID(0)
	if err != nil {
		return db.User{}, err
	}

	return *user, nil
}

// AddUserToLeague adds a user to a specific league
func (s Service) AddUserToLeague(ctx context.Context, leagueID int, userID string) error {
	return s.storage.AddUserToLeague(ctx, leagueID, userID)
}
