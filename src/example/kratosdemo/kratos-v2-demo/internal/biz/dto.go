package biz

type UserDto struct {
	UserName string `mapper:"Name"`
	Age      int
	Id       int
	Address  string
}
