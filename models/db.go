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
	AddTurnip(Price int, UserID int, Name string, Date string) error
	GetTodaysPrices(Date string) ([]*Turnip, error)
	GetTerms() ([]*Term, error)
	GetTerm(Name string) (*Term, error)
	AddTerm(Name string) error
	DeleteTerm(Name string) error
	AddCounter(Name string) error
	CountForUser(TermID int, UserID string) error
	GetForUsers(TermID int) ([]*TermUser, error)
	StoreChatLogID(ChatID int64) error
	AllUsers() ([]*Helth, error)
	GetHelth(ID string) (*Helth, error)
	AddHelth(Sender string, SenderID int) error
	RemoveHelth(SenderID int) error
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

func (db *DB) StoreChatLogID(ChatID int64) error {
	stmt, err := db.Prepare("INSERT INTO log (`ChatID`) VALUES (?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(ChatID)
	if err != nil {
		return err
	}

	return err
}
