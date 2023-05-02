package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

func readFile(path string) (string, error) {
	interval := [2]int{0, math.MaxInt}
	if strings.Contains(path, "#") {
		s := strings.Split(path, "#")

		path = s[0]

		var err error
		for i, v := range strings.Split(s[1], ":") {
			interval[i], err = strconv.Atoi(v)
			if err != nil {
				return "", errors.New("invalid line selectors in file path")
			}
		}
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	str := string(bytes)
	str = strings.Join(strings.Split(str, "\n")[interval[0]-1:interval[1]], "\n")

	return str, nil
}

func TestsCommand(ctx *cli.Context) error {
	inputPath := ctx.Args().Get(0)
	outPath := ctx.Args().Get(1)

	if inputPath == "" || outPath == "" {
		return cli.Exit("Input and output file paths must be defined", 1)
	}

	inputStr, err := readFile(inputPath)
	if err != nil {
		return err
	}

	s := CreateSpinner()
	s.Suffix = " Generating tests for code snippet 🪄\n"
	s.Start()

	result, err := AskGPT(fmt.Sprintf(`Generate tests code for the following code snippet based on what it does in the same language\n\n %s`, inputStr))
	s.Stop()
	if err != nil {
		return err
	}

	err = os.WriteFile(outPath, []byte(result), 0644)
	if err != nil {
		return err
	}

	fmt.Print(boldGreen("\n✅ Test cases generated successfully 🐙\n"))

	return nil
}
