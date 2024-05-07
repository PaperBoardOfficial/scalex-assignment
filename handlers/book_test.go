package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PaperBoardOfficial/scalex-assignment/handlers"
	"github.com/stretchr/testify/assert"
)

func TestAddBook(t *testing.T) {
	t.Run("add valid book", func(t *testing.T) {
		book := handlers.Book{Name: "Test Book", Author: "Test Author", Year: 2022}
		body, _ := json.Marshal(book)
		token, err := handlers.GenerateToken("admin", "adminpass")
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest("POST", "/addBook", bytes.NewBuffer(body))
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: token,
		})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.AddBook)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("add invalid book", func(t *testing.T) {
		book := handlers.Book{Name: "", Author: "Test Author", Year: 2022}
		body, _ := json.Marshal(book)
		token, err := handlers.GenerateToken("admin", "adminpass")
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest("POST", "/addBook", bytes.NewBuffer(body))
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: token,
		})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.AddBook)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestDeleteBook(t *testing.T) {

	t.Run("delete existing book", func(t *testing.T) {
		token, err := handlers.GenerateToken("admin", "adminpass")
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest("DELETE", "/deleteBook?name=Test Book", nil)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: token,
		})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.DeleteBook)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("delete non-existing book", func(t *testing.T) {
		token, err := handlers.GenerateToken("admin", "adminpass")
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest("DELETE", "/deleteBook?name=Non-existing Book", nil)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: token,
		})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.DeleteBook)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
