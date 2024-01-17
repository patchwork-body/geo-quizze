package main

import (
	"log"
	"net/http"
	"patchwork-body/geo-quizze/routes"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	http.HandleFunc("/", routes.Index)
	http.HandleFunc("/quiz", routes.Quiz)
	http.HandleFunc("/modal", routes.Modal)

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
