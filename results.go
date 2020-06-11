// Copyright 2020 Andrew Quinn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofrac

type Result struct {
	Z          complex128
	C          complex128
	Iterations int
	NFactor    float64
}

// Results is a 2D slice of Result objects.
type Results struct {
	results       [][]Result
	maxIterations int
}

// NewResults constructs a 2D slice of Result objects, with outer and inner
// dimensions of rows and cols, respectively.
func NewResults(rows int, cols int, maxIterations int) Results {
	r := Results{maxIterations: maxIterations}
	r.results = make([][]Result, rows)
	for row := range r.results {
		r.results[row] = make([]Result, cols)
	}
	return r
}

// SetResult sets the z, c, and iterations fields of the Result located at
// the coordinates (row, col).
func (r Results) SetResult(row int, col int, z complex128, c complex128, iterations int) {
	r.results[row][col].Z = z
	r.results[row][col].C = c
	r.results[row][col].Iterations = iterations
}

// At retrieves the Result at the coordinates (row, col).
func (r Results) At(row int, col int) *Result {
	return &r.results[row][col]
}

// Dimensions returns the outer and inner dimensions of a Results object.
func (r Results) Dimensions() (rows int, cols int) {
	rows = len(r.results)
	if rows > 0 {
		cols = len(r.results[0])
	}
	return rows, cols
}

func calculateAccumulatedHistogram(r Results) (hist []int) {
	hist = make([]int, r.maxIterations)

	// regular histogram
	for row := range r.results {
		for col := range r.results[row] {
			n := r.results[row][col].Iterations
			hist[n]++
		}
	}

	// accumulate it
	for i, n := range hist {
		if i == 0 {
			continue
		}
		hist[i] = n + hist[i-1]
	}

	return hist
}

func setNFactors(r Results, hist []int) {
	invTotal := 1.0

	if r.maxIterations > 1 {
		lastDivergent := hist[len(hist)-2]
		// non-degenerate case
		if lastDivergent > 0 {
			invTotal = 1.0 / float64(lastDivergent)
		}
	}

	for row := range r.results {
		for col := range r.results[row] {
			result := r.At(row, col)
			i := result.Iterations

			// only escaped results are normalized
			if i < r.maxIterations-1 {
				result.NFactor = float64(hist[i]) * invTotal
			}
		}
	}
}

// Done finalizes a Results object and triggers calculations that depend on
// the entirety of a fractal solution.
func (r Results) Done() {
	hist := calculateAccumulatedHistogram(r)
	setNFactors(r, hist)
}
