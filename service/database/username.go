package database

func (db *appdbimpl) DoLogin(username string) error {
	_, err := db.c.Exec(`INSERT INTO users (username) VALUES (?)`, username)

	if err != nil {
		return err
	}

	return err
}

func (db *appdbimpl) FindUserId(username string) (error, int) {
	var id int
	err := db.c.QueryRow(`SELECT idUser FROM users WHERE username = ?`, username).Scan(&id)

	if err != nil {
		return err, 0
	}

	return nil, id

}

func (db *appdbimpl) CheckUsername(username string) (error, int) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, username).Scan(&count)

	if err != nil {
		return err, -1
	}

	return nil, count
}

func (db *appdbimpl) SetUsername(id int, username string) error {

	var err error

	_, err = db.c.Exec(`UPDATE users SET username = ? WHERE idUser = ?`, id, username)

	return err

}
