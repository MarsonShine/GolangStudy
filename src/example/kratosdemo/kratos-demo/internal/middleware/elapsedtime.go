package middleware

import (
	"math/rand"
	"strconv"

	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
)

// 计算请求时间中间件
type RecordRequestElapsedTime struct {
}

func NewRecordRequestElapsedTime() RecordRequestElapsedTime {
	return RecordRequestElapsedTime{}
}

func (elapsed RecordRequestElapsedTime) ServeHTTP(c *bm.Context) {
	c.Request.Header.Add("ElapsedTime", strconv.FormatInt(rand.Int63(), 10))
}
