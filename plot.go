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

	// SetFracData sets Plotter's FracData member, which stores information
	// about an iterate, which is necessary for some plotting functions.
	SetFracData(fd *FracData)
}

type PlotterBase struct {
	FracData
}

func (pb *PlotterBase) plot(r *Result, pFunc func(r *Result) float64) float64 {
	if r.Iterations == pb.MaxIterations-1 {
		return float64(r.Iterations)
	}
	return pFunc(r)
}

func (pb *PlotterBase) SetFracData(fd *FracData) {
	pb.FracData = *fd
}

// EscapeTimePlotter implements the basic plotting method for fractals, which
// relies solely on the escape time (i.e., pre-bailout iteration count) of an
// iterated calculation.
type EscapeTimePlotter struct {
	PlotterBase
}

func (p EscapeTimePlotter) Plot(r *Result) float64 {
	return float64(r.Iterations)
}

func smooth(val float64, z complex128, d FracDataGetter) float64 {
	mod := cmplx.Abs(z)
	if mod == 0 {
		return val
	}
	lgBase := d.Data().logDegreeInv
	return val + 1 - math.Log(math.Log(mod))*lgBase
}

// SmoothedEscapeTimePlotter maps a Result to a value in a way analogous to
// calculating the Result's electrostatic potential.
type SmoothedEscapeTimePlotter struct {
	PlotterBase
}

func (p SmoothedEscapeTimePlotter) Plot(r *Result) float64 {
	return p.plot(r, func(r *Result) float64 {
		return smooth(float64(r.Iterations), r.Z, &p)
	})
}

// NormalizedEscapeTimePlotter is a version of the escape time plotting method
// in which values are normalized according to the set of diverging results.
type NormalizedEscapeTimePlotter struct {
	PlotterBase
}

func (p NormalizedEscapeTimePlotter) Plot(r *Result) float64 {
	return p.plot(r, func(r *Result) float64 {
		return math.Floor(r.NFactor * float64(p.MaxIterations-2))
	})
}

// NormalizedSmoothedEscapeTimePlotter performs the electrostatic potential
// calculation on the normalized escape time results.
type NormalizedSmoothedEscapeTimePlotter struct {
	NormalizedEscapeTimePlotter
}

func (p NormalizedSmoothedEscapeTimePlotter) Plot(r *Result) float64 {
	return p.plot(r, func(r *Result) float64 {
		return smooth(p.NormalizedEscapeTimePlotter.Plot(r), r.Z, p.Data())
	})
}

// PhasePlotter plots an escaped iterate according to its phase when it
// escapes.
type PhasePlotter struct {
	PlotterBase
}

func (p PhasePlotter) Plot(r *Result) float64 {
	return p.plot(r, func(r *Result) float64 {
		return float64(p.MaxIterations-1) * (cmplx.Phase(r.Z) + math.Pi) / (2 * math.Pi)
	})
}
