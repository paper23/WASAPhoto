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
	DoLogin(username string) error
	FindUserId(username string) (error, int)
	CheckUsername(username string) (error, int)
	SetUsername(id int, username string) error

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
	tables := [4]string{
		`CREATE TABLE IF NOT EXISTS users (
			idUser INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(16) NOT NULL,
			biography VARCHAR(100)
		);`,
		`CREATE TABLE IF NOT EXISTS images (
			idImage INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			idOwner INTEGER NOT NULL,
			dateTime DATETIME NOT NULL,
			url VARCHAR(100) NOT NULL,
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
			idComment INTEGER NOT NULL,
			idUserWriter INTEGER NOT NULL,
			idImageCommented INTEGER NOT NULL,
			text VARCHAR(200) NOT NULL,
			PRIMARY KEY (idComment, idUserWriter, idImageCommented),
			FOREIGN KEY(idUserWriter) REFERENCES users (idUser) ON DELETE CASCADE,
			FOREIGN KEY(idImageCommented) REFERENCES images (idImage) ON DELETE CASCADE
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
