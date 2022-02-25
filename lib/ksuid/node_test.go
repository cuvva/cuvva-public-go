package ksuid

import (
	"context"
	"testing"
)

func BenchmarkGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Generate(context.Background(), "user")
	}
}
