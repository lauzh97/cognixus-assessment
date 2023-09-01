package internal

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id         uuid.UUID
	Email      string
	TodoListId uuid.UUID
	Active     bool
	CreatedOn  time.Time
	UpdatedOn  time.Time
}

type TodoList struct {
	Id        uuid.UUID
	Active    bool
	CreatedOn time.Time
	UpdatedOn time.Time
}

type Item struct {
	Id          uuid.UUID
	TodoListId  uuid.UUID
	Name        string `json:"itemName"`
	Description string `json:"itemDescription"`
	MarkDone    bool   `json:"done"`
	Active      bool
	CreatedOn   time.Time
	UpdatedOn   time.Time
}
