package repository

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type MapsRepo struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewMapsRepository(db *sql.DB, log *logrus.Logger) MapsRepo {
	return MapsRepo{
		db:  db,
		log: log,
	}
}
