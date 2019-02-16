package localcache_test

import (
	"testing"

	. "github.com/template/be/lib/storage/localcache"
)

func BenchmarkSet(b *testing.B) {

	for n := 0; n < b.N; n++ {
		Set("a", 1, 11)
	}

}

func BenchmarkGet(b *testing.B) {

	// Set value
	Set("a", 1, 11)

	for n := 0; n < b.N; n++ {
		Get("a")
	}

}

func BenchmarkDelete(b *testing.B) {

	for n := 0; n < b.N; n++ {
		Delete("a")
	}

}
