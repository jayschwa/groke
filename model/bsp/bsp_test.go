package bsp

import (
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

	group := new(sync.WaitGroup)
	group.Add(len(names))

	for _, name := range names {
		go func(name string) {
			defer group.Done()

			if strings.HasSuffix(name, ".bsp") {
				f, err := os.Open("testdata/" + name)
				if err != nil {
					t.Fatal(err)
				}

				var m *Model
				m, err = Read(f, 0)
				f.Close()

				if err == nil {
					if m == nil {
						t.Fatal(m)
					}

					log.Print(name)
					for _, t := range m.Textures {
						if t.Name != "" {
							switch ds := t.DataSource.(type) {
							case dataSourceInternal:
								encode(ds.Image, name+"_"+t.Name)
							case dataSourceExternal:
								// FIXME
							default:
								log.Fatal("unknown data source type")
							}
						}
					}
				} else {
					log.Fatal(name, " - ", err)
				}
			}
		}(name)
	}

	group.Wait()
}
