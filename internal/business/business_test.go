package internal

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	data "todo/internal/data"
	pb "todo/proto/todo"

	"github.com/google/uuid"
)

var (
	testUserId     = uuid.New()
	testTodoListId = uuid.New()
	testItemId     = uuid.New()
)

func Test_AddNewUser(t *testing.T) {
	testCases := []struct {
		testName    string
		inEmail     string
		wantErr     bool
		expectedErr error
		mockFunc    func()
	}{
		{
			testName:    "Fail - missing email",
			inEmail:     "",
			wantErr:     true,
			expectedErr: errors.New("missing email"),
			mockFunc:    func() {},
		},
		{
			testName:    "Success",
			inEmail:     "test@email.com",
			wantErr:     false,
			expectedErr: nil,
			mockFunc: func() {
				data.AddTodoList = func(ctx context.Context) (uuid.UUID, error) {
					return testTodoListId, nil
				}
				data.AddUser = func(ctx context.Context, email string, todoListId uuid.UUID) (uuid.UUID, error) {
					return testUserId, nil
				}
			},
		},
	}

	// preserve original function
	oriAddTodoList := data.AddTodoList
	oriAddUser := data.AddUser

	for _, tc := range testCases {
		t.Run(tc.testName, func(tt *testing.T) {
			tc.mockFunc()
			_, err := AddNewUser(context.Background(), tc.inEmail)
			if tc.wantErr && errors.Is(err, tc.expectedErr) {
				tt.Errorf("AddNewUser failed, got err: %v, want err: %v", err, tc.expectedErr)
			}
			if !tc.wantErr && err != nil {
				tt.Errorf("AddNewUser failed, not expecting err: %v", err)
			}
		})
	}

	// reset
	data.AddTodoList = oriAddTodoList
	data.AddUser = oriAddUser
}

func Test_CheckUserExists(t *testing.T) {
	testCases := []struct {
		testName    string
		inEmail     string
		expectedOut bool
		wantErr     bool
		expectedErr error
		mockFunc    func()
	}{
		{
			testName:    "Fail - missing email",
			inEmail:     "",
			expectedOut: false,
			wantErr:     true,
			expectedErr: errors.New("missing email"),
			mockFunc:    func() {},
		},
		{
			testName:    "Success - user exist",
			inEmail:     "test@email.com",
			expectedOut: true,
			wantErr:     false,
			expectedErr: nil,
			mockFunc: func() {
				data.GetUser = func(ctx context.Context, email string) (data.User, error) {
					return data.User{
						Id: testUserId,
					}, nil
				}
			},
		},
		{
			testName:    "Success - user do not exist",
			inEmail:     "test2@email.com",
			expectedOut: true,
			wantErr:     false,
			expectedErr: nil,
			mockFunc: func() {
				data.GetUser = func(ctx context.Context, email string) (data.User, error) {
					return data.User{
						Id: uuid.Nil,
					}, nil
				}
			},
		},
	}

	// preserve original function
	oriGetUser := data.GetUser

	for _, tc := range testCases {
		t.Run(tc.testName, func(tt *testing.T) {
			tc.mockFunc()
			out, err := CheckUserExists(context.Background(), tc.inEmail)
			if tc.wantErr && errors.Is(err, tc.expectedErr) {
				tt.Errorf("CheckUserExists failed, got err: %v, want err: %v", err, tc.expectedErr)
			}
			if !tc.wantErr && err != nil {
				tt.Errorf("CheckUserExists failed, not expecting err: %v", err)
			}
			if !reflect.DeepEqual(out, tc.expectedOut) {
				tt.Errorf("CheckUserExists failed, got out: %v, want out: %v", out, tc.expectedOut)
			}
		})
	}

	// reset
	data.GetUser = oriGetUser
}

func Test_AddTodo(t *testing.T) {
	testCases := []struct {
		testName    string
		inEmail     string
		inReq       *pb.AddTodoRequest
		wantErr     bool
		expectedErr error
		mockFunc    func()
	}{
		{
			testName: "Fail - missing email",
			inEmail:  "",
			inReq: &pb.AddTodoRequest{
				ItemName:        "item1",
				ItemDescription: "desc1",
			},
			wantErr:     true,
			expectedErr: errors.New("missing email"),
			mockFunc:    func() {},
		},
		{
			testName: "Fail - missing itemName",
			inEmail:  "test@email.com",
			inReq: &pb.AddTodoRequest{
				ItemName:        "",
				ItemDescription: "desc1",
			},
			wantErr:     true,
			expectedErr: errors.New("missing itemName"),
			mockFunc:    func() {},
		},
		{
			testName: "Fail - missing itemDesc",
			inEmail:  "test@email.com",
			inReq: &pb.AddTodoRequest{
				ItemName:        "item1",
				ItemDescription: "",
			},
			wantErr:     true,
			expectedErr: errors.New("missing itemDescription"),
			mockFunc:    func() {},
		},
		{
			testName: "Fail - user do not exist",
			inEmail:  "test2@email.com",
			inReq: &pb.AddTodoRequest{
				ItemName:        "item1",
				ItemDescription: "desc1",
			},
			wantErr:     true,
			expectedErr: sql.ErrNoRows,
			mockFunc: func() {
				data.GetUser = func(ctx context.Context, email string) (data.User, error) {
					return data.User{
						Id: uuid.Nil,
					}, errors.New(sql.ErrNoRows.Error())
				}
			},
		},
		{
			testName: "Success",
			inEmail:  "test@email.com",
			inReq: &pb.AddTodoRequest{
				ItemName:        "item1",
				ItemDescription: "desc1",
			},
			wantErr:     false,
			expectedErr: sql.ErrNoRows,
			mockFunc: func() {
				data.GetUser = func(ctx context.Context, email string) (data.User, error) {
					return data.User{
						Id: testUserId,
					}, nil
				}
				data.GetTodoListIdByUserId = func(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
					return testTodoListId, nil
				}
				data.AddItem = func(ctx context.Context, userId uuid.UUID, todoListId uuid.UUID, itemName string, itemDescription string) (uuid.UUID, error) {
					return uuid.New(), nil
				}
			},
		},
	}

	// preserve original function
	oriGetUser := data.GetUser
	oriGetTodoListIdByUserId := data.GetTodoListIdByUserId
	oriAddTodo := data.AddItem

	for _, tc := range testCases {
		t.Run(tc.testName, func(tt *testing.T) {
			tc.mockFunc()
			_, err := AddTodo(context.Background(), tc.inEmail, tc.inReq)
			if tc.wantErr && errors.Is(err, tc.expectedErr) {
				tt.Errorf("AddTodo failed, got err: %v, want err: %v", err, tc.expectedErr)
			}
			if !tc.wantErr && err != nil {
				tt.Errorf("AddTodo failed, not expecting err: %v", err)
			}
		})
	}

	// reset
	data.GetUser = oriGetUser
	data.GetTodoListIdByUserId = oriGetTodoListIdByUserId
	data.AddItem = oriAddTodo
}

func Test_MarkTodo(t *testing.T) {
	testCases := []struct {
		testName    string
		inEmail     string
		inReq       *pb.UpdateTodoRequest
		wantErr     bool
		expectedErr error
		mockFunc    func()
	}{
		{
			testName: "Fail - missing email",
			inEmail:  "",
			inReq: &pb.UpdateTodoRequest{
				ItemName: "item1",
			},
			wantErr:     true,
			expectedErr: errors.New("missing email"),
			mockFunc:    func() {},
		},
		{
			testName: "Fail - missing itemName",
			inEmail:  "test@email.com",
			inReq: &pb.UpdateTodoRequest{
				ItemName: "",
			},
			wantErr:     true,
			expectedErr: errors.New("missing itemName"),
			mockFunc:    func() {},
		},
		{
			testName: "Fail - user do not exist",
			inEmail:  "test2@email.com",
			inReq: &pb.UpdateTodoRequest{
				ItemName: "item1",
			},
			wantErr:     true,
			expectedErr: sql.ErrNoRows,
			mockFunc: func() {
				data.GetUser = func(ctx context.Context, email string) (data.User, error) {
					return data.User{
						Id: uuid.Nil,
					}, errors.New(sql.ErrNoRows.Error())
				}
			},
		},
		{
			testName: "Success",
			inEmail:  "test@email.com",
			inReq: &pb.UpdateTodoRequest{
				ItemName: "item1",
			},
			wantErr:     false,
			expectedErr: sql.ErrNoRows,
			mockFunc: func() {
				data.GetUser = func(ctx context.Context, email string) (data.User, error) {
					return data.User{
						Id: testUserId,
					}, nil
				}
				data.GetTodoListIdByUserId = func(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
					return testTodoListId, nil
				}
				data.GetItemByItemName = func(ctx context.Context, todoListId uuid.UUID, itemName string) (data.Item, error) {
					return data.Item{
						Id: testItemId,
					}, nil
				}
				data.UpdateItem = func(ctx context.Context, itemId string, item data.Item) (bool, error) {
					return true, nil
				}
			},
		},
	}

	// preserve original function
	oriGetUser := data.GetUser
	oriGetTodoListIdByUserId := data.GetTodoListIdByUserId
	oriGetItemByItemName := data.GetItemByItemName
	oriUpdateItem := data.UpdateItem

	for _, tc := range testCases {
		t.Run(tc.testName, func(tt *testing.T) {
			tc.mockFunc()
			_, err := MarkTodo(context.Background(), tc.inEmail, tc.inReq)
			if tc.wantErr && errors.Is(err, tc.expectedErr) {
				tt.Errorf("MarkTodo failed, got err: %v, want err: %v", err, tc.expectedErr)
			}
			if !tc.wantErr && err != nil {
				tt.Errorf("MarkTodo failed, not expecting err: %v", err)
			}
		})
	}

	// reset
	data.GetUser = oriGetUser
	data.GetTodoListIdByUserId = oriGetTodoListIdByUserId
	data.GetItemByItemName = oriGetItemByItemName
	data.UpdateItem = oriUpdateItem
}

func Test_DeleteTodo(t *testing.T) {
	testCases := []struct {
		testName    string
		inEmail     string
		inReq       *pb.UpdateTodoRequest
		wantErr     bool
		expectedErr error
		mockFunc    func()
	}{
		{
			testName: "Fail - missing email",
			inEmail:  "",
			inReq: &pb.UpdateTodoRequest{
				ItemName: "item1",
			},
			wantErr:     true,
			expectedErr: errors.New("missing email"),
			mockFunc:    func() {},
		},
		{
			testName: "Fail - missing itemName",
			inEmail:  "test@email.com",
			inReq: &pb.UpdateTodoRequest{
				ItemName: "",
			},
			wantErr:     true,
			expectedErr: errors.New("missing itemName"),
			mockFunc:    func() {},
		},
		{
			testName: "Fail - user do not exist",
			inEmail:  "test2@email.com",
			inReq: &pb.UpdateTodoRequest{
				ItemName: "item1",
			},
			wantErr:     true,
			expectedErr: sql.ErrNoRows,
			mockFunc: func() {
				data.GetUser = func(ctx context.Context, email string) (data.User, error) {
					return data.User{
						Id: uuid.Nil,
					}, errors.New(sql.ErrNoRows.Error())
				}
			},
		},
		{
			testName: "Success",
			inEmail:  "test@email.com",
			inReq: &pb.UpdateTodoRequest{
				ItemName: "item1",
			},
			wantErr:     false,
			expectedErr: sql.ErrNoRows,
			mockFunc: func() {
				data.GetUser = func(ctx context.Context, email string) (data.User, error) {
					return data.User{
						Id: testUserId,
					}, nil
				}
				data.GetTodoListIdByUserId = func(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
					return testTodoListId, nil
				}
				data.GetItemByItemName = func(ctx context.Context, todoListId uuid.UUID, itemName string) (data.Item, error) {
					return data.Item{
						Id: testItemId,
					}, nil
				}
				data.UpdateItem = func(ctx context.Context, itemId string, item data.Item) (bool, error) {
					return true, nil
				}
			},
		},
	}

	// preserve original function
	oriGetUser := data.GetUser
	oriGetTodoListIdByUserId := data.GetTodoListIdByUserId
	oriGetItemByItemName := data.GetItemByItemName
	oriUpdateItem := data.UpdateItem

	for _, tc := range testCases {
		t.Run(tc.testName, func(tt *testing.T) {
			tc.mockFunc()
			_, err := DeleteTodo(context.Background(), tc.inEmail, tc.inReq)
			if tc.wantErr && errors.Is(err, tc.expectedErr) {
				tt.Errorf("DeleteTodo failed, got err: %v, want err: %v", err, tc.expectedErr)
			}
			if !tc.wantErr && err != nil {
				tt.Errorf("DeleteTodo failed, not expecting err: %v", err)
			}
		})
	}

	// reset
	data.GetUser = oriGetUser
	data.GetTodoListIdByUserId = oriGetTodoListIdByUserId
	data.GetItemByItemName = oriGetItemByItemName
	data.UpdateItem = oriUpdateItem
}

func Test_ListTodo(t *testing.T) {
	testCases := []struct {
		testName    string
		inEmail     string
		expectedOut *pb.ListTodoReply
		wantErr     bool
		expectedErr error
		mockFunc    func()
	}{
		{
			testName:    "Fail - missing email",
			inEmail:     "",
			expectedOut: &pb.ListTodoReply{},
			wantErr:     true,
			expectedErr: errors.New("missing email"),
			mockFunc:    func() {},
		},
		{
			testName:    "Fail - user do not exist",
			inEmail:     "test2@email.com",
			expectedOut: &pb.ListTodoReply{},
			wantErr:     true,
			expectedErr: sql.ErrNoRows,
			mockFunc: func() {
				data.GetUser = func(ctx context.Context, email string) (data.User, error) {
					return data.User{
						Id: uuid.Nil,
					}, errors.New(sql.ErrNoRows.Error())
				}
			},
		},
		{
			testName: "Success",
			inEmail:  "test@email.com",
			expectedOut: &pb.ListTodoReply{
				Count: 1,
				Items: []*pb.TodoItem{
					{
						ItemName:        "test1",
						ItemDescription: "desc1",
						Done:            false,
					},
				},
			},
			wantErr:     false,
			expectedErr: sql.ErrNoRows,
			mockFunc: func() {
				data.GetUser = func(ctx context.Context, email string) (data.User, error) {
					return data.User{
						Id: testUserId,
					}, nil
				}
				data.GetTodoListIdByUserId = func(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
					return testTodoListId, nil
				}
				data.ListItem = func(ctx context.Context, todoListId uuid.UUID) ([]data.Item, error) {
					return []data.Item{
						{
							Id:          uuid.New(),
							Name:        "test1",
							Description: "desc1",
							MarkDone:    false,
						},
					}, nil
				}
			},
		},
	}

	// preserve original function
	oriGetUser := data.GetUser
	oriGetTodoListIdByUserId := data.GetTodoListIdByUserId
	oriListItem := data.ListItem

	for _, tc := range testCases {
		t.Run(tc.testName, func(tt *testing.T) {
			tc.mockFunc()
			out, err := ListTodo(context.Background(), tc.inEmail)
			if tc.wantErr && errors.Is(err, tc.expectedErr) {
				tt.Errorf("ListTodo failed, got err: %v, want err: %v", err, tc.expectedErr)
			}
			if !tc.wantErr && err != nil {
				tt.Errorf("ListTodo failed, not expecting err: %v", err)
			}
			if !reflect.DeepEqual(out, tc.expectedOut) {
				tt.Errorf("ListTodo failed, got out: %v, want out: %v", out, tc.expectedOut)
			}
		})
	}

	// reset
	data.GetUser = oriGetUser
	data.GetTodoListIdByUserId = oriGetTodoListIdByUserId
	data.ListItem = oriListItem
}
