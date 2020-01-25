package daemon

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"

	db "github.com/kycklingar/pultimad/storage"
)

func (d *daemon) download(f db.File) error {
	rc, err := f.Download()
	if err != nil {
		return err
	}

	tmpFile, err := createTmpFile(rc)
	rc.Close()
	if err != nil {
		tmpFile.Close()
		return err
	}

	sum, err := sum(tmpFile)
	if err != nil {
		cleanup(tmpFile)
		return err
	}

	fi, err := tmpFile.Stat()
	tmpFile.Close()
	if err = mvFile(fi.Name(), sum); err != nil {
		return err
	}

	fpath := fmt.Sprintf("download/%s", f.Creator)
	os.MkdirAll(fpath, os.ModePerm)

	fname, err := url.PathUnescape(filepath.Base(f.FileURL))

	link(sum, fpath, fname, f.PostID)

	return d.b.CheckFile(f, sum)
}


func createTmpFile(fs io.Reader) (*os.File, error) {
	f, err := os.Create(fmt.Sprint("tmpfile-", rand.Int()))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(f, fs)
	f.Sync()
	return f, err
}

func cleanup(f *os.File) {
	fi, err := f.Stat()
	f.Close()
	if err != nil {
		log.Println(err)
		return
	}

	os.Remove(fi.Name())
}

func sum(f io.ReadSeeker) (string, error) {
	f.Seek(0, io.SeekStart)
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func mvFile(file, sum string) error {
	os.MkdirAll(fmt.Sprintf(".storage/%s", sum[:2]), os.ModePerm)
	f := fmt.Sprintf(".storage/%s/%s", sum[:2], sum)

	var err error
	if _, err = os.Stat(f); err == nil {
		return os.Remove(file)
		return nil
	} else if os.IsNotExist(err) {
		return os.Rename(file, f)
	}

	return err
}

func link(sum, path, fname string, postID int) error {
	filepath := fmt.Sprintf("%s/%d-%s-%s", path, postID, sum[:4], fname)
	return os.Link(fmt.Sprintf(".storage/%s/%s", sum[:2], sum), filepath)
}
