// Copyright 2020 Andrew Quinn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofrac_test

import (
	"github.com/cfdwalrus/gofrac"
	"testing"
)

func TestNewDomain(t *testing.T) {
	isValidSamples := func(x int, y int) bool {
		return x > 0 && y > 0
	}

	tc := []int{-1, 0, 1}
	for _, x := range tc {
		for _, y := range tc {
			d, err := gofrac.NewDomain(0, 0, 1, 1, x, y)
			switch {
			case !isValidSamples(x, y) && err == nil:
				t.Errorf("%T: want: err != nil, got: err == nil", d)
			case isValidSamples(x, y) && err != nil:
				t.Errorf("%T: want: err == nil, got err != nil", d)
			}
		}
	}
}

func TestDomain_At(t *testing.T) {
	maxSamples := 10
	isValidIndices := func(i int, j int) bool {
		return i >= 0 && i < maxSamples && j >= 0 && j < maxSamples
	}

	tc := []int{-1, 0, 5, 9, maxSamples, 100 * maxSamples}
	d, _ := gofrac.NewDomain(0, 0, 1, 1, maxSamples, maxSamples)
	for _, i := range tc {
		for _, j := range tc {
			z, err := d.At(i, j)
			switch {
			case !isValidIndices(i, j) && err == nil:
				// false negative error
				t.Errorf("%T: (i,j) = (%d, %d): want err != nil, got err == nil", d, i, j)
			case isValidIndices(i, j) && err != nil:
				// false positive error
				t.Errorf("%T: (i,j) = (%d, %d): want err == nil, got err != nil", d, i, j)
			case !isValidIndices(i, j) && err != nil:
				// error properly caught
				continue
			default:
				// valid indices yield a complex number interpolated from the inputs of NewDomain
				reZ := float64(i) / float64(maxSamples)
				imZ := float64(j) / float64(maxSamples)
				if real(z) != reZ || 1.0-imag(z) != imZ {
					t.Errorf("%T: (i, j) = (%d, %d): want %f + %fi, got %v", d, i, j, reZ, imZ, z)
				}
			}
		}
	}
}

func TestDomain_Dimensions(t *testing.T) {
	tc := []int{1, 100, 100000}
	for _, xSamples := range tc {
		for _, ySamples := range tc {
			d, _ := gofrac.NewDomain(0, 0, 1, 1, xSamples, ySamples)
			if rows, cols := d.Dimensions(); rows != ySamples || cols != xSamples {
				t.Errorf("%T: want: rows = %d, cols = %d, got: rows = %d, cols = %d", d, ySamples, xSamples, rows, cols)
			}
		}
	}
}
