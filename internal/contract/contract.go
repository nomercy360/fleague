package contract

import (
	"time"
)

type Error struct {
	Message string `json:"message"`
}

type UserAuthResponse struct {
	ID           int       `json:"id"`
	FirstName    *string   `json:"first_name"`
	LastName     *string   `json:"last_name"`
	Username     string    `json:"username"`
	ChatID       int64     `json:"chat_id"`
	LanguageCode *string   `json:"language_code"`
	CreatedAt    time.Time `json:"created_at"`
	Token        string    `json:"token"`
}

type TeamResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ShortName    string `json:"short_name"`
	CrestURL     string `json:"crest_url"`
	Country      string `json:"country"`
	Abbreviation string `json:"abbreviation"`
}

type MatchResponse struct {
	ID         int          `json:"id"`
	Tournament string       `json:"tournament"`
	HomeTeam   TeamResponse `json:"home_team"`
	AwayTeam   TeamResponse `json:"away_team"`
	MatchDate  time.Time    `json:"match_date"`
	Status     string       `json:"status"`
	AwayScore  int          `json:"away_score"`
	HomeScore  int          `json:"home_score"`
}
