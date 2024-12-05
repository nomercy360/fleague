package db

import "context"

func (s *storage) SaveMatch(ctx context.Context, match Match) error {
	query := `
        INSERT INTO matches (id, tournament, home_team, away_team, match_date, status, away_score, home_score)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
        tournament = excluded.tournament,
        home_team = excluded.home_team,
        away_team = excluded.away_team,
        match_date = excluded.match_date,
        status = excluded.status,
        away_score = excluded.away_score,
        home_score = excluded.home_score`

	_, err := s.db.ExecContext(ctx, query,
		match.ID,
		match.Tournament,
		match.HomeTeam,
		match.AwayTeam,
		match.MatchDate,
		match.Status,
		match.AwayScore,
		match.HomeScore,
	)
	return err
}

func (s *storage) GetActiveMatches(ctx context.Context, leagueID *int) ([]Match, error) {
	var query string
	var args []interface{}

	if leagueID != nil {
		query = `
            SELECT
                m.id,
                m.tournament,
                m.home_team,
                m.away_team,
                m.match_date,
                m.status
            FROM matches m
            INNER JOIN league_matches lm ON lm.match_id = m.id
            WHERE lm.league_id = ? AND m.status = 'scheduled'`
		args = append(args, *leagueID)
	} else {
		query = `
            SELECT
                id,
                tournament,
                home_team,
                away_team,
                match_date,
                status
            FROM matches
            WHERE status = 'scheduled'`
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var match Match
		if err := rows.Scan(&match.ID, &match.Tournament, &match.HomeTeam, &match.AwayTeam, &match.MatchDate, &match.Status); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}
