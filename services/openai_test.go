package services_test

import (
	"patchwork-body/geo-quizze/services"
	"testing"

	"github.com/joho/godotenv"
)

func TestGenerateQuestion(t *testing.T) {
	godotenv.Load("../.env")

	// Act
	var question services.ChatResponse
	services.GenerateQuestion("super-hard", &question)

	// Assert
	if len(question.Choices) != 1 {
		t.Errorf("expected lens %d, got %d", 1, len(question.Choices))
	}

	if len(question.Choices[0].Message.Content) == 0 {
		t.Errorf("expected non empty content, got %s", question.Choices[0].Message.Content)
	}
}

func TestGenerateQuestionImage(t *testing.T) {
	godotenv.Load("../.env")

	// Act
	var question services.ChatResponse
	services.GenerateQuestion("super-hard", &question)

	var image services.Image
	services.GenerateQuestionImage(question.Choices[0].Message.Content, &image)

	// Assert
	if len(image.Data) == 0 {
		t.Errorf("expected non empty data, got %s", image.Data)
	}

	if len(image.Data[0].URL) == 0 {
		t.Errorf("expected non empty url, got %s", image.Data[0].URL)
	}
}
