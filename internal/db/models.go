package db

// LeaderboardEntry represents an entry in the leaderboard
type LeaderboardEntry struct {
	UserID   string `db:"user_id"`
	Points   int    `db:"points"`
	SeasonID string `db:"season_id"`
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
