package bsp

import (
	"bytes"
	"image"
	"image/color"
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
		var h []byte
		if offset := Uint32(b[4+i*4:]); offset == 0xffffffff || offset == 0 {
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
		name := string(bytes.ToLower(h[:nameLen]))

		dataOffset := int(Uint32(h[24:]))
		width := int(Uint32(h[16:]))
		height := int(Uint32(h[20:]))

		if dataOffset > 0 {
			var palOffset int
			var dataEnd int
			if palOffset = int(Uint32(h[36:])); palOffset != 0 {
				palOffset += width * height / 64
			} else if palOffset = int(Uint32(h[32:])); palOffset != 0 {
				palOffset += width * height / 16
			} else if palOffset = int(Uint32(h[28:])); palOffset != 0 {
				palOffset += width * height / 4
			} else {
				palOffset = dataOffset + width*height
			}

			palOffset += 2
			dataEnd = palOffset + 256*3

			texs = append(texs, Texture{
				Name: name,
				DataSource: dataSourceInternal{
					h[dataOffset:dataEnd],
					width,
					height,
					hlTexToImage,
				},
			})
		} else {
			texs = append(texs, Texture{
				Name:       name,
				DataSource: dataSourceExternal{},
			})
		}
	}

	return
}

func hlTexToImage(w, h int, data []byte) image.Image {
	palette := make(color.Palette, 0)
	palOffset := len(data) - 256*3

	for i := 0; i < 255; i++ {
		o := i * 3
		palette = append(palette, color.NRGBA{data[palOffset+o+0], data[palOffset+o+1], data[palOffset+o+3], 0xff})
	}
	palette = append(palette, color.NRGBA{0, 0, 0, 0})

	rect := image.Rect(0, 0, w, h)
	return &image.Paletted{
		Pix:     data,
		Stride:  w,
		Rect:    rect,
		Palette: palette,
	}
}
