package zentao

import (
	"errors"
	"go-zentao-task-api/pkg/db"
	"gorm.io/gorm"
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

var es = NewEsAction()

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
	_,err := es.Create(data)
	if err!= nil {
		return 0, err
	}
	return data.ID, nil
}

func (*Action) FindLastLogin(actor string) (*Action, error) {
	var result Action
	if res := db.Orm.Where(&Action{
		Actor:  actor,
		Action: ActionLogin,
	}).Order("id desc").First(&result); res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, res.Error
	}
	return &result, nil
}
