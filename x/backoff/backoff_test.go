package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkNewBackoff(b *testing.B) {
	var back = NewBackoff()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		time.Sleep(back.Next(1))
	}
}

func TestNewBackoff(t *testing.T) {
	var back = NewBackoff()

	// default min is time.Millisecond*10
	delta := time.Millisecond * 10
	for i := 0; i < 10; i++ {
		if i >= 7 {
			// default max is time.Second*1
			assert.Equal(t, back.Next(1), time.Second)
		} else {
			assert.Equal(t, back.Next(1), delta)
		}
		delta *= 2
	}
}
