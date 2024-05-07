package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/golang-jwt/jwt"
)

func Home(w http.ResponseWriter, r *http.Request) {
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

	if !tkn.Valid {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	_, currentFilePath, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFilePath)
	parentDir := filepath.Dir(currentDir)

	var files []string
	files = append(files, filepath.Join(parentDir, "data/regularUser.csv"))
	if claims.Issuer == "admin" {
		files = append(files, filepath.Join(parentDir, "data/adminUser.csv"))
	}

	var books []string
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Println("Failed to open file:", err)
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		csvReader := csv.NewReader(f)
		records, err := csvReader.ReadAll()
		if err != nil {
			fmt.Printf("Failed to open file %s: %v\n", file, err)
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		for _, record := range records {
			books = append(books, record...)
		}
	}

	json.NewEncoder(w).Encode(books)
}
