package backoff

import (
	"math"
	"time"
)

const (
	defaultFactor = 2

	Done = -1
)

var defaultBackoff = newBackoff()

type backOff struct {
	factor             float64
	delayMin, delayMax float64
}

func newBackoff() *backOff {
	return &backOff{
		factor:   defaultFactor,
		delayMin: 10,
		delayMax: 1000,
	}
}

func NewBackoff() *backOff {
	return defaultBackoff.clone()
}

// Next
// Exponential
func (b *backOff) Next(delta int) time.Duration {
	r := b.delayMin * math.Pow(b.factor, float64(delta))
	if r > b.delayMax {
		return b.duration(b.delayMax)
	}

	if r < b.delayMin {
		return b.duration(b.delayMin)
	}

	return b.duration(r)
}

func (b *backOff) duration(t float64) time.Duration {
	return time.Microsecond * time.Duration(t)
}

func (b *backOff) clone() *backOff {
	cb := *b
	return &cb
}
