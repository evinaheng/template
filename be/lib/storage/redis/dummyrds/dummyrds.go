package dummyrds

import (
	"errors"
	"fmt"
	"strings"

	"github.com/template/be/lib/storage/redis"
)

// New module for mocking
func New(config Config) redis.Redis {
	return &dummydis{
		config: config,
	}
}

// AddMock to mocking mapping
func (m Mocker) AddMock(command string, result interface{}, isError bool) {
	m[command] = mock{
		IsError: isError,
		Result:  result,
	}
}

func (c *dummydis) Get(key string) *redis.Result {
	return c.mock("GET " + key)
}

func (c *dummydis) Setex(key string, expireTime int, value interface{}) error {
	return c.mock(fmt.Sprintf("SETEX %s %d %s", key, expireTime, value)).Error
}

func (c *dummydis) Del(keys ...string) error {
	k := strings.Join(keys, " ")
	return c.mock("DEL " + k).Error
}

func (c *dummydis) Expire(key string, seconds int) error {
	return c.mock(fmt.Sprintf("EXPIRE %s %d", key, seconds)).Error
}

func (c *dummydis) Incr(keys ...string) error {
	k := strings.Join(keys, " ")
	return c.mock("INCR " + k).Error
}

func (c *dummydis) Decr(keys ...string) error {
	k := strings.Join(keys, " ")
	return c.mock("DECR " + k).Error
}

func (c *dummydis) HSet(key, field string, value interface{}) error {
	return c.mock(fmt.Sprintf("HSET %s %s %s", key, field, value)).Error
}

func (c *dummydis) HGet(key, field string) *redis.Result {
	return c.mock(fmt.Sprintf("HGET %s %s", key, field))
}

func (c *dummydis) HKeys(hash string) *redis.Result {
	return c.mock(fmt.Sprintf("HKEYS %s", hash))
}

func (c *dummydis) HVals(hash string) *redis.Result {
	return c.mock(fmt.Sprintf("HVALS %s", hash))
}

func (c *dummydis) HGetAll(hash string) *redis.Result {
	return c.mock(fmt.Sprintf("HGETALL %s", hash))
}

func (c *dummydis) mock(command string) *redis.Result {

	res, ok := c.config.MockingMap[command]

	// Mock not found
	if !ok {
		return &redis.Result{
			Error: errors.New("No mocking found for " + command),
		}
	}

	// Mock error
	if res.IsError {
		return &redis.Result{
			Error: fmt.Errorf("%v", res.Result),
		}
	}

	// Mock result
	return &redis.Result{
		Value: res.Result,
	}

}
