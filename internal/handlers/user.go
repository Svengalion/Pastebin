// internal/handlers/user.go
package handlers

import (
	"net/http"

	"github.com/Svengalion/Pastebin/internal/models"
	"github.com/Svengalion/Pastebin/internal/repos"
	"github.com/Svengalion/Pastebin/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repo      repos.UserRepos
	Validator *validator.Validate
	JWTSecret []byte
}

func NewUserHandler(repo repos.UserRepos, jwtSecret []byte) *UserHandler {
	return &UserHandler{
		Repo:      repo,
		Validator: validator.New(),
		JWTSecret: jwtSecret,
	}
}

// RegUser регистрирует нового пользователя
// @Summary Регистрация нового пользователя
// @Description Создаёт нового пользователя с уникальным логином и email
// @Tags users
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Данные пользователя"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} gin.H
// @Failure 409 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users/register [post]
func (h *UserHandler) RegUser(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash your password"})
		return
	}

	user := &models.User{
		Login:    req.Login,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := h.Repo.RegisterUser(user); err != nil {
		switch err {
		case repos.ErrUserEmailAlreadyExist:
			c.JSON(http.StatusConflict, gin.H{"error": "Email already taken"})
		case repos.ErrUserLoginAlreadyExist:
			c.JSON(http.StatusConflict, gin.H{"error": "Login already taken"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
		return
	}

	resp := models.RegisterResponse{
		Id:        user.ID,
		Login:     user.Login,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	c.JSON(http.StatusOK, resp)
}

// LoginUser аутентифицирует пользователя и возвращает JWT токен
// @Summary Аутентификация пользователя
// @Description Аутентифицирует пользователя по email и паролю и возвращает JWT токен
// @Tags users
// @Accept json
// @Produce json
// @Param user body LoginRequest true "Данные для аутентификации"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users/auth [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := h.Validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user, err := h.Repo.GetUserByEmail(req.Email)
	if err != nil {
		if err == repos.ErrUserNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate user"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateJWT(user, h.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	resp := models.LoginResponse{
		Token: token,
	}

	c.JSON(http.StatusOK, resp)
}
