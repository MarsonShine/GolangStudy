package http

import (
	"fmt"
	"kratos-demo/api"
	jb "kratos-demo/internal/jsonpb"
	"net/http"

	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
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
	p := new(api.HelloReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	data, _ := api.DemoSvc.SayHelloURL(c.Context, p)
	c.JSON(data, nil)
}

func getInt64FromProtobuf(c *bm.Context) {
	p := new(api.HelloReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	data, _ := api.DemoSvc.SayHelloURL(c.Context, p)
	c.Render(http.StatusOK, jb.PBJSON{
		Code:    0,
		Message: "0",
		TTL:     1,
		Data:    data,
	})
}
