package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abdotop/octopus"
)

func TestCORSMiddleware(t *testing.T) {
	app := octopus.New()

	app.Use(New(Config{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
		ExposedHeaders:   []string{},
		MaxAge:           86400,
	}))

	app.GET("/test", func(c *octopus.Ctx) {})

	req, err := http.NewRequest("GET", "http://localhost:8888/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.ServeHTTP)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "*"
	if rr.Header().Get("Access-Control-Allow-Origin") != expected {
		t.Errorf("handler returned unexpected Access-Control-Allow-Origin header: got %v want %v", rr.Header().Get("Access-Control-Allow-Origin"), expected)
	}
}
