package core

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockIdentityRepo alignment
type MockIdentityRepo struct{ mock.Mock }
func (m *MockIdentityRepo) Save(ctx context.Context, s *Snippet) error { return m.Called(ctx, s).Error(0) }
func (m *MockIdentityRepo) Search(ctx context.Context, q string) ([]*Snippet, error) {
	args := m.Called(ctx, q)
	return args.Get(0).([]*Snippet), args.Error(1)
}
func (m *MockIdentityRepo) GetAll(ctx context.Context) ([]*Snippet, error) { // Added this!
	args := m.Called(ctx)
	return args.Get(0).([]*Snippet), args.Error(1)
}
func (m *MockIdentityRepo) Close() error { return m.Called().Error(0) }

// MockGraphRepo alignment
type MockGraphRepo struct{ mock.Mock }
func (m *MockGraphRepo) SyncNode(ctx context.Context, s *Snippet) error { return m.Called(ctx, s).Error(0) }
func (m *MockGraphRepo) SaveVersion(ctx context.Context, s *Snippet, v *Version) error { return m.Called(ctx, s, v).Error(0) }
func (m *MockGraphRepo) CiteSnippet(ctx context.Context, c *Citation) error { // Changed signature!
	return m.Called(ctx, c).Error(0) 
}
func (m *MockGraphRepo) Close() error { return m.Called().Error(0) }