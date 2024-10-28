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

func (h *PasteHandler) CreatePaste(c *gin.Context) {
	var req models.CreatePasteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paste := models.Paste{
		Title:   req.Title,
		Content: req.Content,
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

func (h *PasteHandler) GetPaste(c *gin.Context) {
	hash := c.Param("hash")
	if len(hash) != utils.HashSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect hash length"})
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

func (h *PasteHandler) GetAllPastes(c *gin.Context) {
	pastes, err := h.Repo.GetAllPastes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pastes"})
		return
	}
	c.JSON(http.StatusOK, pastes)
}
