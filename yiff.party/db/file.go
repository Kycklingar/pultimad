package db

import (
	"errors"
	"io"
	"log"
	"net/http"

	yp "github.com/kycklingar/pultimad/yiff.party/parser"
)

type File struct {
	FileURL string
	Creator string
	PostID  int
	TypeOf  int
}

var client *http.Client

func (f File) Download() (io.ReadCloser, error) {
	if client == nil {
		client = &http.Client{}
	}
	res, err := client.Get(f.FileURL)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		res.Body.Close()
		return nil, errors.New("not 200")
	}

	return res.Body, nil
}

func (db *DB) Tried(f File) error {
	_, err := db.Exec(
		"UPDATE files SET tries = tries + 1 WHERE post_id = $1 AND typeof = $2 AND file_url = $3",
		f.PostID,
		f.TypeOf,
		f.FileURL,
	)
	return err
}

func (db *DB) CheckFile(f File, sum string) error {
	_, err := db.Exec(
		"UPDATE files SET downloaded = CURRENT_TIMESTAMP, sha256 = $1 WHERE post_id = $2 AND typeof = $3 AND file_url = $4",
		sum,
		f.PostID,
		f.TypeOf,
		f.FileURL,
	)
	return err
}

func (db *DB) StoreSharedFile(sh yp.Shared) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var def func() error
	def = tx.Rollback
	defer func() {
		if err := def(); err != nil {
			log.Println(err)
		}
	}()

	_, err = db.Exec(
		"INSERT INTO posts(creator_id, id, title, body) VALUES($1, $2, $3, $4) ON CONFLICT DO NOTHING",
		sh.CreatorID,
		sh.ID,
		sh.Title,
		sh.Description,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO files(post_id, typeof, file_url) VALUES($1, $2, $3) ON CONFLICT DO NOTHING",
		sh.ID,
		sharedType,
		sh.FileURL,
	)
	if err != nil {
		return err
	}

	def = tx.Commit

	return nil
}

func (db *DB) StorePost(post *yp.Post) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	var def func() error
	def = tx.Rollback
	defer func() {
		if err := def(); err != nil {
			log.Println(err)
		}
	}()

	_, err = tx.Exec(
		"INSERT INTO posts(creator_id, id, title, body) VALUES($1, $2, $3, $4) ON CONFLICT DO NOTHING",
		post.Creator.ID,
		post.ID,
		post.Title,
		post.Body,
	)
	if err != nil {
		return err
	}

	if post.FileURL != "" {
		_, err = tx.Exec(
			"INSERT INTO files(post_id, typeof, file_url) VALUES($1, $2, $3) ON CONFLICT DO NOTHING",
			post.ID,
			postType,
			post.FileURL,
		)
		if err != nil {
			return err
		}
	}

	for _, attachment := range post.Attachments {
		_, err = tx.Exec(
			"INSERT INTO files(post_id, typeof, file_url) VALUES($1, $2, $3) ON CONFLICT DO NOTHING",
			post.ID,
			attachmentType,
			attachment.String(),
		)
		if err != nil {
			return err
		}
	}

	def = tx.Commit

	return nil
}

func (db *DB) GetFiles(limit int) ([]File, error) {
	rows, err := db.Query(
		`SELECT c.name, p.id, a.file_url, a.typeof
		FROM files a
		JOIN posts p
			ON a.post_id = p.id
		JOIN creators c
			ON p.creator_id = c.id
		WHERE a.downloaded IS NULL
		AND a.tries < 5
		ORDER BY a.tries ASC
		LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File

	for rows.Next() {
		var f File
		rows.Scan(&f.Creator, &f.PostID, &f.FileURL, &f.TypeOf)
		files = append(files, f)
	}

	return files, rows.Err()
}
