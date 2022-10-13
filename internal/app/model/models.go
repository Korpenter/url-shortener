package model

// URL represents url record
type URL struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"url"`
	UserID   string `json:"user_id"`
}

type Response struct {
	Result string `json:"result"`
}

type URLItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type BatchReqItem struct {
	CorID   string `json:"correlation_id"`
	OrigURL string `json:"original_url"`
}

type BatchRespItem struct {
	CorID    string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}