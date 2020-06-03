// Copyright 2020 Andrew Quinn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofrac

import "github.com/lucasb-eyer/go-colorful"

// PrettyBands is a color palette containing discrete bands of blue,
// brown, and cream hues.
var PrettyBands = NewUniformBandedPalette(
	colorful.Hsv(24.0, 0.38, 0.33),
	colorful.Hsv(158.0, 0.48, 0.73),
	colorful.Hsv(58.0, 0.72, 0.83),
	colorful.Hsv(58.0, 0.32, 0.95),
	colorful.Hsv(24.0, 0.86, 0.97),
)

// PrettyBands2 is similar to PrettyBands, only it contains some additional
// orange tones.
var PrettyBands2 = NewUniformBandedPalette(
	colorful.Hsv(27.0, 0.75, 0.25),
	colorful.Hsv(188.0, 0.35, 0.82),
	colorful.Hsv(175.0, 0.13, 0.91),
	colorful.Hsv(35.0, 0.17, 0.85),
	colorful.Hsv(52.0, 0.06, 1.00),
)

// BWStripes is a palette of one black and one white color band.
// TODO: Rename to not imply periodicity
var BWStripes = NewUniformBandedPalette(
	colorful.Hsv(0, 0, 0),
	colorful.Hsv(0, 0, 1),
)

// PrettyBlends is an interpolated version of PrettyBands.
var PrettyBlends = BlendedBandedPalette(PrettyBands)

// PrettyBlends2 is an interpolated version of PrettyBands2.
var PrettyBlends2 = BlendedBandedPalette(PrettyBands2)

// BWBlends is an interpolated version of BWStripes.
var BWBlends = BlendedBandedPalette(BWStripes)

// PrettyPeriodic is a periodic version of PrettyBands.
var PrettyPeriodic = PeriodicPalette{
	Period:        1,
	BandedPalette: PrettyBands,
}

// PrettyPeriodic2 is a periodic version of PrettyBands2.
var PrettyPeriodic2 = PeriodicPalette{
	Period:        10,
	BandedPalette: PrettyBands2,
}

// BWPeriodic is a periodic version of BWStripes for fans of zebras and Tim Burton.
// TODO: Rename after renaming BWStripes
var BWPeriodic = PeriodicPalette{
	Period:        1,
	BandedPalette: BWStripes,
}
