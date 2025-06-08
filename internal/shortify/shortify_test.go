package shortify

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/krishnadwypayan/shorturl/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestShortify_Success(t *testing.T) {
	srv := startMockSnowflakeServer()
	defer srv.Close()

	req := model.ShortURLRequest{
		LongURL: "https://example.com",
		Alias:   "",
	}
	res, err := Shortify(req)
	assert.NoError(t, err)
	assert.Equal(t, req.LongURL, res.LongURL)
	assert.True(t, strings.HasPrefix(res.ShortURL, ShortUrl))
	assert.NotEmpty(t, res.ID)
}

func TestShortify_InvalidRequest(t *testing.T) {
	srv := startMockSnowflakeServer()
	defer srv.Close()

	tests := []struct {
		name string
		req  model.ShortURLRequest
	}{
		{"Empty URL", model.ShortURLRequest{LongURL: ""}},
		{"Invalid URL scheme", model.ShortURLRequest{LongURL: "ftp://foo.com"}},
		{"Alias too short", model.ShortURLRequest{LongURL: "http://foo.com", Alias: "ab"}},
		{"Alias too long", model.ShortURLRequest{LongURL: "http://foo.com", Alias: strings.Repeat("a", 21)}},
		{"Alias invalid chars", model.ShortURLRequest{LongURL: "http://foo.com", Alias: "bad!alias"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Shortify(tc.req)
			assert.Error(t, err)
		})
	}
}

func TestShortify_SnowflakeError(t *testing.T) {
	// Service returns error
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "fail", http.StatusInternalServerError)
	}))
	defer srv.Close()
	req := model.ShortURLRequest{LongURL: "http://foo.com"}
	_, err := Shortify(req)
	assert.Error(t, err)
}

func TestShortify_SnowflakeDecodeError(t *testing.T) {
	// Service returns invalid JSON
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "not-json")
	}))
	defer srv.Close()
	req := model.ShortURLRequest{LongURL: "http://foo.com"}
	_, err := Shortify(req)
	assert.Error(t, err)
}

func TestShortify_WithAlias(t *testing.T) {
	req := model.ShortURLRequest{
		LongURL: "https://example.com",
		Alias:   "customAlias",
	}
	res, err := Shortify(req)
	assert.NoError(t, err)
	assert.Equal(t, req.LongURL, res.LongURL)
	assert.Equal(t, ShortUrl+"customAlias", res.ShortURL)
	assert.Equal(t, "customAlias", res.ID)
}

func startMockSnowflakeServer() *httptest.Server {
	// Mock Snowflake service using gin
	router := gin.New()
	router.GET("/generate", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.SnowflakeResponse{ID: "abc123"})
	})
	srv := httptest.NewServer(router)
	// defer srv.Close()

	SnowflakeBaseUrl = srv.URL // Set the Snowflake service URL for testing
	return srv
}
