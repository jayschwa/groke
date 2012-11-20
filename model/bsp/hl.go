package bsp

import (
	"bytes"
	"github.com/ftrvxmtrx/groke/image/hltex"
	"image"
	"io"
	"io/ioutil"
	"unsafe"
)

type hlFace struct {
	Plane     uint16
	Side      uint16
	FirstEdge uint32
	NumEdges  uint16
	TexInfoID uint16
	Light     [4]uint8
	LightMap  uint32
}

const (
	hlLumpEntities = iota
	hlLumpPlanes
	hlLumpTextures
	hlLumpVertices
	hlLumpVisibility
	hlLumpNodes
	hlLumpTextureInformation
	hlLumpFaces
	hlLumpLightmaps
	hlLumpClipNodes
	hlLumpLeaves
	hlLumpMarkSurfaces
	hlLumpEdges
	hlLumpFaceEdgeTables
	hlLumpModels
	hlNumLumps
)

const hlHeaderLen = 4

func hlBSPRead(r io.Reader, flags int, m *Model) (err error) {
	var b []byte

	if rb, ok := r.(*bytes.Buffer); ok {
		b = rb.Bytes()
	} else if rb, err_ := ioutil.ReadAll(r); err_ == nil {
		b = rb
	} else {
		err = err_
		return
	}

	lumps := bspLumpsFrom(b, hlNumLumps)

	// entities
	m.Entities, err = bspReadEntities(lumps[hlLumpEntities].Data(hlHeaderLen, b))
	if err != nil || flags&EntitiesOnly != 0 {
		return
	}

	// textures
	m.Textures, err = hlReadTextures(lumps[hlLumpTextures].Data(hlHeaderLen, b))
	if err != nil {
		return
	}

	// faces
	m.Faces, err = hlReadFaces(b, lumps, m)
	if err != nil {
		return
	}

	return
}

func hlReadFaces(b []byte, lumps []bspLump, m *Model) (out []Face, err error) {
	var (
		edgeIndices []q1EdgeIndex
		faceEdges   []int32
		faces       []hlFace
		planes      []Plane
		texInfos    []q1TexInfo
		verts       []Vector3
	)

	data := make([][]byte, len(lumps))
	for i := range data {
		data[i] = lumps[i].Data(hlHeaderLen, b)
	}

	if planes, err = q1ReadPlanes(data[hlLumpPlanes]); err != nil {
		return
	} else if verts, err = q1ReadVertices(data[hlLumpVertices]); err != nil {
		return
	} else if texInfos, err = q1ReadTexInfo(data[hlLumpTextureInformation]); err != nil {
		return
	} else if edgeIndices, err = q1ReadEdgeIndices(data[hlLumpEdges]); err != nil {
		return
	} else if faceEdges, err = q1ReadFaceEdges(data[hlLumpFaceEdgeTables]); err != nil {
		return
	}

	fb := lumps[hlLumpFaces].Data(hlHeaderLen, b)
	h := sliceHeader(&fb)
	h.Len = int(lumps[hlLumpFaces].Size / 20)
	h.Cap = h.Len
	faces = *(*[]hlFace)(unsafe.Pointer(&h))

	out = make([]Face, 0, len(faces))
	for _, face := range faces {
		v := make([]Vert, 0, int(face.NumEdges)*2)
		fe := faceEdges[face.FirstEdge : int(face.FirstEdge)+int(face.NumEdges)]

		for _, fei := range fe {
			if fei < 0 {
				v = append(v, Vert{
					Pos: verts[edgeIndices[-fei].B],
				})
				v = append(v, Vert{
					Pos: verts[edgeIndices[-fei].A],
				})
			} else {
				v = append(v, Vert{
					Pos: verts[edgeIndices[fei].A],
				})
				v = append(v, Vert{
					Pos: verts[edgeIndices[fei].B],
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
			Verts: v,
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

func hlReadTextures(b []byte) (texs []Texture, err error) {
	numTex := int(Uint32(b))
	texs = make([]Texture, 0, numTex)

	for i := 0; i < cap(texs); i++ {
		if offset := Uint32(b[4+i*4:]); offset == 0xffffffff || offset == 0 {
			texs = append(texs, Texture{
				Name:       "",
				DataSource: dataSourceInternal{},
			})
		} else {
			r := bytes.NewBuffer(b[offset:])

			var im image.Image
			if im, err = hltex.Decode(r); err != nil {
				return
			}
			hlt := im.(*hltex.HLTex)

			var source DataSource
			if hlt.Image.(*image.Paletted).Pix == nil {
				source = dataSourceExternal{}
			} else {
				source = dataSourceInternal{
					hlt.Image,
				}
			}

			texs = append(texs, Texture{
				Name:       hlt.Name,
				DataSource: source,
			})
		}
	}

	return
}
