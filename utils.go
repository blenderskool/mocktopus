package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
)

func GetKeys[T comparable, U any](s map[T]U) []T {
	keys := make([]T, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}

	return keys
}

func SurveyNumberValidator(ans interface{}) error {
	if _, err := strconv.Atoi(ans.(string)); err != nil {
		return errors.New("please enter a number")
	}
	return nil
}

func SurveyNumberTransform(ans interface{}) (newAns interface{}) {
	value, _ := strconv.Atoi(ans.(string))
	return value
}

func CreateSpinner() *spinner.Spinner {
	return spinner.New(spinner.CharSets[14], 25*time.Millisecond)
}

func AskGPT(message string) (string, error) {
	b, err := json.Marshal(map[string]any{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]any{
			{
				"role":    "user",
				"content": message,
			},
		}})
	if err != nil {
		return "", err
	}

	body := bytes.NewBuffer(b)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", body)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+MOCKTOPUS_OPENAI_KEY)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", errors.New("an unexpected error occurred while requesting OpenAI endpoints")
	}

	result := struct {
		Choices []struct {
			Message struct {
				Content string
			}
		}
	}{}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Choices[0].Message.Content, nil
}
