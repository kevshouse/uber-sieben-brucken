package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 🛠️ Step 1: Create 'Mock' versions of our shores
type MockIdentityRepo struct{ mock.Mock }
func (m *MockIdentityRepo) CreateSnippet(ctx context.Context, s *Snippet) error {
		args := m.Called(ctx, s)
		return args.Error(0)
}

type MockGraphRepo struct{ mock.Mock }
func (m *MockGraphRepo) SaveVersion(ctx context.Context, id string, v *Version) error {
		args := m.Called(ctx, id, v)
		return args.Error(0)
}

func (m *MockGraphRepo) CiteSnippet(ctx context.Context, c *Citation) error {
		args := m.Called(ctx, c)
		return args.Error(0)
}

// Step 2: Test the Dual-Handshake
func TestCreateSnippet_Handshake(t *testing.T) {
		mockID := new(MockIdentityRepo)
		mockGraph := new(MockGraphRepo)
		service := NewSnippetService(mockID, mockGraph)

		ctx := context.Background()

		// We expect the brain to talk to BOTH repositories
		// 1. Identity repo handles the Snippet
    	mockID.On("CreateSnippet", ctx, mock.AnythingOfType("*core.Snippet")).Return(nil)
    	// 2. Graph repo handles the Version
    	mockGraph.On("SaveVersion", ctx, mock.Anything, mock.AnythingOfType("*core.Version")).Return(nil)

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