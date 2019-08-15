package models

import "fmt"

type LastFMUser struct {
	UserID     int
	LastfmName string
}

func (db *DB) GetLastFM(UserID int) (*LastFMUser, error) {
	lfm := new(LastFMUser)
	r := db.QueryRow("SELECT UserID, LastfmName FROM lastfm WHERE UserID = ?", UserID)
	err := r.Scan(&lfm.UserID, &lfm.LastfmName)
	if err != nil {
		return nil, err
	}

	return lfm, nil
}

func (db *DB) AddLastFM(UserID int, LastfmName string) error {
	stmt, err := db.Prepare("INSERT INTO lastfm (UserID, LastfmName) VALUES (?, ?)")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = stmt.Exec(UserID, LastfmName)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateLastFM(UserID int, LastfmName string) error {
	stmt, err := db.Prepare("UPDATE lastfm SET LastfmName = ? WHERE UserID = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(LastfmName, UserID)
	if err != nil {
		return err
	}

	return nil
}
