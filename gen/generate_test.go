package gen

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	Generate("test_struct.go")

}

func TestGetTagMap(t *testing.T) {
	tagString := "min(21)"

	tm := getTagMap(tagString)

	if len(*tm) != 1 {
		t.Errorf("tag map length is not 1")
	}

	if (*tm)["min"] != "21" {
		t.Errorf("tag map value is not 21")
	}
}
