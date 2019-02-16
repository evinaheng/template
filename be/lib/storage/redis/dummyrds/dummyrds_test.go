package dummyrds_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/template/be/lib/storage/redis/dummyrds"
)

func TestGet(t *testing.T) {

	m := Mocker{}
	m.AddMock("GET foo", "result", false)

	rds := New(Config{
		MockingMap: m,
	})

	assert.Equal(t, "result", rds.Get("foo").String())
	assert.Equal(t, "", rds.Get("none").String())
}

func TestSetex(t *testing.T) {

	m := Mocker{}
	m.AddMock("SETEX foo 10 bar", "1", false)
	m.AddMock("SETEX err 10 bar", "failed", true)

	rds := New(Config{
		MockingMap: m,
	})

	assert.Nil(t, rds.Setex("foo", 10, "bar"))
	assert.EqualError(t, rds.Setex("err", 10, "bar"), "failed")
}

func TestDel(t *testing.T) {

	m := Mocker{}
	m.AddMock("DEL foo", "1", false)
	m.AddMock("DEL bar", "failed", true)

	rds := New(Config{
		MockingMap: m,
	})

	assert.Nil(t, rds.Del("foo"))
	assert.EqualError(t, rds.Del("bar"), "failed")
}
func TestExpire(t *testing.T) {

	m := Mocker{}
	m.AddMock("EXPIRE foo 10", "1", false)
	m.AddMock("EXPIRE err 10", "failed", true)

	rds := New(Config{
		MockingMap: m,
	})

	assert.Nil(t, rds.Expire("foo", 10))
	assert.EqualError(t, rds.Expire("err", 10), "failed")
}

func TestIncr(t *testing.T) {

	m := Mocker{}
	m.AddMock("INCR foo", "1", false)
	m.AddMock("INCR bar", "failed", true)

	rds := New(Config{
		MockingMap: m,
	})

	assert.Nil(t, rds.Incr("foo"))
	assert.EqualError(t, rds.Incr("bar"), "failed")
}

func TestDecr(t *testing.T) {

	m := Mocker{}
	m.AddMock("DECR foo", "1", false)
	m.AddMock("DECR bar", "failed", true)

	rds := New(Config{
		MockingMap: m,
	})

	assert.Nil(t, rds.Decr("foo"))
	assert.EqualError(t, rds.Decr("bar"), "failed")
}

func TestHSet(t *testing.T) {

	m := Mocker{}
	m.AddMock("HSET foo bar res", "1", false)
	m.AddMock("HSET err bar res", "failed", true)

	rds := New(Config{
		MockingMap: m,
	})

	assert.Nil(t, rds.HSet("foo", "bar", "res"))
	assert.EqualError(t, rds.HSet("err", "bar", "res"), "failed")
}

func TestHGet(t *testing.T) {

	m := Mocker{}
	m.AddMock("HGET foo bar", "1", false)
	m.AddMock("HGET err bar", "failed", true)

	rds := New(Config{
		MockingMap: m,
	})

	assert.Equal(t, 1, rds.HGet("foo", "bar").Int())
	assert.Equal(t, "", rds.HGet("err", "bar").String())
}

func TestHKeys(t *testing.T) {

	m := Mocker{}
	m.AddMock("HKEYS foo", []string{"one", "two"}, false)
	m.AddMock("HKEYS err", "", true)

	rds := New(Config{
		MockingMap: m,
	})

	assert.Len(t, rds.HKeys("foo").StringSlice(), 2)
	assert.Equal(t, "", rds.HKeys("err").String())
}

func TestHVals(t *testing.T) {

	m := Mocker{}
	m.AddMock("HVALS foo", []string{"one", "two"}, false)
	m.AddMock("HVALS err", "", true)

	rds := New(Config{
		MockingMap: m,
	})

	assert.Len(t, rds.HVals("foo").StringSlice(), 2)
	assert.Equal(t, "", rds.HVals("err").String())
}

func TestHGetALl(t *testing.T) {

	m := Mocker{}
	m.AddMock("HGETALL foo", []string{"one", "two"}, false)
	m.AddMock("HGETALL err", "", true)

	rds := New(Config{
		MockingMap: m,
	})

	assert.Len(t, rds.HGetAll("foo").StringSlice(), 2)
	assert.Equal(t, "", rds.HGetAll("err").String())
}
