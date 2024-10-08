package repos

import "errors"

var (
	ErrPasteNotFound         = errors.New("paste not found")
	ErrHashAlreadyExists     = errors.New("hash already exists")
	ErrUserLoginAlreadyExist = errors.New("user with that login already exists")
	ErrUserEmailAlreadyExist = errors.New("user with that email already exists")
	ErrUserNotFound          = errors.New("user not found")
)
