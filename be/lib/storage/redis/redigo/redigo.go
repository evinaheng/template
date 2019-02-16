package redigo

import (
	"time"

	rgo "github.com/gomodule/redigo/redis"
	"github.com/template/be/lib/storage/redis"
)

// New redis redigo module
func New(config Config) redis.Redis {

	// Set default 10 seconds timeout
	if config.Timeout == 0 {
		config.Timeout = 10
	}

	// Open connection to redis server
	return &credis{
		pool: rgo.Pool{
			MaxIdle:     config.MaxIdle,
			IdleTimeout: time.Duration(config.Timeout) * time.Second,
			Dial: func() (rgo.Conn, error) {
				return rgo.Dial(
					"tcp",
					config.Endpoint,
				)
			},
		},
	}

}

func (c *credis) Get(key string) *redis.Result {
	args := append([]interface{}{key})
	return c.cmd("GET", args...)
}

func (c *credis) Setex(key string, expireTime int, value interface{}) error {
	args := append([]interface{}{key, expireTime, value})
	return c.cmd("SETEX", args...).Error
}

func (c *credis) Del(keys ...string) error {
	args := make([]interface{}, len(keys))
	for i, v := range keys {
		args[i] = v
	}
	return c.cmd("DEL", args...).Error
}

func (c *credis) Expire(key string, seconds int) error {
	args := append([]interface{}{key, seconds})
	return c.cmd("EXPIRE", args...).Error
}

func (c *credis) Incr(keys ...string) error {
	args := make([]interface{}, len(keys))
	for i, v := range keys {
		args[i] = v
	}
	return c.cmd("INCR", args...).Error
}

func (c *credis) Decr(keys ...string) error {
	args := make([]interface{}, len(keys))
	for i, v := range keys {
		args[i] = v
	}
	return c.cmd("DECR", args...).Error
}

func (c *credis) HSet(key, field string, value interface{}) error {
	args := append([]interface{}{key, field, value})
	return c.cmd("HSET", args...).Error
}

func (c *credis) HGet(key, field string) *redis.Result {
	args := []interface{}{key, field}
	return c.cmd("HGET", args...)
}

func (c *credis) HKeys(key string) *redis.Result {
	args := []interface{}{key}
	return c.cmd("HKEYS", args...)
}

func (c *credis) HVals(key string) *redis.Result {
	args := []interface{}{key}
	return c.cmd("HVALS", args...)
}

func (c *credis) HGetAll(key string) *redis.Result {
	args := []interface{}{key}
	return c.cmd("HGETALL", args...)
}

func (c *credis) cmd(command string, args ...interface{}) *redis.Result {
	result := &redis.Result{}
	conn := c.pool.Get()
	defer conn.Close()

	data, err := conn.Do(command, args...)
	if err != nil {

		// Retry mechanism
		conn2 := c.pool.Get()
		defer conn2.Close()
		data, err = conn2.Do(command, args...)
		if err != nil {
			result.Error = err
			return result
		}
	}
	result.Value = data

	return result
}
