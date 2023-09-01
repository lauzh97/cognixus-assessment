package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	data "todo/internal/data"
	pb "todo/proto/todo"

	"github.com/google/uuid"
)

// Adds a new record into items table, related to the logged in user
func AddTodo(ctx context.Context, email string, in *pb.AddTodoRequest) (*pb.EmptyReply, error) {
	// validation
	if email == "" {
		return &pb.EmptyReply{}, errors.New("missing email")
	}

	if in.ItemName == "" {
		return &pb.EmptyReply{}, errors.New("missing itemName")
	}

	if in.ItemDescription == "" {
		return &pb.EmptyReply{}, errors.New("missing itemDescription")
	}
	// end validation

	user, err := data.GetUser(ctx, email)
	if err != nil {
		return &pb.EmptyReply{}, err
	}

	// get todoListId
	todoListId, err := data.GetTodoListIdByUserId(ctx, user.Id)
	if err != nil {
		return &pb.EmptyReply{}, err
	}

	// add item
	_, err = data.AddItem(ctx, user.Id, todoListId, in.ItemName, in.ItemDescription)
	if err != nil {
		return &pb.EmptyReply{}, err
	}

	return &pb.EmptyReply{}, nil
}

// Soft delete an existing record into items table, related to the logged in user
func DeleteTodo(ctx context.Context, email string, in *pb.UpdateTodoRequest) (*pb.EmptyReply, error) {
	// validation
	if email == "" {
		return &pb.EmptyReply{}, errors.New("missing email")
	}

	if in.ItemName == "" {
		return &pb.EmptyReply{}, errors.New("missing itemName")
	}
	// end validation

	user, err := data.GetUser(ctx, email)
	if err != nil {
		return &pb.EmptyReply{}, err
	}

	// get todoListId
	todoListId, err := data.GetTodoListIdByUserId(ctx, user.Id)
	if err != nil {
		return &pb.EmptyReply{}, err
	}

	// get item
	item, err := data.GetItemByItemName(ctx, todoListId, in.ItemName)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.EmptyReply{}, errors.New("item do not exist")
		}
		return &pb.EmptyReply{}, err
	}

	// only update "active" column
	item.Active = false

	// update item (soft delete)
	_, err = data.UpdateItem(ctx, item.Id.String(), item)
	if err != nil {
		return &pb.EmptyReply{}, err
	}

	return &pb.EmptyReply{}, nil
}

func ListTodo(ctx context.Context, email string) (*pb.ListTodoReply, error) {
	// validation
	if email == "" {
		return &pb.ListTodoReply{}, errors.New("missing email")
	}
	// end validation

	user, err := data.GetUser(ctx, email)
	if err != nil {
		return &pb.ListTodoReply{}, err
	}

	// get todoListId
	todoListId, err := data.GetTodoListIdByUserId(ctx, user.Id)
	if err != nil {
		return &pb.ListTodoReply{}, err
	}

	items, err := data.ListItem(ctx, todoListId)
	if err != nil {
		return &pb.ListTodoReply{}, err
	}

	// format into json for reply
	var res pb.ListTodoReply
	res.Count = int32(len(items))
	j, err := json.Marshal(items)
	if err != nil {
		return &pb.ListTodoReply{}, err
	}
	err = json.Unmarshal(j, &res.Items)
	if err != nil {
		return &pb.ListTodoReply{}, err
	}

	return &res, nil
}

func MarkTodo(ctx context.Context, email string, in *pb.UpdateTodoRequest) (*pb.EmptyReply, error) {
	// validation
	if email == "" {
		return &pb.EmptyReply{}, errors.New("missing email")
	}

	if in.ItemName == "" {
		return &pb.EmptyReply{}, errors.New("missing itemName")
	}
	// end validation

	user, err := data.GetUser(ctx, email)
	if err != nil {
		return &pb.EmptyReply{}, err
	}

	// get todoListId
	todoListId, err := data.GetTodoListIdByUserId(ctx, user.Id)
	if err != nil {
		return &pb.EmptyReply{}, err
	}

	// get item
	item, err := data.GetItemByItemName(ctx, todoListId, in.ItemName)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.EmptyReply{}, errors.New("item do not exist")
		}
		return &pb.EmptyReply{}, err
	}

	// update value
	item.MarkDone = true
	item.Active = true

	// update item
	_, err = data.UpdateItem(ctx, item.Id.String(), item)
	if err != nil {
		return &pb.EmptyReply{}, err
	}

	return &pb.EmptyReply{}, nil
}

func Ping(ctx context.Context, in *pb.EmptyRequest) (*pb.PingReply, error) {
	return &pb.PingReply{Pong: "Pong"}, nil
}

// Adds a new record into main.user table
func AddNewUser(ctx context.Context, email string) (uuid.UUID, error) {
	// validation
	if email == "" {
		return uuid.Nil, errors.New("missing email")
	}
	// end validation

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
	// validation
	if email == "" {
		return false, errors.New("missing email")
	}
	// end validation

	_, err := data.GetUser(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		return false, errors.New("GetUser failed: " + err.Error())
	}
	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}
