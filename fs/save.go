package fs

import (
	"crypto/sha256"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
)

var StoragePath = ".storage"

// Store the file in the destination, make a hardlink, and return the sha256 sum
func WriteAndLink(destination string, file io.Reader) (string, error) {
	storageFile, err := Write(file)
	if err != nil {
		return "", err
	}

	err = Link(storageFile, destination)

	return filepath.Base(storageFile), err
}

// Write to the storage and return the storage path
func Write(file io.Reader) (string, error) {
	tf, err := tmpfile(file)
	if err != nil {
		return "", err
	}

	filename := tf.Name()

	sum, err := checksum(tf)
	tf.Close()
	if err != nil {
		return "", err
	}

	return moveToStorage(filename, sum)

}

func Truncate(filename string) string {
	if len(filename) > 256 {
		return filename[:203] + "..." + filename[len(filename)-51:]
	}

	return filename
}

func Link(src, dest string) error {
	var err error
	if err = os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	dest = Truncate(dest)

	err = os.Link(src, dest)
	if err != nil {
		if lerr, ok := err.(*os.LinkError); ok && lerr.Err.Error() == "file exists" {
			return nil
		}
	}

	return err
}

func tmpfile(fs io.Reader) (*os.File, error) {
	filename := fmt.Sprint("tmpfile-", rand.Int())
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(f, fs)
	if err != nil {
		f.Close()
		os.Remove(filename)
		return nil, err
	}

	f.Sync()
	f.Seek(0, 0)
	return f, nil
}

func checksum(f io.Reader) (string, error) {
	h := sha256.New()
	_, err := io.Copy(h, f)

	return fmt.Sprintf("%x", h.Sum(nil)), err
}

func moveToStorage(src, sum string) (string, error) {
	var err error
	if err = os.MkdirAll(filepath.Join(StoragePath, sum[:2]), os.ModePerm); err != nil {
		return "", err
	}

	f := filepath.Join(StoragePath, sum[:2], sum)

	if _, err = os.Stat(f); err == nil {
		return f, os.Remove(src)
	} else if os.IsNotExist(err) {
		return f, os.Rename(src, f)
	}

	return f, err
}
