package yp

import (
	"fmt"
	"testing"
)

func TestLoadCreators(t *testing.T) {
	creators, err := LoadCreators()
	for _, creator := range creators {
		fmt.Println(creator.ID, creator.Name)
	}
	if err != nil {
		t.Fatal(err)
	}
}
