package middleware

import (
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
)

func UseRateLimiter(e *bm.Engine) {
	limiter := bm.NewRateLimiter(nil)
	e.Use(limiter.Limit())
}
