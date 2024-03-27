package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
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

func askGPT(message string) (string, error) {
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

func askGemini(message string) (string, error) {
	b, err := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{
				"role": "user",
				"parts": []map[string]any{
					{
						"text": message + "\nDON'T ADD BACKTICKS IN THE RESPONSE.",
					},
				},
			},
		},
		"generationConfig": map[string]any{
			"temperature":     0.9,
			"topK":            1,
			"topP":            1,
			"maxOutputTokens": 2048,
			"stopSequences":   []map[string]any{},
		},
		"safetySettings": []map[string]any{
			{
				"category":  "HARM_CATEGORY_HARASSMENT",
				"threshold": "BLOCK_MEDIUM_AND_ABOVE",
			},
			{
				"category":  "HARM_CATEGORY_HATE_SPEECH",
				"threshold": "BLOCK_MEDIUM_AND_ABOVE",
			},
			{
				"category":  "HARM_CATEGORY_SEXUALLY_EXPLICIT",
				"threshold": "BLOCK_MEDIUM_AND_ABOVE",
			},
			{
				"category":  "HARM_CATEGORY_DANGEROUS_CONTENT",
				"threshold": "BLOCK_MEDIUM_AND_ABOVE",
			},
		},
	})
	if err != nil {
		return "", err
	}

	body := bytes.NewBuffer(b)
	req, err := http.NewRequest("POST", "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.0-pro:generateContent?key="+MOCKTOPUS_GEMINI_KEY, body)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", errors.New("an unexpected error occurred while requesting Gemini endpoints")
	}

	result := struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string
				}
				Role string
			}
		}
	}{}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}

	answer := result.Candidates[0].Content.Parts[0].Text

	// For some reason Gemini is very adamant to markdown formatting
	if strings.HasPrefix(answer, "```") {
		lines := strings.Split(answer, "\n")
		answer = strings.Join(lines[1:len(lines)-1], "\n")
	}
	return answer, nil
}

func AskAI(message string) (string, error) {
	if MOCKTOPUS_OPENAI_KEY != "" {
		return askGPT(message)
	} else if MOCKTOPUS_GEMINI_KEY != "" {
		return askGemini(message)
	}
	panic("Unreachable: Please set relevant API keys as env variables for the AI model to use")
}
