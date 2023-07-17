package main

import (
	"awesomeProject/models"
	"awesomeProject/pkg"
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
	mux := chi.NewRouter()
	mux.Use(
		// CORS
		cors.Handler(cors.Options{
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		}),
		// Logger
		middleware.Logger,
		middleware.Recoverer,
	)

	mux.Route("/task", func(mux chi.Router) {
		// List
		mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
			var todos []models.Todo

			db, err := gorm.Open(sqlite.Open("todos.db"))
			if err != nil {
				panic(err.Error())
			}
			db.Find(&todos)

			pkg.Encode(w, todos)
		})
		// New
		mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
			body, err := pkg.ReadBody(r.Body)
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
		// Get
		mux.Get("/{taskID}", func(w http.ResponseWriter, r *http.Request) {
			taskID := chi.URLParam(r, "taskID")
			var todo models.Todo

			db, err := gorm.Open(sqlite.Open("todos.db"))
			if err != nil {
				panic(err.Error())
			}
			res := db.First(&todo, taskID)
			if res.Error != nil {
				http.NotFound(w, r)
				return
			}

			pkg.Encode(w, todo)
		})
		// Delete
		mux.Delete("/{taskID}", func(w http.ResponseWriter, r *http.Request) {
			taskID := chi.URLParam(r, "taskID")
			var todo models.Todo
			db, err := gorm.Open(sqlite.Open("todos.db"))
			if err != nil {
				panic(err.Error())
			}
			res := db.Delete(&todo, taskID)
			if res.Error != nil {
				http.NotFound(w, r)
				return
			}

			pkg.Encode(w, todo)
		})
	})

	fmt.Printf("Server is working on http://localhost:3100\n")
	log.Fatal(http.ListenAndServe(":3100", mux))
}
