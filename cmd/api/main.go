package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/handler"
	"github.com/user/project/internal/s3"
	"github.com/user/project/internal/service"
	"github.com/user/project/internal/syncer"
	"github.com/user/project/internal/terrors"
	"gopkg.in/yaml.v3"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	Host             string `yaml:"host"`
	Port             int    `yaml:"port"`
	DBPath           string `yaml:"db_path"`
	TelegramBotToken string `yaml:"telegram_bot_token"`
	JWTSecret        string `yaml:"jwt_secret"`
	MetaFetchURL     string `yaml:"meta_fetch_url"`
	AWS              struct {
		AccessKeyID     string `yaml:"access_key_id"`
		SecretAccessKey string `yaml:"secret_access_key"`
		Endpoint        string `yaml:"endpoint"`
		Bucket          string `yaml:"bucket"`
	} `yaml:"aws"`
	AssetsURL   string `yaml:"assets_url"`
	FootballAPI struct {
		APIKey  string `yaml:"api_key"`
		BaseURL string `yaml:"base_url"`
	} `yaml:"football_api"`
}

func ReadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return &cfg, nil
}

func ValidateConfig(cfg *Config) error {
	validate := validator.New()
	return validate.Struct(cfg)
}

func getLoggerMiddleware(logger *slog.Logger) middleware.RequestLoggerConfig {
	return middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}
}

func getServerErrorHandler(e *echo.Echo) func(err error, context2 echo.Context) {
	return func(err error, c echo.Context) {
		var (
			code = http.StatusInternalServerError
			msg  interface{}
		)

		var he *echo.HTTPError
		var terror *terrors.Error
		switch {
		case errors.As(err, &he):
			code = he.Code
			msg = he.Message
		case errors.As(err, &terror):
			code = terror.Code
			msg = terror.Message
		default:
			msg = err.Error()
		}

		if _, ok := msg.(string); ok {
			msg = map[string]interface{}{"error": msg}
		}

		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				err = c.NoContent(code)
			} else {
				err = c.JSON(code, msg)
			}

			if err != nil {
				e.Logger.Error(err)
			}
		}
	}
}

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func getAuthConfig(secret string) echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(_ echo.Context) jwt.Claims {
			return new(contract.JWTClaims)
		},
		SigningKey:             []byte(secret),
		ContinueOnIgnoredError: true,
		ErrorHandler: func(c echo.Context, err error) error {
			var extErr *echojwt.TokenExtractionError
			if !errors.As(err, &extErr) {
				return echo.NewHTTPError(http.StatusUnauthorized, "auth is invalid")
			}

			claims := &contract.JWTClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 30)),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			c.Set("user", token)

			if claims.UID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "auth is invalid")
			}

			return nil
		},
	}
}

func gracefulShutdown(e *echo.Echo, done chan<- bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}

func startSyncer(ctx context.Context, sync *syncer.Syncer) {
	log.Println("Starting initial sync process...")
	if err := sync.SyncMatches(ctx); err != nil {
		log.Printf("Initial sync failed: %v", err)
	}

	if err := sync.ProcessPredictions(ctx); err != nil {
		log.Printf("Initial prediction processing failed: %v", err)
	}

	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Starting sync process...")
			if err := sync.SyncMatches(ctx); err != nil {
				log.Printf("Failed to sync matches: %v", err)
			}

			if err := sync.ProcessPredictions(ctx); err != nil {
				log.Printf("Failed to process predictions: %v", err)
			}
		case <-ctx.Done():
			log.Println("Stopping syncer...")
			return
		}
	}
}

func main() {
	configFilePath := "config.yml"
	configFilePathEnv := os.Getenv("CONFIG_FILE_PATH")
	if configFilePathEnv != "" {
		configFilePath = configFilePathEnv
	}

	cfg, err := ReadConfig(configFilePath)
	if err != nil {
		log.Fatalf("error reading configuration: %v", err)
	}

	if err := ValidateConfig(cfg); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	storage, err := db.ConnectDB(cfg.DBPath)

	if err != nil {
		log.Fatalf("failed to create storage: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	e.Use(middleware.RequestLoggerWithConfig(getLoggerMiddleware(logger)))

	e.HTTPErrorHandler = getServerErrorHandler(e)

	e.Validator = &customValidator{validator: validator.New()}

	apiCfg := service.Config{
		BotToken:  cfg.TelegramBotToken,
		JWTSecret: cfg.JWTSecret,
		AssetsURL: cfg.AssetsURL,
	}

	s3Client, err := s3.NewS3Client(
		cfg.AWS.AccessKeyID, cfg.AWS.SecretAccessKey, cfg.AWS.Endpoint, cfg.AWS.Bucket)

	if err != nil {
		log.Fatalf("Failed to initialize AWS S3 client: %v\n", err)
	}

	svc := service.New(storage, apiCfg, s3Client)

	h := handler.New(svc)

	tmConfig := middleware.TimeoutConfig{
		Timeout: 20 * time.Second,
	}

	e.Use(middleware.TimeoutWithConfig(tmConfig))

	e.POST("/auth/telegram", h.AuthTelegram)

	// Routes
	g := e.Group("/v1")

	authCfg := getAuthConfig(cfg.JWTSecret)

	g.Use(echojwt.WithConfig(authCfg))

	g.GET("/matches", h.ListMatches)
	g.POST("/predictions", h.SavePrediction)
	g.GET("/predictions", h.GetUserPredictions)
	g.GET("/leaderboard", h.GetLeaderboard)
	g.GET("/users/:username", h.GetUserInfo)
	g.GET("/seasons/active", h.GetActiveSeason)
	g.GET("/referrals", h.ListMyReferrals)

	done := make(chan bool, 1)

	go gracefulShutdown(e, done)

	sync := syncer.NewSyncer(storage, cfg.FootballAPI.BaseURL, cfg.FootballAPI.APIKey)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go startSyncer(ctx, sync)

	if err := e.Start(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	<-done
	log.Println("Graceful shutdown complete.")
}
