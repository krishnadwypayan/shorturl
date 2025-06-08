package model

type SnowflakeResponse struct {
	ID string `json:"id"`
}

type ShortURLRequest struct {
	LongURL string `json:"long_url" binding:"required,url"`
	Alias   string `json:"alias,omitempty" binding:"omitempty,alphanum"`
}

type ShortURLResponse struct {
	LongURL  string `json:"long_url"`
	ShortURL string `json:"short_url"`
	ID       string `json:"id"`
}
