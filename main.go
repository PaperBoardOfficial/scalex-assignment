package main

import (
	"log"
	"net/http"

	"github.com/PaperBoardOfficial/scalex-assignment/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/home", handlers.Home)
	http.HandleFunc("/addBook", handlers.AddBook)
	http.HandleFunc("/deleteBook", handlers.DeleteBook)
	http.ListenAndServe(":8080", nil)
}
