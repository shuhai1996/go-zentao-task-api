package zentao

import (
	"github.com/gomodule/redigo/redis"
	"go-zentao-task-api/pkg/gredis"
	"go-zentao-task-api/pkg/util"
)

const (
	UserInfoPrefix = "zentao:userinfo:"
)

func (service *Service) GetUserInfoByToken(token string) (*util.Response, error) {
	conn := gredis.RedisPool.Get()
	defer conn.Close()

	cache, err := redis.Values(conn.Do("HGETALL", UserInfoPrefix+token))
	if err != nil && err != redis.ErrNil {
		return util.ReturnResponse(400, err.Error(), nil), err
	}

	data := service.User

	if err := redis.ScanStruct(cache, data); err != nil {
		return util.ReturnResponse(400, err.Error(), nil), err
	}

	if data.Account == "" {
		return util.ReturnResponse(400, "登录已失效", nil), nil
	}

	conn.Do("expire", UserInfoPrefix+token, 1800) //nolint 鉴权成功后，token重置有效期为30分钟

	return util.ReturnResponse(200, "success", map[string]interface{}{
		"user_id":   data.ID,
		"user_name": data.Nickname,
		"account":   data.Account,
	}), nil
}
