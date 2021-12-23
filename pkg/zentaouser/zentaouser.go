package zentaouser

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"go-zentao-task/model/zentao"
	"go-zentao-task/pkg/config"
)

func Setup() *zentao.User {
	account, pass, err := getUserConfig("")
	if err != nil {
		fmt.Println(err)
	}
	return UserLogin(account, pass)
}

func getUserConfig(store string) (string, string, error) {
	account := config.Get("user.account" + store)
	password := config.Get("user.password" + store)

	if account == "" {
		return "", "", errors.New("用户account不能为空")
	}
	if password == "" {
		return "", "", errors.New("用户password不能为空")
	}
	return account, password, nil
}

func UserLogin(account string, password string) *zentao.User {
	user := zentao.NewUser()
	u, err := user.FindOneByAccount(account)
	if err != nil {
		fmt.Println("用户不存在")
		return nil
	}
	md := md5.Sum([]byte(password))   // 进行md5加密
	pass := hex.EncodeToString(md[:]) //[16]byte转成切片再转成string
	if bytes.Compare([]byte(u.Password), []byte(pass)) != 0 {
		fmt.Println("密码不正确")
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	fmt.Println(u.Realname + "登陆成功！")
	return u
}
