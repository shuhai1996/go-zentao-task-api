package zentao

import (
	"github.com/gin-gonic/gin"
	"go-zentao-task-api/core"
)

type Task struct {
}

type TaskIndexRequest struct {
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"pageSize" binding:"required"`
	Account  string `form:"account"`
}

type UserViewRequest struct {
	ID int `json:"id"`
}

func (*Task) Index(c *core.Context) {
	var r TaskIndexRequest
	if err := c.ShouldBindQuery(&r); err != nil {
		c.Fail(40100, err.Error(), nil)
		return
	}

	if u := c.Request.URL.Query().Get("account"); u != "" {
		service.User.Account = u
	} else {
		user := c.GetStringMap("user_info")
		u = user["account"].(string)
	}

	res := service.GetAllTaskNotDone()
	if res == nil {
		c.Fail(400, "找不到任务", nil)
		return
	}
	var list []map[string]interface{}
	for _, v := range res {
		list = append(list, map[string]interface{}{
			"id":             v.ID,
			"name":           v.Name,
			"status":         v.Status,
			"consumed":       v.Consumed,
			"left":           v.Left,
			"estimate":       v.Estimate,
			"last_edit_time": v.LastEditedDate.Format("2006-01-02 15:04:05"),
		})
	}

	c.Success(gin.H{
		"data": list,
		"_meta": map[string]int{
			"page":     r.PageSize,
			"pageSize": r.Page,
		},
	})
}
