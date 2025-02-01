package db

import (
	"context"
	"time"
)

type Season struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	IsActive  bool      `db:"is_active"`
	Type      string    `db:"type"`
}

const (
	SeasonTypeMonthly  = "monthly"
	SeasonTypeFootball = "football"
)

func (s *Storage) MarkSeasonInactive(ctx context.Context, seasonID string) error {
	query := `
		UPDATE seasons
		SET is_active = 0
		WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, seasonID)
	return err
}

func (s *Storage) CreateSeason(ctx context.Context, season Season) error {
	query := `
		INSERT INTO seasons (id, name, start_date, end_date, is_active, type)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, season.ID, season.Name, season.StartDate, season.EndDate, season.IsActive, season.Type)
	return err
}

func (s *Storage) CountSeasons(ctx context.Context, t string) (int, error) {
	query := "SELECT COUNT(*) FROM seasons WHERE type = ?"
	var count int
	err := s.db.QueryRowContext(ctx, query, t).Scan(&count)
	return count, err
}

func (s *Storage) GetActiveSeasons(ctx context.Context) ([]Season, error) {
	resp := make([]Season, 0)

	query := `
		SELECT
			id,
			name,
			start_date,
			end_date,
			is_active,
			type
		FROM seasons
		WHERE is_active = 1`

	var season Season
	res, err := s.db.QueryContext(ctx, query)
	if err != nil && IsNoRowsError(err) {
		return resp, ErrNotFound
	} else if err != nil {
		return resp, err
	}

	for res.Next() {
		err := res.Scan(
			&season.ID,
			&season.Name,
			&season.StartDate,
			&season.EndDate,
			&season.IsActive,
			&season.Type,
		)
		if err != nil {
			return resp, err
		}
		resp = append(resp, season)
	}

	return resp, nil
}

func (s *Storage) GetActiveSeason(ctx context.Context, t string) (Season, error) {
	query := `
		SELECT
			id,
			name,
			start_date,
			end_date,
			is_active,
			type
		FROM seasons
		WHERE is_active = 1 AND type = ?`

	var season Season
	err := s.db.QueryRowContext(ctx, query, t).Scan(
		&season.ID,
		&season.Name,
		&season.StartDate,
		&season.EndDate,
		&season.IsActive,
		&season.Type,
	)

	if err != nil && IsNoRowsError(err) {
		return Season{}, ErrNotFound
	} else if err != nil {
		return Season{}, err
	}

	return season, nil
}
