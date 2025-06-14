package shortify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/krishnadwypayan/shorturl/internal/logger"
	"github.com/krishnadwypayan/shorturl/internal/model"
	"github.com/krishnadwypayan/shorturl/internal/mongo"
)

var SnowflakeBaseUrl = getSnowflakeBaseUrl()

const (
	SnowflakeGenerateEndpoint = "/generate"
	ShortUrl                  = "http://short.ify/"
	DefaultTTL                = 86400 // 24 hours in seconds
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
	var snowflakeRes model.SnowflakeResponse

	if req.Alias != "" {
		// TODO : Check if the alias already exists in the database
		exists, _ := mongo.CheckAliasExists(req.Alias)
		if exists {
			logger.Error().Msg(fmt.Sprintf("Alias already exists: %s", req.Alias))
			return "", errors.New("alias already exists")
		}

		// If an alias is provided, use it as the ID
		snowflakeRes.ID = req.Alias
	} else {
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

		if err := json.NewDecoder(res.Body).Decode(&snowflakeRes); err != nil {
			logger.Error().Msg(fmt.Sprintf("Failed to decode snowflake response: %v", err))
			return "", errors.New("failed to generate unique ID")
		}
	}

	// set default TTL if not provided
	if req.TTL <= 0 {
		req.TTL = DefaultTTL
	}

	err := mongo.InsertUrlMapping(req, snowflakeRes.ID)
	if err != nil {
		logger.Error().Msg(fmt.Sprintf("Failed to insert URL mapping: %v", err))
		return "", fmt.Errorf("failed to insert URL mapping: %w", err)
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

func getSnowflakeBaseUrl() string {
	if v := os.Getenv("SNOWFLAKE_BASE_URL"); v != "" {
		return v
	}
	return "http://localhost:8080" // default fallback
}
