package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"time"
)

// User represents a user in the system
type User struct {
	ID                 string        `db:"id"`
	FirstName          *string       `db:"first_name"`
	LastName           *string       `db:"last_name"`
	Username           string        `db:"username"`
	AvatarURL          *string       `db:"avatar_url"`
	LanguageCode       *string       `db:"language_code"`
	ChatID             int64         `db:"chat_id"`
	ReferredBy         *string       `db:"referred_by"`
	CreatedAt          time.Time     `db:"created_at"`
	TotalPoints        int           `db:"total_points"`
	TotalPredictions   int           `db:"total_predictions"`
	CorrectPredictions int           `db:"correct_predictions"`
	GlobalRank         int           `db:"global_rank"`
	CurrentWinStreak   int           `db:"current_win_streak"`
	LongestWinStreak   int           `db:"longest_win_streak"`
	FavoriteTeamID     *string       `db:"favorite_team_id"`
	FavoriteTeam       *FavoriteTeam `db:"favorite_team"`
}

type FavoriteTeam Team

func (ft *FavoriteTeam) Scan(src interface{}) error {
	var source []byte
	switch src := src.(type) {
	case []byte:
		source = src
	case string:
		source = []byte(src)
	case nil:
		return nil
	default:
		return errors.New("unsupported type")
	}

	if len(source) == 0 {
		return nil
	}

	if err := json.Unmarshal(source, ft); err != nil {
		return err
	}

	return nil
}

func IsNoRowsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func IsUniqueViolationError(err error) bool {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		return errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique)
	}
	return false
}

func IsForeignKeyViolationError(err error) bool {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		return errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintForeignKey)
	}
	return false
}

func (s *Storage) CreateUser(user User) error {
	query := `
		INSERT INTO users (id, first_name, last_name, username, language_code, chat_id, avatar_url, referred_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query, user.ID, user.FirstName, user.LastName, user.Username, user.LanguageCode, user.ChatID, user.AvatarURL, user.ReferredBy)
	return err
}

func (s *Storage) getUserBy(condition string, value interface{}) (User, error) {
	query := fmt.Sprintf(`
			SELECT id,
				   first_name,
				   last_name,
				   username,
				   language_code,
				   chat_id,
				   created_at,
				   total_points,
				   total_predictions,
				   correct_predictions,
				   avatar_url,
				   referred_by,
				   global_rank,
				   current_win_streak,
				   longest_win_streak,
				   favorite_team
			FROM (SELECT u.id,
						 u.first_name,
						 u.last_name,
						 u.username,
						 u.language_code,
						 u.chat_id,
						 u.created_at,
						 u.total_points,
						 u.total_predictions,
						 u.correct_predictions,
						 u.avatar_url,
						 u.referred_by,
						 u.current_win_streak,
						 u.longest_win_streak,
						 RANK() OVER (ORDER BY total_points DESC) AS global_rank,
						 CASE
							 WHEN u.favorite_team_id IS NOT NULL THEN
								 json_object(
										 'id', t.id,
										 'name', t.name,
										 'short_name', t.short_name,
										 'crest_url', t.crest_url,
										 'country', t.country,
										 'abbreviation', t.abbreviation
								 )
						 END                                  AS favorite_team
				  FROM users u
						   LEFT JOIN teams t ON u.favorite_team_id = t.id) ranked_users
			WHERE %s`, condition)

	var user User
	row := s.db.QueryRowContext(context.Background(), query, value)

	if err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.LanguageCode,
		&user.ChatID,
		&user.CreatedAt,
		&user.TotalPoints,
		&user.TotalPredictions,
		&user.CorrectPredictions,
		&user.AvatarURL,
		&user.ReferredBy,
		&user.GlobalRank,
		&user.CurrentWinStreak,
		&user.LongestWinStreak,
		&user.FavoriteTeam,
	); err != nil && IsNoRowsError(err) {
		return User{}, ErrNotFound
	} else if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Storage) GetUserByChatID(chatID int64) (User, error) {
	return s.getUserBy("chat_id = ?", chatID)
}

func (s *Storage) GetUserByID(id string) (User, error) {
	return s.getUserBy("id = ?", id)
}

func (s *Storage) GetUserByUsername(uname string) (User, error) {
	return s.getUserBy("username = ?", uname)
}

func (s *Storage) UpdateUserPoints(ctx context.Context, userID string, points int, isCorrect bool) error {
	var correctPredictionsIncrement int
	if isCorrect {
		correctPredictionsIncrement = 1
	} else {
		correctPredictionsIncrement = 0
	}

	query := `
		 UPDATE users
        SET total_points = total_points + ?,
            total_predictions = total_predictions + 1,
            correct_predictions = correct_predictions + ?
        WHERE id = ?`

	_, err := s.db.ExecContext(ctx, query, points, correctPredictionsIncrement, userID)
	return err
}

func (s *Storage) UpdateUserPredictionCount(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET total_predictions = total_predictions + 1
		WHERE id = ?`

	_, err := s.db.ExecContext(ctx, query, userID)
	return err
}

func (s *Storage) ListUserReferrals(ctx context.Context, userID string) ([]User, error) {
	query := `
		SELECT id, first_name, last_name, username, language_code, chat_id, created_at, total_points, total_predictions, correct_predictions, avatar_url, referred_by
		FROM users
		WHERE referred_by = ?`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.LanguageCode,
			&user.ChatID,
			&user.CreatedAt,
			&user.TotalPoints,
			&user.TotalPredictions,
			&user.CorrectPredictions,
			&user.AvatarURL,
			&user.ReferredBy,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Storage) UpdateUserInformation(ctx context.Context, user User) error {
	query := `
		UPDATE users
		SET first_name = ?,
		    last_name = ?,
		    username = ?,
		    avatar_url = ?,
		    language_code = ?,
		    favorite_team_id = ?
		WHERE id = ?`

	_, err := s.db.ExecContext(ctx, query, user.FirstName, user.LastName, user.Username, user.AvatarURL, user.LanguageCode, user.FavoriteTeamID, user.ID)
	return err
}

func (s *Storage) UpdateUserStreak(ctx context.Context, userID string, currentStreak int, longestStreak int) error {
	query := `
        UPDATE users
        SET current_win_streak = ?, longest_win_streak = ?
        WHERE id = ?
    `
	_, err := s.db.ExecContext(ctx, query, currentStreak, longestStreak, userID)
	return err
}
