package bsp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
)

type Vector3 [3]float64
type PlaneType byte

const (
	PlaneAxialX = iota
	PlaneAxialY
	PlaneAxialZ
	PlaneNonAxialX
	PlaneNonAxialY
	PlaneNonAxialZ
)

type Model struct {
	Entities string
	Planes   []struct {
		N Vector3
		D float64
		T PlaneType
	}
}

type lumpReader func(*io.SectionReader, *Model) error

type lump struct {
	Read lumpReader
}

type bspHeader struct {
	id    []byte
	lumps []lump
}

var bspHeaders = []*bspHeader{
	&bspHeader{
		[]byte{29, 0, 0, 0},
		[]lump{
			lump{q1BSPReadEntities},
			lump{q1BSPReadPlanes},
			lump{q1BSPReadTextures},
			lump{q1BSPReadVertices},
			lump{q1BSPReadVisibility},
			lump{q1BSPReadNodes},
			lump{q1BSPReadTextureInformation},
			lump{q1BSPReadFaces},
			lump{q1BSPReadLightmaps},
			lump{q1BSPReadClipNodes},
			lump{q1BSPReadLeaves},
			lump{q1BSPReadMarkSurfaces},
			lump{q1BSPReadEdges},
			lump{q1BSPReadFaceEdgeTables},
			lump{q1BSPReadModels},
		},
	},
}

var (
	ErrFormat = errors.New("bsp: invalid bsp format")
)

func (h *bspHeader) Read(r io.ReaderAt) error {
	numLumps := len(h.lumps)
	lmSection := io.NewSectionReader(r, int64(len(h.id)), int64(numLumps*8))
	lumps := make([]struct {
		Offset uint32
		Size   uint32
	}, numLumps)

	if err := binary.Read(lmSection, binary.LittleEndian, &lumps); err != nil {
		return err
	}

	var m Model

	for i, lm := range lumps {
		s := io.NewSectionReader(r, int64(lm.Offset), int64(lm.Size))
		if err := h.lumps[i].Read(s, &m); err != nil {
			return err
		}
	}

	log.Printf("%#v", m)

	return nil
}

func Read(r io.ReaderAt) error {
	for _, bh := range bspHeaders {
		id := make([]byte, len(bh.id))

		if idSize, err := r.ReadAt(id, 0); err != nil {
			return err
		} else if idSize == len(bh.id) && bytes.Compare(id, bh.id) == 0 {
			return bh.Read(r)
		}
	}

	return ErrFormat
}
