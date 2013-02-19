package ar

import (
	"os"
	"testing"
	"time"
)

func TestReadHeader(t *testing.T) {
	f, err := os.Open("./fixtures/hello.a")
	defer f.Close()

	if err != nil {
		t.Errorf(err.Error())
	}
	reader := NewReader(f)
	header, err := reader.Next()
	if err != nil {
		t.Errorf(err.Error())
	}

	expectedName := "hello.txt"
	if header.Name != expectedName {
		t.Errorf("Header name should be %s but is %s", expectedName, header.Name)
	}
	expectedModTime := time.Unix(1361157466, 0)
	if header.ModTime != expectedModTime {
		t.Errorf("ModTime should be %s but is %s", expectedModTime, header.ModTime)
	}
	expectedUid := 501
	if header.Uid != expectedUid {
		t.Errorf("Uid should be %s but is %s", expectedUid, header.Uid)
	}
	expectedGid := 20
	if header.Gid != expectedGid {
		t.Errorf("Gid should be %s but is %s", expectedGid, header.Gid)
	}
}