package db

import "time"

// LeaderboardEntry represents an entry in the leaderboard
type LeaderboardEntry struct {
	ID       int    `db:"id"`
	LeagueID int    `db:"league_id"`
	UserID   string `db:"user_id"`
	Points   int    `db:"points"`
}

// Season represents a season in the competition
type Season struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	IsActive  bool      `db:"is_active"`
}

// Team represents a sports team
type Team struct {
	ID           int    `db:"id"`
	Name         string `db:"name"`
	ShortName    string `db:"short_name"`
	Abbreviation string `db:"abbreviation"`
	CrestURL     string `db:"crest_url"`
	Country      string `db:"country"`
}
