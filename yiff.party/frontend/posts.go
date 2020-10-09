package main

type post struct {
	CID int
	ID int
	Title string
	Body string
	Files files
}

type files struct {
	Post string
	Attachments []string
}

func posts() ([]post, error) {
	rows, err := db.Query(`
		SELECT creator_id, id, title, body
		FROM posts
		`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []post

	for rows.Next() {
		var p post
		err = rows.Scan(&p.CID, &p.ID, &p.Title, &p.Body)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	fs, err := files()

	return mergeFiles(posts, files)
}

func mergeFiles(posts []post, files []file) ([]post, error) {
	
}

type file struct {
	pid int
	typeof int
	url string
	sha256 string
}

func files() ([]file, error{
	rows, err := db.Query(`
		SELECT post_id, typeof, file_url, sha256
		FROM files
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []file

	for rows.Next() {
		var f file
		err = rows.Scan(&f.pid, &f.typeof, &url, &sha256
		if err != nil {
			return nil, err
		}

		files = append(files, f)

	}

	return files, nil
}
