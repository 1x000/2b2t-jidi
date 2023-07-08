package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DATABASE       = "passwd.db"
	ADMIN_EMAIL    = "huaweilaomke402@gmail.com"
	ADMIN_PASSWORD = "awa114514"
)

func checkAdmin(email, password string) bool {
	return email == ADMIN_EMAIL && password == ADMIN_PASSWORD
}

func checkLogin(username, password string) bool {
	db, err := sql.Open("sqlite3", DATABASE)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE username=? AND password=?", username, password).Scan(&userID)
	return err == nil
}

func getUserID(username string) int {
	db, err := sql.Open("sqlite3", DATABASE)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE username=?", username).Scan(&userID)
	if err != nil {
		panic(err)
	}
	return userID
}

func saveImage(image *multipart.FileHeader, baseName string) string {
	baseDir := fmt.Sprintf("%s/images/%s", os.Getenv("GOPATH"), baseName)
	os.MkdirAll(baseDir, os.ModePerm)
	imagePath := fmt.Sprintf("%s/%s", baseDir, filepath.Base(image.Filename))
	f, err := os.Create(imagePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	src, err := image.Open()
	if err != nil {
		panic(err)
	}
	defer src.Close()

	_, err = io.Copy(f, src)
	if err != nil {
		panic(err)
	}
	return imagePath
}

func createDatabase() {
	db, err := sql.Open("sqlite3", DATABASE)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, password TEXT)")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS bases (id INTEGER PRIMARY KEY, name TEXT, x INTEGER, y INTEGER, z INTEGER, image_path TEXT, owner TEXT)")
	if err != nil {
		panic(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		if checkAdmin(username, password) {
			http.Redirect(w, r, "/admin", http.StatusFound)
			return
		} else if checkLogin(username, password) {
			userID := getUserID(username)
			http.Redirect(w, r, fmt.Sprintf("/user/%d", userID), http.StatusFound)
			return
		}
		t.Execute(w, "Invalid username or password")
	} else {
		t.Execute(w, nil)
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Complete admin handler
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Complete user handler
}

func main() {
	createDatabase()
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/user/", userHandler)
	http.ListenAndServe(":4442", nil)
}
