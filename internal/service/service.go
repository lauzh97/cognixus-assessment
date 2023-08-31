package internal

import (
	"context"
	"errors"
	"fmt"
	b "todo/internal/business"
	g "todo/internal/google"
	pb "todo/proto/todo"
)

type TodoServer struct{
	pb.UnimplementedTodoServer
}

var email string

func NewTodoServer(ctx context.Context) pb.TodoServer {
	return &TodoServer{}
}

// Adds a new item into todolist
func (s *TodoServer) AddTodo(ctx context.Context, in *pb.AddTodoRequest) (*pb.EmptyReply, error) {
	if err := s.CheckLogin(ctx); err != nil {return nil, err}
	return b.AddTodo(ctx, email, in)
}

// Soft deletes an item in the todolist
func (s *TodoServer) DeleteTodo(ctx context.Context, in *pb.UpdateTodoRequest) (*pb.EmptyReply, error) {
	if err := s.CheckLogin(ctx); err != nil {return nil, err}
	return b.DeleteTodo(ctx, email, in)
}

// Lists all items of the todolist
func (s *TodoServer) ListTodo(ctx context.Context, in *pb.EmptyRequest) (*pb.ListTodoReply, error) {
	if err := s.CheckLogin(ctx); err != nil {return nil, err}
	return b.ListTodo(ctx, email)
}

// Mark an item as true or completed
func (s *TodoServer) MarkTodo(ctx context.Context, in *pb.UpdateTodoRequest) (*pb.EmptyReply, error) {
	if err := s.CheckLogin(ctx); err != nil {return nil, err}
	return b.MarkTodo(ctx, email, in)
}

// Pong!
func (s *TodoServer) Ping(ctx context.Context, in *pb.EmptyRequest) (*pb.PingReply, error) {
	if err := s.CheckLogin(ctx); err != nil {return nil, err}
	return b.Ping(ctx, in)
}

// Checks if user is logged in using Gmail.
// If user is a new user, then create a new user automatically.
func (s *TodoServer) CheckLogin(ctx context.Context) error {
	email = g.UserDetails.Email
	// user is not logged in
	if email == "" {
		return errors.New("user not logged in. Please log in using http://localhost:8081")
	}

	// user is logged in, check if new user
	userExists, err :=  b.CheckUserExists(ctx, email)
	if err != nil {
		return err
	}

	// new user
	if !userExists {
		_, err = b.AddNewUser(ctx, email)
		if err != nil {
			return err
		}

		fmt.Println("Added new user with email: " + email)
	}

	return nil
}