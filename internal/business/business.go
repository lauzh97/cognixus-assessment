package internal

import (
	"context"
	pb "todo/proto/todo"
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