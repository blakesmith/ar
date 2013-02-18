package ar

import (
	"bytes"
	"testing"
)

func TestGlobalHeaderWrite(t *testing.T) {
	var buf bytes.Buffer
	writer := NewWriter(&buf)
	if err := writer.WriteHeader(new(Header)); err != nil {
		t.Errorf(err.Error())
	}

	globalHeader := buf.Bytes()
	expectedHeader := []byte("!<arch>\n")
	if !bytes.Equal(globalHeader, expectedHeader) {
		t.Errorf("Global header should be %s but is %s", expectedHeader, globalHeader)
	}
}