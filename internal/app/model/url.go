package model

// URL represents url record
type URL struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"url"`
	UserID   string `json:"user_id"`
}
