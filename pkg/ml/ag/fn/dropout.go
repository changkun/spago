// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"github.com/nlpodyssey/spago/pkg/mat"
	"github.com/nlpodyssey/spago/pkg/mat/rand"
	"github.com/nlpodyssey/spago/pkg/mat/rand/bernulli"
)

var _ Function = &Dropout{}

// Dropout is an operator to perform elements dropout with a probability.
type Dropout struct {
	x       Operand
	prob    mat.Float
	q       mat.Float // 1 - p
	randGen *rand.LockedRand[mat.Float]
	mask    mat.Matrix[mat.Float] // filled during the forward
}

// NewDropout returns a new Dropout Function.
func NewDropout(x Operand, p mat.Float, randGen *rand.LockedRand[mat.Float]) *Dropout {
	return &Dropout{
		x:       x,
		prob:    p,
		q:       1.0 - p,
		randGen: randGen,
		mask:    nil,
	}
}

// Forward computes the output of the function.
func (r *Dropout) Forward() mat.Matrix[mat.Float] {
	if r.q > 0.0 {
		r.mask = bernulli.Distribution(r.x.Value().Rows(), r.x.Value().Columns(), r.prob, r.randGen)
		r.mask.ProdScalarInPlace(1.0 / r.q)
	} else {
		r.mask = r.x.Value().ZerosLike()
	}
	return r.x.Value().Prod(r.mask)
}

// Backward computes the backward pass.
func (r *Dropout) Backward(gy mat.Matrix[mat.Float]) {
	if !(r.x.Value().SameDims(gy) || r.x.Value().VectorOfSameSize(gy)) {
		panic("fn: matrices with not compatible size")
	}
	defer mat.ReleaseMatrix(r.mask)
	if r.x.RequiresGrad() {
		gx := gy.Prod(r.mask)
		r.x.PropagateGrad(gx)
	}
}
