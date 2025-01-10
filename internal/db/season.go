package db

import "context"

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
		INSERT INTO seasons (id, name, start_date, end_date, is_active)
		VALUES (?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, season.ID, season.Name, season.StartDate, season.EndDate, season.IsActive)
	return err
}

func (s *Storage) CountSeasons(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM seasons"
	var count int
	err := s.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}
