// Copyright 2020 Andrew Quinn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofrac

import (
	"errors"
)

// DomainReader reads values from a discretization of a bounded 2D space.
type DomainReader interface {
	// At returns the underlying coordinates of a sample (i, j) in a domain as
	//a complex number, where the sample (0, 0) is the top-left corner of the
	// domain. If a non-existent sample is requested, an error is returned.
	At(i int, j int) (loc complex128, err error)

	// Dimensions returns the number of samples taken along each axis of a
	// domain as rows and columns.
	Dimensions() (rows int, cols int)
}

// Domain stores bounds and sampling information over a rectangular 2D surface.
//
// The lower-left corner of the domain is stored in (x0, y0) and reaches to the
// upper-right corner at (x0+xDist, y0+yDist). Along the x and y axes, xs and
// ys samples are taken, respectively. The inverses of xDist and yDist are
// stored in wInv and hInv, respectively, to speed up computation.
type Domain struct {
	x0, y0       float64
	xDist, yDist float64
	xs, ys       int

	wInv float64
	hInv float64
}

func (r *Domain) At(i int, j int) (loc complex128, err error) {
	if i < 0 || i >= r.xs || j < 0 || j >= r.ys {
		return 0, errors.New("gofrac: sample is out of bounds")
	}

	ti := float64(i) * r.wInv
	re := ti*r.xDist + r.x0

	tj := 1.0 - float64(j)*r.hInv
	im := tj*r.yDist + r.y0

	return complex(re, im), nil
}

func (r *Domain) Dimensions() (rows int, cols int) {
	return r.ys, r.xs
}

// NewDomain constructs a rectangular 2D domain.
//
// (x0, y0) is the bottom-left corner of the domain.
// (y0, y1) is the top-right corner of the domain.
// The domain is sampled xSamples and ySamples times along the x and y axes,
// respectively.
func NewDomain(x0, y0, x1, y1 float64, xSamples, ySamples int) (d *Domain, err error) {
	if xSamples <= 0 || ySamples <= 0 {
		return nil, errors.New("gofrac: The number of samples along any axis must be greater than zero")
	}

	return &Domain{
		x0:    x0,
		y0:    y0,
		xs:    xSamples,
		ys:    ySamples,
		xDist: x1 - x0,
		yDist: y1 - y0,
		wInv:  1.0 / float64(xSamples),
		hInv:  1.0 / float64(ySamples),
	}, nil
}
