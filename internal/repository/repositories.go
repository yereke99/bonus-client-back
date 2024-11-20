package repository

import (
	"database/sql"
)

type IAuthRepository interface {
}

type Repositories struct {
	AuthRepository IAuthRepository
}

func NewRepository(db *sql.DB) *Repositories {

	return &Repositories{
		AuthRepository: NewAuthRepository(db),
	}
}
