// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linear

import (
	"github.com/nlpodyssey/spago/ag"
	"github.com/nlpodyssey/spago/losses"
	"github.com/nlpodyssey/spago/mat"
	"github.com/nlpodyssey/spago/nn"
	"github.com/nlpodyssey/spago/nn/activation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModel_Forward(t *testing.T) {
	t.Run("float32", testModelForward[float32])
	t.Run("float64", testModelForward[float64])
}

func testModelForward[T mat.DType](t *testing.T) {
	model := newTestModel[T]()
	g := ag.NewGraph[T]()

	// == Forward

	x := g.NewVariable(mat.NewVecDense([]T{-0.8, -0.9, -0.9, 1.0}), true)

	actProc := nn.Reify(activation.New[T](ag.OpTanh), g)
	proc := nn.Reify(model, g)
	y := nn.ToNode[T](actProc.Forward(proc.Forward(x)...)) // TODO: test linear only

	assert.InDeltaSlice(t, []T{-0.39693, -0.79688, 0.0, 0.70137, -0.18775}, y.Value().Data(), 1.0e-05)

	// == Backward

	gold := g.NewVariable(mat.NewVecDense([]T{0.0, 0.5, -0.4, -0.9, 0.9}), false)
	loss := losses.MSE(g, y, gold, false)
	g.Backward(loss)

	assert.InDeltaSlice(t, []T{0.0126, -2.07296, 1.07476, -0.14158}, x.Grad().Data(), 0.005)

	assert.InDeltaSlice(t, []T{
		0.26751, 0.30095, 0.30095, -0.33439,
		0.37867, 0.42601, 0.42601, -0.47334,
		-0.32, -0.36, -0.36, 0.4,
		-0.65089, -0.73226, -0.73226, 0.81362,
		0.83952, 0.94446, 0.94446, -1.04940,
	}, model.W.Grad().Data(), 1.0e-05)

	assert.InDeltaSlice(t, []T{
		-0.33439, -0.47334, 0.4, 0.81362, -1.0494,
	}, model.B.Grad().Data(), 1.0e-05)
}

func newTestModel[T mat.DType]() *Model[T] {
	model := New[T](4, 5)
	model.W.Value().SetData([]T{
		0.5, 0.6, -0.8, -0.6,
		0.7, -0.4, 0.1, -0.8,
		0.7, -0.7, 0.3, 0.5,
		0.8, -0.9, 0.0, -0.1,
		0.4, 1.0, -0.7, 0.8,
	})
	model.B.Value().SetData([]T{0.4, 0.0, -0.3, 0.8, -0.4})
	return model
}
