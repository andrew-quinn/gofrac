// Copyright 2020 Andrew Quinn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofrac_test

import (
	"github.com/cfdwalrus/gofrac"
	"math"
	"testing"
)

func TestEscapeTimePlotter_Plot(t *testing.T) {
	var p gofrac.EscapeTimePlotter
	for i := 0; i < 32; i++ {
		it := 1 << uint(i)
		// plotting should only depend on the Iterations field
		var fakeResult = &gofrac.Result{
			Z:          42 + 17i,
			C:          13 + 37i,
			Iterations: it,
			NFactor:    0.12345,
		}
		want := float64(fakeResult.Iterations)
		got := p.Plot(fakeResult)
		if got != want {
			t.Errorf("%T: want: %0.2f, got: %0.2f", p, want, got)
		}
	}
}

// TODO: Make this not be terrible
func TestSmoothedEscapeTimePlotter_Plot(t *testing.T) {
	maxIt := 10

	f := gofrac.FracData{
		Radius:        4,
		MaxIterations: maxIt,
	}
	f.SetDegree(2)

	var p gofrac.SmoothedEscapeTimePlotter
	p.SetFracData(&f)

	converged := gofrac.Result{Iterations: maxIt - 1}
	want := float64(maxIt - 1)
	got := p.Plot(&converged)
	if got != want {
		t.Errorf("%T: want: %0.2f, got: %0.2f", p, want, got)
	}
}

func TestNormalizedEscapeTimePlotter_Plot(t *testing.T) {
	maxIt := 11
	//s := 1 / float64(maxIt-1)
	r := gofrac.NewResults(1, maxIt, maxIt)
	rows, cols := r.Dimensions()
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			r.SetResult(row, col, 0, 0, col)
		}
	}
	r.Done()

	f := gofrac.FracData{}
	f.SetRadius(0)
	f.SetDegree(2)
	f.SetMaxIterations(maxIt)
	var p gofrac.NormalizedEscapeTimePlotter
	p.SetFracData(&f)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			want := float64(col)
			got := p.Plot(r.At(row, col))
			if math.Abs(got-want) > 0.00001 {
				t.Errorf("%T: want: %0.5f, got: %0.5f", p, want, got)
			}
		}
	}
}

// TODO: Soooo, there might be some work to do here
func TestNormalizedSmoothedEscapeTimePlotter_Plot(t *testing.T) {

}
