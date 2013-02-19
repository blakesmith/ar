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
}