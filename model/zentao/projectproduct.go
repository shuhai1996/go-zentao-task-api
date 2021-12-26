package zentao

import (
	"errors"
	"go-zentao-task-api/pkg/db"
	"gorm.io/gorm"
)

type ProjectProduct struct {
	Project int `json:"project"`
	Product int `json:"product"`
	Branch  int `json:"branch"`
	Plan    int `json:"plan"`
}

func (ProjectProduct) TableName() string {
	return "zt_projectproduct"
}

func NewProjectProduct() *ProjectProduct {
	return &ProjectProduct{}
}

func (*ProjectProduct) FindOneByProject(project int) (*ProjectProduct, error) {
	var result ProjectProduct
	if res := db.Orm.Where(&ProjectProduct{
		Project: project,
	}).First(&result); res != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, res.Error
	}
	return &result, nil
}
