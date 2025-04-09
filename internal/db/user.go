package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"github.com/user/project/internal/nanoid"
	"time"
)

// User represents a user in the system
type User struct {
	ID                 string    `db:"id"`
	FirstName          *string   `db:"first_name"`
	LastName           *string   `db:"last_name"`
	Username           string    `db:"username"`
	AvatarURL          *string   `db:"avatar_url"`
	LanguageCode       *string   `db:"language_code"`
	ChatID             int64     `db:"chat_id"`
	ReferredBy         *string   `db:"referred_by"`
	CreatedAt          time.Time `db:"created_at"`
	TotalPredictions   int       `db:"total_predictions"`
	CorrectPredictions int       `db:"correct_predictions"`
	CurrentWinStreak   int       `db:"current_win_streak"`
	LongestWinStreak   int       `db:"longest_win_streak"`
	FavoriteTeamID     *string   `db:"favorite_team_id"`
	FavoriteTeam       *Team     `db:"favorite_team"`
	Badges             []Badge   `db:"badges"`
	SubscriptionActive bool      `db:"subscription_active"`
	SubscriptionExpiry time.Time `db:"subscription_expiry"`
}

type Badge struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Icon      string    `json:"icon"`
	AwardedAt time.Time `json:"awarded_at"`
}

func UnmarshalJSONToSlice[T any](src interface{}) ([]T, error) {
	var source []byte

	switch s := src.(type) {
	case []byte:
		source = s
	case string:
		source = []byte(s)
	case nil:
		return []T{}, nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", s)
	}

	var result []T
	if err := json.Unmarshal(source, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return result, nil
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
	query := `SELECT 
				u.id,
				u.first_name,
				u.last_name,
				u.username,
				u.language_code,
				u.chat_id,
				u.created_at,
				u.total_predictions,
				u.correct_predictions,
				u.avatar_url,
				u.referred_by,
				u.current_win_streak,
				u.longest_win_streak,
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
				END AS favorite_team,
				json_group_array(
					DISTINCT json_object(
						'id', b.id,
						'name', b.name,
						'awarded_at', strftime('%Y-%m-%dT%H:%M:%SZ', ub.awarded_at),
						'color', b.color,
						'icon', b.icon
					)
				) FILTER (WHERE b.id IS NOT NULL) AS badges
			FROM users u
			LEFT JOIN teams t ON u.favorite_team_id = t.id
			LEFT JOIN user_badges ub ON u.id = ub.user_id
			LEFT JOIN badges b ON ub.badge_id = b.id
			WHERE ` + fmt.Sprintf("%s GROUP BY u.id", condition)

	var user User
	row := s.db.QueryRowContext(context.Background(), query, value)
	var badgeJSON string
	var favoriteTeamJSON *string

	if err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.LanguageCode,
		&user.ChatID,
		&user.CreatedAt,
		&user.TotalPredictions,
		&user.CorrectPredictions,
		&user.AvatarURL,
		&user.ReferredBy,
		&user.CurrentWinStreak,
		&user.LongestWinStreak,
		&favoriteTeamJSON,
		&badgeJSON,
	); err != nil && IsNoRowsError(err) {
		return User{}, ErrNotFound
	} else if err != nil {
		return User{}, err
	}

	var err error
	user.Badges, err = UnmarshalJSONToSlice[Badge](badgeJSON)
	if err != nil {
		return user, err
	}

	if favoriteTeamJSON != nil {
		var team Team
		if err := json.Unmarshal([]byte(*favoriteTeamJSON), &team); err != nil {
			return user, err
		}
		user.FavoriteTeam = &team
	}

	return user, nil
}

func (s *Storage) GetUserByChatID(chatID int64) (User, error) {
	return s.getUserBy("u.chat_id = ?", chatID)
}

func (s *Storage) GetUserByID(id string) (User, error) {
	return s.getUserBy("u.id = ?", id)
}

func (s *Storage) GetUserByUsername(uname string) (User, error) {
	return s.getUserBy("u.username = ?", uname)
}

func (s *Storage) UpdateUserPoints(ctx context.Context, userID string, isCorrect bool) error {
	var correctPredictionsIncrement int
	if isCorrect {
		correctPredictionsIncrement = 1
	} else {
		correctPredictionsIncrement = 0
	}

	query := `
		 UPDATE users
        SET total_predictions = total_predictions + 1,
            correct_predictions = correct_predictions + ?
        WHERE id = ?`

	_, err := s.db.ExecContext(ctx, query, correctPredictionsIncrement, userID)
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
		SELECT id, first_name, last_name, username, language_code, chat_id, created_at, total_predictions, correct_predictions, avatar_url, referred_by
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

func (s *Storage) GetAllUsers(ctx context.Context) ([]User, error) {
	query := `
		SELECT id, first_name, last_name, username, language_code, chat_id, created_at, total_predictions, correct_predictions, avatar_url, referred_by, current_win_streak, longest_win_streak, favorite_team_id
		FROM users
		WHERE total_predictions > 0`

	rows, err := s.db.QueryContext(ctx, query)
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
			&user.TotalPredictions,
			&user.CorrectPredictions,
			&user.AvatarURL,
			&user.ReferredBy,
			&user.CurrentWinStreak,
			&user.LongestWinStreak,
			&user.FavoriteTeamID,
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

func (s *Storage) GetAllUsersWithFavoriteTeam(ctx context.Context) ([]User, error) {
	query := `
        SELECT id, first_name, last_name, username, language_code, chat_id, favorite_team_id
        FROM users
        WHERE favorite_team_id IS NOT NULL
    `
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.LanguageCode, &user.ChatID, &user.FavoriteTeamID); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *Storage) FollowUser(ctx context.Context, followerID, followingID string) error {
	query := `
		INSERT INTO user_followers (follower_id, following_id)
		VALUES (?, ?)
	`
	_, err := s.db.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		if IsUniqueViolationError(err) {
			return fmt.Errorf("already following user: %w", err)
		}
		if IsForeignKeyViolationError(err) {
			return fmt.Errorf("one of the users does not exist: %w", err)
		}
		return err
	}
	return nil
}

// UnfollowUser allows a user to unfollow another user
func (s *Storage) UnfollowUser(ctx context.Context, followerID, followingID string) error {
	query := `
		DELETE FROM user_followers
		WHERE follower_id = ? AND following_id = ?
	`
	result, err := s.db.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("not following this user")
	}
	return nil
}

// IsFollowing checks if a user is following another user
func (s *Storage) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	query := `
		SELECT 1 FROM user_followers
		WHERE follower_id = ? AND following_id = ?
	`
	var exists int
	err := s.db.QueryRowContext(ctx, query, followerID, followingID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetFollowers retrieves a list of users following the given user
func (s *Storage) GetFollowers(ctx context.Context, userID string) ([]User, error) {
	query := `
		SELECT u.id, u.first_name, u.last_name, u.username, u.avatar_url
		FROM user_followers uf
		JOIN users u ON uf.follower_id = u.id
		WHERE uf.following_id = ?
	`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.AvatarURL); err != nil {
			return nil, err
		}
		followers = append(followers, user)
	}
	return followers, rows.Err()
}

// GetFollowing retrieves a list of users the given user is following
func (s *Storage) GetFollowing(ctx context.Context, userID string) ([]User, error) {
	query := `
		SELECT u.id, u.first_name, u.last_name, u.username, u.avatar_url
		FROM user_followers uf
		JOIN users u ON uf.following_id = u.id
		WHERE uf.follower_id = ?
	`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.AvatarURL); err != nil {
			return nil, err
		}
		following = append(following, user)
	}
	return following, rows.Err()
}

func (s *Storage) RecordUserLogin(ctx context.Context, userID string) error {
	query := `
        INSERT INTO user_logins (id, user_id, login_time)
        VALUES (?, ?, ?)
    `
	id := nanoid.Must() // Assuming you have a nanoid package for unique IDs
	_, err := s.db.ExecContext(ctx, query, id, userID, time.Now())
	return err
}

func (s *Storage) HasLoggedInToday(ctx context.Context, userID string) (bool, error) {
	query := `
        SELECT COUNT(*)
        FROM user_logins
        WHERE user_id = ?
        AND DATE(login_time) = DATE('now')
    `
	var count int
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
