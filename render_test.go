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
	maxIterations int
}

func (p fakePlotter) Plot(r *gofrac.Result) float64 {
	return float64(r.Iterations % 2)
}

func (p *fakePlotter) SetMaxIterations(n int) {
	p.maxIterations = n
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
	r := gofrac.NewResults(rows, cols)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			r.SetResult(row, col, 0, 0, row)
		}
	}
	mockSampleColor = func(val float64, _ int) color.Color {
		if val == 0.0 {
			return color.Black
		}
		return color.White
	}

	bitmap := gofrac.Render(&r, &fakePlotter{10}, fakePalette{})

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
