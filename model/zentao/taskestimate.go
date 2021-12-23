package zentao

import (
	"github.com/jinzhu/gorm"
	"go-zentao-task/pkg/db"
	"time"
)

type TaskEstimate struct {
	ID       int     `json:"id"`
	Task     int     `json:"task"`
	Date     string  `json:"date"`
	Left     float64 `json:"left"`
	Consumed float64 `json:"consumed"`
	Account  string  `json:"account"`
	Work     string  `json:"work"`
}

func (TaskEstimate) TableName() string {
	return "zt_taskestimate"
}

func (TaskEstimate) GetToday(account string) ([]TaskEstimate, error) {
	var result []TaskEstimate
	date := time.Now().Format("2006-01-02")
	if err := db.Orm.Where(&TaskEstimate{
		Date:    date,
		Account: account,
	}).Find(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return result, nil
}

func NewTaskEstimate() *TaskEstimate {
	return &TaskEstimate{}
}

func (*TaskEstimate) Create(task int, left float64, consumed float64, account string, work string) (int, error) {
	data := &TaskEstimate{
		Task:     task,
		Date:     time.Now().Format("2006-01-02"),
		Left:     left,
		Consumed: consumed,
		Account:  account,
		Work:     work,
	}
	op := db.Orm.Create(data)
	if op.Error != nil {
		return 0, op.Error
	}
	return data.ID, nil
}
