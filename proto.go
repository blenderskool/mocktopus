package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/emicklei/proto"
	"github.com/urfave/cli/v2"
)

type ProtoDefs = map[string]*proto.Message

func getAllDefinitions(path string) (ProtoDefs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Parse the file as proto
	parser := proto.NewParser(file)
	definition, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	defs := ProtoDefs{}
	proto.Walk(definition, proto.WithMessage(func(m *proto.Message) {
		defs[m.Name] = m
	}))

	return defs, nil
}

func getDefStr(defs ProtoDefs, name string) string {
	dependencies := []string{}
	fieldsStrParts := []string{}
	def := defs[name]

	// Walk through each field in the definition
	walk(def, proto.WithNormalField(func(field *proto.NormalField) {
		// Check if the type of this field is defined in this file
		if dependentDef, ok := defs[field.Type]; ok {
			// If the type is user-defined, then recursively construct the string of this defined type
			dependencies = append(dependencies, getDefStr(defs, dependentDef.Name))
		}

		baseStr := fmt.Sprintf("%s %s=%d;", field.Type, field.Name, field.Sequence)
		if field.Repeated {
			baseStr += "repeated " + baseStr
		}

		fieldsStrParts = append(fieldsStrParts, baseStr)
	}))

	fieldsStr := strings.Join(fieldsStrParts, "\n")
	defStr := fmt.Sprintf("message %s {\n%s\n}", def.Name, fieldsStr)

	dependencies = append(dependencies, defStr)
	return strings.Join(dependencies, "\n\n")
}

func walk(m *proto.Message, handlers ...proto.Handler) {
	for _, e := range m.Elements {
		for _, handler := range handlers {
			handler(e)
		}
	}
}

func prepareQuestions(defs ProtoDefs, isCode bool) []*survey.Question {
	qs := []*survey.Question{
		{
			Name: "definition",
			Prompt: &survey.Select{
				Message: "Which definition do you want mock data for?",
				Options: GetKeys(defs),
			},
			Validate: survey.Required,
		},
	}

	if !isCode {
		qs = append(qs, &survey.Question{
			Name: "count",
			Prompt: &survey.Input{
				Message: "Number of records to generate?",
				Default: "1",
			},
			Validate:  SurveyNumberValidator,
			Transform: SurveyNumberTransform,
		})
	}

	return qs
}

func ProtoCommand(ctx *cli.Context) error {
	inputPath := ctx.Args().Get(0)
	outPath := ctx.Args().Get(1)
	isCode := ctx.Bool("code")

	if inputPath == "" || outPath == "" {
		return cli.Exit("Input and output file paths must be defined", 1)
	}

	if !strings.HasSuffix(inputPath, ".proto") {
		return cli.Exit("Input file must be a .proto file", 1)
	}

	s := CreateSpinner()
	s.Suffix = " Scanning for definitions\n"
	s.Start()
	time.Sleep(time.Second)
	s.Stop()

	defs, err := getAllDefinitions(inputPath)
	if err != nil {
		return err
	}

	if len(defs) == 0 {
		return cli.Exit("No definitions found, exiting...", 0)
	}

	fmt.Printf(boldGreen("%d definitions found\n"), len(defs))

	answers := struct {
		Definition string
		Count      int
	}{}

	err = survey.Ask(prepareQuestions(defs, isCode), &answers)

	if err != nil {
		return err
	}

	defStr := getDefStr(defs, answers.Definition)

	s.Suffix = " Generating mock data for proto definition 🪄\n"
	if isCode {
		s.Suffix = " Generating code for generating mock data for proto definition 🪄\n"
	}
	s.Start()

	var result string
	if isCode {
		result, err = AskAI(
			fmt.Sprintf(`Generate JS code with "@faker-js/faker" library to create mock data for the "%s" proto definition. proto definition in object format. Use only UUID for id fields and working image urls if needed\n\n%s`, answers.Definition, defStr),
		)
	} else {
		result, err = AskAI(
			fmt.Sprintf(`Generate valid JSON array with %d unique items and each item satisfying the "%s" proto definition. Use only UUID for id fields and working image urls if needed\n\n%s`, answers.Count, answers.Definition, defStr),
		)
	}

	s.Stop()
	if err != nil {
		return err
	}

	err = os.WriteFile(outPath, []byte(result), 0644)
	if err != nil {
		return err
	}

	if isCode {
		fmt.Print(boldGreen("\n✅ Code for mock data generated successfully 🐙\n"))
	} else {
		fmt.Print(boldGreen("\n✅ Mock data generated successfully 🐙\n"))
	}
	return nil
}
