package database

// insert a new photo given its id, its owner id and the date of upload
func (db *appdbimpl) InsertPhoto(idOwner int, date string, url string) (error, int) {
	res, err := db.c.Exec(`INSERT INTO images (idOwner, dateTime, url) VALUES (?, ?, ?)`, idOwner, date, url)

	if err != nil {
		return err, -1
	}

	var idImage int64
	idImage, err = res.LastInsertId()
	if err != nil {
		return err, -1
	}

	return err, int(idImage)
}

func (db *appdbimpl) SelectImgUrl(idImage int) (error, string) {
	var url string
	err := db.c.QueryRow(`SELECT url FROM images WHERE idImage = ?`, idImage).Scan(&url)

	if err != nil {
		return err, ""
	}

	return nil, url
}

func (db *appdbimpl) SelectImgDate(idImage int) (error, string) {
	var date string
	err := db.c.QueryRow(`SELECT dateTime FROM images WHERE idImage = ?`, idImage).Scan(&date)

	if err != nil {
		return err, ""
	}

	return nil, date
}

func (db *appdbimpl) FindImage(idImage int, idOwner int) (error, int) {
	var count int
	err := db.c.QueryRow(`SELECT count(*) FROM images WHERE idImage = ? AND idOwner = ?`, idImage, idOwner).Scan(&count)

	if err != nil {
		return err, -1
	}

	return err, count
}

func (db *appdbimpl) DeleteImage(idImage int, idOwner int) error {
	_, err := db.c.Exec(`DELETE FROM images WHERE idImage = ? AND idOwner = ?`, idImage, idOwner)

	return err
}

func (db *appdbimpl) CheckPhotoOwnership(idImage int, idOwner int) (error, int) {

	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM images WHERE idImage = ? AND idOwner = ?`, idImage, idOwner).Scan(&count)

	if err != nil {
		return err, -1
	}

	return err, count

}
