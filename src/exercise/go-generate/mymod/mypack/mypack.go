package mypack

//go:generate mygenerate arg1 "multiword arg"
func PackFunc() string {
	return "mymod/mypack.PackFunc"
}
