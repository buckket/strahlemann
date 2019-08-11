package database

import (
	"database/sql"
	"fmt"
	"github.com/buckket/strahlemann/utils"
	_ "github.com/mattn/go-sqlite3"
)

type StrahlemannDatabase struct {
	*sql.DB
}

type Post struct {
	ID        int64
	LastTweet int64
	Content   string
	Position  int
	Complete  bool
}

type Status struct {
	AllEntries      int
	DoneEntries     int
	LastTweet       int64
	CurrentPosition int
	CurrentLenght   int
	NextTweet       string
}

func (db *StrahlemannDatabase) CreateSchema() error {
	schema := `
CREATE TABLE IF NOT EXISTS blog (
	id			INTEGER PRIMARY KEY AUTOINCREMENT,
	last_tweet	INTEGER,
	content		TEXT,
	position	INTEGER,
	complete	BOOLEAN
);
`
	_, err := db.Exec(schema)
	return err
}

func (db *StrahlemannDatabase) InsertPost(post *Post) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insertPost := `INSERT INTO blog (last_tweet, content, position, complete) VALUES (?, ?, ?, ?)`
	_, err = tx.Exec(insertPost, post.LastTweet, post.Content, post.Position, post.Complete)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *StrahlemannDatabase) UpdatePost(post *Post) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updatePost := `UPDATE blog SET last_tweet = ?, position = ?, complete = ? WHERE id = ?`
	_, err = tx.Exec(updatePost, post.LastTweet, post.Position, post.Complete, post.ID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *StrahlemannDatabase) GetNextPost() (*Post, error) {
	post := Post{}

	qry := `SELECT id, last_tweet, content, position FROM blog WHERE complete = FALSE ORDER BY id ASC`
	err := db.QueryRow(qry).Scan(&post.ID, &post.LastTweet, &post.Content, &post.Position)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (db *StrahlemannDatabase) GetStatus() (*Status, error) {
	status := Status{}

	qry := `SELECT COUNT(*) FROM blog`
	err := db.QueryRow(qry).Scan(&status.AllEntries)
	if err != nil {
		return nil, err
	}

	qry = `SELECT COUNT(*) FROM blog WHERE complete = TRUE`
	err = db.QueryRow(qry).Scan(&status.DoneEntries)
	if err != nil {
		return nil, err
	}

	if status.AllEntries == status.DoneEntries {
		return &status, nil
	}

	post, err := db.GetNextPost()
	if err != nil {
		return nil, err
	}

	status.LastTweet = post.LastTweet
	status.CurrentPosition = post.Position
	status.CurrentLenght = len(post.Content)
	status.NextTweet, _ = utils.ExtractTweet(post.Content[post.Position:])

	return &status, nil
}

func New(target string) (*StrahlemannDatabase, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?&_fk=true", target))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &StrahlemannDatabase{db}, nil
}
