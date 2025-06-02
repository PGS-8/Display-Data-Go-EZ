package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type user struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Age         int    `json:"age"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Nationality string `json:"nationality"`
}

var current_db *sql.DB

// Task 2 : Get All users in the database using GET Method
func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, err := current_db.Query("SELECT * FROM users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var users []user

		for rows.Next() {
			var u user
			if err := rows.Scan(&u.Id, &u.Name, &u.Age, &u.Email, &u.Password, &u.Nationality); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, u)
		}

		userJSON, err := json.Marshal(users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(users)

		w.Header().Set("Content-type", "application/json")
		w.Write(userJSON)
	}
}

// Task 3 : Get an user from specific ID using GET Method
func handleUser(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, "users/")
	query_id := urlPathSegment[len(urlPathSegment)-1]

	switch r.Method {
	case http.MethodGet:
		var id int
		var name string
		var age int
		var email string
		var password string
		var nationality string

		err := current_db.QueryRow("SELECT * FROM users WHERE id = ?", query_id).Scan(
			&id, &name, &age, &email, &password, &nationality)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		user := user{
			Id:          id,
			Name:        name,
			Age:         age,
			Email:       email,
			Password:    password,
			Nationality: nationality,
		}

		userJSON, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(user)

		w.Header().Set("Content-type", "application/json")
		w.Write(userJSON)
	}
}

// =========================================

// Task 1 : Connecting API with MySQL
func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Get env vars
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Fail to connect the database")
		return
	}

	fmt.Println("Connect the database successfully")
	current_db = db
}

// Task 4 : Have Middleware and CORS
func corsMiddleware(next http.Handler, path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware Start")
		fmt.Println("Location: " + path)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		next.ServeHTTP(w, r)
		fmt.Println("Middleware Finish")
	})
}

func main() {
	usersHandler := http.HandlerFunc(handleUsers)
	userHandler := http.HandlerFunc(handleUser)
	http.Handle("/users", corsMiddleware(usersHandler, "/users"))
	http.Handle("/users/", corsMiddleware(userHandler, "/users/:id"))
	http.ListenAndServe(":5500", nil)
}
