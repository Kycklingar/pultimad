package daemon

import (
	"fmt"
	"path/filepath"

	db "github.com/kycklingar/pultimad/storage"
	"github.com/kycklingar/pultimad/storage/fs"
)

func (d *daemon) download(f db.File) error {
	rc, err := f.Download()
	if err != nil {
		return err
	}

	sPath, err := fs.Write(rc)
	rc.Close()
	if err != nil {
		return err
	}

	sum := filepath.Base(sPath)
	filename := filepath.Base(f.FileURL)
	dest := filepath.Join(d.downloadPath, f.Creator, fmt.Sprintf("%d-%s-%s", f.PostID, sum[:4], filename))

	err = fs.Link(sPath, dest)
	if err != nil {
		return err
	}

	return d.b.CheckFile(f, sum)
}
