package gofrac

import (
	"math"
	"math/cmplx"
)

type Plotter interface {
	Plot(r *Result) float64
}

type EscapeTimePlotter struct{}

func (p EscapeTimePlotter) Plot(r *Result) float64 {
	return float64(r.iterations)
}

var invLog2 = math.Log2E

func smoothVal(val float64, z complex128) float64 {
	lz := math.Log(cmplx.Abs(z))
	nu := math.Log(lz*invLog2) * invLog2
	return val + 1 - nu
}

type SmoothedEscapeTimePlotter struct{}

func (p SmoothedEscapeTimePlotter) Plot(r *Result) float64 {
	if r.iterations == glob.maxIterations-1 {
		return float64(r.iterations)
	}
	return smoothVal(float64(r.iterations), r.z)
}

type NormalizedEscapeTimePlotter struct{}

func (p NormalizedEscapeTimePlotter) Plot(r *Result) float64 {
	if r.iterations == glob.maxIterations-1 {
		return float64(r.iterations)
	}
	return r.nFactor * float64(glob.maxIterations-1)
}

type NormalizedSmoothedEscapeTimePlotter struct {
	NormalizedEscapeTimePlotter
}

func (p NormalizedSmoothedEscapeTimePlotter) Plot(r *Result) float64 {
	val := p.NormalizedEscapeTimePlotter.Plot(r)
	if int(val) == glob.maxIterations-1 {
		return val
	}

	return smoothVal(val, r.z) - 1 + math.Log(math.Log(8000.0)) * invLog2
}
