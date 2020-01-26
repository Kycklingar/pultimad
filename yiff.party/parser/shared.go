package yp

import (
	"bytes"
	"encoding/json"
	"io"
)

type Shared struct {
	CreatorID   int
	Title       string
	Description string
	FileURL     string `json:"file_url"`
	ID          int
	Created     int
}

func parseSharedFilesJson(jsonStream io.Reader) ([]Shared, error) {
	var tmp = struct {
		SharedFiles []Shared `json:"shared_files"`
	}{}

	var buf = new(bytes.Buffer)
	_, err := buf.ReadFrom(jsonStream)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(buf.Bytes(), &tmp); err != nil {
		return nil, err
	}

	return tmp.SharedFiles, nil
}
