package syncer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/user/project/internal/db"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/user/project/internal/contract"
)

func generateMatchReminderText(user db.User, homeTeam db.Team, awayTeam db.Team) string {
	messages := map[string]string{
		"ru": fmt.Sprintf(
			"⚽ Играет ваш любимый клуб!\n\nКак насчет проверить свое футбольное чутье?",
		),
		"en": fmt.Sprintf(
			"⚽ Your favorite team is playing!\n\nHow about testing your football intuition?",
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

type RecapMatch struct {
	HomeTeam           string  `json:"homeTeam"`
	AwayTeam           string  `json:"awayTeam"`
	HomeCrest          string  `json:"homeCrest"`
	AwayCrest          string  `json:"awayCrest"`
	Score              string  `json:"score"`
	PredictedOutcome   *string `json:"predictedOutcome,omitempty"`
	PredictedHomeScore *int    `json:"predictedHomeScore,omitempty"`
	PredictedAwayScore *int    `json:"predictedAwayScore,omitempty"`
	IsCorrect          bool    `json:"isCorrect"`
	Points             int     `json:"points"`
	Popularity         float64 `json:"popularity"`
}

type WeeklyRecapImageRequest struct {
	WeekStartDate       time.Time    `json:"weekStartDate"`
	TotalPoints         int          `json:"totalPoints"`
	LeaderboardPosition int          `json:"leaderboardPosition"`
	Matches             []RecapMatch `json:"matches"`
}

func (s *Syncer) SendWeeklyRecap(ctx context.Context) error {
	users, err := s.storage.GetAllUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	year, week := time.Now().AddDate(0, 0, -7).ISOWeek() // Смотрим прошлую неделю
	weekNum := fmt.Sprintf("%d-W%d", year, week)
	startOfWeek := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	for _, user := range users {
		alreadySent, err := s.storage.HasNotificationBeenSent(ctx, user.ID, "recap", weekNum)
		if err != nil {
			log.Printf("Failed to check notification log for user %s: %v", user.ID, err)
			continue
		}
		if alreadySent {
			continue
		}

		predictions, err := s.storage.GetPredictionsByUserID(ctx, user.ID,
			db.WithOnlyCompleted(),
			db.WithStartTime(startOfWeek),
			db.WithEndTime(endOfWeek),
			db.WithLimit(5),
		)

		if err != nil {
			log.Printf("Failed to fetch predictions for user %s: %v", user.ID, err)
			continue
		}

		// if predictions are empty, skip the user
		if len(predictions) == 0 {
			continue
		}

		// Получаем позицию в лидерборде
		leaderboardPos, totalPoints, err := s.storage.GetUserMonthlyRank(ctx, user.ID)
		if err != nil {
			log.Printf("Failed to fetch monthly rank for user %s: %v", user.ID, err)
			continue
		}

		matches := make([]RecapMatch, 0, len(predictions))
		for _, pred := range predictions {
			match, err := s.storage.GetMatchByID(ctx, pred.MatchID)
			if err != nil {
				log.Printf("Failed to fetch match %s: %v", pred.MatchID, err)
				continue
			}

			matches = append(matches, RecapMatch{
				HomeTeam:           match.HomeTeam.ShortName,
				AwayTeam:           match.AwayTeam.ShortName,
				HomeCrest:          match.HomeTeam.CrestURL,
				AwayCrest:          match.AwayTeam.CrestURL,
				Score:              fmt.Sprintf("%d:%d", *match.HomeScore, *match.AwayScore),
				PredictedOutcome:   pred.PredictedOutcome,
				PredictedHomeScore: pred.PredictedHomeScore,
				PredictedAwayScore: pred.PredictedAwayScore,
				IsCorrect:          pred.PointsAwarded > 0,
				Points:             pred.PointsAwarded,
				Popularity:         match.Popularity,
			})
		}

		// sort by match popularity from low to high
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Popularity < matches[j].Popularity
		})

		imgRequest := WeeklyRecapImageRequest{
			WeekStartDate:       startOfWeek,
			TotalPoints:         totalPoints,
			LeaderboardPosition: leaderboardPos,
			Matches:             matches,
		}

		imgData, err := fetchWeeklyRecapImage(s.cfg.ImagePreviewURL, imgRequest)
		if err != nil {
			log.Printf("Failed to fetch recap image for user %s: %v", user.ID, err)
			continue
		}

		// Формируем текстовое сообщение
		lang := "en"
		if user.LanguageCode != nil {
			lang = *user.LanguageCode
		}
		message := generateWeeklyRecapText(lang, totalPoints, leaderboardPos, len(matches))

		buttonText := "View Details"
		if lang == "ru" {
			buttonText = "Подробнее"
		}

		// Отправляем уведомление
		err = s.notifier.SendPhotoNotification(contract.SendNotificationParams{
			Image:      imgData,
			ChatID:     user.ChatID,
			Message:    bot.EscapeMarkdown(message),
			WebAppURL:  fmt.Sprintf("%s/weekly-recap?week=%s", s.cfg.WebAppURL, weekNum),
			ButtonText: buttonText,
		})
		if err == nil {
			err := s.storage.LogNotification(ctx, user.ID, "recap", weekNum)
			if err != nil {
				log.Printf("Failed to log notification for user %s: %v", user.ID, err)
			}
			log.Printf("Sent weekly recap to user %s for week %s", user.ID, weekNum)
		} else {
			log.Printf("Failed to send recap to user %s: %v", user.ID, err)
		}
	}

	return nil
}

func generateWeeklyRecapText(lang string, totalPoints, leaderboardPos, matchCount int) string {
	messages := map[string]string{
		"ru": fmt.Sprintf(
			"Привет! Итоги недели:\n- Прогнозов сделано: %d\n- Очков набрал: %d\n- Место в лидерборде: #%d",
			matchCount, totalPoints, leaderboardPos,
		),
		"en": fmt.Sprintf(
			"Hello! Weekly recap:\n- Predictions made: %d\n- Points earned: %d\n- Leaderboard position: #%d",
			matchCount, totalPoints, leaderboardPos,
		),
	}
	if text, exists := messages[lang]; exists {
		return text
	}
	return messages["en"]
}

type ImageRequest struct {
	Tournament string    `json:"tournament"`
	HomeTeam   string    `json:"homeTeam"`
	AwayTeam   string    `json:"awayTeam"`
	MatchDate  time.Time `json:"matchDate"`
	HomeCrest  string    `json:"homeCrest"`
	AwayCrest  string    `json:"awayCrest"`
}

func fetchImage(baseURL, endpoint string, body interface{}) ([]byte, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	u.Path += endpoint

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
		log.Printf("Failed to download image from %s: %s", endpoint, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to download image from %s, got status code: %d", endpoint, resp.StatusCode)
		return nil, fmt.Errorf("failed to download image from %s, status code: %d", endpoint, resp.StatusCode)
	}

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read image data from %s: %s", endpoint, err)
		return nil, err
	}

	return imgData, nil
}

// fetchWeeklyRecapImage использует общую функцию для получения изображения рекапа
func fetchWeeklyRecapImage(baseURL string, request WeeklyRecapImageRequest) ([]byte, error) {
	return fetchImage(baseURL, "/api/weekly-recap", request)
}

func fetchPreviewImage(baseURL string, match db.Match, homeTeam db.Team, awayTeam db.Team) ([]byte, error) {
	body := ImageRequest{
		Tournament: match.Tournament,
		HomeTeam:   homeTeam.ShortName,
		AwayTeam:   awayTeam.ShortName,
		MatchDate:  match.MatchDate,
		HomeCrest:  homeTeam.CrestURL,
		AwayCrest:  awayTeam.CrestURL,
	}
	return fetchImage(baseURL, "/api/football-card", body)
}
