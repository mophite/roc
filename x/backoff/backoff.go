// Copyright (c) 2021 roc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

package backoff

import (
    "math"
    "time"
)

//Truncated Binary Exponential Back—off,TBEB
//After the collision site stops sending, it will not send data again immediately,
//but back off for a random time,
//reducing the probability of collision during retransmission.
const defaultFactor = 2

var defaultBackoff = newBackoff()

type backOff struct {
    factor, delayMin, delayMax float64
    Attempts                   int
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

func WitchBackoff(factor, delayMin, delayMax float64, attempts int) *backOff {
    b := defaultBackoff.clone()
    b.factor = factor
    b.delayMax = delayMax
    b.delayMin = delayMin
    b.Attempts = attempts
    return b
}

// Next
// Exponential
func (b *backOff) Next(delta int) time.Duration {
    r := b.delayMin * math.Pow(b.factor, float64(b.Attempts))
    b.Attempts += delta
    if r > b.delayMax {
        return b.duration(b.delayMax)
    }

    if r < b.delayMin {
        return b.duration(b.delayMin)
    }

    return b.duration(r)
}

func (b *backOff) duration(t float64) time.Duration {
    return time.Millisecond * time.Duration(t)
}

func (b *backOff) clone() *backOff {
    cb := *b
    return &cb
}
