package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
)

type mockService struct{}

// 1. Matches CreateSnippet in service.go
func (m *mockService) CreateSnippet(ctx context.Context, title, ownerID, content string) (*core.Snippet, error) {
    return &core.Snippet{Title: title}, nil
}

// 2. Matches CiteSnippet in service.go
func (m *mockService) CiteSnippet(ctx context.Context, sourceID, targetID, contextStr string) error {
    return nil
}

// 3. Matches SearchSnippets in service.go
func (m *mockService) SearchSnippets(ctx context.Context, query string) ([]*core.Snippet, error) {
    return []*core.Snippet{}, nil
}

func TestCreateSnippetHandler(t *testing.T) {
	// 1. Setup
	// Initialise the handler with the mock service
	// Adjust 'NewHandler' to match your actual handler constructor
	svc := &mockService{}
	h := NewHandler(svc)

	// 2. Create dummy JSON payload
	jsonBody := `{"title": "Test Snippet", "owner_id": "user_123", "content": "fmt.Println(\"Hello\")"}`
	req, err := http.NewRequest("POST", "/snippets", bytes.NewBufferString(jsonBody))
	if err != nil {
    t.Fatalf("failed to create request: %v", err) // Added format string here
	}
	req.Header.Set("Content-Type", "application/json")

	// 3. Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// 4. Call the handler function directly
	h.CreateSnippet(rr, req)

	// 5. Assert the response status code and body
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", rr.Code)
	}

}