// Copyright 2020 Andrew Quinn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofrac_test

import (
	"github.com/cfdwalrus/gofrac"
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"testing"
)

type testCase struct {
	val   float64
	color color.Color
}

func cmp(t *testing.T, palette gofrac.ColorSampler, tc testCase, maxIterations int) {
	want, _ := colorful.MakeColor(tc.color)
	got := palette.SampleColor(tc.val, maxIterations)
	if want != got {
		t.Errorf("%T: val=%0.1f, want: %v, got: %v", palette, tc.val, want, got)
	}
}

func TestSpectralPalette_SampleColor(t *testing.T) {
	// full spectrum starting at red
	rainbow := gofrac.SpectralPalette{Sweep: 360.0}
	maxIt := 361
	rainbowTC := []testCase{
		{val: 0.0, color: color.RGBA{0xff, 0x00, 0x00, 0xff}},
		{val: 120.0, color: color.RGBA{0x00, 0xff, 0x00, 0xff}},
		{val: 240.0, color: color.RGBA{0x00, 0x00, 0xff, 0xff}},
		{val: 360.0, color: color.RGBA{0x00, 0x00, 0x00, 0xff}},
	}

	for _, tc := range rainbowTC {
		cmp(t, rainbow, tc, maxIt)
	}
}

func TestBandedPalette_SampleColor(t *testing.T) {
	bands := gofrac.NewUniformBandedPalette(
		color.RGBA{0xff, 0x00, 0x00, 0xff},
		color.RGBA{0x00, 0xff, 0x00, 0xff},
		color.RGBA{0x00, 0x00, 0xff, 0xff},
	)

	maxIt := 100
	bandsTC := []testCase{
		{val: 0.0, color: color.RGBA{0xff, 0x00, 0x00, 0xff}},
		{val: 1.0, color: color.RGBA{0xff, 0x00, 0x00, 0xff}},

		{val: 33.0, color: color.RGBA{0x00, 0xff, 0x00, 0xff}},
		{val: 65.0, color: color.RGBA{0x00, 0xff, 0x00, 0xff}},

		{val: 66.0, color: color.RGBA{0x00, 0x00, 0xff, 0xff}},

		{val: 99.0, color: color.RGBA{0x00, 0x00, 0x00, 0xff}},
	}

	for _, tc := range bandsTC {
		cmp(t, bands, tc, maxIt)
	}
}

func TestBlendedBandedPalette_SampleColor(t *testing.T) {
	blends := gofrac.NewUniformBlendedBandedPalette(
		color.RGBA{0xff, 0x00, 0x00, 0xff},
		color.RGBA{0x00, 0xff, 0x00, 0xff},
		color.RGBA{0x00, 0x00, 0xff, 0xff},
	)

	maxIT := 6
	blendsTC := []testCase{
		{val: 0.0, color: color.RGBA{0xff, 0x00, 0x00, 0xff}},
		{val: 2.0, color: color.RGBA{0x00, 0xff, 0x00, 0xff}},
		{val: 4.0, color: color.RGBA{0x00, 0x00, 0xff, 0xff}},
		{val: 5.0, color: color.RGBA{0x00, 0x00, 0x00, 0xff}},
	}
	for _, tc := range blendsTC {
		cmp(t, blends, tc, maxIT)
	}
}

func TestPeriodicPalette_SampleColor(t *testing.T) {
	bw := gofrac.PeriodicPalette{
		BandedPalette: gofrac.NewUniformBandedPalette(
			color.RGBA{0x00, 0x00, 0x00, 0xff},
			color.RGBA{0xff, 0xff, 0xff, 0xff},
		),
		Period: 1,
	}

	maxIt := 22
	bwTC := make([]testCase, maxIt-1)
	// even black, odd white
	for i := range bwTC {
		if i%2 == 0 {
			bwTC[i] = testCase{val: float64(i), color: color.RGBA{0x00, 0x00, 0x00, 0xff}}
		} else {
			bwTC[i] = testCase{val: float64(i), color: color.RGBA{0xff, 0xff, 0xff, 0xff}}
		}
	}
	bwTC[len(bwTC)-1].color = color.RGBA{0x00, 0x00, 0x00, 0xff} // convergent color is always black -- that might become variable in the future

	for _, tc := range bwTC {
		cmp(t, bw, tc, maxIt)
	}
}
