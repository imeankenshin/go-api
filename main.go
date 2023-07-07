package main

import (
	"awesomeProject/models"
	"awesomeProject/util"
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

func Migration() error {
	db, err := gorm.Open(sqlite.Open("todos.db"))
	if err != nil {
		panic(err.Error())
	}
	return db.AutoMigrate(&models.Todo{})
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
		var todos []models.Todo

		db, err := gorm.Open(sqlite.Open("todos.db"))
		if err != nil {
			panic(err.Error())
		}
		db.Find(&todos)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		_ = json.NewEncoder(w).Encode(todos)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := util.ReadBody(r.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		db, err := gorm.Open(sqlite.Open("todos.db"))
		if err != nil {
			panic(err.Error())
		}

		result := db.Create(&models.Todo{
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
