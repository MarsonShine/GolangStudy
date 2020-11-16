package contracts

import "gormdemo/src/models"

// TODO 未实现
type UserService interface {
	Create(user *models.User)
	Get(id uint) *models.User
	Update(user *models.User) bool
	Delete(id uint) bool
}
