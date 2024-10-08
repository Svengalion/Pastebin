package repos

import (
	"errors"

	"github.com/Svengalion/Pastebin/internal/models"
	"gorm.io/gorm"
)

type UserRepos interface {
	RegisterUser(user *models.User) (err error)
	AuthUser(login string, password string) (paste *models.User, err error)
}

type userRepos struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) UserRepos {
	return &userRepos{db}
}

func (r *userRepos) RegisterUser(user *models.User) (err error) {
	if err := r.db.Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			existingUser := &models.User{}
			if err := r.db.Where("login = ?", user.Login).First(existingUser).Error; err == nil {
				return ErrUserLoginAlreadyExist
			}
			if err := r.db.Where("email = ?", user.Email).First(existingUser).Error; err == nil {
				return ErrUserEmailAlreadyExist
			}
		}
		return err
	}
	return nil
}

func (r *userRepos) AuthUser(login string, password string) (user *models.User, err error) {
	if err := r.db.First(&login, "login = ?", login).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
