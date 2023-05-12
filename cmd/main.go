package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"protocoll/internal"
)

func run() int {
	// application CLI design
	app := &cli.App{
		Name:  "protocoll",
		Usage: "Converts proto files into a Postman/Insomnia collection",
		Commands: []*cli.Command{
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "Generate a collection",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "folder",
						Usage:    "Folder containing proto files",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "name",
						Usage:    "Collection name",
						Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					generator := internal.Generator{}
					return generator.Generate(cCtx.String("folder"), cCtx.String("name"))
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
