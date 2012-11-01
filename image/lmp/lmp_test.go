package lmp

import (
	"groke/archive/wad"
	"image"
	"image/png"
	"os"
	"testing"
)

func encode(name string, m image.Image) (err error) {
	var f *os.File

	if f, err = os.Create(name); err != nil {
		return
	}

	defer f.Close()
	err = png.Encode(f, m)

	return
}

func TestLMP(t *testing.T) {
	w, err := wad.OpenReader("test.wad")

	if err != nil {
		t.Fatal(err)
	}

	for _, f := range w.File {
		r, err := f.Open()
		if err != nil {
			t.Error(f.Name, err)
		}

		im, format, err := image.Decode(r)
		if err != nil {
			t.Error(f.Name, err)
		} else if format != "lmp" {
			t.Error(f.Name, "not lmp")
		} else if err := encode(f.Name+".png", im); err != nil {
			t.Error(f.Name, err)
		}
	}
}
