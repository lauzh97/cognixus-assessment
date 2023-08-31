package internal

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

var DB *sql.DB

func AddTodoList(ctx context.Context) (uuid.UUID, error) {
	id := uuid.New()

	query := `INSERT INTO main.todoList(id, createdOn) values($1,$2);`
	_, err := DB.Exec(query, id, time.Now())
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func AddUser(ctx context.Context, email string, todoListId uuid.UUID) (uuid.UUID, error) {
	id := uuid.New()

	query := `INSERT INTO main.user(id, email, todoListId, createdOn) values($1,$2,$3,$4);`
	_, err := DB.Exec(query, id, email, todoListId, time.Now())
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func GetUser(ctx context.Context, email string) (User, error) {
	query := `SELECT * FROM main.user WHERE email = $1;`
	row := DB.QueryRow(query, email)

	var user User
	err := row.Scan(
		&user.Id, 
		&user.Email, 
		&user.TodoListId, 
		&user.Active, 
		&user.CreatedOn, 
		&user.UpdatedOn,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}