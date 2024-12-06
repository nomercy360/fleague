package db

import "context"

func (s *storage) GetTeamByName(ctx context.Context, name string) (Team, error) {
	query := `
		SELECT
			id,
			name,
			short_name,
			crest_url,
			country,
			abbreviation
		FROM teams
		WHERE name = ?`

	var team Team
	err := s.db.QueryRowContext(ctx, query, name).Scan(
		&team.ID,
		&team.Name,
		&team.ShortName,
		&team.CrestURL,
		&team.Country,
		&team.Abbreviation,
	)

	if err != nil && IsNoRowsError(err) {
		return Team{}, ErrNotFound
	} else if err != nil {
		return Team{}, err
	}

	return team, nil
}

func (s *storage) GetTeamByID(ctx context.Context, id int) (Team, error) {
	query := `
		SELECT
			id,
			name,
			short_name,
			crest_url,
			country,
			abbreviation
		FROM teams
		WHERE id = ?`

	var team Team
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&team.ID,
		&team.Name,
		&team.ShortName,
		&team.CrestURL,
		&team.Country,
		&team.Abbreviation,
	)

	if err != nil && IsNoRowsError(err) {
		return Team{}, ErrNotFound
	} else if err != nil {
		return Team{}, err
	}

	return team, nil
}
