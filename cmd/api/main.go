package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/handler"
	"github.com/user/project/internal/service"
	"github.com/user/project/internal/syncer"
	"gopkg.in/yaml.v3"
	"log"
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
	FootballAPIKey   string `yaml:"football_api_key"`
}

var (
	footballAPIBaseURL = "https://api.football-data.org/v4"
)

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

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}

func startSyncer(ctx context.Context, sync *syncer.Syncer) {
	// Run the syncer immediately on startup
	log.Println("Starting initial sync process...")
	if err := sync.SyncMatches(ctx); err != nil {
		log.Printf("Initial sync failed: %v", err)
	}

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Starting sync process...")
			if err := sync.SyncMatches(ctx); err != nil {
				log.Printf("Failed to sync matches: %v", err)
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

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	svc := service.New(storage, cfg.TelegramBotToken)

	h := handler.New(svc)

	r.Get("/health", h.Health)
	r.Post("/auth/telegram", h.AuthTelegram)

	r.Route("/v1", func(r chi.Router) {
		setupAPIEndpoints(r, h)
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	done := make(chan bool, 1)

	go gracefulShutdown(server, done)

	sync := syncer.NewSyncer(storage, footballAPIBaseURL, cfg.FootballAPIKey, "CL")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go startSyncer(ctx, sync)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("http server error: %v", err)
	}

	<-done
	log.Println("Graceful shutdown complete.")
}

func setupAPIEndpoints(r chi.Router, h *handler.Handler) {
	r.Get("/matches", h.ListMatches)
}
