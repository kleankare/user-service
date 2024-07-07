package ports

import "time"

type CacheRepository interface {
	Get(key string, value interface{}) error
	Set(key string, value interface{}, duration time.Duration) error
	Delete(key string) error
}
