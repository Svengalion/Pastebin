// internal/handlers/user.go
package handlers

import (
	"log"
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

func (h *UserHandler) RegUser(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Ошибка привязки JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		log.Printf("Ошибка валидации: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка хэширования пароля: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось хэшировать ваш пароль"})
		return
	}

	user := &models.User{
		Login:    req.Login,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := h.Repo.RegisterUser(user); err != nil {
		log.Printf("Ошибка регистрации пользователя: %v", err)
		switch err {
		case repos.ErrUserEmailAlreadyExist:
			c.JSON(http.StatusConflict, gin.H{"error": "Email уже используется"})
		case repos.ErrUserLoginAlreadyExist:
			c.JSON(http.StatusConflict, gin.H{"error": "Логин уже используется"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать пользователя"})
		}
		return
	}

	resp := models.RegisterResponse{
		Id:        user.ID,
		Login:     user.Login,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	c.JSON(http.StatusCreated, resp) // Изменено на StatusCreated
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
