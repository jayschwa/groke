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

	err = Read(f, 0)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}

	f.Close()
}
