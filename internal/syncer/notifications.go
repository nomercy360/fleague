package syncer

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"log"
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

func (s *Syncer) SendMatchNotification(ctx context.Context) error {
	users, err := s.storage.GetAllUsersWithFavoriteTeam(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch users with favorite teams: %w", err)
	}

	for _, user := range users {
		if user.FavoriteTeamID == nil {
			continue
		}

		matches, err := s.storage.GetTodayMatchesForTeam(ctx, *user.FavoriteTeamID)
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

			message := fmt.Sprintf("📅 Your favorite team %s is playing today against %s at %s.",
				homeTeam.Name, awayTeam.Name, match.MatchDate.Format("15:04"))

			err = s.notifier.SendTextNotification(contract.SendNotificationParams{
				ChatID:     user.ChatID,
				Message:    bot.EscapeMarkdown(message),
				WebAppURL:  fmt.Sprintf("%s/m_%s", s.webAppURL, match.ID),
				ButtonText: "Make Prediction",
			})

			if err == nil {
				err := s.storage.LogNotification(ctx, user.ID, "match", match.ID)
				if err != nil {
					return err
				}
			} else {
				log.Printf("Failed to send notification to user %s: %v", user.ID, err)
			}
		}
	}

	return nil
}
