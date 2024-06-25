package gen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cscoding21/csgen"
)

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

func TestGetFile(t *testing.T) {
	//---test the passed in filename
	fileName1 := "file.test.go"
	gf := csgen.GetFile(fileName1)
	if !strings.HasSuffix(gf, fileName1) || !strings.HasPrefix(gf, "/") {
		t.Errorf("file name is not correct")
	}

	//---spoof the generator logic to get the file name
	fileName2 := "file.test2.go"
	os.Setenv("GOFILE", fileName2)
	gf = csgen.GetFile()

	if !strings.HasSuffix(gf, fileName2) || !strings.HasPrefix(gf, "/") {
		t.Errorf("file name is not correct")
	}
}

func TestCheckFileHasValidatorTags(t *testing.T) {
	file, err := filepath.Abs("../tests/data.go")
	if err != nil {
		t.Error(err)
	}

	hasTags, err := CheckFileHasValidatorTags(file)
	if err != nil {
		t.Error(err)
	}

	if !hasTags {
		t.Errorf("expected %s to have validator tags", file)
	}

	file, err = filepath.Abs("../tests/data2.go")
	if err != nil {
		t.Error(err)
	}

	hasTags, err = CheckFileHasValidatorTags(file)
	if err != nil {
		t.Error(err)
	}

	if hasTags {
		t.Errorf("expected %s NOT to have validator tags", file)
	}
}
