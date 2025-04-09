package db

import (
	"context"
	"database/sql"
	"errors"
)

func (s *Storage) GetLeaderboard(ctx context.Context, seasonID string) ([]LeaderboardEntry, error) {
	query := `
        SELECT
            season_id,
            user_id,
            points
        FROM leaderboards
        WHERE season_id = ?
        ORDER BY points DESC LIMIT 100`

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

func (s *Storage) UpdateUserLeaderboardPoints(ctx context.Context, userID, seasonID string, points int) error {
	query := `
		INSERT INTO leaderboards (season_id, user_id, points)
		VALUES (?, ?, ?)
		ON CONFLICT (season_id, user_id)
		DO UPDATE SET points = points + ?`

	_, err := s.db.ExecContext(ctx, query, seasonID, userID, points, points)
	return err
}

type Rank struct {
	SeasonID   string `db:"season_id" json:"season_id"`
	Position   int    `db:"position" json:"position"`
	Points     int    `db:"points" json:"points"`
	SeasonType string `db:"season_type" json:"season_type"`
}

func (s *Storage) GetUserRank(ctx context.Context, userID string) ([]Rank, error) {
	query := `
		WITH active_seasons AS (
			SELECT id, type FROM seasons WHERE is_active = 1
		), ranked_leaderboard AS (
			SELECT
				l.season_id,
				l.user_id,
				l.points,
				RANK() OVER (PARTITION BY l.season_id ORDER BY l.points DESC) AS position,
				s.type
			FROM leaderboards l
			JOIN active_seasons s ON l.season_id = s.id
		)
		SELECT season_id, position, points, type AS season_type
		FROM ranked_leaderboard
		WHERE user_id = ?`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ranks := make([]Rank, 0)
	for rows.Next() {
		var rank Rank
		if err := rows.Scan(
			&rank.SeasonID,
			&rank.Position,
			&rank.Points,
			&rank.SeasonType,
		); err != nil {
			return nil, err
		}
		ranks = append(ranks, rank)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ranks, nil
}

func (s *Storage) GetUserMonthlyRank(ctx context.Context, userID string) (position int, points int, err error) {
	query := `
		WITH ranked_leaderboard AS (
			SELECT
				l.user_id,
				l.points,
				RANK() OVER (ORDER BY l.points DESC) AS position
			FROM leaderboards l
			JOIN seasons s ON l.season_id = s.id
			WHERE s.is_active = 1 AND s.type = 'monthly'
		)
		SELECT position, points
		FROM ranked_leaderboard
		WHERE user_id = ?
	`

	err = s.db.QueryRowContext(ctx, query, userID).Scan(&position, &points)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, nil
		}
		return 0, 0, err
	}

	return position, points, nil
}
