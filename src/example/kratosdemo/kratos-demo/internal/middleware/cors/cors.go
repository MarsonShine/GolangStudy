package cors

import (
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
)

var handFunc bm.HandlerFunc

// 计算请求时间中间件
type CorsMiddleware struct{}

func NewCors() CorsMiddleware {
	return CorsMiddleware{}
}

// func (cors CorsMiddleware) ServeHTTP(c *bm.Context) {

// }

func (core CorsMiddleware) UseCros(engine *bm.Engine) {
	handFunc = bm.CORS([]string{"localhost"})
	engine.Use(handFunc)
}
