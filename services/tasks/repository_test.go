package tasks

import (
	"context"
	"fmt"
	"rest_api_sample/pkg/db"
	"rest_api_sample/pkg/registry"
	"testing"
)

func TestGetRepository(t *testing.T) {
	r := GetRepository()
	if r == nil {
		t.Error("repo is null")
	}
}

func TestRepository_Create(t *testing.T) {

	var task = Task{
		Title:    "test1",
		Sprint:   "bar",
		Estimate: "1",
		Status:   "in-progress",
		Assignee: "foo",
	}

	res, err := GetRepository().Create(task)
	if err != nil {
		t.Errorf("error in create task, %s", err)
	}

	if res.ID == "" {
		t.Error("task create without id")
	}

	if res.Title != task.Title || res.Status != task.Status ||
		res.Estimate != task.Estimate || res.Sprint != task.Sprint || res.Assignee != task.Assignee {
		t.Error("task is not same as original")
	}
}

func init() {

	fmt.Println("test tasks repo init")
	dbService, err := db.New("/tmp/test_tasks")
	if err != nil {
		panic(err)
	}

	ctx := context.WithValue(context.Background(), db.ContextKey, dbService)
	registry.Run(ctx)
}
