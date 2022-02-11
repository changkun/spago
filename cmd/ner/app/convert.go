// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

import (
	"github.com/nlpodyssey/spago/pkg/mat"
	"github.com/nlpodyssey/spago/pkg/nlp/sequencelabeler"
	"github.com/urfave/cli/v2"
)

func newConvertCommandFor[T mat.DType](app *NERApp) *cli.Command {
	return &cli.Command{
		Name:        "convert",
		Usage:       "Run the " + programName + " to convert a pre-processed Flair model.",
		Description: "Run the " + programName + " converter.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "model-folder",
				Destination: &app.modelFolder,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "model-name",
				Destination: &app.modelName,
				Required:    true,
			},
		},
		Action: func(c *cli.Context) error {
			sequencelabeler.Convert[T](app.modelFolder, app.modelName)
			return nil
		},
	}
}
