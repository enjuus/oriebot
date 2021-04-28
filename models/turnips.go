package models

import "fmt"

type Turnip struct {
	ID     int
	UserID int
	Name   string
	Price  int
	Date   string
}

func (db *DB) GetTodaysPrices(Date string) ([]*Turnip, error) {
	query := fmt.Sprintf("SELECT turnips.* FROM turnips INNER JOIN (" +
		"SELECT UserID, MAX(ID) AS maxID FROM turnips GROUP BY UserID) AS tmp " +
		"ON turnips.UserID = tmp.UserID AND turnips.ID = tmp.maxID WHERE `Date` = '%s'", Date)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prices := make([]*Turnip, 0)
	for rows.Next() {
		t := new(Turnip)
		err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.Price, &t.Date)
		if err != nil {
			return nil, err
		}
		prices = append(prices, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prices, nil
}

func (db *DB) AddTurnip(Price int, UserID int, Name string, Date string) error {
	stmt, err := db.Prepare("INSERT INTO turnips (UserID, `Name`, Price, `Date`) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(UserID, Name, Price, Date)
	if err != nil {
		return err
	}

	return err
}
