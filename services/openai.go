package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// https://platform.openai.com/docs/api-reference/images/create
type ImagePayload struct {
	Prompt         string  `json:"prompt"`
	ResponseFormat *string `json:"response_format,omitempty"`
	Model          *string `json:"model,omitempty"`
	N              *int    `json:"n,omitempty"`
	Quality        *string `json:"quality,omitempty"`
	Size           *string `json:"size,omitempty"`
	Style          *string `json:"style,omitempty"`
}

type ImageData struct {
	URL string `json:"url"`
}

type Image struct {
	Data []ImageData `json:"data"`
}

const API_URL = "https://api.openai.com/v1"

func GenerateQuestionImage(question string, result *Image) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	if apiKey == "" {
		fmt.Println("Please set OPENAI_API_KEY")
		os.Exit(1)
	}

	model := "dall-e-3"
	responseFormat := "url"
	n := 1
	quality := "hd"
	size := "1024x1024"
	style := "vivid"

	data := ImagePayload{
		Prompt:         question,
		Model:          &model,
		ResponseFormat: &responseFormat,
		N:              &n,
		Quality:        &quality,
		Size:           &size,
		Style:          &style,
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", API_URL+"/images/generations", body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&result)
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// https://platform.openai.com/docs/api-reference/chat/create
type ChatPayload struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

type ChatResponseChoiceMessage struct {
	Content string `json:"content"`
}

type ChatResponseChoice struct {
	Message ChatResponseChoiceMessage `json:"message"`
}

type ChatResponse struct {
	Choices []ChatResponseChoice `json:"choices"`
}

func GenerateQuestion(difficulty string, result *ChatResponse) {
	validDifficulties := map[string]bool{
		"super-easy": true,
		"easy":       true,
		"medium":     true,
		"hard":       true,
		"super-hard": true,
	}

	if !validDifficulties[difficulty] {
		fmt.Println("Invalid difficulty level. It must be 'easy', 'medium', or 'hard'.")
		return
	}

	apiKey := os.Getenv("OPENAI_API_KEY")

	if apiKey == "" {
		fmt.Println("Please set OPENAI_API_KEY")
		os.Exit(1)
	}

	payload := ChatPayload{
		Messages: []Message{
			{
				Role: "system",
				Content: `You're a game host on a quiz show. You're asking the contestants a questions about geography.
					There're 5 difficulty levels: super-easy, easy, medium, hard, and super-hard.
					Depending on the difficulty level, the question will be different.
					Each time when you're will be asking to generate a question, you'll be given a difficulty level.
					You need to generate a question about geography with the given difficulty level.
					For example, if the difficulty level is "easy", you can generate a question like "What is the name of the tallest mountain in Africa?".
				`,
			},
			{
				Role: "system",
				Content: `There're a list of questions you're already asked.
					[]
					You can't ask the same question twice.
					You can't ask a question about something that doesn't exist or made up.
				`,
			},
			{
				Role:    "system",
				Content: "Create a question about " + difficulty + " geography.",
			},
		},
		Model: "gpt-4",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", API_URL+"/chat/completions", body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&result)
}

func ValidateUserAnswer(answer string, question string) bool {
	apiKey := os.Getenv("OPENAI_API_KEY")

	if apiKey == "" {
		fmt.Println("Please set OPENAI_API_KEY")
		os.Exit(1)
	}

	payload := ChatPayload{
		Messages: []Message{
			{
				Role: "system",
				Content: fmt.Sprintf(
					`is the "%s" the valid answer for the question: "%s", answer only "yes" or "no" in lowercase without punctuation`,
					answer,
					question,
				),
			},
		},
		Model: "gpt-4",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", API_URL+"/chat/completions", body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	var result ChatResponse
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Choices[0].Message.Content == "yes"
}
