package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	telegram "github.com/go-telegram/bot"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/user/project/internal/api"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/notification"
	"github.com/user/project/internal/s3"
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
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	DBPath            string `yaml:"db_path"`
	TelegramBotToken  string `yaml:"telegram_bot_token"`
	JWTSecret         string `yaml:"jwt_secret"`
	OGImagePreviewSVC string `yaml:"og_img_preview_svc"`
	WebAppURL         string `yaml:"web_app_url"`
	AWS               struct {
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
	OpenAIKey         string `yaml:"openai_key"`
	TelegramChannelID int64  `yaml:"telegram_channel_id"`
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
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
			}

			if uid := api.GetContextUserID(c); uid != "" {
				attrs = append(attrs, slog.String("uid", uid))
			}

			if v.Error != nil {
				attrs = append(attrs, slog.String("err", v.Error.Error()))
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR", attrs...)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST", attrs...)
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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}

func startSyncer(ctx context.Context, sync *syncer.Syncer) {
	log.Println("Starting initial sync process...")
	if err := sync.SyncTeams(ctx); err != nil {
		log.Printf("Initial team sync failed: %v", err)
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

			if err := sync.ManageSeasons(ctx); err != nil {
				log.Printf("Initial season sync failed: %v", err)
			}
		case <-ctx.Done():
			log.Println("Stopping syncer...")
			return
		}
	}
}

func startNotificationJob(ctx context.Context, sync *syncer.Syncer) {
	//send when app starts
	if err := sync.SendMatchNotification(ctx); err != nil {
		log.Printf("Failed to send match notifications: %v", err)
	}

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatalf("Failed to load Moscow timezone: %v", err)
	}

	for {
		now := time.Now().In(location) // Get current time in Moscow timezone
		nextRun := time.Date(now.Year(), now.Month(), now.Day(), 06, 02, 0, 0, location)

		// If it's already past 10 AM MSK today, schedule for tomorrow
		if now.After(nextRun) {
			nextRun = nextRun.Add(24 * time.Hour)
		}

		waitDuration := time.Until(nextRun)
		log.Printf("Next notification job scheduled at: %v (Moscow Time)", nextRun)

		timer := time.NewTimer(waitDuration)

		select {
		case <-timer.C:
			log.Println("Running notification job at 10 AM Moscow Time...")

			if err := sync.SendMatchNotification(ctx); err != nil {
				log.Printf("Failed to send match notifications: %v", err)
			}
		case <-ctx.Done():
			log.Println("Stopping notification job...")
			timer.Stop()
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

	apiCfg := api.Config{
		BotToken:  cfg.TelegramBotToken,
		JWTSecret: cfg.JWTSecret,
		AssetsURL: cfg.AssetsURL,
		OpenAIKey: cfg.OpenAIKey,
	}

	s3Client, err := s3.NewS3Client(
		cfg.AWS.AccessKeyID, cfg.AWS.SecretAccessKey, cfg.AWS.Endpoint, cfg.AWS.Bucket)

	if err != nil {
		log.Fatalf("Failed to initialize AWS S3 client: %v\n", err)
	}

	a := api.New(storage, apiCfg, s3Client)

	tmConfig := middleware.TimeoutConfig{
		Timeout: 20 * time.Second,
	}

	e.Use(middleware.TimeoutWithConfig(tmConfig))

	e.POST("/auth/telegram", a.TelegramAuth)

	// Routes
	g := e.Group("/v1")

	authCfg := getAuthConfig(cfg.JWTSecret)

	g.Use(echojwt.WithConfig(authCfg))

	g.GET("/matches", a.ListMatches)
	g.GET("/matches/:id", a.GetMatchByID)
	g.POST("/predictions", a.SavePrediction)
	g.GET("/predictions", a.GetUserPredictions)
	g.GET("/leaderboard", a.GetLeaderboard)
	g.GET("/users/:username", a.GetUserInfo)
	g.GET("/seasons/active", a.GetActiveSeasons)
	g.GET("/referrals", a.ListMyReferrals)
	g.GET("/teams", a.ListTeams)
	g.PUT("/users", a.UpdateUser)
	g.GET("/match/popular", a.GetTodayMostPopularMatch)

	done := make(chan bool, 1)

	go gracefulShutdown(e, done)

	bot, err := telegram.New(cfg.TelegramBotToken)

	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	notifier := notification.NewTelegramNotifier(bot)

	syncerCfg := syncer.Config{
		APIBaseURL:      cfg.FootballAPI.BaseURL,
		APIKey:          cfg.FootballAPI.APIKey,
		WebAppURL:       cfg.WebAppURL,
		OpenAIKey:       cfg.OpenAIKey,
		ImagePreviewURL: cfg.OGImagePreviewSVC,
		ChannelChatID:   cfg.TelegramChannelID,
	}

	sync := syncer.NewSyncer(storage, notifier, syncerCfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go startSyncer(ctx, sync)

	// go startNotificationJob(ctx, sync)

	if err := e.Start(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	<-done
	log.Println("Graceful shutdown complete.")
}
