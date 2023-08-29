package internal

import (
	"context"
	b "todo/internal/business"
	pb "todo/proto/todo"
)

type TodoServer struct{
	pb.UnimplementedTodoServer
}

func (s *TodoServer) AddTodo(ctx context.Context, in *pb.AddTodoRequest) (*pb.BasicReply, error) {
	return b.AddTodo(ctx, in)
}

func (s *TodoServer) DeleteTodo(ctx context.Context, in *pb.UpdateTodoRequest) (*pb.BasicReply, error) {
	return b.DeleteTodo(ctx, in)
}

func (s *TodoServer) ListTodo(ctx context.Context, in *pb.ListTodoRequest) (*pb.ListTodoReply, error) {
	return b.ListTodo(ctx, in)
}

func (s *TodoServer) MarkTodo(ctx context.Context, in *pb.UpdateTodoRequest) (*pb.BasicReply, error) {
	return b.MarkTodo(ctx, in)
}

func (s *TodoServer) Ping(ctx context.Context, in *pb.EmptyRequest) (*pb.PingReply, error) {
	return b.Ping(ctx, in)
}
