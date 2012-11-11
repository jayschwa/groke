package bsp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
)

type Vector3 [3]float64
type PlaneType uint8
type TexFlags int

type Plane struct {
	N Vector3
	D float64
	T PlaneType
}

type EdgeVIndex struct {
	Ai int
	Bi int
}

type Edge struct {
	A Vector3
	B Vector3
}

type Texture struct {
	Name       string
	Width      int
	Height     int
	dataOffset int64
	r          *io.SectionReader
}

type TexInfo struct {
	S       Vector3
	T       Vector3
	Ds      float64
	Dt      float64
	Texture *Texture
	Flags   TexFlags
}

type Entity map[string]string

type Model struct {
	Edges        []Edge
	EdgeVIndices []EdgeVIndex
	Entities     []Entity
	Planes       []Plane
	Textures     []Texture
	TexInfos     []TexInfo
	Verts        []Vector3
}

type lumpReader func(*io.SectionReader, *Model) error

type lump struct {
	Read lumpReader
	Flag int
}

type bspHeader struct {
	id    []byte
	lumps []lump
}

var pt2String = []string{
	"PlaneAxialX",
	"PlaneAxialY",
	"PlaneAxialZ",
	"PlaneNonAxialX",
	"PlaneNonAxialY",
	"PlaneNonAxialZ",
}

const (
	PlaneAxialX = iota
	PlaneAxialY
	PlaneAxialZ
	PlaneNonAxialX
	PlaneNonAxialY
	PlaneNonAxialZ
)

const (
	TexAnimated = TexFlags(1 << iota)
)

const (
	NoEntities = 1 << iota
	NoTextures
	NoLightmaps
)

var bspHeaders = []*bspHeader{
	&bspHeader{
		[]byte{29, 0, 0, 0},
		[]lump{
			lump{q1BSPReadEntities, NoEntities},
			lump{q1BSPReadPlanes, 0},
			lump{q1BSPReadTextures, NoTextures},
			lump{q1BSPReadVertices, 0},
			lump{q1BSPReadVisibility, 0},
			lump{q1BSPReadNodes, 0},
			lump{q1BSPReadTextureInformation, NoTextures},
			lump{q1BSPReadFaces, 0},
			lump{q1BSPReadLightmaps, NoLightmaps},
			lump{q1BSPReadClipNodes, 0},
			lump{q1BSPReadLeaves, 0},
			lump{q1BSPReadMarkSurfaces, 0},
			lump{q1BSPReadEdges, 0},
			lump{q1BSPReadFaceEdgeTables, 0},
			lump{q1BSPReadModels, 0},
		},
	},
}

var (
	ErrFormat = errors.New("bsp: invalid bsp format")
)

func (t PlaneType) String() string {
	return pt2String[t]
}

func (h *bspHeader) Read(r io.ReaderAt, flags int) error {
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
		if flags&h.lumps[i].Flag == 0 {
			s := io.NewSectionReader(r, int64(lm.Offset), int64(lm.Size))
			if err := h.lumps[i].Read(s, &m); err != nil {
				return err
			}
		}
	}

	log.Printf("%#v", m.Textures)

	return nil
}

func Read(r io.ReaderAt, flags int) error {
	for _, bh := range bspHeaders {
		id := make([]byte, len(bh.id))

		if idSize, err := r.ReadAt(id, 0); err != nil {
			return err
		} else if idSize == len(bh.id) && bytes.Compare(id, bh.id) == 0 {
			return bh.Read(r, flags)
		}
	}

	return ErrFormat
}
