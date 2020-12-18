// go build -buildmode=plugin  plugin 目前只支持 linux 和 mac 操作系统 不支持 windows
package plugin1

func Invoke() string {
	str := "this is plugin 1"
	return str
}
