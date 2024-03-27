package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

func PersonaCommand(ctx *cli.Context) error {
	answers := struct {
		Product string
	}{}

	survey.Ask([]*survey.Question{
		{
			Name: "product",
			Prompt: &survey.Input{
				Message: "Describe your product:",
			},
			Validate: survey.Required,
		},
	}, &answers)

	s := CreateSpinner()
	s.Suffix = " Generating user personas for the product 🪄\n"
	s.Start()

	result, err := AskAI(fmt.Sprintf(`Create a few user personas with name alliterations and different backgrounds for %s. Also add behavior, needs and wants, demographics to each persona`, answers.Product))
	s.Stop()
	if err != nil {
		return err
	}

	fmt.Println(boldGreen("\n✅ User personas generated successfully 🐙\n"))
	fmt.Println(result)

	return nil
}
