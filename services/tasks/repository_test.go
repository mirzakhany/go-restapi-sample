package tasks

import (
	"context"
	"os"
	"testing"

	"github.com/mirzakhany/rest_api_sample/pkg/db"
)

func TestMain(m *testing.M) {

	dbService, err := db.New("/tmp/test_tasks")
	if err != nil {
		panic(err)
	}
	ctx := context.WithValue(context.Background(), db.ContextKey, dbService)
	New(ctx)

	code := m.Run()

	_ = os.Remove("/tmp/test_tasks")
	os.Exit(code)
}

func TestGetRepository(t *testing.T) {
	r := GetRepository()
	if r == nil {
		t.Error("repo is nil")
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

func TestRepository_GetOne(t *testing.T) {

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

	savedTask, err := GetRepository().GetOne(res.ID)
	if err != nil {
		t.Errorf("error in loading task, %s", err)
	}
	if savedTask.Title != task.Title || savedTask.Status != task.Status ||
		savedTask.Estimate != task.Estimate || savedTask.Sprint != task.Sprint || savedTask.Assignee != task.Assignee {
		t.Error("task is not same as original")
	}
}

func TestRepository_Update(t *testing.T) {

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

	task.Title = "new-title-value"

	updatedTask, err := GetRepository().Update(res.ID, task)
	if err != nil {
		t.Errorf("error in update task, %s", err)
	}

	if updatedTask.Title != task.Title || updatedTask.Status != task.Status ||
		updatedTask.Estimate != task.Estimate || updatedTask.Sprint != task.Sprint ||
		updatedTask.Assignee != task.Assignee {
		t.Error("task is not same as last updated task")
	}
}

func TestRepository_Delete(t *testing.T) {

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

	err = GetRepository().Delete(res.ID)
	if err != nil {
		t.Errorf("error in delete task, %s", err)
	}

}

func TestRepository_GetAll(t *testing.T) {

	emptyBucket()

	var tasks = []Task{{
		Title:    "test1",
		Sprint:   "bar1",
		Estimate: "1",
		Status:   "in-progress",
		Assignee: "bar",
	}, {
		Title:    "test2",
		Sprint:   "bar",
		Estimate: "2",
		Status:   "in-backlog",
		Assignee: "foo",
	}, {
		Title:    "test3",
		Sprint:   "bar2",
		Estimate: "4",
		Status:   "done",
		Assignee: "bar",
	}}

	repo := GetRepository()

	for _, task := range tasks {
		_, err := repo.Create(task)
		if err != nil {
			t.Errorf("error in create task, %s", err)
		}
	}

	savedTasks, err := repo.GetAll()

	if err != nil {
		t.Errorf("error in loading tasks, %s", err)
	}
	if len(savedTasks) != len(tasks) {
		t.Errorf("result count in not same as input %d != %d", len(savedTasks), len(tasks))
	}
}

func emptyBucket() {
	repo := GetRepository()
	savedTasks, _ := repo.GetAll()
	for _, t := range savedTasks {
		err := repo.Delete(t.ID)
		if err != nil {
			panic(err)
		}
	}
}
