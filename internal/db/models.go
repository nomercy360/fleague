package db

import "time"

// LeaderboardEntry represents an entry in the leaderboard
type LeaderboardEntry struct {
	UserID   string `db:"user_id"`
	Points   int    `db:"points"`
	SeasonID string `db:"season_id"`
}

// Season represents a season in the competition
type Season struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	IsActive  bool      `db:"is_active"`
}

// Team represents a sports team
type Team struct {
	ID           string `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	ShortName    string `db:"short_name" json:"short_name"`
	Abbreviation string `db:"abbreviation" json:"abbreviation"`
	CrestURL     string `db:"crest_url" json:"crest_url"`
	Country      string `db:"country" json:"country"`
}
