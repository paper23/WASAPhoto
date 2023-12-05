package database

// insert a new photo given its id, its owner id and the date of upload
func (db *appdbimpl) LikePhoto(idLiker int, idImage int) error {
	_, err := db.c.Exec(`INSERT INTO likes (idLiker, idImageLiked) VALUES (?, ?)`, idLiker, idImage)

	return err
}

func (db *appdbimpl) CheckLike(idLiker int, idImage int) (error, int) {

	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM likes WHERE idLiker = ? AND idImageLiked = ?`, idLiker, idImage).Scan(&count)

	if err != nil {
		return err, -1
	}

	return err, count
}

func (db *appdbimpl) UnlikePhoto(idLiker int, idImage int) error {
	_, err := db.c.Exec(`DELETE FROM likes WHERE idLiker = ? AND idImageLiked = ?`, idLiker, idImage)

	return err
}
