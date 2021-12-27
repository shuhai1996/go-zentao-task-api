package zentao

import (
	"fmt"
	"go-zentao-task-api/model/zentao"
)

type TaskService struct {
	T *zentao.Task
}

var Task *TaskService

func InitializeTaskService() *TaskService {
	task := zentao.NewTask()
	Task = &TaskService{T:task}
	return Task
}

// GetAllTaskByAccount 根据account获取任务
func (ts *TaskService) GetAllTaskByAccount(account string, status string) []zentao.Task {
	tasks, err := ts.T.FindAll(account, status, 0, "asc")
	if err != nil {
		fmt.Println("获取任务异常")
		return nil
	}
	return tasks
}