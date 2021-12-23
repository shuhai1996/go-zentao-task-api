package zentao

import (
	"github.com/jinzhu/gorm"
	"go-zentao-task/pkg/db"
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
	if err := db.Orm.Where(&ProjectProduct{
		Project: project,
	}).First(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &result, nil
}
