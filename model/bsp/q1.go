package bsp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"unsafe"
)

type q1EdgeIndex struct {
	A uint16
	B uint16
}

type q1Face struct {
	Plane     uint16
	Side      uint16
	Edge      uint32
	NumEdges  uint16
	TexInfoID uint16
	LightType uint8
	LightBase uint8
	Light     [2]uint8
	LightMap  uint32
}

type q1Lump struct {
	Offset uint32
	Size   uint32
}

type q1Plane struct {
	N [3]float32
	D float32
	T uint32
}

type q1TexInfo struct {
	S     [3]float32
	Ds    float32
	T     [3]float32
	Dt    float32
	TexID uint32
	Anim  uint32
}

const (
	q1LumpEntities = iota
	q1LumpPlanes
	q1LumpTextures
	q1LumpVertices
	q1LumpVisibility
	q1LumpNodes
	q1LumpTextureInformation
	q1LumpFaces
	q1LumpLightmaps
	q1LumpClipNodes
	q1LumpLeaves
	q1LumpMarkSurfaces
	q1LumpEdges
	q1LumpFaceEdgeTables
	q1LumpModels
)

func q1BSPRead(r io.Reader, flags int, m *Model) (err error) {
	var b []byte

	if rb, ok := r.(*bytes.Buffer); ok {
		b = rb.Bytes()
	} else if rb, err_ := ioutil.ReadAll(r); err_ == nil {
		b = rb
	} else {
		err = err_
		return
	}

	var lumps []q1Lump
	{
		h := sliceHeader(&b)
		h.Len = 15
		h.Cap = 15
		lumps = *(*[]q1Lump)(unsafe.Pointer(&h))
	}

	// entities
	m.Entities, err = q1ReadEntities(lumps[q1LumpEntities].Data(b))
	if err != nil || flags&EntitiesOnly != 0 {
		return
	}

	// textures
	m.Textures, err = q1ReadTextures(lumps[q1LumpTextures].Data(b))
	if err != nil {
		return
	}

	// faces
	m.Faces, err = q1ReadFaces(b, lumps, m)
	if err != nil {
		return
	}

	return
}

func q1ReadPlanes(b []byte) (planes []Plane, err error) {
	h := sliceHeader(&b)
	h.Len = len(b) / 2
	h.Cap = h.Len
	planes32 := *(*[]q1Plane)(unsafe.Pointer(&h))
	planes = make([]Plane, 0, len(planes32))

	for i := 0; i < cap(planes); i++ {
		p := planes32[i]
		planes = append(planes, Plane{
			N: Vector3{
				float64(p.N[0]),
				float64(p.N[1]),
				float64(p.N[2]),
			},
			D: float64(p.D),
			T: PlaneType(p.T),
		})
	}

	return
}

func q1ReadEdgeIndices(b []byte) (edgeIndices []q1EdgeIndex, err error) {
	h := sliceHeader(&b)
	h.Len = len(b) / 4
	h.Cap = h.Len
	edgeIndices = *(*[]q1EdgeIndex)(unsafe.Pointer(&h))
	return
}

func q1ReadVertices(b []byte) (vertices []Vector3, err error) {
	h := sliceHeader(&b)
	h.Len = len(b) / 12
	h.Cap = h.Len
	verts32 := *(*[][3]float32)(unsafe.Pointer(&h))
	vertices = make([]Vector3, 0, len(verts32))

	for i := 0; i < cap(vertices); i++ {
		vertices = append(vertices, Vector3{
			float64(verts32[i][0]),
			float64(verts32[i][1]),
			float64(verts32[i][2]),
		})
	}

	return
}

func q1ReadFaceEdges(b []byte) (faceEdges []int16, err error) {
	h := sliceHeader(&b)
	h.Len = len(b) / 2
	h.Cap = h.Len
	faceEdges = *(*[]int16)(unsafe.Pointer(&h))
	return
}

func q1ReadTextureInformation(b []byte) (texInfos []q1TexInfo, err error) {
	h := sliceHeader(&b)
	h.Len = len(b) / 40
	h.Cap = h.Len
	texInfos = *(*[]q1TexInfo)(unsafe.Pointer(&h))
	return
}

func q1ReadFaces(b []byte, lumps []q1Lump, m *Model) (out []Face, err error) {
	var (
		edgeIndices []q1EdgeIndex
		faceEdges   []int16
		faces       []q1Face
		planes      []Plane
		texInfos    []q1TexInfo
		verts       []Vector3
	)

	if planes, err = q1ReadPlanes(lumps[q1LumpPlanes].Data(b)); err != nil {
		return
	} else if verts, err = q1ReadVertices(lumps[q1LumpVertices].Data(b)); err != nil {
		return
	} else if texInfos, err = q1ReadTextureInformation(lumps[q1LumpTextureInformation].Data(b)); err != nil {
		return
	} else if edgeIndices, err = q1ReadEdgeIndices(lumps[q1LumpEdges].Data(b)); err != nil {
		return
	} else if faceEdges, err = q1ReadFaceEdges(lumps[q1LumpFaceEdgeTables].Data(b)); err != nil {
		return
	}

	fb := lumps[q1LumpFaces].Data(b)
	h := sliceHeader(&fb)
	h.Len = int(lumps[q1LumpFaces].Size / 20)
	h.Cap = h.Len
	faces = *(*[]q1Face)(unsafe.Pointer(&h))

	out = make([]Face, 0, len(faces))
	for i := 0; i < cap(out); i++ {
		face := faces[i]
		edges := make([]Edge, 0, int(face.NumEdges))
		fe := faceEdges[face.Edge : int(face.Edge)+cap(edges)]

		for _, fei := range fe {
			if fei < 0 {
				edges = append(edges, Edge{
					verts[edgeIndices[-fei].B],
					verts[edgeIndices[-fei].A],
				})
			} else {
				edges = append(edges, Edge{
					verts[edgeIndices[fei].A],
					verts[edgeIndices[fei].B],
				})
			}
		}

		ti := &texInfos[face.TexInfoID]
		var flags TexFlags
		if ti.Anim != 0 {
			flags |= TexAnimated
		}
		s := Vector3{float64(ti.S[0]), float64(ti.S[1]), float64(ti.S[2])}
		t := Vector3{float64(ti.T[0]), float64(ti.T[1]), float64(ti.T[2])}

		out = append(out, Face{
			Edges: edges,
			Front: face.Side == 0,
			Plane: &planes[face.Plane],
			TexInfo: TexInfo{
				S:       s,
				T:       t,
				Ds:      float64(ti.Ds),
				Dt:      float64(ti.Ds),
				Texture: &m.Textures[ti.TexID],
				Flags:   flags,
			},
		})
	}

	return
}

func q1ReadEntities(b []byte) (ents []Entity, err error) {
	ents = make([]Entity, 0, 64)
	ent := make(Entity)
	inBlock := 0

	for i := 0; i < len(b); {
		c := b[i]
		i++

		if c == '{' {
			inBlock++
		} else if c == '}' {
			if inBlock == 1 {
				ents = append(ents, ent)
				ent = make(Entity)
			}

			inBlock--
		} else if c == '"' && inBlock == 1 {
			keyIndex := bytes.IndexByte(b[i:], '"')
			if keyIndex < 0 {
				err = fmt.Errorf("key not closed with doublequote")
				break
			}
			key := stringFrom(b[i : i+keyIndex])
			i += keyIndex + 1

			for i < len(b) {
				c = b[i]
				i++

				if c == ' ' || c == '\t' {
					continue
				} else if c == '"' {
					valueIndex := bytes.IndexByte(b[i:], '"')
					if valueIndex < 0 {
						err = fmt.Errorf("key not closed with doublequote")
						break
					}
					value := stringFrom(b[i : i+valueIndex])
					i += valueIndex + 1

					ent[key] = value
					break
				} else {
					err = fmt.Errorf("bsp: unexpected char %q at pos %d", c, i)
				}
			}
		} else if c != ' ' && c != '\t' && c != '\r' && c != '\n' && c != 0 {
			err = fmt.Errorf("bsp: unexpected char %q at pos %d", c, i)
			return
		}
	}

	return
}

func q1ReadTextures(b []byte) (texs []Texture, err error) {
	numTex := Uint32(b)
	texs = make([]Texture, 0, numTex)

	for i := 0; i < cap(texs); i++ {
		offset := Uint32(b[4+i*4:])
		h := b[offset:]

		nameLen := bytes.IndexByte(h[:16], 0)
		if nameLen < 0 || nameLen > 16 {
			nameLen = 16
		}

		dataOffset := Uint32(h[24:])
		width := Uint32(h[16:])
		height := Uint32(h[20:])
		texs = append(texs, Texture{
			Name:   string(bytes.ToLower(h[:nameLen])),
			Width:  width,
			Height: height,
			Data:   h[dataOffset : dataOffset+width*height],
		})
	}

	return
}

func Uint32(b []byte) int {
	return int(binary.LittleEndian.Uint32(b))
}

func (lump q1Lump) String() string {
	return fmt.Sprintf("{0x%x, 0x%x}", lump.Offset, lump.Size)
}

func (lump q1Lump) Data(b []byte) []byte {
	offset := int(lump.Offset - 4)
	return b[offset : offset+int(lump.Size)]
}

func sliceHeader(raw *[]byte) reflect.SliceHeader {
	return *(*reflect.SliceHeader)(unsafe.Pointer(raw))
}

func stringHeader(raw *[]byte) reflect.StringHeader {
	return *(*reflect.StringHeader)(unsafe.Pointer(raw))
}

func stringFrom(b []byte) string {
	h := reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&b[0])),
		Len:  len(b),
	}

	return *(*string)(unsafe.Pointer(&h))
}
