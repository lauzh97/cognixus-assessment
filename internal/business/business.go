package internal

import (
	"context"
	"database/sql"
	"errors"
	data "todo/internal/data"
	pb "todo/proto/todo"

	"github.com/google/uuid"
)

func AddTodo(ctx context.Context, in *pb.AddTodoRequest) (*pb.BasicReply, error) {
	return nil, nil
}

func DeleteTodo(ctx context.Context, in *pb.UpdateTodoRequest) (*pb.BasicReply, error) {
	return nil, nil
}

func ListTodo(ctx context.Context, in *pb.ListTodoRequest) (*pb.ListTodoReply, error) {
	return nil, nil
}

func MarkTodo(ctx context.Context, in *pb.UpdateTodoRequest) (*pb.BasicReply, error) {
	return nil, nil
}

func Ping(ctx context.Context, in *pb.EmptyRequest) (*pb.PingReply, error) {
	return &pb.PingReply{Pong: "Pong"}, nil
}

// Adds a new record into main.user table
func AddNewUser(ctx context.Context, email string) (uuid.UUID, error) {
	// add a new todolist for the new user
	todoListId, err := data.AddTodoList(ctx)
	if err != nil {
		return uuid.Nil, errors.New("AddTodoList failed: " + err.Error())
	}

	// add a new user
	userId, err := data.AddUser(ctx, email, todoListId)
	if err != nil {
		return uuid.Nil, errors.New("AddUser failed: " + err.Error())
	}

	return userId, nil
}

// Checks if user exists in the main.user table.
func CheckUserExists(ctx context.Context, email string) (bool, error) {
	_, err := data.GetUser(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		return false, errors.New("GetUser failed: " + err.Error())
	}
	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}