package db

import "context"

func (s *Storage) GetLeaderboard(ctx context.Context, seasonID string) ([]LeaderboardEntry, error) {
	query := `
        SELECT
            season_id,
            user_id,
            points
        FROM leaderboards
        WHERE season_id = ?
        ORDER BY points DESC LIMIT 30`

	rows, err := s.db.QueryContext(ctx, query, seasonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []LeaderboardEntry
	for rows.Next() {
		var entry LeaderboardEntry
		if err := rows.Scan(&entry.SeasonID, &entry.UserID, &entry.Points); err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return leaderboard, nil
}

func (s *Storage) GetActiveSeason(ctx context.Context) (Season, error) {
	query := `
		SELECT
			id,
			name,
			start_date,
			end_date,
			is_active
		FROM seasons
		WHERE is_active = 1`

	var season Season
	if err := s.db.QueryRowContext(ctx, query).Scan(&season.ID, &season.Name, &season.StartDate, &season.EndDate, &season.IsActive); err != nil {
		return Season{}, err
	}

	return season, nil
}

func (s *Storage) UpdateUserLeaderboardPoints(ctx context.Context, userID, seasonID string, points int) error {
	query := `
		INSERT INTO leaderboards (season_id, user_id, points)
		VALUES (?, ?, ?)
		ON CONFLICT (season_id, user_id)
		DO UPDATE SET points = points + ?`

	_, err := s.db.ExecContext(ctx, query, seasonID, userID, points, points)
	return err
}
