// Copyright 2020 Andrew Quinn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofrac

import (
	"image/color"
	"runtime"
	"sync"
)

// bitmap stores a 2D field of color.Color that can be used to generate images.
type bitmap [][]color.Color

// NewBitmap constructs a slice of length r of slices of color.Colors, each of
// which are of length c.
func NewBitmap(r int, c int) bitmap {
	b := make(bitmap, r)
	for r := range b {
		b[r] = make([]color.Color, c)
	}
	return b
}

// Render combines the fractal iteration results with a plotting method and
// generates a bitmap according to the color palette provided.
func Render(results *Results, plotter Plotter, palette ColorSampler) bitmap {
	rows, cols := results.Dimensions()
	bitmap := NewBitmap(rows, cols)

	rowJobs := make(chan int, rows)

	numWorkers := runtime.NumCPU()
	wg := sync.WaitGroup{}
	for worker := 0; worker < numWorkers; worker++ {
		wg.Add(1)
		go func() {
			for row := range rowJobs {
				for col := 0; col < cols; col++ {
					result := results.At(row, col)
					val := plotter.Plot(result)
					bitmap[row][col] = palette.SampleColor(val, results.maxIterations)
				}
			}
			wg.Done()
		}()
	}

	for row := 0; row < rows; row++ {
		rowJobs <- row
	}

	close(rowJobs)
	wg.Wait()

	return bitmap
}
