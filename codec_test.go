package proto

import (
	"bytes"
	"testing"
)

func TestReadBody(t *testing.T) {
	t.Skip("skip as no proto")
	s := &struct {
		Name string
	}{}
	c := NewCodec()
	b := bytes.NewReader(nil)
	err := c.ReadBody(b, s)
	if err != nil {
		t.Fatal(err)
	}
}
