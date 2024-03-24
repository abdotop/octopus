package octopus

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	app := New()

	go func() {
		log.Fatalln(app.Run(":8080"))

	}()

	// Wait for the server to start
	time.Sleep(time.Second * 2)

	status, body, err := pingURL("http://localhost:8080")
	fmt.Println(status, body, err)
	if err != nil {
		t.Fatalf("Failed to ping URL: %v", err)
	}

	// Check the status and body
	t.Logf("Status: %d, Body: %s", status, body)
}

func pingURL(url string) (int, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}

	return resp.StatusCode, string(body), nil
}
