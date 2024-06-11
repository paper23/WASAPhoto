package database

// insert a new photo given its id, its owner id and the date of upload
func (db *appdbimpl) CommentPhoto(idUserWriter int, idImage int, text string) (error, int) {

	res, err := db.c.Exec(`INSERT INTO comments (idUserWriter, idImageCommented, text) VALUES (?, ?, ?)`, idUserWriter, idImage, text)

	if err != nil {
		return err, -1
	}

	var idComment int64
	idComment, err = res.LastInsertId()
	if err != nil {
		return err, -1
	}

	return err, int(idComment)
}

func (db *appdbimpl) FindComment(idComment int) (error, int) {

	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM comments WHERE idComment = ?`, idComment).Scan(&count)

	if err != nil {
		return err, -1
	}

	return err, count

}

func (db *appdbimpl) SelectCommentText(idComment int) (error, string) {
	var text string
	err := db.c.QueryRow(`SELECT text FROM comments WHERE idComment = ?`, idComment).Scan(&text)

	if err != nil {
		return err, ""
	}

	return err, text
}

func (db *appdbimpl) CheckCommentOwnership(idComment int, idUserWriter int) (error, bool) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM comments WHERE idComment = ? AND idUserWriter = ?`, idComment, idUserWriter).Scan(&count)

	if err != nil {
		return err, false
	}

	if count > 0 {
		return err, true
	} else {
		return err, false
	}
}

func (db *appdbimpl) UncommentPhoto(idComment int) error {
	_, err := db.c.Exec(`DELETE FROM comments WHERE idComment = ?`, idComment)

	return err
}

func (db *appdbimpl) CountComments(idImage int) (error, int) {

	var commentsCounter int
	commentsCounter = 0

	err := db.c.QueryRow(`SELECT COUNT (*) FROM comments WHERE idImageCommented = ?`, idImage).Scan(&commentsCounter)
	if err != nil {
		return err, -1
	}

	return nil, commentsCounter
}

func (db *appdbimpl) GetComments(idImage int) (error, []CommentWrapper) {

	var comments []CommentWrapper

	rows, err := db.c.Query(`SELECT * FROM comments WHERE idImageCommented = ?`, idImage)
	if err != nil {
		return err, nil
	}

	for rows.Next() {
		var comment CommentWrapper

		err := rows.Scan(&comment.CommentData.IdComment, &comment.CommentData.IdUserWriter, &comment.CommentData.IdImage, &comment.CommentData.Text)
		if err != nil {
			return err, nil
		}

		err = db.c.QueryRow(`SELECT username FROM users WHERE idUser = ?`, comment.CommentData.IdUserWriter).Scan(&comment.Username)
		if err != nil {
			return err, nil
		}

		comments = append(comments, comment)
	}

	err = rows.Err()

	if err != nil {
		return err, nil
	}

	return nil, comments

}
