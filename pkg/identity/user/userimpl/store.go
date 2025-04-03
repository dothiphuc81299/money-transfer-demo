package userimpl

import (
	"money-transfer-demo/pkg/infra/storage/postgres"

	"gorm.io/gorm"
)

type store struct {
	db *gorm.DB
}

func NewStore(database postgres.DBConnector) *store {
	return &store{db: database.GetDB()}
}

func (s *store) createUser(