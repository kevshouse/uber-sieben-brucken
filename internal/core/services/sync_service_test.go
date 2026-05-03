package services_test

import (
	"context"
	"testing"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
	"github.com/kevshouse/uber-sieben-brucken/internal/core/services"

	"github.com/stretchr/testify/mock"
)

func TestSyncService_CreateGenesis(t *testing.T) {
	t.Run("it should create a snippet and its genesis version securely", func(t *testing.T) {
		// 1. Setup Testify Mocks
		idRepo := new(core.MockIdentityRepo)
		graphRepo := new(core.MockGraphRepo)

		// 2. Define Expectations
		idRepo.On("Save", mock.Anything, mock.AnythingOfType("*core.Snippet")).Return(nil)
		graphRepo.On("SaveVersion", mock.Anything, mock.AnythingOfType("*core.Snippet"), mock.AnythingOfType("*core.Version")).Return(nil)

		srv := services.NewSyncService(idRepo, graphRepo)

		// 3. Act
		snippet, err := srv.CreateGenesis(context.Background(), "Genesis Test", "owner-123", "fmt.Println('Genesis')")

		// 4. Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if snippet.Title != "Genesis Test" {
			t.Errorf("Expected title 'Genesis Test', got %s", snippet.Title)
		}
		if snippet.CreatedAt.IsZero() {
			t.Error("Temporal Failure: CreatedAt was not set")
		}

		idRepo.AssertExpectations(t)
		graphRepo.AssertExpectations(t)
	})
}

// Refactoring your older tests to also use the new mocks
func TestSyncService_Sync(t *testing.T) {
	t.Run("it should reach the graph shore during sync", func(t *testing.T) {
		idRepo := new(core.MockIdentityRepo)
		graphRepo := new(core.MockGraphRepo)
		
		// Expect SyncNode to be called
		graphRepo.On("SyncNode", mock.Anything, mock.AnythingOfType("*core.Snippet")).Return(nil)
		
		srv := services.NewSyncService(idRepo, graphRepo)
		err := srv.Sync(context.Background(), &core.Snippet{ID: "123"})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		graphRepo.AssertExpectations(t)
	})
}