package bsp

import (
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	f, err := os.Open("test.bsp")
	if err != nil {
		t.Fatal(err)
	}

	var m *Model
	m, err = Read(f, 0)
	if err != nil {
		f.Close()
		t.Fatal(err)
	} else if m == nil {
		f.Close()
		t.Fatal(m)
	}

	f.Close()
}
