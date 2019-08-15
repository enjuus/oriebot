package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Datastore interface {
	AllQuotes() ([]*Quote, error)
	AddQuote(Message string, Sender string, SenderFirstName string, SenderLastName string, SenderID int) error
	GetQuote(ID string) (*Quote, error)
	CountQuotes() int
	GetLastFM(UserID int) (*LastFMUser, error)
	AddLastFM(UserID int, LastfmName string) error
	UpdateLastFM(UserID int, LastfmName string) error
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
