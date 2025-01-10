// internal/repos/user_repos.go
package repos

import (
	"errors"

	"github.com/Svengalion/Pastebin/internal/models"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type UserRepos interface {
	RegisterUser(user *models.User) (err error)
	GetUserByLogin(login string) (user *models.User, err error)
	GetUserByEmail(email string) (user *models.User, err error)
}

type userRepos struct {
	db *gorm.DB
}

func NewUserRepos(db *gorm.DB) UserRepos {
	return &userRepos{db}
}

func (r *userRepos) RegisterUser(user *models.User) (err error) {
	if err := r.db.Create(user).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				if pgErr.ConstraintName == "users_email_key" {
					return ErrUserEmailAlreadyExist
				}
				if pgErr.ConstraintName == "users_login_key" {
					return ErrUserLoginAlreadyExist
				}
			}
		}
		return err
	}
	return nil
}

func (r *userRepos) GetUserByLogin(login string) (user *models.User, err error) {
	if err := r.db.Where("login = ?", login).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepos) GetUserByEmail(email string) (user *models.User, err error) {
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
