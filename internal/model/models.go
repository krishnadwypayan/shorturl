package model

type SnowflakeResponse struct {
	ID string `json:"id"`
}

// ShortURLRequest represents the payload for creating a shortened URL.
// It contains the original long URL, an optional custom alias, and an optional time-to-live (TTL) in seconds.
//
// Fields:
//   - LongURL: The original URL to be shortened. Must be a valid URL and is required.
//   - Alias: An optional custom alias for the shortened URL. Must be alphanumeric if provided.
//   - TTL: Optional time-to-live for the shortened URL in seconds. Must be numeric if provided.
type ShortURLRequest struct {
	LongURL string `json:"long_url" binding:"required,url"`
	Alias   string `json:"alias,omitempty" binding:"omitempty,alphanum"`
	TTL     int64  `json:"ttl,omitempty" binding:"omitempty,numeric"`
}

type ShortURLResponse struct {
	LongURL  string `json:"long_url"`
	ShortURL string `json:"short_url"`
	ID       string `json:"id"`
}
