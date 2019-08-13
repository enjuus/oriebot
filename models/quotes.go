package models

type Quote struct {
	ID              int32
	Message         string
	Sender          string
	SenderFirstName string
	SenderLastName  string
	SenderID	int
}

func (db *DB) CountQuotes() int {
	row, err := db.Query("SELECT COUNT(*) as count FROM quotes")
	defer row.Close()

	if err != nil {
		return 0
	}
	var count int
	err = row.Scan(&count)

	return count

}

func (db *DB) AllQuotes() ([]*Quote, error) {
	rows, err := db.Query("SELECT * FROM quotes")
	defer rows.Close()

	quotes := make([]*Quote, 0)
	for rows.Next() {
		qt := new(Quote)
		err := rows.Scan(&qt.ID, &qt.Sender, &qt.SenderFirstName, &qt.SenderLastName, &qt.SenderID, &qt.Message)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, qt)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return quotes, nil
}

func (db *DB) GetQuote(ID string) (*Quote, error) {
	qt := new(Quote)
	r := db.QueryRow("SELECT ID, Message, Sender, SenderFirstName, SenderLastName, SenderID FROM quotes WHERE ID = ?", ID)
	err := r.Scan(&qt.ID, &qt.Message, &qt.Sender, &qt.SenderFirstName, &qt.SenderLastName, &qt.SenderID)
	if err != nil {
		return nil, err
	}

	return qt, nil
}

func (db *DB) AddQuote(Message string, Sender string, SenderFirstName string, SenderLastName string, SenderID int) error {
	stmt, err := db.Prepare("INSERT INTO quotes (Message, Sender, SenderFirstName, SenderLastName, SenderID) values (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(Message, Sender, SenderFirstName, SenderLastName, SenderID)
	if err != nil {
		return err
	}

	return nil
}
