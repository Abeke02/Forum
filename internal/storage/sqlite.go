package storage

import (
	"database/sql"
	"fmt"
)

const reactionTable = `CREATE TABLE IF NOT EXISTS reaction (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	id_object INT,
	action TEXT,
	creator TEXT,
	object TEXT,
	FOREIGN KEY (creator) REFERENCES user(username)
);`

const userTable = `CREATE TABLE IF NOT EXISTS user (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT UNIQUE,
	username TEXT UNIQUE,
	hashPassword TEXT,
	session_token TEXT DEFAULT NULL,
	expiresAt DATETIME DEFAULT NULL
);`

const postTable = `CREATE TABLE IF NOT EXISTS post (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	creator TEXT,
	title TEXT,
	description TEXT,
	likes INT,
	dislikes INT,
	created_at DATE DEFAULT (datetime('now','localtime')),
	FOREIGN KEY (creator) REFERENCES user(username)
);`

const categoriesTable = `CREATE TABLE IF NOT EXISTS categories (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	tag TEXT,
	id_post INT
);`

const commentTable = `CREATE TABLE IF NOT EXISTS comment (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	id_post INT,
	creator TEXT,
	comment TEXT,
	likes INT,
	dislikes INT,
	created_at DATE DEFAULT (datetime('now','localtime')),
	FOREIGN KEY (creator) REFERENCES user(username)
);`

var tables = []string{userTable, postTable, categoriesTable, commentTable, reactionTable}

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "forum.db?_foreign_keys=1")
	if err != nil {
		return nil, fmt.Errorf("storage: init db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("storage: init db: %w", err)
	}
	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			return nil, fmt.Errorf("storage: create tables: %w", err)
		}
	}
	return db, nil
}
