package Week02

import (
	"database/sql"
	"errors"
	"fmt"
	pkgerrors "github.com/pkg/errors"
	"testing"
)

type Task struct {
	ID string
	name string
}


func FindTaskById(task_id string) (*Task, error) {
	//...
	return nil, FakeNoDataError()
}

func FakeNoDataError() error {
	return sql.ErrNoRows
}

type GetTaskRequest struct {
	ID string  `json:"id" binding:"required"`
}

func QueryTaskById(param *GetTaskRequest) (*Task, error) {
	task, err := FindTaskById(param.ID)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "QueryTaskById:")
	}
	return task, nil
}

func TestErrorNoRows(t *testing.T) {
	param := GetTaskRequest{
		ID:"haili007",
	}
	task, err := QueryTaskById(&param)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("data not found")
		}
		fmt.Printf("%+v", err)
	}
	fmt.Println(task)
}

