package db

import "time"

// Match represents a sports match
type Match struct {
	ID         int       `db:"id"`
	Tournament string    `db:"tournament"`
	HomeTeam   string    `db:"home_team"`
	AwayTeam   string    `db:"away_team"`
	MatchDate  time.Time `db:"match_date"`
	Status     string    `db:"status"`
	HomeScore  *int      `db:"home_score"` // Nullable, set after match completion
	AwayScore  *int      `db:"away_score"` // Nullable, set after match completion
}

// Prediction represents a user's prediction for a match
type Prediction struct {
	ID              int       `db:"id"`
	UserID          string    `db:"user_id"`
	MatchID         int       `db:"match_id"`
	PredictedWinner string    `db:"predicted_winner"` // home, away, or draw
	PointsAwarded   int       `db:"points_awarded"`
	CreatedAt       time.Time `db:"created_at"`
}

// League represents a user-created league
type League struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"` // Optional field
	OwnerID     string    `db:"owner_id"`
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
}

// LeagueMember represents a user's membership in a league
type LeagueMember struct {
	ID       int       `db:"id"`
	LeagueID int       `db:"league_id"`
	UserID   string    `db:"user_id"`
	JoinedAt time.Time `db:"joined_at"`
}

// LeagueMatch represents the association of a match with a league
type LeagueMatch struct {
	ID       int `db:"id"`
	LeagueID int `db:"league_id"`
	MatchID  int `db:"match_id"`
}

// LeaderboardEntry represents an entry in the leaderboard
type LeaderboardEntry struct {
	ID       int    `db:"id"`
	LeagueID int    `db:"league_id"`
	UserID   string `db:"user_id"`
	Points   int    `db:"points"`
}

// Referral represents a referral made by a user
type Referral struct {
	ID          int       `db:"id"`
	ReferrerID  string    `db:"referrer_id"`
	ReferredID  string    `db:"referred_id"`
	CreatedAt   time.Time `db:"created_at"`
	RewardGiven bool      `db:"reward_given"`
}

// Season represents a season in the competition
type Season struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	IsActive  bool      `db:"is_active"`
}

// SeasonMatch represents the association of a match with a season
type SeasonMatch struct {
	ID       int `db:"id"`
	SeasonID int `db:"season_id"`
	MatchID  int `db:"match_id"`
}
