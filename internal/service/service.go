package internal

import (
	"context"
	"errors"
	g "todo/internal/google"
	b "todo/internal/business"
	pb "todo/proto/todo"
)

type TodoServer struct{
	pb.UnimplementedTodoServer
}

func NewTodoServer(ctx context.Context) pb.TodoServer {
	return &TodoServer{}
}

func (s *TodoServer) AddTodo(ctx context.Context, in *pb.AddTodoRequest) (*pb.BasicReply, error) {
	if err := s.CheckLogin(); err != nil {return nil, err}
	return b.AddTodo(ctx, in)
}

func (s *TodoServer) DeleteTodo(ctx context.Context, in *pb.UpdateTodoRequest) (*pb.BasicReply, error) {
	if err := s.CheckLogin(); err != nil {return nil, err}
	return b.DeleteTodo(ctx, in)
}

func (s *TodoServer) ListTodo(ctx context.Context, in *pb.ListTodoRequest) (*pb.ListTodoReply, error) {
	if err := s.CheckLogin(); err != nil {return nil, err}
	return b.ListTodo(ctx, in)
}

func (s *TodoServer) MarkTodo(ctx context.Context, in *pb.UpdateTodoRequest) (*pb.BasicReply, error) {
	if err := s.CheckLogin(); err != nil {return nil, err}
	return b.MarkTodo(ctx, in)
}

func (s *TodoServer) Ping(ctx context.Context, in *pb.EmptyRequest) (*pb.PingReply, error) {
	if err := s.CheckLogin(); err != nil {return nil, err}
	return b.Ping(ctx, in)
}

func (s *TodoServer) CheckLogin() error {
	if g.UserDetails.Email == "" {
		return errors.New("user not logged in. Please log in using http://localhost:8081")
	}

	return nil
}