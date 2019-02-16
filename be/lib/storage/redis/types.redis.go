package redis

// A Redis offers a standard interface for caching mechanism
type Redis interface {
	Get(string) *Result
	Setex(string, int, interface{}) error
	Expire(string, int) error
	Del(...string) error
	HSet(string, string, interface{}) error
	HGet(string, string) *Result
	HKeys(hash string) *Result
	HVals(hash string) *Result
	HGetAll(hash string) *Result
	Incr(...string) error
	Decr(...string) error
}

// Result struct
type Result struct {
	Value interface{}
	Error error
}
