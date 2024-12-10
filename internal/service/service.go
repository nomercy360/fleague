package service

import (
	"context"
	"errors"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"sort"
)

// storager interface for database operations
type storager interface {
	Health() (db.HealthStats, error)
	GetLeaderboard(ctx context.Context, leagueID int) ([]db.LeaderboardEntry, error)
	AddPrediction(ctx context.Context, prediction db.Prediction) error
	GetActiveMatches(ctx context.Context) ([]db.Match, error)
	GetUserByChatID(chatID int64) (*db.User, error)
	CreateUser(user db.User) error
	GetTeamByID(ctx context.Context, teamID int) (db.Team, error)
	GetUserPredictionByMatchID(ctx context.Context, uid, matchID int) (*db.Prediction, error)
	SavePrediction(ctx context.Context, prediction db.Prediction) error
	GetMatchByID(ctx context.Context, matchID int) (db.Match, error)
	GetPredictionsByUserID(ctx context.Context, uid int) ([]db.Prediction, error)
}

const userIDContextKey = "user_id"

func GetUserIDFromContext(ctx context.Context) int {
	uid, ok := ctx.Value(userIDContextKey).(int)
	if !ok {
		return 0
	}

	return uid
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

func toMatchResponse(match db.Match, homeTeam db.Team, awayTeam db.Team) contract.MatchResponse {
	return contract.MatchResponse{
		ID:         match.ID,
		Tournament: match.Tournament,
		MatchDate:  match.MatchDate,
		Status:     match.Status,
		HomeTeam: contract.TeamResponse{
			ID:           match.HomeTeamID,
			Name:         homeTeam.Name,
			ShortName:    homeTeam.ShortName,
			CrestURL:     homeTeam.CrestURL,
			Country:      homeTeam.Country,
			Abbreviation: homeTeam.Abbreviation,
		},
		AwayTeam: contract.TeamResponse{
			ID:           match.AwayTeamID,
			Name:         awayTeam.Name,
			ShortName:    awayTeam.ShortName,
			CrestURL:     awayTeam.CrestURL,
			Country:      awayTeam.Country,
			Abbreviation: awayTeam.Abbreviation,
		},
		HomeScore: match.HomeScore,
		AwayScore: match.AwayScore,
	}
}

// GetActiveMatches fetches active matches for a league or all matches
func (s Service) GetActiveMatches(ctx context.Context) ([]contract.MatchResponse, error) {
	res, err := s.storage.GetActiveMatches(ctx)

	if err != nil {
		return nil, err
	}

	uid := GetUserIDFromContext(ctx)

	var matches []contract.MatchResponse
	for _, match := range res {
		homeTeam, err := s.storage.GetTeamByID(ctx, match.HomeTeamID)
		if err != nil {
			return nil, err
		}

		awayTeam, err := s.storage.GetTeamByID(ctx, match.AwayTeamID)
		if err != nil {
			return nil, err
		}

		prediction, err := s.storage.GetUserPredictionByMatchID(ctx, uid, match.ID)
		if err != nil && !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}

		resp := toMatchResponse(match, homeTeam, awayTeam)

		if prediction != nil {
			resp.Prediction = prediction
		}

		matches = append(matches, resp)
	}

	// sort matches by date
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].MatchDate.Before(matches[j].MatchDate)
	})

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
