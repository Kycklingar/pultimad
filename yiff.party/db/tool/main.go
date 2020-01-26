package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: tool <connStr> <path>")
	}

	var connStr, path string

	connStr = os.Args[1]
	path = os.Args[2]

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT c.name, f.post_id, f.file_url, f.sha256 FROM files f JOIN posts p ON f.post_id = p.id JOIN creators c ON p.creator_id = c.id  WHERE f.sha256 IS NOT NULL")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	type file struct {
		creator string
		postid  int
		fileUrl string
		sha256  string
	}
	for rows.Next() {
		var f file
		if err = rows.Scan(&f.creator, &f.postid, &f.fileUrl, &f.sha256); err != nil {
			panic(err)
		}

		fmt.Println(f.creator, f.postid, f.sha256, f.fileUrl)

		fp := filePath(f.postid, path, f.creator, fileName(f.fileUrl), f.sha256)

		mkdirs(fp)
		if err = os.Link(storageFpath(f.sha256), fp); err != nil {
			log.Println(err)
		}
	}

}

func storageFpath(sum string) string {
	return fmt.Sprintf(".storage/%s/%s", sum[:2], sum)
}

func fileName(fileUrl string) string {
	fname, err := url.PathUnescape(filepath.Base(fileUrl))
	if err != nil {
		panic(err)
	}
	return fname
}

func filePath(pID int, path, creator, fileName, sha256 string) string {
	return fmt.Sprintf("%s/%s/%d-%s-%s", path, creator, pID, sha256[:4], fileName)
}

func mkdirs(fp string) {
	d, _ := filepath.Split(fp)
	os.MkdirAll(d, os.ModePerm)
}
