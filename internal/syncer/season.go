package syncer

import (
	"context"
	"errors"
	"fmt"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/nanoid"
	"log"
	"time"
)

func (s *Syncer) ManageSeasons(ctx context.Context) error {
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	activeSeason, err := s.storage.GetActiveSeason(ctx)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return fmt.Errorf("failed to get active season: %w", err)
	}

	newSeasonRequired := false
	if errors.Is(err, db.ErrNotFound) {
		newSeasonRequired = true
	} else {
		if !(activeSeason.StartDate.Equal(firstOfMonth) && activeSeason.EndDate.Equal(lastOfMonth)) {
			newSeasonRequired = true
		}
	}

	if newSeasonRequired {
		if !errors.Is(err, db.ErrNotFound) {
			err := s.storage.MarkSeasonInactive(ctx, activeSeason.ID)
			if err != nil {
				return fmt.Errorf("failed to mark previous season inactive: %w", err)
			}
		}

		seasonCount, err := s.storage.CountSeasons(ctx)
		if err != nil {
			return fmt.Errorf("failed to count existing seasons: %w", err)
		}

		newSeasonName := fmt.Sprintf("S%d", seasonCount+1)

		newSeason := db.Season{
			ID:        nanoid.Must(),
			Name:      newSeasonName,
			StartDate: firstOfMonth,
			EndDate:   lastOfMonth,
			IsActive:  true,
		}

		if err := s.storage.CreateSeason(ctx, newSeason); err != nil {
			return fmt.Errorf("failed to create new season: %w", err)
		}

		log.Printf("New season created: %s (%s - %s)", newSeason.Name, newSeason.StartDate, newSeason.EndDate)
	}

	return nil
}
