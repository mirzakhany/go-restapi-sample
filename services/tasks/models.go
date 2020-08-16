package tasks

// Task is db model for single task
type Task struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Sprint   string `json:"sprint"`
	Estimate string `json:"estimate"`
	Status   string `json:"status"`
	Assignee string `json:"assignee"`
}

// TaskRequest create or update task request
type TaskRequest struct {
	Title    string `form:"title" json:"title" xml:"title" binding:"required"`
	Sprint   string `form:"sprint" json:"sprint" xml:"sprint" binding:"required"`
	Estimate string `form:"estimate" json:"estimate" xml:"estimate" binding:"required"`
	Status   string `form:"status" json:"status" xml:"status" binding:"required"`
	Assignee string `form:"assignee" json:"assignee" xml:"assignee" binding:"required"`
}
