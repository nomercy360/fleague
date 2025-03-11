package syncer

import (
	"context"
	"errors"
	"fmt"
	telegram "github.com/go-telegram/bot"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"log"
)

func (s *Syncer) ProcessPredictions(ctx context.Context) error {
	matches, err := s.storage.GetCompletedMatchesWithoutCompletedPredictions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get completed matches: %w", err)
	}

	seasons, err := s.storage.GetActiveSeasons(ctx)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return fmt.Errorf("failed to get active season: %w", err)
	} else if errors.Is(err, db.ErrNotFound) {
		return fmt.Errorf("no active season found")
	}

	for _, match := range matches {
		// Retrieve all predictions for the current match
		predictions, err := s.storage.GetPredictionsForMatch(ctx, match.ID)
		if err != nil {
			log.Printf("Failed to fetch predictions for match %s: %v", match.ID, err)
			continue
		}

		for _, prediction := range predictions {
			// Ensure match scores are available
			if match.AwayScore == nil || match.HomeScore == nil {
				log.Printf("Skipping prediction for match %s with missing scores", match.ID)
				continue
			}

			// Calculate base points based on prediction correctness
			basePoints := calculateBasePoints(match, prediction)
			isExactCorrect := basePoints == 7
			isOutcomeCorrect := basePoints == 3

			// Determine if the prediction was correct
			isCorrect := isExactCorrect || isOutcomeCorrect

			// Fetch the user to update streaks
			user, err := s.storage.GetUserByID(prediction.UserID)
			if err != nil {
				log.Printf("Failed to fetch user %s: %v", prediction.UserID, err)
				continue
			}

			// Initialize bonus points
			bonusPoints := 0

			if isCorrect {
				// Increment the user's current streak
				user.CurrentWinStreak += 1

				// Update the longest streak if necessary
				if user.CurrentWinStreak > user.LongestWinStreak {
					user.LongestWinStreak = user.CurrentWinStreak
				}

				// Calculate bonus based on the new streak
				bonusPoints = calculateBonus(user.CurrentWinStreak)
			} else {
				// Reset the user's current streak
				user.CurrentWinStreak = 0
			}

			// Calculate total points (base + bonus)
			totalPoints := basePoints + bonusPoints

			// Update the prediction result with total points
			err = s.storage.UpdatePredictionResult(ctx, prediction.MatchID, prediction.UserID, totalPoints)
			if err != nil {
				log.Printf("Failed to update prediction result for match %s, user %s: %v", prediction.MatchID, prediction.UserID, err)
				continue
			}

			for _, season := range seasons {
				err = s.storage.UpdateUserLeaderboardPoints(ctx, prediction.UserID, season.ID, totalPoints)
				if err != nil {
					log.Printf("Failed to update leaderboard for user %s: %v", prediction.UserID, err)
					continue
				}
			}

			err = s.storage.UpdateUserPoints(ctx, prediction.UserID, isCorrect)
			if err != nil {
				log.Printf("Failed to update user points for user %s: %v", prediction.UserID, err)
				continue
			}

			err = s.storage.UpdateUserStreak(ctx, user.ID, user.CurrentWinStreak, user.LongestWinStreak)
			if err != nil {
				log.Printf("Failed to update streak for user %s: %v", user.ID, err)
				continue
			}

			// go s.notifyUser(ctx, user, user.CurrentWinStreak, bonusPoints)
		}
	}

	return nil
}

func calculateBonus(currentStreak int) int {
	switch {
	case currentStreak >= 11:
		return 10
	case currentStreak >= 7:
		return 5
	case currentStreak >= 4:
		return 2
	default:
		return 0
	}
}

func calculateBasePoints(match db.Match, prediction db.Prediction) int {
	awayScore := *match.AwayScore
	homeScore := *match.HomeScore

	// Exact score prediction
	if prediction.PredictedHomeScore != nil && prediction.PredictedAwayScore != nil {
		predictedHomeScore := *prediction.PredictedHomeScore
		predictedAwayScore := *prediction.PredictedAwayScore

		if homeScore == predictedHomeScore && awayScore == predictedAwayScore {
			return 7
		}
	}

	// Outcome prediction
	if prediction.PredictedOutcome != nil {
		outcome := *prediction.PredictedOutcome

		if outcome == db.MatchOutcomeDraw && homeScore == awayScore {
			return 3
		}
		if outcome == db.MatchOutcomeHome && homeScore > awayScore {
			return 3
		}
		if outcome == db.MatchOutcomeAway && awayScore > homeScore {
			return 3
		}
	}

	return 0
}

func (s *Syncer) notifyUser(ctx context.Context, user db.User, streak int, bonusPoints int) {
	if streak < 4 && bonusPoints == 0 {
		return
	}

	message := fmt.Sprintf("🎉 Congratulations! You've achieved a streak of %d correct predictions and earned an extra %d points!", streak, bonusPoints)
	err := s.notifier.SendTextNotification(contract.SendNotificationParams{
		ChatID:  user.ChatID,
		Message: telegram.EscapeMarkdown(message),
	})

	if err != nil {
		log.Printf("Failed to send notification to user %s: %v", user.ID, err)
	}
}
