package database

//insert a new photo given its id, its owner id and the date of upload
func (db *appdbimpl) InsertPhoto(idOwner int, date string, url string) (error, int64) {
	res, err := db.c.Exec(`INSERT INTO images (idOwner, dateTime, url) VALUES (?, ?, ?)`, idOwner, date, url)

	if err != nil {
		return err, -1
	}

	var idImage int64
	idImage, err = res.LastInsertId()
	if err != nil {
		return err, -1
	}

	return err, idImage
}