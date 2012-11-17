package bsp

import (
	"bytes"
	"io"
	"io/ioutil"
	"unsafe"
)

type q2Face struct {
	Plane          uint16
	Side           uint16
	FirstEdge      uint32
	NumEdges       uint16
	TexInfoID      uint16
	LightMapStyles [4]uint8
	LightMap       uint32
}

type q2Node struct {
	Plane      uint32
	FrontChild int32
	BackChild  int32
	BBoxMin    [3]uint16
	BBoxMax    [3]uint16
	FirstFace  uint16
	NumFaces   uint16
}

type q2Plane q1Plane

type q2TexInfo struct {
	S       [3]float32
	Ds      float32
	T       [3]float32
	Dt      float32
	Flags   uint32
	Value   uint32
	Texture [32]byte
	Next    uint32
}

const (
	q2LumpEntities = iota
	q2LumpPlanes
	q2LumpVertices
	q2LumpVisibility
	q2LumpNodes
	q2LumpTextureInformation
	q2LumpFaces
	q2LumpLightmaps
	q2LumpLeaves
	q2LumpLeafFaceTable
	q2LumpLeafBrushTable
	q2LumpEdges
	q2LumpFaceEdgeTables
	q2LumpModels
	q2LumpBrushes
	q2LumpBrushSides
	q2LumpPop
	q2LumpAreas
	q2LumpAreaPortals
	q2NumLumps
)

const q2HeaderLen = 8

func q2BSPRead(r io.Reader, flags int, m *Model) (err error) {
	var b []byte

	if rb, ok := r.(*bytes.Buffer); ok {
		b = rb.Bytes()
	} else if rb, err_ := ioutil.ReadAll(r); err_ == nil {
		b = rb
	} else {
		err = err_
		return
	}

	h := sliceHeader(&b)
	h.Len = q2NumLumps
	h.Cap = q2NumLumps
	lumps := *(*[]bspLump)(unsafe.Pointer(&h))

	// entities
	m.Entities, err = bspReadEntities(lumps[q2LumpEntities].Data(q2HeaderLen, b))
	if err != nil || flags&EntitiesOnly != 0 {
		return
	}

	// faces
	m.Faces, err = q2ReadFaces(b, lumps, m)
	if err != nil {
		return
	}

	return
}

func q2ReadPlanes(b []byte) (planes []Plane, err error) {
	// same as in Quake
	return q1ReadPlanes(b)
}

func q2ReadTexInfo(b []byte) (texInfos []q2TexInfo, err error) {
	h := sliceHeader(&b)
	h.Len = len(b) / 76
	h.Cap = h.Len
	texInfos = *(*[]q2TexInfo)(unsafe.Pointer(&h))
	return
}

func q2ReadFaces(b []byte, lumps []bspLump, m *Model) (out []Face, err error) {
	var (
		edgeIndices []q1EdgeIndex
		faceEdges   []int16
		faces       []q2Face
		planes      []Plane
		texInfos    []q2TexInfo
		verts       []Vector3
	)

	data := make([][]byte, len(lumps))
	for i := range data {
		d := lumps[i].Data(q2HeaderLen, b)
		data[i] = d
	}

	// many are the same as in Quake
	if planes, err = q1ReadPlanes(data[q2LumpPlanes]); err != nil {
		return
	} else if verts, err = q1ReadVertices(data[q2LumpVertices]); err != nil {
		return
	} else if texInfos, err = q2ReadTexInfo(data[q2LumpTextureInformation]); err != nil {
		return
	} else if edgeIndices, err = q1ReadEdgeIndices(data[q2LumpEdges]); err != nil {
		return
	} else if faceEdges, err = q1ReadFaceEdges(data[q2LumpFaceEdgeTables]); err != nil {
		return
	}

	textures := make(map[[32]byte]*Texture)
	m.Textures = make([]Texture, 0)
	for _, ti := range texInfos {
		if _, ok := textures[ti.Texture]; !ok {
			nameLen := bytes.IndexByte(ti.Texture[:], 0)
			if nameLen < 0 || nameLen > 32 {
				nameLen = 32
			}
			m.Textures = append(m.Textures, Texture{
				DataSource: dataSourceExternal{},
				Name:       string(ti.Texture[:nameLen]),
			})
			textures[ti.Texture] = &m.Textures[len(m.Textures)-1]
		}
	}

	// set Next field
	for _, ti := range texInfos {
		if ti.Next < uint32(len(texInfos)) {
			textures[ti.Texture].Next = textures[texInfos[ti.Next].Texture]
		}
	}

	fb := lumps[q2LumpFaces].Data(q2HeaderLen, b)
	h := sliceHeader(&fb)
	h.Len = int(lumps[q2LumpFaces].Size / 20)
	h.Cap = h.Len
	faces = *(*[]q2Face)(unsafe.Pointer(&h))

	out = make([]Face, 0, len(faces))
	for i := 0; i < cap(out); i++ {
		face := faces[i]
		edges := make([]Edge, 0, int(face.NumEdges))
		fe := faceEdges[face.FirstEdge : int(face.FirstEdge)+cap(edges)]

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

		s := qVector3(ti.S)
		t := qVector3(ti.T)

		out = append(out, Face{
			Plane: &planes[face.Plane],
			TexInfo: TexInfo{
				S:       s,
				T:       t,
				Ds:      float64(ti.Ds),
				Dt:      float64(ti.Ds),
				Texture: textures[ti.Texture],
				Flags:   flags,
			},
		})
	}

	return
}
