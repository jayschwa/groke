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

type Face struct {
	Edges   []Edge
	Front   bool
	Plane   *Plane
	TexInfo *TexInfo
}

type Model struct {
	Edges        []Edge
	EdgeVIndices []EdgeVIndex
	Entities     []Entity
	Faces        []Face
	Planes       []Plane
	Textures     []Texture
	TexInfos     []TexInfo
	Verts        []Vector3
}

type lumpReader func(*io.SectionReader, *Model) error

type lump struct {
	// i is an order in which the section is loaded.
	i int
	// name is the name of the section.
	Name string
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
			lump{0, "entities", q1BSPReadEntities, NoEntities},
			lump{1, "planes", q1BSPReadPlanes, 0},
			lump{2, "textures", q1BSPReadTextures, NoTextures},
			lump{3, "vertices", q1BSPReadVertices, 0},
			lump{4, "visibility", q1BSPReadVisibility, 0},
			lump{5, "nodes", q1BSPReadNodes, 0},
			lump{6, "texture info", q1BSPReadTextureInformation, NoTextures},
			lump{8, "lightmaps", q1BSPReadLightmaps, NoLightmaps},
			lump{9, "clipnodes", q1BSPReadClipNodes, 0},
			lump{10, "leaves", q1BSPReadLeaves, 0},
			lump{11, "mark surfaces", q1BSPReadMarkSurfaces, 0},
			lump{12, "edges", q1BSPReadEdges, 0},
			lump{13, "face edge tables", q1BSPReadFaceEdgeTables, 0},
			lump{14, "models", q1BSPReadModels, 0},
			lump{7, "faces", q1BSPReadFaces, 0},
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

	for _, hl := range h.lumps {
		lm := lumps[hl.i]
		if flags&hl.Flag == 0 {
			s := io.NewSectionReader(r, int64(lm.Offset), int64(lm.Size))
			if err := hl.Read(s, &m); err != nil {
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
