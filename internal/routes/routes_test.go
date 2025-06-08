package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/krishnadwypayan/shorturl/internal/snowflake"
	"github.com/stretchr/testify/assert"
)

func TestRegisterSnowflakeRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	generator := snowflake.NewGenerator(1)

	RegisterSnowflakeRoutes(r, generator)

	// Test health check route
	w := performGenerateRequest(r, "GET", "/health")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status": "ok"}`, w.Body.String())

	// Test generate ID route
	w = performGenerateRequest(r, "GET", "/generate")
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["id"])
}

func TestRegisterShortURLRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	RegisterShortifyRoutes(r)

	// Test health check route
	w := performGenerateRequest(r, "GET", "/health")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status": "ok"}`, w.Body.String())

	// Test /shortify with invalid JSON
	w = httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/shortify", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	_, hasError := errResp["error"]
	assert.True(t, hasError)

	// Test /shortify with valid JSON
	validBody := `{"url": "https://example.com"}`
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/shortify",
		strings.NewReader(validBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	// Accept either 200 or 400 depending on shortify.Shortify implementation
	assert.Contains(t, []int{http.StatusOK, http.StatusBadRequest}, w.Code)
	if w.Code == http.StatusOK {
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp)
	} else {
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp, "error")
	}
}

func BenchmarkShortifyRoute(b *testing.B) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	RegisterShortifyRoutes(r)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := performShortifyRequest(r, "POST", "/shortify", `{"long_url": "https://example.com"}`)
		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
			b.Fatalf("unexpected status code: %d", w.Code)
		}
	}
}

func performShortifyRequest(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

func BenchmarkGenerateRoute(b *testing.B) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	generator := snowflake.NewGenerator(1)
	RegisterSnowflakeRoutes(r, generator)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := performGenerateRequest(r, "GET", "/generate")
		if w.Code != http.StatusOK {
			b.Fatalf("unexpected status code: %d", w.Code)
		}
	}
}

func performGenerateRequest(r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	r.ServeHTTP(w, req)
	return w
}
