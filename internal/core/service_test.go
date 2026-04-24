package core_test

import (
	"context"
	"testing"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 🛠️ Step 1: Create 'Mock' versions of our shores
type MockIdentityRepo struct{ mock.Mock }
func (m *MockIdentityRepo) CreateSnippet(ctx context.Context, s *core.Snippet) error {
		args := m.Called(ctx, s)
		return args.Error(0)
}

type MockGraphRepo struct{ mock.Mock }
func (m *MockGraphRepo) SaveVersion(ctx context.Context, s *core.Snippet, v *core.Version) error {
		args := m.Called(ctx, s, v) // we pass the whole snippet 's'.
		return args.Error(0)
}


func (m *MockIdentityRepo) Search(ctx context.Context, query string) ([]*core.Snippet, error) {
    args := m.Called(ctx, query)
    return args.Get(0).([]*core.Snippet), args.Error(1)
}

func (m *MockGraphRepo) CiteSnippet(ctx context.Context, c *core.Citation) error {
		args := m.Called(ctx, c)
		return args.Error(0)
}


// Step 2: Test the Dual-Handshake
func TestCreateSnippet_Handshake(t *testing.T) {
		mockID := new(MockIdentityRepo)
		mockGraph := new(MockGraphRepo)
		service := core.NewSnippetService(mockID, mockGraph)

		ctx := context.Background()

		// We expect the brain to talk to BOTH repositories
		// 1. Identity repo handles the Snippet
    	mockID.On("CreateSnippet", ctx, mock.AnythingOfType("*core.Snippet")).Return(nil)
    	// 2. Graph repo handles the Version
    	// Both arguments are now pointers to structs
		mockGraph.On("SaveVersion", ctx, mock.AnythingOfType("*core.Snippet"), mock.AnythingOfType("*core.Version")).Return(nil)

		// Execute the action
		snippet, err := service.CreateSnippet(ctx, "Test Bridge", "user_123", "fmt.Println('Hello Euler')")

		// Verify the result
		assert.NoError(t, err)
		assert.NotNil(t, snippet)
		assert.Equal(t, "Test Bridge", snippet.Title)

		// Confirm the handshake occurred
		mockID.AssertExpectations(t)
		mockGraph.AssertExpectations(t)
}

func TestCiteSnippet_LegacyCompatibility(t *testing.T) {
    mockID := new(MockIdentityRepo)
    mockGraph := new(MockGraphRepo)
    service := core.NewSnippetService(mockID, mockGraph)
    ctx := context.Background()

    sourceID := "modern-id"
    targetID := "legacy-id-null-title"
    contextStr := "Referencing legacy Euler notes"

    // We expect the service to attempt the citation regardless of the target's internal state
    mockGraph.On("CiteSnippet", ctx, mock.MatchedBy(func(c *core.Citation) bool {
        return c.SourceID == sourceID && c.TargetID == targetID
    })).Return(nil)

    // Action
    err := service.CiteSnippet(ctx, sourceID, targetID, contextStr)

    // Assert
    assert.NoError(t, err)
    mockGraph.AssertExpectations(t)
}
		
