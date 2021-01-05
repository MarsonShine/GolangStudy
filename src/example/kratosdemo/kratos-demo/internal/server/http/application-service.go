package http

import (
	"fmt"
	"kratos-demo/internal/model"
	"math/rand"

	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
)

// 获取用户信息
func getUserHandler(c *bm.Context) {
	id, success := c.Params.Get("id")
	if !success {
		c.JSONMap(map[string]interface{}{
			"message": "id 不正确",
			"success": success,
			"elapsed": fmt.Sprintf("来自中间件的值=%s", c.Request.Header.Get("ElapsedTime")),
		}, nil)
	} else {
		c.JSONMap(map[string]interface{}{
			"data":    fmt.Sprintf("id = %s", id),
			"success": success,
			"elapsed": fmt.Sprintf("来自中间件的值=%s", c.Request.Header.Get("ElapsedTime")),
		}, nil)
	}
}

// 返回int64类型
func getInt64(c *bm.Context) {
	k := &model.Article{
		ID:      rand.Int63(),
		Content: "这是Content",
		Author:  "marsonshine",
	}
	c.JSON(k, nil)

	// c.Render(http.StatusOK, render.PB{
	// 	Code: 0,
	// 	Message: "0",
	// 	TTL: 1,
	// 	Data: ,
	// })
}

func getInt64FromProtobuf(c *bm.Context) {

}
