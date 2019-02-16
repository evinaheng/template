package redigo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/template/be/lib/storage/redis/redigo"
)

func TestGet(t *testing.T) {

	cfg := Config{
		Endpoint: "null",
	}
	c := New(cfg)

	res := c.Get("test").String()
	assert.Equal(t, "", res)
}

func TestSetex(t *testing.T) {

	cfg := Config{
		Endpoint: "null",
	}
	c := New(cfg)

	assert.EqualError(t, c.Setex("test", 100, 1), "dial tcp: address null: missing port in address")
}
func TestDel(t *testing.T) {

	cfg := Config{
		Endpoint: "null",
	}
	c := New(cfg)

	assert.EqualError(t, c.Del("test"), "dial tcp: address null: missing port in address")
}

func TestExpire(t *testing.T) {

	cfg := Config{
		Endpoint: "null",
	}
	c := New(cfg)

	assert.EqualError(t, c.Expire("test", 123), "dial tcp: address null: missing port in address")
}

func TestIncr(t *testing.T) {

	cfg := Config{
		Endpoint: "null",
	}
	c := New(cfg)

	assert.EqualError(t, c.Incr("test"), "dial tcp: address null: missing port in address")
}
func TestDecr(t *testing.T) {

	cfg := Config{
		Endpoint: "null",
	}
	c := New(cfg)

	assert.EqualError(t, c.Decr("test"), "dial tcp: address null: missing port in address")
}
func TestHSet(t *testing.T) {

	cfg := Config{
		Endpoint: "null",
	}
	c := New(cfg)

	assert.EqualError(t, c.HSet("test", "bar", 100), "dial tcp: address null: missing port in address")
}
func TestHGet(t *testing.T) {

	cfg := Config{
		Endpoint: "null",
	}
	c := New(cfg)

	assert.Equal(t, "", c.HGet("test", "bar").String())
}

func TestHKeys(t *testing.T) {

	cfg := Config{
		Endpoint: "null",
	}
	c := New(cfg)

	assert.Equal(t, "", c.HKeys("test").String())
}

func TestHVals(t *testing.T) {

	cfg := Config{
		Endpoint: "null",
	}
	c := New(cfg)

	assert.Equal(t, "", c.HVals("test").String())
}

func TestHGetAll(t *testing.T) {

	cfg := Config{
		Endpoint: "null",
	}
	c := New(cfg)

	assert.Equal(t, []string(nil), c.HGetAll("test").StringSlice())
}
