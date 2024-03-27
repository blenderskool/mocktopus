package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var boldRed = color.New(color.Bold, color.FgHiRed).SprintFunc()
var boldGreen = color.New(color.Bold, color.FgGreen).SprintFunc()
var MOCKTOPUS_OPENAI_KEY = os.Getenv("MOCKTOPUS_OPENAI_KEY")
var MOCKTOPUS_GEMINI_KEY = os.Getenv("MOCKTOPUS_GEMINI_KEY")

func main() {
	if MOCKTOPUS_OPENAI_KEY == "" && MOCKTOPUS_GEMINI_KEY == "" {
		fmt.Println(
			`Please add API keys for either one of the models as an env variable with its respective name:
 * OpenAI: "MOCKTOPUS_OPENAI_KEY"
 * Gemini: "MOCKTOPUS_GEMINI_KEY"`,
		)
		return
	}

	app := &cli.App{
		Name:    "mocktopus",
		Usage:   fmt.Sprintf("🐙 %s CLI tool to generate mocks for anything!", boldGreen("GPT powered")),
		Version: "1.0.0",
		Authors: []*cli.Author{
			{
				Name: "Akash Hamirwasia",
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "proto",
				Description: "generate mock data for complex structures by analyzing proto definitions",
				Usage:       "proto <source> <destination>",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "code",
						Aliases: []string{"c"},
						Usage:   "generate code for generating mock data",
					},
				},
				Action: ProtoCommand,
			},
			{
				Name:        "placeholder",
				Description: "generate mock data from natural descriptions",
				Action:      PlaceholderCommand,
			},
			{
				Name:        "tests",
				Description: "generate test cases for code snippets",
				Usage:       "tests <source> <destination>",
				Action:      TestsCommand,
			},
			{
				Name:        "persona",
				Description: "generate user personas for a product",
				Action:      PersonaCommand,
			},
		},
		ExitErrHandler: func(ctx *cli.Context, err error) {
			if err == nil {
				return
			}

			fmt.Println(boldRed("⚠️ Error occurred while running ", ctx.Command.Name, " command:"))
			fmt.Println(err)
		},
	}

	app.Run(os.Args)
}
