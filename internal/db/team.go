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

func (s *storage) GetTeamByID(ctx context.Context, id string) (Team, error) {
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

func (s *storage) SaveTeam(ctx context.Context, team Team) error {
	query := `
		INSERT INTO teams (id, name, short_name, crest_url, country, abbreviation)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET
		name = excluded.name,
		short_name = excluded.short_name,
		crest_url = excluded.crest_url,
		country = excluded.country,
		abbreviation = excluded.abbreviation
	`

	_, err := s.db.ExecContext(ctx, query,
		team.ID,
		team.Name,
		team.ShortName,
		team.CrestURL,
		team.Country,
		team.Abbreviation,
	)

	return err
}
