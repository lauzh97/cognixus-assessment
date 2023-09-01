package internal

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

var DB *sql.DB

var AddItem = func(ctx context.Context, userId uuid.UUID, todoListId uuid.UUID, itemName string, itemDescription string) (uuid.UUID, error) {
	id := uuid.New()

	query := `INSERT INTO main.item(id, todoListId, name, description) VALUES ($1,$2,$3,$4);`
	_, err := DB.Exec(query, id, todoListId, itemName, itemDescription)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

var UpdateItem = func(ctx context.Context, itemId string, item Item) (bool, error) {
	query := `UPDATE main.item SET name=$1, description=$2, markDone=$3, active=$4, updatedOn=$5 WHERE id=$6;`
	_, err := DB.Exec(query, item.Name, item.Description, item.MarkDone, item.Active, time.Now(), item.Id)
	if err != nil {
		return false, err
	}

	return true, nil
}

var GetItemByItemName = func(ctx context.Context, todoListId uuid.UUID, itemName string) (Item, error) {
	query := `SELECT id, todoListId, name, description, markDone FROM main.item WHERE todoListId=$1 AND name=$2 AND active=true`
	row := DB.QueryRow(query, todoListId, itemName)

	var item Item
	err := row.Scan(
		&item.Id,
		&item.TodoListId,
		&item.Name,
		&item.Description,
		&item.MarkDone,
	)
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

var ListItem = func(ctx context.Context, todoListId uuid.UUID) ([]Item, error) {
	query := `SELECT name, description, markDone FROM main.item WHERE todoListId=$1 AND active=true`
	rows, err := DB.Query(query, todoListId)
	if err != nil {
		return nil, err
	}

	var items []Item
	for rows.Next() {
		var item Item
		err = rows.Scan(&item.Name, &item.Description, &item.MarkDone)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

var GetTodoListIdByUserId = func(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	query := `SELECT todoListId FROM main.user WHERE id=$1`
	row := DB.QueryRow(query, userId)

	var todoListId uuid.UUID
	err := row.Scan(&todoListId)
	if err != nil {
		return uuid.Nil, err
	}

	return todoListId, nil
}

var AddTodoList = func(ctx context.Context) (uuid.UUID, error) {
	id := uuid.New()

	query := `INSERT INTO main.todoList(id) VALUES($1);`
	_, err := DB.Exec(query, id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

var AddUser = func(ctx context.Context, email string, todoListId uuid.UUID) (uuid.UUID, error) {
	id := uuid.New()

	query := `INSERT INTO main.user(id, email, todoListId) VALUES($1,$2,$3);`
	_, err := DB.Exec(query, id, email, todoListId)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

var GetUser = func(ctx context.Context, email string) (User, error) {
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
