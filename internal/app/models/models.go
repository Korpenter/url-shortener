package models

// URL represents a URL object that contains information about a shortened URL.
type URL struct {
	// ShortURL is the shortened version of the URL.
	ShortURL string `json:"short_url"`
	// LongURL is the original version of the URL.
	LongURL string `json:"url"`
	// UserID is the ID of the user who created the shortened URL.
	UserID string `json:"user_id"`
	// Deleted indicates whether the URL has been deleted or not.
	Deleted bool `json:"deleted"`
}

// Response represents a shortened URL sent in response to Users' request.
type Response struct {
	// Result is the shortened version of the URL.
	Result string `json:"result"`
}

// URLItem represents a URL item with shortened and original URLs.
type URLItem struct {
	// ShortURL is the shortened version of the URL.
	ShortURL string `json:"short_url"`
	// OriginalURL is the original, long version of the URL.
	OriginalURL string `json:"original_url"`
}

// DeleteURLItem represents an item containing information about a URL to be deleted.
type DeleteURLItem struct {
	// UserID is the ID of the user who is trying to delete the URL.
	UserID string `json:"user_id"`
	// ShortURL is the shortened version of the URL.
	ShortURL string `json:"short_url"`
}

// BatchReqItem represents an item in a batch request for creating shortened URLs.
type BatchReqItem struct {
	// CorID is the correlation ID for the original URL.
	CorID string `json:"correlation_id"`
	// OrigURL is the original, long version of the URL.
	OrigURL string `json:"original_url"`
}

// BatchRespItem represents an item in a batch response containing shortened URLs.
type BatchRespItem struct {
	// CorID is the correlation ID for the request.
	CorID string `json:"correlation_id"`
	// ShortURL is the shortened version of the URL.
	ShortURL string `json:"short_url"`
}
