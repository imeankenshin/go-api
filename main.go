package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func Migration() error {
	db, err := gorm.Open(sqlite.Open("todos.db"))
	if err != nil {
		panic(err.Error())
	}
	return db.AutoMigrate(&Todo{})
}

func main() {
	r := chi.NewRouter()
	r.Use(
		// CORS
		cors.Handler(cors.Options{
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		}),
		// Logger
		middleware.Logger,
		middleware.Recoverer,
	)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		var todos []Todo

		db, err := gorm.Open(sqlite.Open("todos.db"))
		if err != nil {
			panic(err.Error())
		}
		db.Find(&todos)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		_ = json.NewEncoder(w).Encode(todos)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		db, err := gorm.Open(sqlite.Open("todos.db"))
		if err != nil {
			panic(err.Error())
		}
		result := db.Create(&Todo{
			Title:       body["title"],
			Description: body["description"],
			Done:        false,
		})
		if result.Error != nil {
			panic(result.Error.Error())
		}
		fmt.Fprintf(w, "Success!")
	})

	println("Server is working on http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", r))
}
