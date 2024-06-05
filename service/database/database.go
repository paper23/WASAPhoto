/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// user
	DoLogin(username string) error
	FindUserId(username string) (error, int)
	CheckUsername(username string) (error, int)
	SetUsername(id int, username string) error
	SelectUser(id int) (error, int, string, string)
	FindUserById(id int) (error, int)
	FollowUser(idUser int, idUserToFollow int) error
	BanUser(idUser int, idUserToBan int) error
	CheckBan(idUser int, idUserToCheck int) (error, int)
	SbanUser(idUser int, idUserToSban int) error
	CheckFollowing(idUser int, idUserToUnfollow int) (error, int)
	UnfollowUser(idUser int, idUserToUnfolow int) error
	FindUserBio(idUser int) (error, string)
	FindUsername(idUser int) (error, string)
	CountFollowing(idUser int) (error, int)
	CountFollowers(idUser int) (error, int)

	// image
	InsertPhoto(idOwner int, date string, file []byte) (error, int)
	SelectImgDate(idImage int) (error, string)
	FindImage(idImage int, idOwner int) (error, int)
	DeleteImage(idImage int, idOwner int) error
	CommentPhoto(idUserWriter int, idImage int, text string) (error, int)
	FindComment(idComment int) (error, int)
	SelectCommentText(idComment int) (error, string)
	CheckCommentOwnership(idComment int, idUserWriter int) (error, bool)
	UncommentPhoto(idComment int) error
	LikePhoto(idLiker int, idImage int) error
	CheckLike(idLiker int, idImage int) (error, int)
	UnlikePhoto(idLiker int, idImage int) error
	CheckPhotoOwnership(idImage int, idOwner int) (error, int)
	GetUserImagesId(idOwner int) (error, []int)
	GetUserImagesFile(idImages []int) (error, []byte)

	// stream
	GetStream(idUser int) (error, []int, []int, []string)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	_, errPramga := db.Exec(`PRAGMA foreign_keys= ON`)
	if errPramga != nil {
		return nil, fmt.Errorf("error setting pragmas: %w", errPramga)
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = createDatabase(db)
		if err != nil {
			return nil, fmt.Errorf("error creating database: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

func createDatabase(db *sql.DB) error {
	tables := [6]string{
		`CREATE TABLE IF NOT EXISTS users (
			idUser INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(16) NOT NULL UNIQUE,
			biography VARCHAR(100)
		);`,
		`CREATE TABLE IF NOT EXISTS images (
			idImage INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			idOwner INTEGER NOT NULL,
			dateTime DATETIME NOT NULL,
			file BLOB NOT NULL,
			FOREIGN KEY(idOwner) REFERENCES users (idUser) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS likes (
			idLiker INTEGER NOT NULL,
			idImageLiked INTEGER NOT NULL,
			PRIMARY KEY (idLiker, idImageLiked),
			FOREIGN KEY(idLiker) REFERENCES users (idUser) ON DELETE CASCADE,
			FOREIGN KEY(idImageLiked) REFERENCES images (idImage) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS comments (
			idComment INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			idUserWriter INTEGER NOT NULL,
			idImageCommented INTEGER NOT NULL,
			text VARCHAR(200) NOT NULL,
			FOREIGN KEY(idUserWriter) REFERENCES users (idUser) ON DELETE CASCADE,
			FOREIGN KEY(idImageCommented) REFERENCES images (idImage) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS follows(
			idFollower INTEGER NOT NULL,
			idFollowed INTEGER NOT NULL,
			PRIMARY KEY (idFollower,idFollowed),
			FOREIGN KEY(idFollower) REFERENCES users (idUser) ON DELETE CASCADE,
			FOREIGN KEY(idFollowed) REFERENCES users (idUser) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS bans(
			idUser INTEGER NOT NULL,
			idBanned INTEGER NOT NULL,
			PRIMARY KEY (idUser,idBanned),
			FOREIGN KEY(idUser) REFERENCES users (idUser) ON DELETE CASCADE,
			FOREIGN KEY(idBanned) REFERENCES users (idUser) ON DELETE CASCADE
			);`,
	}

	for i := 0; i < len(tables); i++ {
		sqlStmt := tables[i]
		_, err := db.Exec(sqlStmt)

		if err != nil {
			return err
		}
	}
	return nil

}
