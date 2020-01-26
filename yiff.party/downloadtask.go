package yp

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/kycklingar/pultimad/daemon"
	"github.com/kycklingar/pultimad/fs"
	"github.com/kycklingar/pultimad/yiff.party/db"
)

type dlTask struct {
	file db.File
	db   *db.DB
}

func (t dlTask) Domain() string { return ypDomain }

func (t dlTask) Description() string { return fmt.Sprintf("%s - Downloading file %s", t.file.Creator, t.file.FileURL) }

func (t dlTask) Do() []daemon.Taskif {
	rc, err := t.file.Download()
	if err != nil {
		log.Println(err)
		t.db.Tried(t.file)
		return nil
	}

	sPath, err := fs.Write(rc)
	rc.Close()
	if err != nil {
		log.Println(err)
		return nil
	}

	sum := filepath.Base(sPath)
	filename := filepath.Base(t.file.FileURL)
	dest := filepath.Join("download", ypDomain, t.file.Creator, fmt.Sprintf("%d-%s-%s", t.file.PostID, sum[:4], filename))

	err = fs.Link(sPath, dest)
	if err != nil {
		log.Println(err)
		return nil
	}

	err = t.db.CheckFile(t.file, sum)
	if err != nil {
		log.Println(err)
	}

	return nil
}
