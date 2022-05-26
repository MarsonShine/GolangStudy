package p1

import (
	"fmt"
	"go-cip/10-cacheModel/demo4/p2"
	"go-cip/10-cacheModel/demo4/trace"
)

var V1_p1 = trace.Trace("init v1_p1", p2.V1_p2)
var V2_p1 = trace.Trace("init v2_p1", p2.V2_p2)

func init() {
	fmt.Println("init func in p1")
}
