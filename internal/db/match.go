package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Match represents a sports match
type Match struct {
	ID         string      `db:"id" json:"id"`
	Tournament string      `db:"tournament" json:"tournament"`
	HomeTeamID string      `db:"home_team_id" json:"home_team_id"`
	AwayTeamID string      `db:"away_team_id" json:"away_team_id"`
	MatchDate  time.Time   `db:"match_date" json:"match_date"`
	Status     string      `db:"status" json:"status"`
	HomeScore  *int        `db:"home_score" json:"home_score"` // Nullable, set after match completion
	AwayScore  *int        `db:"away_score" json:"away_score"` // Nullable, set after match completion
	HomeOdds   *float64    `db:"home_odds" json:"home_odds"`
	DrawOdds   *float64    `db:"draw_odds" json:"draw_odds"`
	AwayOdds   *float64    `db:"away_odds" json:"away_odds"`
	HomeTeam   Team        `db:"-" json:"home_team"`
	AwayTeam   Team        `db:"-" json:"away_team"`
	Prediction *Prediction `db:"-" json:"prediction,omitempty"`
	Popularity float64     `db:"popularity" json:"popularity"`
}

const (
	MatchStatusScheduled = "scheduled"
	MatchStatusCompleted = "completed"
	MatchStatusOngoing   = "ongoing"
)

func UnmarshalJSONToStruct[T any](src interface{}) (T, error) {
	var source []byte
	var zeroValue T

	switch s := src.(type) {
	case []byte:
		source = s
	case string:
		source = []byte(s)
	case nil:
		return zeroValue, nil
	default:
		return zeroValue, fmt.Errorf("unsupported type: %T", s)
	}

	var result T
	if err := json.Unmarshal(source, &result); err != nil {
		return zeroValue, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return result, nil
}

func (s *Storage) GetLastMatchesByTeamID(ctx context.Context, teamID string, limit int) ([]Match, error) {
	query := `
        SELECT m.id, m.tournament, m.home_team_id, m.away_team_id, m.match_date, m.status, m.home_score, m.away_score, m.popularity,
               json_object('id', home_team.id, 'name', home_team.name, 'short_name', home_team.short_name, 'crest_url', home_team.crest_url, 'country', home_team.country, 'abbreviation', home_team.abbreviation) as home_team,
               json_object('id', away_team.id, 'name', away_team.name, 'short_name', away_team.short_name, 'crest_url', away_team.crest_url, 'country', away_team.country, 'abbreviation', away_team.abbreviation) as away_team
        FROM matches m
        JOIN teams home_team ON home_team.id = home_team_id
        JOIN teams away_team ON away_team.id = away_team_id
        WHERE (home_team_id = ? OR away_team_id = ?)
        AND status = ?
        ORDER BY match_date DESC
        LIMIT ?
    `
	rows, err := s.db.QueryContext(ctx, query, teamID, teamID, MatchStatusCompleted, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var match Match
		var homeTeam, awayTeam interface{}
		if err := rows.Scan(
			&match.ID,
			&match.Tournament,
			&match.HomeTeamID,
			&match.AwayTeamID,
			&match.MatchDate,
			&match.Status,
			&match.HomeScore,
			&match.AwayScore,
			&match.Popularity,
			&homeTeam,
			&awayTeam,
		); err != nil {
			return nil, err
		}

		homeTeamStruct, err := UnmarshalJSONToStruct[Team](homeTeam)
		if err != nil {
			return nil, err
		}

		awayTeamStruct, err := UnmarshalJSONToStruct[Team](awayTeam)
		if err != nil {
			return nil, err
		}

		match.HomeTeam = homeTeamStruct
		match.AwayTeam = awayTeamStruct

		matches = append(matches, match)
	}

	return matches, nil
}

func (s *Storage) SaveMatch(ctx context.Context, match Match) error {
	query := `
        INSERT INTO matches (id, tournament, home_team_id, away_team_id, match_date, status, away_score, home_score, home_odds, draw_odds, away_odds, popularity)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
        tournament = excluded.tournament,
        home_team_id = excluded.home_team_id,
        away_team_id = excluded.away_team_id,
        match_date = excluded.match_date,
        status = excluded.status,
        away_score = excluded.away_score,
        home_score = excluded.home_score,
        home_odds = excluded.home_odds,
        draw_odds = excluded.draw_odds,
        away_odds = excluded.away_odds,
        popularity = excluded.popularity`

	_, err := s.db.ExecContext(ctx, query,
		match.ID,
		match.Tournament,
		match.HomeTeamID,
		match.AwayTeamID,
		match.MatchDate,
		match.Status,
		match.AwayScore,
		match.HomeScore,
		match.HomeOdds,
		match.DrawOdds,
		match.AwayOdds,
		match.Popularity,
	)
	return err
}

func (s *Storage) GetActiveMatches(ctx context.Context, userID string) ([]Match, error) {
	var query string
	var args []interface{}

	query = `
		SELECT
			m.id,
			m.tournament,
			m.home_team_id,
			m.away_team_id,
			m.match_date,
			m.status,
			m.home_odds,
			m.draw_odds,
			m.away_odds,
			m.popularity,
			json_object('id', t1.id, 'name', t1.name, 'short_name', t1.short_name, 'crest_url', t1.crest_url, 'country', t1.country, 'abbreviation', t1.abbreviation) as home_team,
			json_object('id', t2.id, 'name', t2.name, 'short_name', t2.short_name, 'crest_url', t2.crest_url, 'country', t2.country, 'abbreviation', t2.abbreviation) as away_team,
			CASE
				WHEN p.user_id IS NOT NULL THEN
					json_object(
						'user_id', p.user_id,
						'match_id', p.match_id,
						'predicted_outcome', p.predicted_outcome,
						'predicted_home_score', p.predicted_home_score,
						'predicted_away_score', p.predicted_away_score,
						'points_awarded', p.points_awarded,
						'created_at', CASE WHEN p.created_at IS NOT NULL THEN strftime('%Y-%m-%dT%H:%M:%SZ', p.created_at) ELSE NULL END,
						'updated_at', CASE WHEN p.updated_at IS NOT NULL THEN strftime('%Y-%m-%dT%H:%M:%SZ', p.updated_at) ELSE NULL END,
						'completed_at', CASE WHEN p.completed_at IS NOT NULL THEN strftime('%Y-%m-%dT%H:%M:%SZ', p.completed_at) ELSE NULL END
					)
				ELSE NULL
			END as prediction
		FROM matches m
		JOIN teams t1 ON m.home_team_id = t1.id
		JOIN teams t2 ON m.away_team_id = t2.id
		LEFT JOIN predictions p ON m.id = p.match_id AND p.user_id = ?
		WHERE m.status = 'scheduled' AND m.match_date BETWEEN datetime('now') AND datetime('now', '+7 days')
		ORDER BY m.popularity DESC`

	args = append(args, userID)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var match Match
		var homeTeam, awayTeam, prediction interface{}
		if err := rows.Scan(
			&match.ID,
			&match.Tournament,
			&match.HomeTeamID,
			&match.AwayTeamID,
			&match.MatchDate,
			&match.Status,
			&match.HomeOdds,
			&match.DrawOdds,
			&match.AwayOdds,
			&match.Popularity,
			&homeTeam,
			&awayTeam,
			&prediction,
		); err != nil {
			return nil, err
		}

		homeTeamStruct, err := UnmarshalJSONToStruct[Team](homeTeam)
		if err != nil {
			return nil, err
		}

		awayTeamStruct, err := UnmarshalJSONToStruct[Team](awayTeam)
		if err != nil {
			return nil, err
		}

		match.HomeTeam = homeTeamStruct
		match.AwayTeam = awayTeamStruct

		if prediction != nil {
			predictionStruct, err := UnmarshalJSONToStruct[Prediction](prediction)
			if err != nil {
				return nil, err
			}
			match.Prediction = &predictionStruct
		} else {
			match.Prediction = nil
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

func (s *Storage) GetMatchByID(ctx context.Context, id string) (Match, error) {
	query := `
		SELECT
			m.id,
			m.tournament,
			m.home_team_id,
			m.away_team_id,
			m.match_date,
			m.status,
			m.home_score,
			m.away_score,
			m.home_odds,
			m.draw_odds,
			m.away_odds,
			m.popularity,
			json_object('id', t1.id, 'name', t1.name, 'short_name', t1.short_name, 'crest_url', t1.crest_url, 'country', t1.country, 'abbreviation', t1.abbreviation) as home_team,
			json_object('id', t2.id, 'name', t2.name, 'short_name', t2.short_name, 'crest_url', t2.crest_url, 'country', t2.country, 'abbreviation', t2.abbreviation) as away_team
		FROM matches m
		JOIN teams t1 ON m.home_team_id = t1.id
		JOIN teams t2 ON m.away_team_id = t2.id
		WHERE m.id = ?`

	var match Match
	var homeTeam, awayTeam interface{}
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
		&match.HomeOdds,
		&match.DrawOdds,
		&match.AwayOdds,
		&match.Popularity,
		&homeTeam,
		&awayTeam,
	); err != nil && IsNoRowsError(err) {
		return Match{}, ErrNotFound
	} else if err != nil {
		return Match{}, err
	}

	homeTeamStruct, err := UnmarshalJSONToStruct[Team](homeTeam)
	if err != nil {
		return Match{}, err
	}

	awayTeamStruct, err := UnmarshalJSONToStruct[Team](awayTeam)
	if err != nil {
		return Match{}, err
	}

	match.HomeTeam = homeTeamStruct
	match.AwayTeam = awayTeamStruct

	return match, nil
}

func (s *Storage) GetCompletedMatchesWithoutCompletedPredictions(ctx context.Context) ([]Match, error) {
	query := `
		SELECT m.id, m.home_score, m.away_score
		FROM matches m
		WHERE m.status = 'completed' AND (SELECT COUNT(*) FROM predictions p WHERE p.match_id = m.id AND p.completed_at IS NOT NULL) = 0
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var match Match
		if err := rows.Scan(&match.ID, &match.HomeScore, &match.AwayScore); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func (s *Storage) GetMatchesForTeam(ctx context.Context, teamID string, hoursAhead int) ([]Match, error) {
	query := `
        SELECT id, tournament, home_team_id, away_team_id, match_date, status, popularity
        FROM matches
        WHERE (home_team_id = ? OR away_team_id = ?)
        AND status = ?
        AND match_date BETWEEN DATETIME('now') AND DATETIME('now', ? || ' hours')
    `

	rows, err := s.db.QueryContext(ctx, query, teamID, teamID, MatchStatusScheduled, hoursAhead)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var match Match
		if err := rows.Scan(
			&match.ID,
			&match.Tournament,
			&match.HomeTeamID,
			&match.AwayTeamID,
			&match.MatchDate,
			&match.Status,
			&match.Popularity,
		); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	return matches, nil
}

type PredictionStats struct {
	Home float64 `json:"home"`
	Draw float64 `json:"draw"`
	Away float64 `json:"away"`
}

func (s *Storage) GetPredictionStats(ctx context.Context, matchID string) (PredictionStats, error) {
	query := `
        SELECT
            SUM(CASE WHEN predicted_outcome = 'home' THEN 1 ELSE 0 END) AS home,
            SUM(CASE WHEN predicted_outcome = 'draw' THEN 1 ELSE 0 END) AS draw,
            SUM(CASE WHEN predicted_outcome = 'away' THEN 1 ELSE 0 END) AS away,
            COUNT(*) AS total
        FROM predictions
        WHERE match_id = ?
    `
	var home, draw, away, total float64
	err := s.db.QueryRowContext(ctx, query, matchID).Scan(&home, &draw, &away, &total)
	if err != nil {
		return PredictionStats{}, err
	}

	if total == 0 {
		return PredictionStats{}, nil
	}

	return PredictionStats{
		Home: (home / total) * 100,
		Draw: (draw / total) * 100,
		Away: (away / total) * 100,
	}, nil
}

func (s *Storage) GetTodayMostPopularMatch(ctx context.Context) (Match, error) {
	query := `
		SELECT m.id, m.tournament, m.home_team_id, m.away_team_id, m.match_date, m.status, m.home_score, m.away_score, m.popularity,
			   json_object('id', home_team.id, 'name', home_team.name, 'short_name', home_team.short_name, 'crest_url', home_team.crest_url, 'country', home_team.country, 'abbreviation', home_team.abbreviation) as home_team,
			   json_object('id', away_team.id, 'name', away_team.name, 'short_name', away_team.short_name, 'crest_url', away_team.crest_url, 'country', away_team.country, 'abbreviation', away_team.abbreviation) as away_team
		FROM matches m
		JOIN teams home_team ON home_team.id = home_team_id
		JOIN teams away_team ON away_team.id = away_team_id
		WHERE datetime(m.match_date) BETWEEN datetime('now', 'localtime') AND datetime('now', '+24 hours', 'localtime')
		ORDER BY m.match_date ASC
		LIMIT 1
	`

	var match Match
	var homeTeam, awayTeam interface{}
	row := s.db.QueryRowContext(ctx, query)

	if err := row.Scan(
		&match.ID,
		&match.Tournament,
		&match.HomeTeamID,
		&match.AwayTeamID,
		&match.MatchDate,
		&match.Status,
		&match.HomeScore,
		&match.AwayScore,
		&match.Popularity,
		&homeTeam,
		&awayTeam,
	); err != nil && IsNoRowsError(err) {
		return Match{}, ErrNotFound
	}

	homeTeamStruct, err := UnmarshalJSONToStruct[Team](homeTeam)
	if err != nil {
		return Match{}, err
	}

	awayTeamStruct, err := UnmarshalJSONToStruct[Team](awayTeam)
	if err != nil {
		return Match{}, err
	}

	match.HomeTeam = homeTeamStruct
	match.AwayTeam = awayTeamStruct

	return match, nil
}
