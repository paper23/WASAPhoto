package database

// insert a new user given its username (set biography = "")
func (db *appdbimpl) DoLogin(username string) error {
	_, err := db.c.Exec(`INSERT INTO users (username, biography) VALUES (?, "")`, username)

	if err != nil {
		return err
	}

	return err
}

// find user's id given its username, returns the id
func (db *appdbimpl) FindUserId(username string) (error, int) {
	var id int
	err := db.c.QueryRow(`SELECT idUser FROM users WHERE username = ?`, username).Scan(&id)

	if err != nil {
		return err, 0
	}

	return nil, id

}

// count how many users have the same username given the username, returns the counter
func (db *appdbimpl) CheckUsername(username string) (error, int) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, username).Scan(&count)

	if err != nil {
		return err, -1
	}

	return nil, count
}

// update a user's username given its id and the new username
func (db *appdbimpl) SetUsername(id int, username string) error {

	var err error

	_, err = db.c.Exec(`UPDATE users SET username = ? WHERE idUser = ?`, username, id)

	return err

}

// returns all the informations about a user given its id
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

func (db *appdbimpl) FindUserById(id int) (error, int) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE idUser = ?`, id).Scan(&count)

	if err != nil {
		return err, -1
	}

	return err, count
}

func (db *appdbimpl) FollowUser(idUser int, idUserToFollow int) error {
	_, err := db.c.Exec(`INSERT INTO follows (idFollower, idFollowed) VALUES (?, ?)`, idUser, idUserToFollow)

	return err
}

func (db *appdbimpl) CheckBan(idUser int, idUserToCheck int) (error, int) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM bans WHERE idUser = ? AND idBanned = ?`, idUser, idUserToCheck).Scan(&count)

	if err != nil {
		return err, -1
	}

	return err, count
}

func (db *appdbimpl) SbanUser(idUser int, idUserToSban int) error {
	_, err := db.c.Exec(`DELETE FROM bans WHERE idUser = ? AND idBanned = ?`, idUser, idUserToSban)

	return err
}

func (db *appdbimpl) CheckFollowing(idUser int, idUserToCheck int) (error, int) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM follows WHERE idFollower = ? AND idFollowed = ?`, idUser, idUserToCheck).Scan(&count)

	if err != nil {
		return err, -1
	}

	return err, count
}

func (db *appdbimpl) UnfollowUser(idUser int, idUserToUnfollow int) error {
	_, err := db.c.Exec(`DELETE FROM follows WHERE idFollower = ? AND idFollowed = ?`, idUser, idUserToUnfollow)

	return err
}

func (db *appdbimpl) BanUser(idUser int, idUserToBan int) error {

	var err error
	err = db.UnfollowUser(idUser, idUserToBan)
	if err != nil {
		return err
	}

	err = db.UnfollowUser(idUserToBan, idUser)
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`INSERT INTO bans (idUser, idBanned) VALUES (?, ?)`, idUser, idUserToBan)

	return err
}

func (db *appdbimpl) FindUserBio(idUser int) (error, string) {
	var bio string
	err := db.c.QueryRow(`SELECT biography FROM users WHERE idUser = ?`, idUser).Scan(&bio)

	if err != nil {
		return err, ""
	}

	return err, bio
}

func (db *appdbimpl) FindUsername(id int) (error, string) {
	var username string
	err := db.c.QueryRow(`SELECT username FROM users WHERE idUser = ?`, id).Scan(&username)

	if err != nil {
		return err, ""
	}

	return err, username
}

func (db *appdbimpl) CountFollowing(idUser int) (error, int) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM follows WHERE idFollower = ?`, idUser).Scan(&count)

	if err != nil {
		return err, -1
	}

	return err, count
}

func (db *appdbimpl) CountFollowers(idUser int) (error, int) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM follows WHERE idFollowed = ?`, idUser).Scan(&count)

	if err != nil {
		return err, -1
	}

	return err, count
}

func (db *appdbimpl) GetBanned(idUser int) (error, []User) {

	var bannedList []User

	rows, err := db.c.Query(`SELECT idBanned FROM bans WHERE idUser = ?`, idUser)
	if err != nil {
		return err, nil
	}

	for rows.Next() {
		var idBanned User

		err := rows.Scan(&idBanned.IdUser)
		if err != nil {
			return err, nil
		}

		bannedList = append(bannedList, idBanned)
	}
	err = rows.Err()

	if err != nil {
		return err, nil
	}

	for index := range bannedList {

		err = db.c.QueryRow(`SELECT username FROM users WHERE idUser = ?`, bannedList[index].IdUser).Scan(&bannedList[index].Username)
		if err != nil {
			return err, nil
		}
	}

	return nil, bannedList

}
