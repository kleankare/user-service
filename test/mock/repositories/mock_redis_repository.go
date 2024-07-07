package repositories

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockRedisRepository struct {
	mock.Mock
}

func (m *MockRedisRepository) Get(key string, value interface{}) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockRedisRepository) Set(key string, value interface{}, duration time.Duration) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockRedisRepository) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}
