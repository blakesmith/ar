package ar

import (
	"log"
	"os"
	"testing"
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

	expected := "hello.txt"
	if header.Name != expected {
		log.Println([]byte(header.Name))
		t.Errorf("Header name should be %s but is %s", expected, header.Name)
	}
}