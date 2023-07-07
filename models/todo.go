package models

type Todo struct {
	Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}
