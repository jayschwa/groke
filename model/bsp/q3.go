package bsp

import (
	"bytes"
	"io"
	"io/ioutil"
	"unsafe"
)

type q3Effect struct {
	Name  [64]byte
	Brush uint32
	Side  uint32
}

type q3Face struct {
	Texture        uint32
	Effect         uint32
	Type           q3FaceType
	FirstVertex    uint32
	NumVerts       uint32
	FirstMeshVert  uint32
	NumMeshVerts   uint32
	LightMap       uint32
	LightMapX      uint32
	LightMapY      uint32
	LightMapW      uint32
	LightMapH      uint32
	LightMapOrigin [3]float32
	LightMapS      [3]float32
	LightMapT      [3]float32
	Normal         [3]float32
	PatchW         uint32
	PatchH         uint32
}

type q3Plane struct {
	N [3]float32
	D float32
}

type q3TexInfo struct {
	Name     [64]byte
	SurFlags q3SurfFlags
	Contents q3Contents
}

type q3Vertex struct {
	Position       [3]float32
	SurfCoords     [2]float32
	LightMapCoords [2]float32
	Normal         [3]float32
	Color          [4]uint8
}

type q3Contents uint32

const (
	q3ContentsSolid = q3Contents(1 << iota)
	_
	_
	q3ContentsLava
	q3ContentsSlime
	q3ContentsWater
	q3ContentsFog
	q3ContentsNotTeam1
	q3ContentsNotTeam2
	q3ContentsNoBotClip
	_
	_
	_
	q3ContentsAreaPortal
	q3ContentsPlayerClip
	q3ContentsMonsterClip
	q3ContentsTeleporter
	q3ContentsJumpPad
	q3ContentsClusterPortal
	q3ContentsDoNotEnter
	q3ContentsBotclip
	q3ContentsMover
	q3ContentsOrigin
	q3ContentsBody
	q3ContentsCorpse
	q3ContentsDetail
	q3ContentsStructural
	q3ContentsTranslucent
	q3ContentsTrigger
	q3ContentsNoDrop
)

type q3SurfFlags uint32

const (
	q3SurfNoDamage = q3SurfFlags(1 << iota)
	q3SurfSlick
	q3SurfSky
	q3SurfLadder
	q3SurfNoImpact
	q3SurfNoMarks
	q3SurfFlesh
	q3SurfNoDraw
	q3SurfHint
	q3SurfSkip
	q3SurfNoLightmap
	q3SurfPointLight
	q3SurfMetalSteps
	q3SurfNoSteps
	q3SurfNonSolid
	q3SurfLightFilter
	q3SurfAlphaShadow
	q3SurfNoDLight
	q3SurfDust
)

type q3FaceType uint32

const (
	_ = q3FaceType(iota)
	q3FacePolygon
	q3FacePatch
	q3FaceMesh
	q3FaceBillboard
)

const (
	q3LumpEntities = iota
	q3LumpTextureInformation
	q3LumpPlanes
	q3LumpNodes
	q3LumpLeaves
	q3LumpLeafFaceTable
	q3LumpLeafBrushTable
	q3LumpModels
	q3LumpBrushes
	q3LumpBrushSides
	q3LumpVertices
	q3LumpMeshVertices
	q3LumpEffects
	q3LumpFaces
	q3LumpLightmaps
	q3LumpLightVols
	q3LumpVisibility
	q3NumLumps
)

const q3HeaderLen = 8

func q3BSPRead(r io.Reader, flags int, m *Model) (err error) {
	var b []byte

	if rb, ok := r.(*bytes.Buffer); ok {
		b = rb.Bytes()
	} else if rb, err_ := ioutil.ReadAll(r); err_ == nil {
		b = rb
	} else {
		err = err_
		return
	}

	lumps := bspLumpsFrom(b, q3NumLumps)

	// entities
	m.Entities, err = bspReadEntities(lumps[q3LumpEntities].Data(q3HeaderLen, b))
	if err != nil || flags&EntitiesOnly != 0 {
		return
	}

	// faces
	m.Faces, err = q3ReadFaces(b, lumps, m)
	if err != nil {
		return
	}

	m.Triangle = true

	return
}

func q3ReadPlanes(b []byte) (planes []Plane, err error) {
	h := sliceHeader(&b)
	h.Len = len(b) / 16
	h.Cap = h.Len
	planes32 := *(*[]q3Plane)(unsafe.Pointer(&h))
	planes = make([]Plane, 0, len(planes32))

	for i := 0; i < cap(planes); i++ {
		p := planes32[i]
		planes = append(planes, Plane{
			N: qVector3(p.N),
			D: float64(p.D),
			T: PlaneNoType,
		})
	}

	return
}

func q3ReadTexInfo(b []byte) (texInfos []q3TexInfo, err error) {
	h := sliceHeader(&b)
	h.Len = len(b) / 72
	h.Cap = h.Len
	texInfos = *(*[]q3TexInfo)(unsafe.Pointer(&h))
	return
}

func q3ReadMeshVertices(b []byte) (meshVerts []uint32, err error) {
	h := sliceHeader(&b)
	h.Len = len(b) / 4
	h.Cap = h.Len
	meshVerts = *(*[]uint32)(unsafe.Pointer(&h))
	return
}

func q3ReadVertices(b []byte) (vertices []Vector3, err error) {
	h := sliceHeader(&b)
	h.Len = len(b) / 44
	h.Cap = h.Len
	verts32 := *(*[]q3Vertex)(unsafe.Pointer(&h))
	vertices = make([]Vector3, 0, len(verts32)) // FIXME -- general case

	for i := 0; i < cap(vertices); i++ {
		vertices = append(vertices, qVector3(verts32[i].Position))
	}

	return
}

func q3ReadFaces(b []byte, lumps []bspLump, m *Model) (out []Face, err error) {
	var (
		faces     []q3Face
		planes    []Plane
		texInfos  []q3TexInfo
		meshVerts []uint32
		verts     []Vector3
	)

	data := make([][]byte, len(lumps))
	for i := range data {
		d := lumps[i].Data(q3HeaderLen, b)
		data[i] = d
	}

	// many are the same as in Quake
	if planes, err = q3ReadPlanes(data[q3LumpPlanes]); err != nil {
		return
	} else if verts, err = q3ReadVertices(data[q3LumpVertices]); err != nil {
		return
	} else if texInfos, err = q3ReadTexInfo(data[q3LumpTextureInformation]); err != nil {
		return
	} else if meshVerts, err = q3ReadMeshVertices(data[q3LumpMeshVertices]); err != nil {
		return
	}

	fb := lumps[q3LumpFaces].Data(q3HeaderLen, b)
	h := sliceHeader(&fb)
	h.Len = int(lumps[q3LumpFaces].Size / 104)
	h.Cap = h.Len
	faces = *(*[]q3Face)(unsafe.Pointer(&h))

	// FIXME -- only polygons for now
	out = make([]Face, 0)
	for _, face := range faces {
		var v []Vert

		if face.Type == q3FacePolygon || face.Type == q3FaceMesh {
			v = make([]Vert, 0)
			for i := uint32(0); i < face.NumMeshVerts; i++ {
				v = append(v, Vert{
					Pos: verts[face.FirstVertex+meshVerts[face.FirstMeshVert+i]],
				})
			}
		} else {
			continue
		}
		out = append(out, Face{
			Verts:   v,
			Plane:   nil,
			TexInfo: TexInfo{},
		})
	}

	_ = planes
	_ = texInfos
	_ = meshVerts

	return
}
