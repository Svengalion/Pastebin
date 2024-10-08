package repos

import (
	"errors"

	"github.com/Svengalion/Pastebin/internal/models"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type PasteRepos interface {
	CreatePaste(paste *models.Paste) (err error)
	GetPaste(hash string) (paste *models.Paste, err error)
}

type pasteRepos struct {
	db *gorm.DB
}

func NewPasteRepos(db *gorm.DB) PasteRepos {
	return &pasteRepos{db}
}

func (r *pasteRepos) CreatePaste(paste *models.Paste) (err error) {
	if err := r.db.Create(paste).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrHashAlreadyExists
		}
		return err
	}
	return nil
}

func (r *pasteRepos) GetPaste(hash string) (paste *models.Paste, err error) {
	if err := r.db.First(&paste, "hash = ?", hash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPasteNotFound
		}
		return nil, err
	}
	return paste, nil
}
