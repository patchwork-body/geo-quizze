package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"patchwork-body/geo-quizze/services"
)

type QuestionResponse struct {
	Question string `json:"question"`
	ImageUrl string `json:"imageUrl"`
}

func Questions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var question services.ChatResponse
	services.GenerateQuestion("super-hard", &question)

	var image services.Image
	services.GenerateQuestionImage(question.Choices[0].Message.Content, &image)

	fmt.Println(image)

	body, err := json.Marshal(QuestionResponse{
		Question: question.Choices[0].Message.Content,
		ImageUrl: image.Data[0].URL,
	})

	if err != nil {
		panic(err)
	}

	w.Write(body)
}
