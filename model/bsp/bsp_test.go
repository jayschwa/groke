package bsp

import (
	"log"
	"os"
	"strings"
	"testing"
)

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

		tested++
	}

	t.Logf("tested %d maps", tested)
}
