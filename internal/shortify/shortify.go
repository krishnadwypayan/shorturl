package shortify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/krishnadwypayan/shorturl/internal/logger"
	"github.com/krishnadwypayan/shorturl/internal/model"
)

var SnowflakeBaseUrl = "http://localhost:8080"

const (
	SnowflakeGenerateEndpoint = "/generate"
	ShortUrl                  = "http://short.ify/"
)

func Shortify(req model.ShortURLRequest) (model.ShortURLResponse, error) {
	// Validate the request
	if err := validateShortURLRequest(req); err != nil {
		logger.Error().Msg(fmt.Sprintf("Invalid short URL request: %v", err))
		return model.ShortURLResponse{}, err
	}

	// Generate a unique ID for the short URL
	id, err := generateUniqueID(req)
	if err != nil {
		logger.Error().Msg(fmt.Sprintf("Failed to generate unique ID: %v", err))
		return model.ShortURLResponse{}, err
	}

	// Create the short URL response
	res := model.ShortURLResponse{
		LongURL:  req.LongURL,
		ShortURL: ShortUrl + id,
		ID:       id,
	}

	return res, nil
}

func generateUniqueID(req model.ShortURLRequest) (string, error) {
	if req.Alias != "" {
		// If an alias is provided, use it as the ID
		return req.Alias, nil
	}

	res, err := http.Get(SnowflakeBaseUrl + SnowflakeGenerateEndpoint)
	if err != nil {
		logger.Error().Msg(fmt.Sprintf("Failed to call snowflake service: %v", err))
		return "", errors.New("failed to generate unique ID")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logger.Error().Msg(fmt.Sprintf("Snowflake service returned status: %d", res.StatusCode))
		return "", errors.New("failed to generate unique ID")
	}

	var snowflakeRes model.SnowflakeResponse
	if err := json.NewDecoder(res.Body).Decode(&snowflakeRes); err != nil {
		logger.Error().Msg(fmt.Sprintf("Failed to decode snowflake response: %v", err))
		return "", errors.New("failed to generate unique ID")
	}
	return snowflakeRes.ID, nil
}

func validateShortURLRequest(req model.ShortURLRequest) error {
	if req.LongURL == "" {
		return errors.New("long URL is required")
	}
	if !strings.HasPrefix(req.LongURL, "http://") && !strings.HasPrefix(req.LongURL, "https://") {
		return errors.New("long URL must start with http:// or https://")
	}
	if req.Alias != "" {
		if len(req.Alias) < 3 || len(req.Alias) > 20 {
			return errors.New("alias must be 3-20 characters")
		}
		if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(req.Alias) {
			return errors.New("alias contains invalid characters")
		}
	}
	return nil
}
