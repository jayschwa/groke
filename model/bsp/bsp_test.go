package bsp

import (
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	f, err := os.Open("slaughtr.bsp")
	if err != nil {
		t.Fatal(err)
	}

	err = Read(f)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}

	f.Close()
}
