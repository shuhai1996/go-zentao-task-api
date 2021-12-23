package example

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go-zentao-task/pkg/db"
	"time"
)

type DemoGin struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	ClientCode   string    `json:"client_code"`
	ClientSecret string    `json:"-"`
	Status       int       `json:"status"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
}

func (DemoGin) TableName() string {
	return "demo_gin"
}

func FindAll(clientCode string) ([]DemoGin, error) {
	var result []DemoGin
	if err := db.Orm.Where(&DemoGin{
		ClientCode: clientCode,
		Status:     1,
	}).Find(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return result, nil
}

func FindOne(clientCode string) (*DemoGin, error) {
	// NOTE When query with struct, GORM will only query with those fields has non-zero value
	// that means if your field’s value is 0, '', false or other zero values
	// it won’t be used to build query conditions

	// SELECT * FROM demo_gin WHERE client_code = "client_code";
	var result DemoGin
	if err := db.Orm.Where(&DemoGin{
		ClientCode: clientCode,
		Status:     0,
	}).First(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &result, nil
}

func Update(clientCode, clientSecret string) (int64, error) {
	var demogin DemoGin
	// op := db.Orm.Table("demo_gin").Where("client_code = ?", clientCode).Updates
	op := db.Orm.Model(&demogin).Where("client_code = ?", clientCode).Updates(map[string]interface{}{
		"client_secret": clientSecret,
		"update_time":   time.Now().Format("2006-01-02 15:04:05"),
	})
	if op.Error != nil {
		return 0, op.Error
	}
	return op.RowsAffected, nil
}

func Create() (int, error) {
	data := &DemoGin{
		Name:         "name",
		ClientCode:   "client_code",
		ClientSecret: "client_secret",
		Status:       1,
		CreateTime:   time.Now(),
	}
	op := db.Orm.Create(data)
	if op.Error != nil {
		return 0, op.Error
	}
	return data.ID, nil
}

func Delete(id int) (int64, error) {
	var data DemoGin
	db.Orm.Where(&DemoGin{ID: id}).First(&data)
	if data.ID == 0 {
		return 0, errors.New("record not found")
	}

	op := db.Orm.Delete(&data)
	if op.Error != nil {
		return 0, op.Error
	}
	return op.RowsAffected, nil
}

func RawSql() (*DemoGin, error) {
	var result DemoGin
	if err := db.Orm.Raw("select * from demo_gin where id = ?", 1).Scan(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &result, nil
}

type WireExample1 struct {
}

func NewWireExample1() *WireExample1 {
	return &WireExample1{}
}

func (*WireExample1) A() string {
	return "WireExample1 Model: A"
}

type WireExample2 struct {
}

func NewWireExample2() *WireExample2 {
	return &WireExample2{}
}

func (*WireExample2) A() string {
	return "WireExample2 Model: A"
}
