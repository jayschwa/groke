package bsp

import (
	"image"
	"image/png"
	"log"
	"os"
	"strings"
	"testing"
)

func encode(m image.Image, name string) (err error) {
	var f *os.File

	if f, err = os.Create("testdata/" + name + ".png"); err != nil {
		return
	}

	defer f.Close()
	err = png.Encode(f, m)

	return
}

func TestRead(t *testing.T) {
	f, err := os.Open("testdata")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var names []string
	names, err = f.Readdirnames(0)
	if err != nil {
		t.Fatal(err)
	}

	var tested int
	for _, name := range names {
		if !strings.HasSuffix(name, ".bsp") {
			continue
		}

		log.Print(name)
		f, err = os.Open("testdata/" + name)
		if err != nil {
			t.Fatal(err)
		}

		var m *Model
		m, err = Read(f, 0)
		f.Close()

		if err != nil {
			continue
			//t.Fatal(err)
		} else if m == nil {
			t.Fatal(m)
		}

		for _, t := range m.Textures {
			if t.Name != "" {
				encode(t.Image(), name+"_"+t.Name)
			}
		}

		tested++
	}

	t.Logf("tested %d maps", tested)
}
