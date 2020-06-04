// Copyright 2020 Andrew Quinn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofrac_test

import (
	"github.com/cfdwalrus/gofrac"
	"testing"
)

// mocks

type fakeDomain struct{}

func (d fakeDomain) At(i int, j int) (z complex128, err error) {
	return 0, nil
}

var dimensionsMock func() (int, int)

func (d fakeDomain) Dimensions() (rows int, cols int) {
	return dimensionsMock()
}

type fakeFrac struct{}

var mockFrac func(complex128) *gofrac.Result

func (f fakeFrac) Frac(loc complex128) *gofrac.Result {
	return mockFrac(loc)
}

// real stuff

func TestFracIt(t *testing.T) {
	mockFrac = func(complex128) *gofrac.Result {
		return &gofrac.Result{}
	}

	for _, it := range []int{-1, 0} {
		// malformed domain
		dimensionsMock = func() (int, int) {
			return it, it
		}
		_, err := gofrac.FracIt(fakeDomain{}, fakeFrac{}, 1)
		if err == nil {
			t.Errorf("Error not caught for malformed domain")
		}
	}

	// well-formed domain
	dimensionsMock = func() (int, int) {
		return 10, 10
	}

	// bad iterations argument
	for _, it := range []int{-1, 0} {
		_, err := gofrac.FracIt(fakeDomain{}, fakeFrac{}, it)
		if err == nil {
			t.Errorf("Error not caught for bad iteration count")
		}
	}

	// degenerate case
	degenerateIterations := 5
	degenerateResult := gofrac.Result{
		Z:          0,
		C:          0,
		Iterations: degenerateIterations - 1,
		NFactor:    0,
	}
	mockFrac = func(complex128) *gofrac.Result {
		return &degenerateResult
	}
	rDegenerate, err := gofrac.FracIt(fakeDomain{}, fakeFrac{}, degenerateIterations)
	if err != nil {
		t.Error(err)
	}
	for _, row := range *rDegenerate {
		for _, result := range row {
			if result != degenerateResult {
				t.Errorf("%T: Unexpected Result for degenerate case. wanted: %v, got: %v", rDegenerate, degenerateResult, result)
			}
		}
	}

	// normal operation
	normalResult := gofrac.Result{
		Z:          0,
		C:          0,
		Iterations: 4, // all diverge in less than maxIterations-1 time
		NFactor:    1, // since all results have the same value for Iterations
	}
	mockFrac = func(complex128) *gofrac.Result {
		return &normalResult
	}
	r, err := gofrac.FracIt(fakeDomain{}, fakeFrac{}, 10)
	if err != nil {
		t.Error(err)
	}
	for _, row := range *r {
		for _, result := range row {
			if result != normalResult {
				t.Errorf("%T: Unexpected Result for normal case. wanted: %v, got: %v", rDegenerate, degenerateResult, result)
			}
		}
	}
}
