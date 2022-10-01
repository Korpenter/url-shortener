package model

type Response struct {
	Result string `json:"result"`
}

type ResponseURLItem struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"url"`
}
