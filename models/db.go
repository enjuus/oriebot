package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Datastore interface {
	AllQuotes() ([]*Quote, error)
	AddQuote(Message string, Sender string, SenderFirstName string, SenderLastName string) error
	GetQuote(ID string) (*Quote, error)
	CountQuotes() int
}

type DB struct {
	*sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}