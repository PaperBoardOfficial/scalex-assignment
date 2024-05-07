package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PaperBoardOfficial/scalex-assignment/handlers"
)

func TestLogin(t *testing.T) {
	reqBody := bytes.NewBuffer([]byte(`{"username":"admin","password":"adminpass"}`))
	req, err := http.NewRequest("POST", "/login", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	cookie := rr.Result().Cookies()[0]
	if cookie.Name != "token" {
		t.Errorf("handler returned unexpected cookie: got %v want %v",
			cookie.Name, "token")
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	reqBody := bytes.NewBuffer([]byte(`{"username":"admin","password":"wrongpass"}`))
	req, err := http.NewRequest("POST", "/login", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestLoginMissingRequestBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
