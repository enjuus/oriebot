package models

type Helth struct {
	ID       int32
	Sender   string
	SenderID int
	Start    string
	Stop     string
}

func (db *DB) AllUsers() ([]*Helth, error) {
	rows, err := db.Query("SELECT * FROM helth")
	defer rows.Close()

	helths := make([]*Helth, 0)
	for rows.Next() {
		ht := new(Helth)
		err := rows.Scan(&ht.ID, &ht.Sender, &ht.SenderID, &ht.Start, &ht.Stop)
		if err != nil {
			return nil, err
		}
		helths = append(helths, ht)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return helths, nil
}

func (db *DB) GetHelth(ID string) (*Helth, error) {
	ht := new(Helth)
	r := db.QueryRow("SELECT ID, Sender, SenderID, Start, Stop FROM helth WHERE ID = ?", ID)
	err := r.Scan(&ht.ID, &ht.Sender, &ht.SenderID, &ht.Start, &ht.Stop)
	if err != nil {
		return nil, err
	}

	return ht, nil
}

func (db *DB) AddHelth(Sender string, SenderID int) error {
	stmt, err := db.Prepare("INSERT INTO helth (Sender, SenderID, Start, Stop) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(Sender, SenderID, "07:00AM", "06:30PM")
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) RemoveHelth(SenderID int) error {
	stmt, err := db.Prepare("DELETE FROM helth WHERE SenderID = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(SenderID)
	if err != nil {
		return err
	}

	return nil
}
