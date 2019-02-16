package redis

import (
	rgo "github.com/gomodule/redigo/redis"
	"github.com/template/be/lib/convert"
)

// Int result convertion type
func (r *Result) Int() int {
	if r.Error != nil {
		return 0
	}

	return convert.ToInt(r.Value)

}

// Int64 result convertion type
func (r *Result) Int64() int64 {
	if r.Error != nil {
		return 0
	}

	return convert.ToInt64(r.Value)
}

// String result convertion type
func (r *Result) String() string {
	if r.Error != nil {
		return ""
	}

	return convert.ToString(r.Value)

}

// Bytes result convertion type
func (r *Result) Bytes() []byte {
	if r.Error != nil {
		return nil
	}

	res := convert.ToByteArr(r.Value)
	if len(res) == 0 {
		return nil
	}

	return res

}

// ByteSlice result convertion type
func (r *Result) ByteSlice() [][]byte {
	if r.Error != nil {
		return nil
	}

	val, _ := rgo.ByteSlices(r.Value, nil)
	return val
}

// StringSlice result convertion type
func (r *Result) StringSlice() []string {
	if r.Error != nil {
		return nil
	}

	if val, ok := r.Value.([]string); ok {
		return val
	}

	val, _ := rgo.Strings(r.Value, nil)
	return val
}

// Float64 result convertion type
func (r *Result) Float64() float64 {
	if r.Error != nil {
		return 0
	}

	return convert.ToFloat64(r.Value)
}
