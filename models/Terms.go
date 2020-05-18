package models

import (
	"database/sql"
	"fmt"
	"strings"
)

type Term struct {
	ID    int
	Name  string
	Count int
}

type TermUser struct {
	TermID int
	UserID string
	Count  int
}

func (db *DB) GetTerms() ([]*Term, error) {
	query := fmt.Sprintf("SELECT * FROM terms")
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	terms := make([]*Term, 0)
	for rows.Next() {
		t := new(Term)
		err := rows.Scan(&t.ID, &t.Name, &t.Count)
		if err != nil {
			return nil, err
		}
		terms = append(terms, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return terms, nil
}

func (db *DB) AddTerm(Name string) error {
	stmt, err := db.Prepare("INSERT INTO terms (`Name`, `Count`) VALUES (?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(strings.TrimSpace(Name), 0)
	if err != nil {
		return err
	}

	return err
}

func (db *DB) DeleteTerm(Name string) error {
	stmt, err := db.Prepare("DELETE FROM terms WHERE Name = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(Name)
	if err != nil {
		return err
	}

	return err
}

func (db *DB) AddCounter(Name string) error {
	stmt, err := db.Prepare("UPDATE terms SET `Count` = `Count` + 1 WHERE Name = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(Name)
	if err != nil {
		return err
	}

	return err
}

func (db *DB) GetTerm(Name string) (*Term, error) {
	t := new(Term)
	r := db.QueryRow("SELECT * FROM terms WHERE Name = ?", Name)
	err := r.Scan(&t.ID, &t.Name, &t.Count)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (db *DB) CountForUser(TermID int, UserID string) error {
	tu := new(TermUser)
	r := db.QueryRow("SELECT * FROM termuser WHERE TermID = ? AND UserID = ?", TermID, UserID)
	err := r.Scan(&tu.TermID, &tu.UserID, &tu.Count)
	if err != nil && sql.ErrNoRows != nil {
		stmt, err := db.Prepare("INSERT INTO termuser (`TermID`, `UserID`, `Count`) VALUES (?, ?, 0)")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(TermID, UserID)
		if err != nil {
			return err
		}
	}
	stmt, err := db.Prepare("UPDATE termuser SET `Count` = `Count` + 1 WHERE TermID = ? AND UserID = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(TermID, UserID)
	if err != nil {
		return err
	}

	return err
}

func (db *DB) GetForUsers(TermID int) ([]*TermUser, error) {
	userTerms := make([]*TermUser, 0)

	rows, err := db.Query("SELECT * FROM termuser WHERE TermID = ?", TermID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		tu := new(TermUser)
		err := rows.Scan(&tu.TermID, &tu.UserID, &tu.Count)
		if err != nil {
			return nil, err
		}
		userTerms = append(userTerms, tu)
	}

	return userTerms, nil
}
