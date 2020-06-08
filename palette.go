// Copyright 2020 Andrew Quinn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofrac

import (
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"math"
)

// ColorSampler converts a floating point value to a color.Color in a color
// palette.
type ColorSampler interface {
	// TODO: Make "blackout" color configurable
	// SampleColor returns the color.Color of a palette corresponding to a
	// floating point value given by val.
	SampleColor(val float64, maxIterations int) color.Color
}

func isConvergent(val float64, maxIterations int) bool {
	return int(val) == maxIterations-1 || maxIterations <= 1
}

var black, _ = colorful.MakeColor(color.Black)

// SpectralPalette contains a range of hues from portion of the HSV color space
// where S and V are both 1.0. The palette starts at Offset degrees and
// travels Sweep degrees around the HSV space.
type SpectralPalette struct {
	Sweep  float64
	Offset float64
}

func (p SpectralPalette) SampleColor(val float64, maxIterations int) color.Color {
	if isConvergent(val, maxIterations) {
		return black
	}

	t := val / float64(maxIterations-1)
	h := t*p.Sweep + p.Offset
	return colorful.Hsl(h, 1.0, 0.5)
}

// BandedPalette is an alias for color.Palette, which itself is an alias for
// a slice of color.Color. It represents a collection of discrete color
// bands.
type BandedPalette color.Palette

// NewUniformBandedPalette constructs a BandedPalette containing all of the
// color.Color items given in colors.
func NewUniformBandedPalette(colors ...color.Color) BandedPalette {
	return colors
}

func (p BandedPalette) SampleColor(val float64, maxIterations int) color.Color {
	if isConvergent(val, maxIterations) {
		return black
	}

	i := int(float64(len(p)) * val / float64(maxIterations-1))
	clr, _ := colorful.MakeColor(p[i])
	return clr
}

// BlendedBandedPalette is a collection of interpolated color.Color items.
type BlendedBandedPalette color.Palette

// NewUniformBlendedBandedPalette constructs a BlendedBandedPalette containing
// all of the color.Color items given in colors.
func NewUniformBlendedBandedPalette(colors ...color.Color) BlendedBandedPalette {
	return BlendedBandedPalette(NewUniformBandedPalette(colors...))
}

func (p BlendedBandedPalette) SampleColor(val float64, maxIterations int) color.Color {
	if isConvergent(val, maxIterations) {
		return black
	}

	lastIdx := len(p) - 1
	t := val / float64(maxIterations-2)
	scaledVal := t * float64(lastIdx)

	lo := math.Floor(scaledVal)
	loIdx := int(lo)
	blendT := scaledVal - lo

	hi := math.Ceil(scaledVal)
	hiIdx := int(hi)

	if hiIdx > lastIdx {
		hiIdx = lastIdx
	}

	c1, _ := colorful.MakeColor(p[loIdx])
	c2, _ := colorful.MakeColor(p[hiIdx])

	c := c1.BlendLab(c2, blendT)
	clr, _ := colorful.MakeColor(c)
	return clr
}

// PeriodicPalette is a cyclic palette of discrete color bands whose width is
// given by Period.
type PeriodicPalette struct {
	BandedPalette
	Period int
}

func (p PeriodicPalette) SampleColor(val float64, maxIterations int) color.Color {
	if isConvergent(val, maxIterations) {
		return black
	}

	idx := (int(val) / p.Period) % len(p.BandedPalette)
	clr, _ := colorful.MakeColor(p.BandedPalette[idx])
	return clr
}
