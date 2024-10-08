package handlers

import (
	"errors"
	"net/http"

	"github.com/Svengalion/Pastebin/internal/models"
	"github.com/Svengalion/Pastebin/internal/repos"
	"github.com/Svengalion/Pastebin/internal/utils"
	"github.com/gin-gonic/gin"
)

type PasteHandler struct {
	Repo repos.PasteRepos
}

func NewPasteHandler(repo repos.PasteRepos) *PasteHandler {
	return &PasteHandler{Repo: repo}
}

// CreatePaste godoc
// @Summary Создание новой пасты
// @Description Создаёт новую пасту с уникальным хэшем
// @Tags pastes
// @Accept json
// @Produce json
// @Param paste body models.CreatePasteRequest true "Паста для создания"
// @Success 201 {object} models.Paste
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /pastes/new_paste [post]
func (h *PasteHandler) CreatePaste(c *gin.Context) {
	var paste models.Paste
	if err := c.ShouldBindJSON(&paste); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := 0; i < 10; i++ {
		hash, err := utils.GenerateHash()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		paste.Hash = hash

		err = h.Repo.CreatePaste(&paste)
		if err != nil {
			if errors.Is(err, repos.ErrHashAlreadyExists) {
				continue
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, paste)
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate hash try again"})
}

// GetPaste godoc
// @Summary Получение пасты по хэшу
// @Description Получает пасту по её уникальному хэшу
// @Tags pastes
// @Accept json
// @Produce json
// @Param hash path string true "Хэш пасты"
// @Success 200 {object} models.Paste
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /pastes/{hash} [get]
func (h *PasteHandler) GetPaste(c *gin.Context) {
	hash := c.Param("hash")
	if len(hash) != utils.HashSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorect hash length"})
		return
	}

	paste, err := h.Repo.GetPaste(hash)
	if err != nil {
		if errors.Is(err, repos.ErrPasteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Paste not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, paste)
}
