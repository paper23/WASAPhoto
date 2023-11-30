package database

//insert a new user given its username (set biography = "")
func (db *appdbimpl) DoLogin(username string) error {
	_, err := db.c.Exec(`INSERT INTO users (username, biography) VALUES (?, "")`, username)

	if err != nil {
		return err
	}

	return err
}

//find user's id given its username, returns the id
func (db *appdbimpl) FindUserId(username string) (error, int) {
	var id int
	err := db.c.QueryRow(`SELECT idUser FROM users WHERE username = ?`, username).Scan(&id)

	if err != nil {
		return err, 0
	}

	return nil, id

}

//count how many users have the same username given the username, returns the counter
func (db *appdbimpl) CheckUsername(username string) (error, int) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, username).Scan(&count)

	if err != nil {
		return err, -1
	}

	return nil, count
}

//update a user's username given its id and the new username
func (db *appdbimpl) SetUsername(id int, username string) error {

	var err error

	_, err = db.c.Exec(`UPDATE users SET username = ? WHERE idUser = ?`, username, id)

	return err

}

//returns all the informations about a user given its id
func (db *appdbimpl) SelectUser(id int) (error, int, string, string) {
	var err error

	var idUser int
	var username string
	var bio string

	err = db.c.QueryRow(`SELECT * FROM users WHERE idUser = ?`, id).Scan(&idUser, &username, &bio)

	if err != nil {
		return err, idUser, username, bio
	}

	return err, idUser, username, bio
}