package lmp

import (
	"github.com/ftrvxmtrx/groke/archive/wad"
	"github.com/ftrvxmtrx/tga"
	"image"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
)

func encode(m image.Image, name string) (err error) {
	var f *os.File

	if f, err = os.Create("testdata/" + name + ".tga"); err != nil {
		return
	}

	defer f.Close()
	err = tga.Encode(f, m)

	return
}

func TestLMP(t *testing.T) {
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

	group := new(sync.WaitGroup)
	group.Add(len(names))

	for _, name := range names {
		go func(name string) {
			defer group.Done()

			if strings.HasSuffix(name, ".wad") {
				log.Print(name)

				w, err := wad.OpenReader("testdata/" + name)
				if err == nil {
					for _, f := range w.File {
						r, err := f.Open()
						if err != nil {
							t.Error(f.Name, err)
						}

						im, err := Decode(r)
						if err != nil {
							t.Error(f.Name, err)
						} else if err := encode(im, name+"_"+f.Name); err != nil {
							t.Error(f.Name, err)
						}
					}
				}
			}
		}(name)
	}

	group.Wait()
}
