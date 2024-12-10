package db

import "context"

func (s *storage) SaveMatch(ctx context.Context, match Match) error {
	query := `
        INSERT INTO matches (id, tournament, home_team_id, away_team_id, match_date, status, away_score, home_score)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
        tournament = excluded.tournament,
        home_team_id = excluded.home_team_id,
        away_team_id = excluded.away_team_id,
        match_date = excluded.match_date,
        status = excluded.status,
        away_score = excluded.away_score,
        home_score = excluded.home_score`

	_, err := s.db.ExecContext(ctx, query,
		match.ID,
		match.Tournament,
		match.HomeTeamID,
		match.AwayTeamID,
		match.MatchDate,
		match.Status,
		match.AwayScore,
		match.HomeScore,
	)
	return err
}

func (s *storage) GetActiveMatches(ctx context.Context) ([]Match, error) {
	var query string
	var args []interface{}

	query = `
		SELECT
			m.id,
			m.tournament,
			m.home_team_id,
			m.away_team_id,
			m.match_date,
			m.status
		FROM matches m WHERE m.status = 'scheduled'`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var match Match
		if err := rows.Scan(&match.ID, &match.Tournament, &match.HomeTeamID, &match.AwayTeamID, &match.MatchDate, &match.Status); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

func (s *storage) GetMatchByID(ctx context.Context, id int) (Match, error) {
	query := `
		SELECT
			m.id,
			m.tournament,
			m.home_team_id,
			m.away_team_id,
			m.match_date,
			m.status,
			m.home_score,
			m.away_score
		FROM matches m WHERE m.id = ?`

	var match Match
	row := s.db.QueryRowContext(ctx, query, id)

	if err := row.Scan(
		&match.ID,
		&match.Tournament,
		&match.HomeTeamID,
		&match.AwayTeamID,
		&match.MatchDate,
		&match.Status,
		&match.HomeScore,
		&match.AwayScore,
	); err != nil && IsNoRowsError(err) {
		return Match{}, ErrNotFound
	} else if err != nil {
		return Match{}, err
	}

	return match, nil
}
