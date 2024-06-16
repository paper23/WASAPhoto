package database

import (
	"database/sql"
)

// return a list of images as requested from GetMyStream endpoint
func (db *appdbimpl) GetStream(idUser int) (error, []int, []int, []string) {

	var imagesIdList []int
	var usersIdList []int
	var usernamesList []string
	var rows *sql.Rows
	var err error

	query := `SELECT images.idImage, users.idUser, users.username
	FROM images
	JOIN users ON images.idOwner = users.idUser
	JOIN follows ON users.idUser = follows.idFollowed
	WHERE follows.idFollower = ?
	ORDER BY images.dateTime DESC
	LIMIT 100;`
	rows, err = db.c.Query(query, idUser)

	if err != nil {
		return err, nil, nil, nil
	}

	for rows.Next() {
		var idImage int
		var idUser int
		var username string
		err = rows.Scan(&idImage, &idUser, &username)
		if err != nil {
			return err, nil, nil, nil
		}
		imagesIdList = append(imagesIdList, idImage)
		usersIdList = append(usersIdList, idUser)
		usernamesList = append(usernamesList, username)
	}

	err = rows.Err()

	if err != nil {
		return err, nil, nil, nil
	}

	return err, imagesIdList, usersIdList, usernamesList

}

// return a list of images as requested from GetMyStream endpoint
func (db *appdbimpl) GetStream2(idUser int) (error, []Stream2) {

	var userStream []Stream2

	query := `SELECT users.username, images.*
	FROM images
	JOIN users ON images.idOwner = users.idUser
	JOIN follows ON users.idUser = follows.idFollowed
	WHERE follows.idFollower = ?
	ORDER BY images.dateTime DESC
	LIMIT 100;`

	rows, err := db.c.Query(query, idUser)
	if err != nil {
		return err, nil
	}

	for rows.Next() {

		var tmp Stream2

		err := rows.Scan(&tmp.Username, &tmp.Image.IdImage, &tmp.Image.IdOwner, &tmp.Image.DateTime, &tmp.Image.File)
		if err != nil {
			return err, nil
		}

		userStream = append(userStream, tmp)
	}

	err = rows.Err()
	if err != nil {
		return err, nil
	}

	return nil, userStream

}
