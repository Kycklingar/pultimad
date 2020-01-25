package fs

// Store the file in the destination, make a hardlink, and return the sha256 sum
func WriteAndLink(destination string, file io.Reader) (string, error) {
	storageFile, err := Write(file)
	if err != nil {
		return err
	}

	err = os.Link(storageFile, destination)

	return filepath.Base(storageFile), err
}

// Write to the storage
func Write(file io.Reader) (string, error) {
	tf, err := tmpfile(file)
	if err != nil
	{
		return err
	}

	filename := tf.Name()

	sum, err := checksum(tf)
	tf.Close()
	if err != nil{
		return err
	}

	return moveToStorage(filename, sum)

}

func tmpfile(fs io.Reader) (*os.File, error) {
	filename := fmt.Sprintf("tmpfile-", rand.Int())
	f, err := os.Create(filename)
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(f, fs)
	f.Sync()
	return f, err
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

	f := filePath.Join(StoragePath, sum[:2], sum)

	if _, err = os.Stat(f); err == nil {
		return f, os.Remove(src)
	} else if os.IsNotExist(err) {
		return f, os.Rename(src, f)
	}

	return f, err
}

