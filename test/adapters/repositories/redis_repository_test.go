package test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/kleankare/user-service/internal/adapters/repositories"
	"github.com/kleankare/user-service/internal/core/ports"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type redisRepositoryTestSuite struct {
	suite.Suite
	mockRedis *miniredis.Miniredis
	client    *redis.Client
	cacheRepo ports.CacheRepository
}

func (suite *redisRepositoryTestSuite) SetupSuite() {
	var err error
	suite.mockRedis, err = miniredis.Run()
	require.NoError(suite.T(), err)

	redisConfig := redis.Options{
		Addr: suite.mockRedis.Addr(),
	}
	suite.client = redis.NewClient(&redisConfig)
	_, err = suite.client.Ping(context.Background()).Result()
	require.NoError(suite.T(), err)

	suite.cacheRepo = repositories.NewRedisRepository(suite.client)
}

func (suite *redisRepositoryTestSuite) TearDownSuite() {
	suite.mockRedis.Close()
	suite.client.Close()
}

func TestRedisCacheRepository(t *testing.T) {
	suite.Run(t, new(redisRepositoryTestSuite))
}

func (suite *redisRepositoryTestSuite) TestGet_Success() {
	key := "test-key"
	expectedValue := map[string]interface{}{"foo": "bar"}
	data, _ := json.Marshal(expectedValue)
	suite.mockRedis.Set(key, string(data))

	var value map[string]interface{}
	err := suite.cacheRepo.Get(key, &value)
	suite.NoError(err)
	suite.Equal(expectedValue, value)
}

func (suite *redisRepositoryTestSuite) TestGet_CacheMiss() {
	key := "non-existent-key"
	var value string
	err := suite.cacheRepo.Get(key, &value)
	suite.EqualError(err, fmt.Sprintf("cache miss for key %q", key))
	suite.Empty(value)
}

func (suite *redisRepositoryTestSuite) TestGet_UnmarshalError() {
	key := "test-key"
	suite.mockRedis.Set(key, "invalid-json")

	var value map[string]interface{}
	err := suite.cacheRepo.Get(key, &value)
	suite.EqualError(err, fmt.Sprintf("failed to unmarshal cache value for key %q: invalid character 'i' looking for beginning of value", key))
	suite.Nil(value)
}

func (suite *redisRepositoryTestSuite) TestSet_Success() {
	key := "test-key"
	value := map[string]interface{}{"foo": "bar"}
	duration := 10 * time.Second

	err := suite.cacheRepo.Set(key, value, duration)

	suite.NoError(err)
	cachedValue, err := suite.mockRedis.Get(key)
	suite.NoError(err)

	var cachedValueMap map[string]interface{}
	err = json.Unmarshal([]byte(cachedValue), &cachedValueMap)
	suite.NoError(err)
	suite.Equal(value, cachedValueMap)
}

func (suite *redisRepositoryTestSuite) TestSet_MarshalError() {
	key := "test-key"
	value := make(chan int) // Unmarshalable value
	duration := 10 * time.Second

	err := suite.cacheRepo.Set(key, value, duration)

	suite.EqualError(err, fmt.Sprintf("failed to marshal cache value for key %q: %v", key, "json: unsupported type: chan int"))
}

func (suite *redisRepositoryTestSuite) TestDelete_Success() {
	key := "test-key"
	value := map[string]interface{}{"foo": "bar"}
	duration := 10 * time.Second

	err := suite.cacheRepo.Set(key, value, duration)
	suite.NoError(err)

	err = suite.cacheRepo.Delete(key)
	suite.NoError(err)

	_, err = suite.mockRedis.Get(key)
	suite.Error(err)
}

func (suite *redisRepositoryTestSuite) TestDelete_NonExistentKey() {
	key := "non-existent-key"

	err := suite.cacheRepo.Delete(key)

	suite.NoError(err)
}
