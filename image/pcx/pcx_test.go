package pcx

import (
	"github.com/ftrvxmtrx/tga"
	"image"
	"io"
	"io/ioutil"
	"os"
	"strings"
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

func TestPCX(t *testing.T) {
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, fi := range files {
		name := fi.Name()
		if strings.HasSuffix(name, ".pcx") {
			var r io.Reader
			r, err = os.Open("testdata/" + name)
			if err != nil {
				t.Error(name, err)
			}

			var im image.Image
			if im, err = Decode(r); err != nil {
				t.Fatal(err)
			} else if err = encode(im, name+".tga"); err != nil {
				t.Error(name, err)
			}
		}
	}
}
