package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

func PlaceholderCommand(ctx *cli.Context) error {
	answers := struct {
		Placeholder string
		Count       int
	}{}

	survey.Ask([]*survey.Question{
		{
			Name: "placeholder",
			Prompt: &survey.Input{
				Message: "What do you want placeholder for?",
			},
			Validate: survey.Required,
		},
		{
			Name: "count",
			Prompt: &survey.Input{
				Message: "Number of records to generate?",
			},
			Validate:  survey.ComposeValidators(survey.Required, SurveyNumberValidator),
			Transform: SurveyNumberTransform,
		},
	}, &answers)

	s := CreateSpinner()
	s.Suffix = " Generating mock placeholder data 🪄\n"
	s.Start()

	result, err := AskAI(fmt.Sprintf(`Generate %d placeholder data for %s`, answers.Count, answers.Placeholder))
	s.Stop()
	if err != nil {
		return err
	}

	fmt.Println(boldGreen("\n✅ Mock data generated successfully 🐙\n"))
	fmt.Println(result)

	return nil
}
