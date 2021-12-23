package zentao

import (
	"github.com/jinzhu/gorm"
	"go-zentao-task/pkg/db"
)

type User struct {
	ID       int    `json:"id"`
	Dept     string `json:"dept"`
	Type     string `json:"type"`
	Account  string `json:"account"`
	Password string `json:"password"`
	Realname string `json:"realname"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Deleted  int    `json:"deleted"`
}

func (User) TableName() string {
	return "zt_user"
}

func NewUser() *User {
	return &User{}
}

func (*User) FindOneByAccount(account string) (*User, error) {
	var result User
	if err := db.Orm.Where(&User{
		Account: account,
		Deleted: 0,
	}).First(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &result, nil
}
