package zentao

import (
	"github.com/jinzhu/gorm"
	"go-zentao-task/pkg/db"
	"time"
)

type Action struct {
	ID         int       `json:"id"`
	ObjectID   int       `db:"objectID" gorm:"column:objectID"`
	ObjectType string    `db:"objectType" gorm:"column:objectType"`
	Product    string    `json:"product"`
	Project    int       `json:"project"`
	Execution  int       `json:"execution"`
	Actor      string    `json:"actor"`
	Action     string    `json:"action"`
	Extra      string    `json:"extra"`
	Date       time.Time `json:"date"`
}

const ActionLogin = "login"
const ActionStart = "started"
const ActionRecorde = "recordestimate"

func (Action) TableName() string {
	return "zt_action"
}

func NewAction() *Action {
	return &Action{}
}

func (*Action) Create(objectID int, objectType string, product string, project int, execution int, actor string, action string, extra string) (int, error) {
	data := &Action{
		ObjectID:   objectID,
		ObjectType: objectType,
		Product:    product,
		Project:    project,
		Execution:  execution,
		Actor:      actor,
		Action:     action,
		Extra:      extra,
		Date:       time.Now(),
	}
	op := db.Orm.Create(data)
	if op.Error != nil {
		return 0, op.Error
	}
	return data.ID, nil
}

func (*Action) FindLastLogin(actor string) (*Action, error) {
	var result Action
	if err := db.Orm.Where(&Action{
		Actor:  actor,
		Action: ActionLogin,
	}).Order("id desc").First(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &result, nil
}
