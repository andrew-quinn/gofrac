// Copyright 2020 Andrew Quinn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofrac_test

import (
	"github.com/cfdwalrus/gofrac"
	"image/color"
	"testing"
)

// mocks

type fakePlotter struct {
	gofrac.PlotterBase
}

func (p fakePlotter) Plot(r *gofrac.Result) float64 {
	return float64(r.Iterations % 2)
}

type fakePalette struct{}

var mockSampleColor func(float64, int) color.Color

func (p fakePalette) SampleColor(val float64, maxIterations int) color.Color {
	return mockSampleColor(val, maxIterations)
}

// real stuff

func TestRender(t *testing.T) {
	rows := 10
	cols := 6
	fakeResults := gofrac.NewResults(rows, cols, rows)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			fakeResults.SetResult(row, col, 0, 0, row)
		}
	}

	mockSampleColor = func(val float64, _ int) color.Color {
		if val == 0.0 {
			return color.Black
		}
		return color.White
	}

	// feed Render fakePlotter, which maps fakeResults to alternating {0, 1}
	// values, then map those values to black and white, respectively, with
	// fakePalette
	bitmap := gofrac.Render(&fakeResults, &fakePlotter{}, fakePalette{})

	for row := 0; row < rows; row++ {
		c := mockSampleColor(float64(row%2), 0)
		for col := 0; col < cols; col++ {
			want := c
			got := bitmap[row][col]
			if want != got {
				t.Errorf("Render: want: %v, got: %v", want, got)
			}
		}
	}
}
