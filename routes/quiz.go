package routes

import (
	"html/template"
	"net/http"
	"patchwork-body/geo-quizze/services"
)

type QuizUserAnswer struct {
	Answer    string
	IsCorrect bool
}

type QuizTemplateParams struct {
	Id          string
	Question    string
	ImageUrl    string
	UserAnswers *[]QuizUserAnswer
}

type QuizAnswersHistoryTemplateParams struct {
	UserAnswers *[]QuizUserAnswer
}

func Quiz(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		var question services.ChatResponse
		services.GenerateQuestion("easy", &question)

		var image services.Image
		services.GenerateQuestionImage(question.Choices[0].Message.Content, &image)

		tmplParams := QuizTemplateParams{
			Id:          "1",
			Question:    question.Choices[0].Message.Content,
			ImageUrl:    image.Data[0].URL,
			UserAnswers: &[]QuizUserAnswer{},
		}

		tmpl := template.Must(template.ParseFiles("views/index.html", "views/auth.html", "views/quiz.html"))

		err := tmpl.ExecuteTemplate(w, "index.html", tmplParams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		http.Redirect(w, r, "/quiz", http.StatusSeeOther)

	case http.MethodPatch:
		userAnswer := r.FormValue("answer")
		question := r.FormValue("question")

		tmpl := template.Must(template.ParseFiles("views/answers-history.html"))

		isCorrect := services.ValidateUserAnswer(userAnswer, question)

		var userAnswers []QuizUserAnswer

		userAnswers = append(userAnswers, QuizUserAnswer{
			Answer:    userAnswer,
			IsCorrect: isCorrect,
		})

		tmplParams := QuizAnswersHistoryTemplateParams{
			UserAnswers: &userAnswers,
		}

		err := tmpl.Execute(w, tmplParams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
