package tasks

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"rest_api_sample/pkg/db"
	"rest_api_sample/pkg/registry"
)

// BucketName repository bucket name
const BucketName = "tasks"

var repo *Repository

type Repository struct {
	DBService *db.Service
}

func New(ctx context.Context) *Repository {
	dbService, ok := ctx.Value(db.ContextKey).(*db.Service)
	if !ok {
		log.Panic("could not get database connection pool from context")
	}

	_, err := dbService.CreateBucket(BucketName)
	if err != nil {
		log.Panicf("create bucket %s failed: %s", BucketName, err.Error())
	}

	return &Repository{DBService: dbService}
}

func GetRepository() *Repository {
	return repo
}

func (r *Repository) Create(task Task) (Task, error) {
	taskID := uuid.New().String()
	task.ID = taskID
	err := r.DBService.SetJson(taskID,BucketName, task)
	return task, err
}

func (r *Repository) Update(taskID string, task Task) (Task, error) {

	ok, err := r.DBService.IsExist(taskID, BucketName)
	if err != nil{
		return Task{}, err
	}

	if !ok{
		return Task{}, fmt.Errorf("task not found: %s", taskID)
	}

	err = r.DBService.SetJson(taskID, BucketName, task)
	if err != nil{
		return Task{}, err
	}

	return task, nil
}

func (r *Repository) Delete(taskID string) error {

	ok, err := r.DBService.IsExist(taskID, BucketName)
	if err != nil{
		return err
	}

	if !ok{
		return fmt.Errorf("task not found: %s", taskID)
	}

	return r.DBService.Delete(taskID, BucketName)
}

func (r *Repository) GetOne(taskID string) (Task, error) {

	var task Task
	err := r.DBService.GetJson(taskID, BucketName, &task)
	if err != nil{
		return Task{}, err
	}
	return task, nil
}

func (r *Repository) GetAll() ([]Task, error) {
	var tasks []Task
	err := r.DBService.GetJsonList(BucketName, &tasks)
	if err != nil{
		return nil, err
	}
	return tasks, nil
}


func init() {

	fmt.Println("tasks repo init")
	// make sure that our bucket is exit
	registry.Register(func(ctx context.Context) error {
		repo = New(ctx)
		return nil
	}, 5, true)
}