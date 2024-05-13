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

	t.Log(tm)

}
