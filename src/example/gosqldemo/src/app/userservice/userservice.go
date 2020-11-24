package userservice

import (
	"fmt"
	"gosqldemo/src/domain"
	"gosqldemo/src/repository"
)

type UserAppService struct {
}

var r = &repository.UserRepository{}

func NewUserService() *UserAppService {
	r.New()
	return &UserAppService{}
}

func (us *UserAppService) GetAllUsers() []domain.User {
	// 这里做业务逻辑，判断，验证等
	return r.GetUserAll()
}

func (us *UserAppService) GetUserById(id int) domain.User {
	if id < 1 {
		panic("用户不存在")
	}
	return r.GetUser(id)
}

func (us *UserAppService) CreateUser(user domain.User) error {
	if len(user.Name) == 0 {
		panic("用户姓名必填")
	}
	// if len(*user.Email) == 0 {
	// 	panic("email必填")
	// }
	return r.CreateUser(user)
}

func (us *UserAppService) DeleteUserById(userID int) (bool, error) {
	if userID < 1 {
		return false, fmt.Errorf("参数错误")
	}
	return r.DeleteUser(userID)
}
