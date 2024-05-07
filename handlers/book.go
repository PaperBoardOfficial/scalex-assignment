package handlers

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type Book struct {
	Name   string `json:"name"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	tknStr := c.Value
	claims := &jwt.StandardClaims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		http.Error(w, "Error parsing token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !tkn.Valid || claims.Issuer != "admin" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var book Book
	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if book.Name == "" || book.Author == "" || book.Year < 0 || book.Year > time.Now().Year() {
		http.Error(w, "Invalid book data", http.StatusBadRequest)
		return
	}

	file, err := os.OpenFile("../data/regularUser.csv", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	defer writer.Flush()

	writer.Write([]string{book.Name, book.Author, strconv.Itoa(book.Year)})
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	tknStr := c.Value
	claims := &jwt.StandardClaims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		http.Error(w, "Error parsing token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !tkn.Valid || claims.Issuer != "admin" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	bookName := r.URL.Query().Get("name")
	if bookName == "" {
		http.Error(w, "Invalid book name", http.StatusBadRequest)
		return
	}

	filePath := "../data/regularUser.csv"
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	var newRecords [][]string
	bookFound := false
	for _, record := range records {
		if strings.EqualFold(record[0], bookName) {
			bookFound = true
			continue
		}
		newRecords = append(newRecords, record)
	}

	if !bookFound {
		http.Error(w, "Book not found", http.StatusBadRequest)
		return
	}

	file, err = os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(newRecords)
	if err != nil {
		http.Error(w, "Failed to write to file", http.StatusInternalServerError)
		return
	}
}
