package service

import (
	"context"
	"errors"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/s3"
	"sort"
)

// storager interface for database operations
type storager interface {
	Health() (db.HealthStats, error)
	GetLeaderboard(ctx context.Context, leagueID int) ([]db.LeaderboardEntry, error)
	AddPrediction(ctx context.Context, prediction db.Prediction) error
	GetActiveMatches(ctx context.Context) ([]db.Match, error)
	GetUserByChatID(chatID int64) (*db.User, error)
	GetUserByID(id int) (*db.User, error)
	GetUserByUsername(uname string) (*db.User, error)
	CreateUser(user db.User) error
	GetTeamByID(ctx context.Context, teamID int) (db.Team, error)
	GetUserPredictionByMatchID(ctx context.Context, uid, matchID int) (*db.Prediction, error)
	SavePrediction(ctx context.Context, prediction db.Prediction) error
	GetMatchByID(ctx context.Context, matchID int) (db.Match, error)
	GetPredictionsByUserID(ctx context.Context, uid int, onlyCompleted bool) ([]db.Prediction, error)
	GetActiveSeason(ctx context.Context) (db.Season, error)
	UpdateUserPredictionCount(ctx context.Context, userID int) error
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
	s3Client *s3.Client
}

// New creates a new Service instance
func New(storage storager, s3Client *s3.Client, botToken string) *Service {
	return &Service{
		storage:  storage,
		botToken: botToken,
		s3Client: s3Client,
	}
}

// Health checks the database health
func (s Service) Health() (db.HealthStats, error) {
	return s.storage.Health()
}

// GetLeaderboard fetches the leaderboard for a specific league
func (s Service) GetLeaderboard(ctx context.Context) ([]contract.LeaderboardEntry, error) {
	season, err := s.storage.GetActiveSeason(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.storage.GetLeaderboard(ctx, season.ID)

	if err != nil {
		return nil, err
	}

	var leaderboard []contract.LeaderboardEntry
	for _, entry := range res {
		user, err := s.storage.GetUserByID(entry.UserID)
		if err != nil && !errors.Is(err, db.ErrNotFound) {
			return nil, err
		} else if err != nil {
			continue
		}

		userProfile := contract.UserProfile{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			AvatarURL: user.AvatarURL,
		}

		leaderboard = append(leaderboard, contract.LeaderboardEntry{
			User:     userProfile,
			UserID:   entry.UserID,
			Points:   entry.Points,
			SeasonID: entry.SeasonID,
		})
	}

	// sort leaderboard by points
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].Points > leaderboard[j].Points
	})

	return leaderboard, nil
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

// GetActiveSeason fetches the active season
func (s Service) GetActiveSeason(ctx context.Context) (contract.SeasonResponse, error) {
	season, err := s.storage.GetActiveSeason(ctx)
	if err != nil {
		return contract.SeasonResponse{}, err
	}

	return contract.SeasonResponse{
		ID:        season.ID,
		Name:      season.Name,
		IsActive:  season.IsActive,
		StartDate: season.StartDate,
		EndDate:   season.EndDate,
	}, nil
}
