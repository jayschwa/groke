package bsp

import (
	"bytes"
	"github.com/ftrvxmtrx/groke/image/lmp"
	"image"
	"io"
	"io/ioutil"
	"unsafe"
)

type q1EdgeIndex struct {
	A uint16
	B uint16
}

type q1Face struct {
	Plane     uint16
	Side      uint16
	FirstEdge uint32
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
	q1NumLumps
)

const q1HeaderLen = 4

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

	h := sliceHeader(&b)
	h.Len = q1NumLumps
	h.Cap = q1NumLumps
	lumps := *(*[]bspLump)(unsafe.Pointer(&h))

	// entities
	m.Entities, err = bspReadEntities(lumps[q1LumpEntities].Data(q1HeaderLen, b))
	if err != nil || flags&EntitiesOnly != 0 {
		return
	}

	// textures
	m.Textures, err = q1ReadTextures(lumps[q1LumpTextures].Data(q1HeaderLen, b))
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
			N: qVector3(p.N),
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
		vertices = append(vertices, qVector3(verts32[i]))
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

func q1ReadTexInfo(b []byte) (texInfos []q1TexInfo, err error) {
	h := sliceHeader(&b)
	h.Len = len(b) / 40
	h.Cap = h.Len
	texInfos = *(*[]q1TexInfo)(unsafe.Pointer(&h))
	return
}

func q1ReadFaces(b []byte, lumps []bspLump, m *Model) (out []Face, err error) {
	var (
		edgeIndices []q1EdgeIndex
		faceEdges   []int16
		faces       []q1Face
		planes      []Plane
		texInfos    []q1TexInfo
		verts       []Vector3
	)

	data := make([][]byte, len(lumps))
	for i := range data {
		data[i] = lumps[i].Data(q1HeaderLen, b)
	}

	if planes, err = q1ReadPlanes(data[q1LumpPlanes]); err != nil {
		return
	} else if verts, err = q1ReadVertices(data[q1LumpVertices]); err != nil {
		return
	} else if texInfos, err = q1ReadTexInfo(data[q1LumpTextureInformation]); err != nil {
		return
	} else if edgeIndices, err = q1ReadEdgeIndices(data[q1LumpEdges]); err != nil {
		return
	} else if faceEdges, err = q1ReadFaceEdges(data[q1LumpFaceEdgeTables]); err != nil {
		return
	}

	fb := lumps[q1LumpFaces].Data(q1HeaderLen, b)
	h := sliceHeader(&fb)
	h.Len = int(lumps[q1LumpFaces].Size / 20)
	h.Cap = h.Len
	faces = *(*[]q1Face)(unsafe.Pointer(&h))

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
		if ti.Anim != 0 {
			flags |= TexAnimated
		}
		s := qVector3(ti.S)
		t := qVector3(ti.T)

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

func q1ReadTextures(b []byte) (texs []Texture, err error) {
	numTex := int(Uint32(b))
	texs = make([]Texture, 0, numTex)

	for i := 0; i < cap(texs); i++ {
		var h []byte
		if offset := Uint32(b[4+i*4:]); offset == 0xffffffff {
			texs = append(texs, Texture{
				Name:       "",
				DataSource: dataSourceInternal{},
			})
			continue
		} else {
			h = b[offset:]
		}

		nameLen := bytes.IndexByte(h[:16], 0)
		if nameLen < 0 || nameLen > 16 {
			nameLen = 16
		}

		dataOffset := int(Uint32(h[24:]))
		width := int(Uint32(h[16:]))
		height := int(Uint32(h[20:]))
		texs = append(texs, Texture{
			Name: string(bytes.ToLower(h[:nameLen])),
			DataSource: dataSourceInternal{
				h[dataOffset : dataOffset+width*height],
				width,
				height,
				lmpToImage,
			},
		})
	}

	return
}

func lmpToImage(w, h int, data []byte) image.Image {
	rect := image.Rect(0, 0, w, h)
	return &image.Paletted{
		Pix:     data,
		Stride:  w,
		Rect:    rect,
		Palette: lmp.Palette,
	}
}
