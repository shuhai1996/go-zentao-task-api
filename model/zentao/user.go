package zentao

import (
	"errors"
	"go-zentao-task-api/pkg/db"
	"gorm.io/gorm"
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
	if res := db.Orm.Where(&User{
		Account: account,
		Deleted: 0,
	}).First(&result); res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("用户不存在")
	}
	return &result, nil
}
