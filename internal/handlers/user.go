package handlers

import (
	"errors"
	"net/http"

	"github.com/Svengalion/Pastebin/internal/models"
	"github.com/Svengalion/Pastebin/internal/repos"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Repo repos.UserRepos
}

func NewUserHandler(repo repos.UserRepos) *UserHandler {
	return &UserHandler{Repo: repo}
}

// RegUser godoc
// @Summary Регистрация нового пользователя
// @Description Создаёт нового пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param paste body models.User true "Данные пользователя"
// @Success 201 {object} models.User
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users/registration [post]
func (h *UserHandler) RegUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.Repo.RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// AuthUser godoc
// @Summary Авторизация пользователя
// @Description Пока что возвращает пользователя по логину
// @Tags pastes
// @Accept json
// @Produce json
// @Param login path string true "Логин пользователя"
// @Success 200 {object} models.User
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users/auth/{login, password} [get]
func (h *UserHandler) AuthUser(c *gin.Context) {
	login, password := c.Param("login"), c.Param("password")

	user, err := h.Repo.AuthUser(login, password)
	if err != nil {
		if errors.Is(err, repos.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
