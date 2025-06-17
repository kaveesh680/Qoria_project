package ping

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	Handler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	// Check status code
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, res.StatusCode)
	}

	// Check Content-Type
	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected Content-Type application/json; got %s", contentType)
	}

	// Check body
	expectedBody := `"ping"`
	buf := make([]byte, len(expectedBody))
	_, err := res.Body.Read(buf)
	if err != nil && err.Error() != "EOF" {
		t.Fatalf("failed to read response body: %v", err)
	}

	actual := strings.TrimSpace(string(buf))
	if actual != expectedBody {
		t.Errorf("expected body %s; got %s", expectedBody, actual)
	}
}
