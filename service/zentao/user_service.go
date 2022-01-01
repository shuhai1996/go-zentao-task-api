package zentao

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"go-zentao-task-api/model/zentao"
	"go-zentao-task-api/pkg/gredis"
	"go-zentao-task-api/pkg/util"
	"strconv"
)

type UserService struct {
	U *zentao.User
	A *zentao.Action
}

var User *UserService

func InitializeUserService() *UserService {
	user := zentao.NewUser()
	action := zentao.NewAction()
	User = &UserService{U:user, A: action}
	return User
}


func (us *UserService) UserLogin(account string, password string) (*util.Response, error) {
	conn := gredis.RedisPool.Get() //获取redis 连接
	defer conn.Close()
	u, err := us.U.FindOneByAccount(account)
	if err != nil {
		return nil, err
	}
	md := md5.Sum([]byte(password))   // 进行md5加密
	pass := hex.EncodeToString(md[:]) //[16]byte转成切片再转成string
	if bytes.Compare([]byte(u.Password), []byte(pass)) != 0 {
		return nil, errors.New("密码不正确")
	}
	token := fmt.Sprintf("%x", md5.Sum(uuid.NewV4().Bytes()))
	if _, err := conn.Do("hmset", redis.Args{}.Add(UserInfoPrefix+token).AddFlat(u)...); err != nil {
		return nil, err
	}
	conn.Do("expire", UserInfoPrefix+token, 1800) //nolint 无操作 30分钟后过期
	fmt.Println(u.Realname + "登陆成功！")
	//创建操作记录
	us.A.Create(u.ID, "user", ","+strconv.Itoa(0)+",", 0, 0, u.Account, zentao.ActionLogin, "")
	return util.ReturnResponse(200, "success", map[string]string{
		"access_token": token,
	}), err
}

func (us *UserService) Logout(token string) (*util.Response, error) {
	conn := gredis.RedisPool.Get()
	defer conn.Close()

	_, err := conn.Do("del", UserInfoPrefix+token) //删除token
	if err != nil && err != redis.ErrNil {
		return util.ReturnResponse(400, err.Error(), nil), err
	}

	return util.ReturnResponse(200, "success", nil), nil
}
