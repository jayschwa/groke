package bsp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"reflect"
	"unsafe"
)

type q1Lump struct {
	Offset uint32
	Size   uint32
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

func Uint32(b []byte) int {
	return int(binary.LittleEndian.Uint32(b))
}

func (lump q1Lump) String() string {
	return fmt.Sprintf("{0x%x, 0x%x}", lump.Offset, lump.Size)
}

func (lump q1Lump) Data(b []byte) []byte {
	offset := int(lump.Offset - 4)
	log.Printf("lump [%d:%d]", offset, offset+int(lump.Size))
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
	{
		// edge indices
		var edgeIndices []q1EdgeIndex
		eb := lumps[q1LumpEdges].Data(b)
		h := sliceHeader(&eb)
		h.Len = int(lumps[q1LumpEdges].Size / 4)
		h.Cap = h.Len
		edgeIndices = *(*[]q1EdgeIndex)(unsafe.Pointer(&h))

		// vertices
		var verts []Vector3
		{
			vb := lumps[q1LumpVertices].Data(b)
			h := sliceHeader(&vb)
			h.Len = int(lumps[q1LumpVertices].Size / 12)
			h.Cap = h.Len
			verts32 := *(*[][3]float32)(unsafe.Pointer(&h))
			verts = make([]Vector3, len(verts32))
			for i := 0; i < len(verts); i++ {
				verts[i][0] = float64(verts32[i][0])
				verts[i][1] = float64(verts32[i][1])
				verts[i][2] = float64(verts32[i][2])
			}
		}

		// face edge tables
		var faceEdges []int16
		{
			fb := lumps[q1LumpFaceEdgeTables].Data(b)
			h := sliceHeader(&fb)
			h.Len = int(lumps[q1LumpFaceEdgeTables].Size / 2)
			h.Cap = h.Len
			faceEdges = *(*[]int16)(unsafe.Pointer(&h))
		}

		// planes
		var planes []Plane
		{
			pb := lumps[q1LumpPlanes].Data(b)
			h := sliceHeader(&pb)
			h.Len = int(lumps[q1LumpPlanes].Size / 2)
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
		}

		// texture information
		var texInfos []q1TexInfo
		{
			tb := lumps[q1LumpTextureInformation].Data(b)
			h := sliceHeader(&tb)
			h.Len = int(lumps[q1LumpTextureInformation].Size / 40)
			h.Cap = h.Len
			texInfos = *(*[]q1TexInfo)(unsafe.Pointer(&h))
		}

		// faces
		var faces []q1Face
		{
			fb := lumps[q1LumpFaces].Data(b)
			h := sliceHeader(&fb)
			h.Len = int(lumps[q1LumpFaces].Size / 20)
			h.Cap = h.Len
			faces = *(*[]q1Face)(unsafe.Pointer(&h))
		}

		m.Faces = make([]Face, 0, len(faces))
		for i := 0; i < cap(m.Faces); i++ {
			face := faces[i]
			fe := faceEdges[face.Edge : int(face.Edge)+int(face.NumEdges)]

			edges := make([]Edge, int(face.NumEdges))
			for e := 0; e < len(edges); e++ {
				fei := fe[e]
				if fei < 0 {
					edges[e] = Edge{
						verts[edgeIndices[-fei].B],
						verts[edgeIndices[-fei].A],
					}
				} else {
					edges[e] = Edge{
						verts[edgeIndices[fei].A],
						verts[edgeIndices[fei].B],
					}
				}
			}

			ti := &texInfos[face.TexInfoID]
			var flags TexFlags
			if ti.Anim != 0 {
				flags |= TexAnimated
			}
			s := Vector3{float64(ti.S[0]), float64(ti.S[1]), float64(ti.S[2])}
			t := Vector3{float64(ti.T[0]), float64(ti.T[1]), float64(ti.T[2])}

			m.Faces = append(m.Faces, Face{
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
	}

	return
}

func q1ReadEntities(b []byte) (ents []Entity, err error) {
	ents = make([]Entity, 256)
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
			data:   h[dataOffset : dataOffset+width*height],
		})
	}

	return
}
