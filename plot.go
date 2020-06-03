// Copyright 2020 Andrew Quinn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofrac

import (
	"math"
	"math/cmplx"
)

// Plotter can map the output of a fractal calculation to a floating point
// value in interesting ways.
type Plotter interface {
	// Plot maps a Result object onto a floating point number.
	Plot(r *Result) float64
}

// EscapeTimePlotter implements the basic plotting method for fractals, which
// relies solely on the escape time (i.e., pre-bailout iteration count) of an
// iterated calculation.
type EscapeTimePlotter struct{}

func (p EscapeTimePlotter) Plot(r *Result) float64 {
	return float64(r.Iterations)
}

var invLog2 = math.Log2E

func smoothVal(val float64, z complex128) float64 {
	lz := math.Log(cmplx.Abs(z))
	nu := math.Log(lz*invLog2) * invLog2
	return val + 1 - nu
}

// SmoothedEscapeTimePlotter maps a Result to a value in a way analogous to
// calculating the Result's electrostatic potential.
type SmoothedEscapeTimePlotter struct{}

func (p SmoothedEscapeTimePlotter) Plot(r *Result) float64 {
	if r.Iterations == glob.maxIterations-1 {
		return float64(r.Iterations)
	}
	return smoothVal(float64(r.Iterations), r.Z)
}

// NormalizedEscapeTimePlotter is a version of the escape time plotting method
// in which values are normalized according to the set of diverging results.
type NormalizedEscapeTimePlotter struct{}

func (p NormalizedEscapeTimePlotter) Plot(r *Result) float64 {
	if r.Iterations == glob.maxIterations-1 {
		return float64(r.Iterations)
	}
	return r.NFactor * float64(glob.maxIterations-1)
}

// NormalizedSmoothedEscapeTimePlotter performs the electrostatic potential
// calculation on the normalized escape time results.
type NormalizedSmoothedEscapeTimePlotter struct {
	NormalizedEscapeTimePlotter
}

func (p NormalizedSmoothedEscapeTimePlotter) Plot(r *Result) float64 {
	val := p.NormalizedEscapeTimePlotter.Plot(r)
	if int(val) == glob.maxIterations-1 {
		return val
	}

	return smoothVal(val, r.Z) - 1 + math.Log(math.Log(8000.0))*invLog2
}
