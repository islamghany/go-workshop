package data

import (
	"database/sql"
	"errors"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrEditConflict   = errors.New("edit conflict")
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Users  UserModel
	Tokens TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Tokens: TokenModel{DB: db},
		Users:  UserModel{DB: db}, // Initialize a new UserModel instance.
	}
}
