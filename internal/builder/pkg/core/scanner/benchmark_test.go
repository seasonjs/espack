package scanner

import (
	"strings"
	"testing"
)

//goos: darwin
//goarch: amd64
//pkg: github.com/seasonjs/espack/internal/builder/pkg/core/scanner
//cpu: VirtualApple @ 2.50GHz
//BenchmarkBuilder
//BenchmarkBuilder-10      	 8131868	       141.4 ns/op
//BenchmarkNewString
//BenchmarkNewString-10    	 8294235	       136.9 ns/op

func benchmark(b *testing.B, f func(int, []byte) string) {
	buf := make([]byte, 1024)
	str := strings.NewReader(`
		// this is an Comment 
		function a(){}
		`)
	str.Read(buf)
	for i := 0; i < b.N; i++ {
		f(10000, buf)
	}
}

func BenchmarkBuilder(b *testing.B) {
	benchmark(b, func(i int, s []byte) string {
		var builder strings.Builder
		builder.Write(s)
		return builder.String()

	})
}
func BenchmarkNewString(b *testing.B) {
	benchmark(b, func(i int, bytes []byte) string {
		var _ strings.Builder
		return string(bytes)
	})
}
