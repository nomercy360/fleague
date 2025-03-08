package api

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/s3"
	"github.com/user/project/internal/terrors"
	"net/http"
)

// storager interface for database operations
type storager interface {
	Health() (db.HealthStats, error)
	GetLeaderboard(ctx context.Context, seasonID string) ([]db.LeaderboardEntry, error)
	AddPrediction(ctx context.Context, prediction db.Prediction) error
	GetActiveMatches(ctx context.Context, userID string) ([]db.Match, error)
	GetUserByChatID(chatID int64) (db.User, error)
	GetUserByID(id string) (db.User, error)
	GetUserByUsername(uname string) (db.User, error)
	CreateUser(user db.User) error
	GetTeamByID(ctx context.Context, teamID string) (db.Team, error)
	GetUserPredictionByMatchID(ctx context.Context, uid, matchID string) (db.Prediction, error)
	SavePrediction(ctx context.Context, prediction db.Prediction) error
	GetMatchByID(ctx context.Context, matchID string) (db.Match, error)
	GetPredictionsByUserID(ctx context.Context, uid string, onlyCompleted bool) ([]db.Prediction, error)
	GetActiveSeasons(ctx context.Context) ([]db.Season, error)
	UpdateUserPredictionCount(ctx context.Context, userID string) error
	ListUserReferrals(ctx context.Context, userID string) ([]db.User, error)
	UpdateUserPoints(ctx context.Context, userID string, points int, isCorrect bool) error
	ListTeams(ctx context.Context) ([]db.Team, error)
	UpdateUserInformation(ctx context.Context, user db.User) error
	GetUserRank(ctx context.Context, userID string) ([]db.Rank, error)
	GetLastMatchesByTeamID(ctx context.Context, teamID string, limit int) ([]db.Match, error)
	GetPredictionStats(ctx context.Context, userID string) (db.PredictionStats, error)
	GetTodayMostPopularMatch(ctx context.Context) (db.Match, error)
}

type API struct {
	storage  storager
	cfg      Config
	s3Client *s3.Client
}

type Config struct {
	JWTSecret string
	BotToken  string
	AssetsURL string
	OpenAIKey string
}

func New(storage storager, cfg Config, s3Client *s3.Client) *API {
	return &API{
		storage:  storage,
		cfg:      cfg,
		s3Client: s3Client,
	}
}

func (a API) Health(c echo.Context) error {
	stats, err := a.storage.Health()
	if err != nil {
		return terrors.InternalServer(err, "failed to get health stats")
	}

	return c.JSON(http.StatusOK, stats)
}
