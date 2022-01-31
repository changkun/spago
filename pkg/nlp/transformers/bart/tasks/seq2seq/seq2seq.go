// Copyright 2021 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq2seq

import (
	"github.com/nlpodyssey/spago/pkg/ml/ag"
	"github.com/nlpodyssey/spago/pkg/ml/nn"
	"github.com/nlpodyssey/spago/pkg/nlp/tokenizers/sentencepiece"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bart/config"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bart/head/conditionalgeneration"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bart/loader"
)

var _ nn.Model = &BartForConditionalGeneration{}

// BartForConditionalGeneration contains the Model and the Tokenizer
// used for conditional generation tasks.
// For example, Machine Translation and Summarization.
type BartForConditionalGeneration struct {
	*conditionalgeneration.Model
	Tokenizer *sentencepiece.Tokenizer
}

// LoadModel loads a BartForConditionalGeneration from file.
func LoadModel(modelPath string) (*BartForConditionalGeneration, error) {
	model, err := loader.Load(modelPath)
	if err != nil {
		return nil, err
	}
	tokenizer, err := sentencepiece.NewFromModelFolder(modelPath, false) // TODO: lowercase from config
	if err != nil {
		return nil, err
	}
	return &BartForConditionalGeneration{
		Model:     model.(*conditionalgeneration.Model),
		Tokenizer: tokenizer,
	}, nil
}

// Generate generates new texts starting from the input.
func (t *BartForConditionalGeneration) Generate(text string) (string, error) {
	g := ag.NewGraph(ag.IncrementalForward(false))
	defer g.Clear()

	proc := nn.ReifyForInference(t.Model, g)
	bartConfig := proc.BART.Config

	tokens := t.Tokenizer.Tokenize(text)
	tokenIDs := t.Tokenizer.TokensToIDs(tokens)

	tokenIDs = append(tokenIDs, bartConfig.EosTokenID)

	rawGeneratedIDs := proc.Generate(tokenIDs)
	generatedIDs := t.stripBadTokens(rawGeneratedIDs, bartConfig)

	generatedTokens := t.Tokenizer.IDsToTokens(generatedIDs)
	generatedText := t.Tokenizer.Detokenize(generatedTokens)

	return generatedText, nil
}

func (t *BartForConditionalGeneration) stripBadTokens(ids []int, bartConfig config.Config) []int {
	result := make([]int, 0, len(ids))
	for _, id := range ids {
		if id == bartConfig.EosTokenID || id == bartConfig.PadTokenID || id == bartConfig.BosTokenID ||
			id == bartConfig.DecoderStartTokenID {
			continue
		}
		result = append(result, id)
	}
	return result
}
