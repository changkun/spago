// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nn

import (
	"github.com/nlpodyssey/spago/ag"
	"github.com/nlpodyssey/spago/mat"
)

// Param is the interface for a Model parameter.
type Param[T mat.DType] interface {
	ag.Node[T]

	// Name returns the params name (can be empty string).
	Name() string
	// SetName set the params name (can be empty string).
	SetName(name string)
	// Type returns the params type (weights, biases, undefined).
	Type() ParamsType
	// SetType set the params type (weights, biases, undefined).
	SetType(pType ParamsType)
	// SetRequiresGrad set whether the param requires gradient, or not.
	SetRequiresGrad(value bool)
	// ReplaceValue replaces the value of the parameter and clears the support structure.
	ReplaceValue(value mat.Matrix[T])
	// ApplyDelta updates the value applying the delta.
	ApplyDelta(delta mat.Matrix[T])
	// Payload returns the optimizer support structure (can be nil).
	Payload() *Payload[T]
	// SetPayload is a thread safe operation to set the given Payload on the
	// receiver Param.
	SetPayload(payload *Payload[T])
	// ClearPayload clears the support structure.
	ClearPayload()
}
