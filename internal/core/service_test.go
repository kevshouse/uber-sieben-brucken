package core_test

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockIdentityRepo struct{ mock.Mock }
func (m *MockIdentityRepo) Save(ctx context.Context, s any) error { return m.Called(ctx, s).Error(0) }
func (m *MockIdentityRepo) CreateSnippet(ctx context.Context, s any) error { return m.Called(ctx, s).Error(0) }
func (m *MockIdentityRepo) Search(ctx context.Context, q string) (any, error) {
	args := m.Called(ctx, q)
	return args.Get(0), args.Error(1)
}
func (m *MockIdentityRepo) Close() error { return m.Called().Error(0) }

type MockGraphRepo struct{ mock.Mock }
func (m *MockGraphRepo) SyncNode(ctx context.Context, s any) error { return m.Called(ctx, s).Error(0) }
func (m *MockGraphRepo) SaveVersion(ctx context.Context, s, v any) error { return m.Called(ctx, s, v).Error(0) }
func (m *MockGraphRepo) CiteSnippet(ctx context.Context, c any) error { return m.Called(ctx, c).Error(0) }
func (m *MockGraphRepo) Close() error { return m.Called().Error(0) }