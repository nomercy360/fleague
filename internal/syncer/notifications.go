package syncer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/user/project/internal/db"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/user/project/internal/contract"
)

func (s *Syncer) SendWeeklyRecap(ctx context.Context) error {
	if time.Now().Weekday() != time.Monday {
		log.Println("Skipping weekly recap; today is not Monday.")
	}

	users, err := s.storage.GetAllUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	year, week := time.Now().ISOWeek()
	weekNum := fmt.Sprintf("%d-W%d", year, week)

	for _, user := range users {
		recapData, err := s.storage.GetWeeklyRecap(ctx, user.ID)
		if err != nil {
			log.Printf("Failed to fetch recap data for user %s: %v", user.ID, err)
			continue
		}

		alreadySent, err := s.storage.HasNotificationBeenSent(ctx, user.ID, "recap", weekNum)
		if err != nil {
			log.Printf("Failed to check notification log for user %s: %v", user.ID, err)
			continue
		}
		if alreadySent {
			continue
		}

		message := fmt.Sprintf(
			"📊 Weekly Recap:\n- Predictions: %d\n- Wins: %d\n- Losses: %d\n- Win Streak: %d",
			recapData.TotalPredictions, recapData.Wins, recapData.Losses, recapData.CurrentStreak,
		)

		err = s.notifier.SendTextNotification(contract.SendNotificationParams{
			ChatID:  user.ChatID,
			Message: bot.EscapeMarkdown(message),
		})
		if err == nil {
			err := s.storage.LogNotification(ctx, user.ID, "recap", weekNum)
			if err != nil {
				log.Printf("Failed to log notification for user %s: %v", user.ID, err)
			}
		} else {
			log.Printf("Failed to send recap to user %s: %v", user.ID, err)
		}
	}

	return nil
}

func generateMatchReminderText(user db.User, homeTeam db.Team, awayTeam db.Team) string {
	messages := map[string]string{
		"ru": fmt.Sprintf(
			"⚽ Играет ваш любимый клуб!\n\nСделайте прогноз и зарабатывайте очки! 🎯",
		),
		"en": fmt.Sprintf(
			"⚽ Your favorite team is playing!\n\nMake your prediction and earn points! 🎯",
		),
	}

	lang := "en"
	if user.LanguageCode != nil {
		lang = *user.LanguageCode
	}

	if text, exists := messages[lang]; exists {
		return text
	}
	return messages["en"]
}

func (s *Syncer) SendMatchNotification(ctx context.Context) error {
	users, err := s.storage.GetAllUsersWithFavoriteTeam(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch users with favorite teams: %w", err)
	}

	for _, user := range users {
		if user.FavoriteTeamID == nil {
			continue
		}

		matches, err := s.storage.GetMatchesForTeam(ctx, *user.FavoriteTeamID, 14)
		if err != nil {
			log.Printf("Failed to fetch matches for user %s: %v", user.ID, err)
			continue
		}

		for _, match := range matches {
			alreadySent, err := s.storage.HasNotificationBeenSent(ctx, user.ID, "match", match.ID)
			if err != nil || alreadySent {
				continue
			}

			homeTeam, err := s.storage.GetTeamByID(ctx, match.HomeTeamID)
			if err != nil {
				return err
			}

			awayTeam, err := s.storage.GetTeamByID(ctx, match.AwayTeamID)
			if err != nil {
				return err
			}

			imgData, err := fetchPreviewImage(s.cfg.ImagePreviewURL, match, homeTeam, awayTeam)
			if err != nil {
				log.Printf("Failed to fetch preview image for match %s: %v", match.ID, err)
				continue
			}

			err = s.notifier.SendPhotoNotification(contract.SendNotificationParams{
				Image:      imgData,
				ChatID:     user.ChatID,
				Message:    bot.EscapeMarkdown(generateMatchReminderText(user, homeTeam, awayTeam)),
				WebAppURL:  fmt.Sprintf("%s/matches/%s", s.cfg.WebAppURL, match.ID),
				ButtonText: "Make your prediction",
			})

			if err == nil {
				err := s.storage.LogNotification(ctx, user.ID, "match", match.ID)
				if err != nil {
					return err
				}

				log.Printf("Sent notification to user %s about match %s", user.ID, match.ID)
			} else {
				log.Printf("Failed to send notification to user %s: %v", user.ID, err)
			}
		}
	}

	// special notification for chanel about the most popular match
	mostPopularMatch, err := s.storage.GetTodayMostPopularMatch(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch most popular match: %w", err)
	}

	// check if the match is not already sent
	alreadySent, err := s.storage.HasNotificationBeenSent(ctx, fmt.Sprintf("channel:%d", s.cfg.ChannelChatID), "match", mostPopularMatch.ID)
	if err != nil || alreadySent {
		log.Printf("Match %s already sent to channel", mostPopularMatch.ID)
		return nil
	}

	imgData, err := fetchPreviewImage(s.cfg.ImagePreviewURL, mostPopularMatch, mostPopularMatch.HomeTeam, mostPopularMatch.AwayTeam)
	if err != nil {
		log.Printf("Failed to fetch preview image for most popular match: %v", err)
	}

	text := fmt.Sprintf("%s vs %s сегодня в %s, не забудьте сделать прогноз!", mostPopularMatch.HomeTeam.ShortName, mostPopularMatch.AwayTeam.ShortName, mostPopularMatch.MatchDate.Format("15:04"))
	if err = s.notifier.SendPhotoNotification(contract.SendNotificationParams{
		Image:      imgData,
		ChatID:     s.cfg.ChannelChatID,
		Message:    bot.EscapeMarkdown(text),
		BotWebApp:  fmt.Sprintf("%s?startapp=m_%s", s.cfg.BotWebApp, mostPopularMatch.ID),
		ButtonText: "Сделать ставочку",
	}); err != nil {
		log.Printf("Failed to send notification to channel: %v", err)
	}

	if err == nil {
		err := s.storage.LogNotification(ctx, fmt.Sprintf("channel:%d", s.cfg.ChannelChatID), "match", mostPopularMatch.ID)
		if err != nil {
			return err
		}

		log.Printf("Sent notification to channel about match %s", mostPopularMatch.ID)
	}

	return nil
}

type ImageRequest struct {
	Tournament string    `json:"tournament"`
	HomeTeam   string    `json:"homeTeam"`
	AwayTeam   string    `json:"awayTeam"`
	MatchDate  time.Time `json:"matchDate"`
	HomeCrest  string    `json:"homeCrest"`
	AwayCrest  string    `json:"awayCrest"`
}

func fetchPreviewImage(baseUrl string, match db.Match, homeTeam db.Team, awayTeam db.Team) ([]byte, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	u.Path += "/api/football-card"

	body := ImageRequest{
		Tournament: match.Tournament,
		HomeTeam:   homeTeam.ShortName,
		AwayTeam:   awayTeam.ShortName,
		MatchDate:  match.MatchDate,
		HomeCrest:  homeTeam.CrestURL,
		AwayCrest:  awayTeam.CrestURL,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to download image: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to download image, got status code: %d", resp.StatusCode)
		return nil, errors.New(fmt.Sprintf("failed to download image, status code: %d", resp.StatusCode))
	}

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read image data: %s", err)
		return nil, err
	}

	return imgData, nil
}
