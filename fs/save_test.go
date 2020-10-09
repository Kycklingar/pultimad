package fs

import (
	"os"
	"strings"
	"testing"
)

func TestWrite(t *testing.T) {
	StoragePath = ".test-storage"
	defer os.RemoveAll(StoragePath)
	var strStream = strings.NewReader("Hello World!\n")

	fp, err := Write(strStream)
	if err != nil {
		t.Fatal(err)
	}

	exp := StoragePath + "/03/03ba204e50d126e4674c005e04d82e84c21366780af1f43bd54a37816b6ab340"
	if fp != exp {
		t.Fatalf("Expected '%s', got %s", exp, fp)
	}
}

func TestWriteAndLink(t *testing.T) {
	StoragePath = ".test-storage"
	tmpPath := "tmp-test"

	defer os.RemoveAll(StoragePath)
	defer os.RemoveAll(tmpPath)

	var strStream = strings.NewReader("Hello World!\n")

	sum1, err := WriteAndLink(tmpPath+"/1", strStream)
	if err != nil {
		t.Fatal(err)
	}

	strStream.Seek(0, 0)
	sum2, err := WriteAndLink(tmpPath+"/2", strStream)
	if err != nil {
		t.Fatal(err)
	}

	strStream.Seek(0, 0)
	sum3, err := WriteAndLink(tmpPath+"/1", strStream)
	if err != nil {
		t.Fatal(err)
	}

	if sum1 != sum2 && sum2 != sum3 {
		t.Fatalf("Sums not matching %s, %s, %s", sum1, sum2, sum3)
	}
}

func TestWriteMaxLength(t *testing.T) {
	StoragePath = ".test-storage"
	tmpPath := "tmp-test"

	var filename string

	for i := 0; i < 240; i++ {
		filename += "a"
	}

	filename += "EndOfTheFileName.txt"

	var ss = strings.NewReader("Hello World!\n")

	_, err := WriteAndLink(tmpPath+"/"+filename, ss)
	if err != nil {
		t.Fatal(err)
	}

}
