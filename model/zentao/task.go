package zentao

import (
	"github.com/jinzhu/gorm"
	"go-zentao-task/pkg/db"
	"time"
)

type Task struct {
	ID             int       `json:"id"`
	Project        int       `json:"project"`
	Parent         int       `json:"parent"`
	Module         int       `json:"module"`
	Status         string    `json:"status"`
	Type           string    `json:"type"`
	Name           string    `json:"name"`
	Estimate       string    `json:"estimate"`
	Consumed       float64   `json:"consumed"`
	Left           float64   `json:"left"`
	Execution      int       `json:"execution"`
	AssignedTo     string    `json:"assignedTo" gorm:"column:assignedTo"`
	FromBug        int       `json:"fromBug" gorm:"column:fromBug"`
	FinishedBy     string    `json:"finishedBy" gorm:"column:finishedBy"`
	FinishedDate   time.Time `json:"finishedDate" gorm:"column:finishedDate"`
	LastEditedDate time.Time `json:"lastEditedDate" gorm:"column:lastEditedDate"`
	Deleted        int       `json:"deleted"`
}

func (Task) TableName() string {
	return "zt_task"
}

func NewTask() *Task {
	return &Task{}
}

const TypeDev = "devel"           // 开发
const TypeDiscuss = "discuss"     //讨论，如"开会"等
const StatusDo = "doing"          // 开发中
const StatusWait = "wait"         //等待中
const StatusFinished = "finished" // 完成

func (*Task) FindOne(id int) (*Task, error) {
	var result Task
	if err := db.Orm.Where(&Task{
		ID: id,
	}).First(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &result, nil
}

func (*Task) FindAll(assignedTo, status string, left float64, leftSort string) ([]Task, error) {
	var result []Task
	odb := db.Orm.Where("assignedTo = ?", assignedTo)

	if status != "" {
		odb = odb.Where("status = ?", status)
	} else {
		odb = odb.Where("status = 'doing' or status = 'wait'")
	}
	if left > 0 {
		odb = odb.Where("`left` > ? ", 0)
	}

	if err := odb.Order("status desc").Order("`left` " + leftSort).Order("id desc").Find(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return result, nil
}

func (*Task) UpdateOne(task int, estimate float64, left float64, actor string, finishDate time.Time, status string) (int64, error) {
	op := db.Orm.Model(&Task{}).Where(&Task{ID: task}).Updates(map[string]interface{}{
		"consumed":       estimate,
		"left":           left,
		"finishedBy":     actor,
		"finishedDate":   finishDate,
		"lastEditedDate": time.Now(),
		"status":         status,
	})
	if op.Error != nil {
		return 0, op.Error
	}
	return op.RowsAffected, nil
}
