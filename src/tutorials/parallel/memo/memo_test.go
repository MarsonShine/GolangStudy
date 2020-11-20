package memo_test

import (
	"testing"
	"tutorials/parallel/memo"
	"tutorials/parallel/memotest"
)

var httpGetBody = memotest.HTTPGetBody

func Test(t *testing.T) {
	m := memo.New(httpGetBody)
	memotest.Sequential(t, m)
}
