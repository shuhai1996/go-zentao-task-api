package zentao

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-zentao-task-api/core"
	"go-zentao-task-api/service/zentao"
	"time"
)

type Actions struct {
	*Es
}

type Es struct {
}


type actionIndexRequest struct {
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"pageSize" binding:"required"`
	Account  string `form:"account"`
}

func (*Actions) Index(c *core.Context) {
	var r actionIndexRequest
	if err := c.ShouldBindQuery(&r); err != nil {
		c.Fail(40100, err.Error(), nil)
		return
	}

	total, res :=zentao.GetActionList(r.Page,r.PageSize,``)
	if total == 0 {
		c.Fail(400, "没有行为", nil)
		return
	}
	var list []map[string]interface{}
	for _, v := range res {
		m:= make(map[string]interface{})
		t,_:=json.Marshal(v.Source)
		json.Unmarshal(t, &m)
		TransDateString(m)
		list = append(list, m)
	}

	c.Success(gin.H{
		"data": list,
		"_meta": map[string]int{
			"page": r.Page,
			"pageSize": r.PageSize,
			"total": int(total),
		},
	})
}

func (*Actions) View(c *core.Context) {
	id := c.Param("id")
	res, err :=zentao.FindOne(id)

	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}
	if res == nil {
		c.Fail(400, "找不到该行为", nil)
		return
	}
	m:= make(map[string]interface{})
	va,_:= json.Marshal(res.Source)
	json.Unmarshal(va, &m)
	TransDateString(m)
	c.Success(gin.H{
		"data": m,
	})
}

func TransDateString(m map[string]interface{})  {
	if m["date"] !=nil {
		t:= m["date"].(float64)
		m["create_date"] = time.Unix(int64(t), 0).Format("2006-01-02 15:04:05")
	}
}

func (*Actions) Delete(c *core.Context)  {
	id := c.Param("id")
	res, err :=zentao.DeleteOne(id)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}
	c.Success(gin.H{
		"data": res,
	})
}

