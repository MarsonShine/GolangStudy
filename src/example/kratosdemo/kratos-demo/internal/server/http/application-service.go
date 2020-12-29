package http

import (
	"fmt"

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
