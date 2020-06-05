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

	// SetMaxIterations sets the maximum iteration count.
	SetMaxIterations(n int)
}

type plotterData struct {
	maxIterations int
}

func (it *plotterData) SetMaxIterations(n int) {
	it.maxIterations = n
}

// EscapeTimePlotter implements the basic plotting method for fractals, which
// relies solely on the escape time (i.e., pre-bailout iteration count) of an
// iterated calculation.
type EscapeTimePlotter struct {
	plotterData
}

func (p EscapeTimePlotter) Plot(r *Result) float64 {
	return float64(r.Iterations)
}

var invLog2 = math.Log2E

// this currently only works for quadratic fractals
func smooth(val float64, z complex128, maxIt int) float64 {
	h := math.Log(cmplx.Abs(z)) / math.Log(float64(maxIt-1))
	return val - math.Log(h-1)*invLog2
}

// SmoothedEscapeTimePlotter maps a Result to a value in a way analogous to
// calculating the Result's electrostatic potential.
type SmoothedEscapeTimePlotter struct {
	plotterData
}

func (p SmoothedEscapeTimePlotter) Plot(r *Result) float64 {
	if r.Iterations == p.maxIterations-1 {
		return float64(r.Iterations)
	}
	return smooth(float64(r.Iterations), r.Z, p.maxIterations)
}

// NormalizedEscapeTimePlotter is a version of the escape time plotting method
// in which values are normalized according to the set of diverging results.
type NormalizedEscapeTimePlotter struct {
	plotterData
}

func (p NormalizedEscapeTimePlotter) Plot(r *Result) float64 {
	if r.Iterations == p.maxIterations-1 {
		return float64(r.Iterations)
	}
	return math.Floor(r.NFactor * float64(p.maxIterations-2))
}

// NormalizedSmoothedEscapeTimePlotter performs the electrostatic potential
// calculation on the normalized escape time results.
type NormalizedSmoothedEscapeTimePlotter struct {
	NormalizedEscapeTimePlotter
}

func (p NormalizedSmoothedEscapeTimePlotter) Plot(r *Result) float64 {
	val := p.NormalizedEscapeTimePlotter.Plot(r)
	if int(val) == p.NormalizedEscapeTimePlotter.maxIterations-1 {
		return val
	}

	return smooth(val, r.Z, p.maxIterations)
}
