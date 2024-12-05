package db

import "context"

func (s *storage) GetLeaderboard(ctx context.Context, leagueID int) ([]LeaderboardEntry, error) {
	query := `
        SELECT
            id,
            league_id,
            user_id,
            points
        FROM leaderboards
        WHERE league_id = ?
        ORDER BY points DESC`

	rows, err := s.db.QueryContext(ctx, query, leagueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []LeaderboardEntry
	for rows.Next() {
		var entry LeaderboardEntry
		if err := rows.Scan(&entry.ID, &entry.LeagueID, &entry.UserID, &entry.Points); err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return leaderboard, nil
}
