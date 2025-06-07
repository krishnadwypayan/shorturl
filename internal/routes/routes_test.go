package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
	w := performRequest(r, "GET", "/health")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status": "ok"}`, w.Body.String())

	// Test generate ID route
	w = performRequest(r, "GET", "/generate")
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["id"])
}

func BenchmarkGenerateRoute(b *testing.B) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	generator := snowflake.NewGenerator(1)
	RegisterSnowflakeRoutes(r, generator)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := performRequest(r, "GET", "/generate")
		if w.Code != http.StatusOK {
			b.Fatalf("unexpected status code: %d", w.Code)
		}
	}
}

func performRequest(r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	r.ServeHTTP(w, req)
	return w
}
