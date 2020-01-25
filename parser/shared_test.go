package yp

import (
	"strings"
	"testing"
)

const jsonStr = `
{
    "posts": "test",
    "shared_files": [
        {
            "description": "Windows installer for version 0.16.1, downloaded from \"Download Information\" Patreon post",
            "file_name": "Bonfire_Windows_x64_0161f1.exe",
            "file_url": "https://yiff.party/shared_data/f9e4e3d851cb0fce2a5321eb2df7e276d030f668/Bonfire_Windows_x64_0161f1.exe",
            "id": 11744965863,
            "title": "Bonfire 0.16.1 Windows 64-Bit Installer",
            "uploaded": 1519192458
        },
        {
            "description": "Version 0.36.2",
            "file_name": "Bonfire_Windows_x64_0362.exe",
            "file_url": "https://yiff.party/shared_data/72183b60ed3e2b635d726d393616070a44fda845/Bonfire_Windows_x64_0362.exe",
            "id": 12666465633,
            "title": "Bonfire 0.36.2",
            "uploaded": 1572870427
        }
    ]
}
`

func TestShared(t *testing.T) {
	var str = strings.NewReader(jsonStr)
	files, err := parseSharedFilesJson(str)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 2 {
		t.Fatal("incorrect length")
	}

	if files[0].ID != 11744965863 {
		t.Fatal("incorrect id")
	}
}
