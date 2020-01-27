package db

import "testing"

func TestDownload(t *testing.T) {
	var file File
	file.FileURL = "https://yiff.party/patreon_data/2533136/31252381/hirespack-2019kinktober.png"
	rc, err := file.Download()
	if err != nil {
		t.Fatal(err)
	}

	if rc == nil {
		t.Fatal("rc is nil")
	}
}
