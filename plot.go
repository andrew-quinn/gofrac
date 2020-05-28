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

var invLog2 = 1.0 / math.Log(2.0)

type SmoothedEscapeTimePlotter struct{}

func (p SmoothedEscapeTimePlotter) Plot(r *Result) float64 {
	if r.iterations < glob.maxIterations-1 {
		lz := math.Log(cmplx.Abs(r.z))
		nu := math.Log(lz*invLog2) * invLog2
		return float64(r.iterations+1) - nu
	} else {
		return float64(r.iterations)
	}

}
