package mypack

//go:generate go run gen.go arg1 arg2 $GOVERSION
//最后的$GOVERSION表示取得环境变量，对应go env GOVERSION 的值，如果要显示源字符串（即$GOVERSION）则需要改成：go:generate go run gen.go arg1 arg2 ${DOLLAR}GOVERSION
//DOLLAR的解释，详见：https://eli.thegreenplace.net/2021/a-comprehensive-guide-to-go-generate/

func PackFunc() string {
	return "insourcegenerator/mypack.PackFunc"
}
