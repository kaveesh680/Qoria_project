package ping

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncodePingResponse(t *testing.T) {
	w := httptest.NewRecorder()
	response := "ping"

	err := EncodePingResponse(context.Background(), w, response)
	if err != nil {
		t.Fatalf("EncodePingResponse returned error: %v", err)
	}

	res := w.Result()
	defer res.Body.Close()

	// Check status code
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, res.StatusCode)
	}

	// Check content type
	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json; got %s", ct)
	}

	// Check response body
	var decoded string
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if decoded != "ping" {
		t.Errorf("expected body 'ping'; got '%s'", decoded)
	}
}
