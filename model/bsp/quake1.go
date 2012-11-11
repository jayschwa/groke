package bsp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	_ = int8(-iota)
	bspQ1ContentsEmpty
	bspQ1ContentsSolid
	bspQ1ContentsWater
	bspQ1ContentsSlime
	bspQ1ContentsLava
	bspQ1ContentsSky
	bspQ1ContentsOrigin
	bspQ1ContentsClip
	bspQ1ContentsCurrent0
	bspQ1ContentsCurrent90
	bspQ1ContentsCurrent180
	bspQ1ContentsCurrent270
	bspQ1ContentsCurrentUp
	bspQ1ContentsCurrentDown
)

func q1BSPReadEntities(r *io.SectionReader, m *Model) error {
	b := bufio.NewReader(r)
	ent := make(Entity)

	for inBlock := 0; ; {
		if c, err := b.ReadByte(); err == io.EOF {
			break
		} else if err != nil {
			return err
		} else if c == '{' {
			inBlock++
		} else if c == '}' {
			if inBlock == 1 {
				m.Entities = append(m.Entities, ent)
				ent = make(Entity)
			}

			inBlock--
		} else if c == '"' && inBlock == 1 {
			key, err := b.ReadString('"')
			if err != nil {
				return err
			}

			for {
				if c, err := b.ReadByte(); err != nil {
					return err
				} else if c == ' ' || c == '\t' {
					continue
				} else if c == '"' {
					value, err := b.ReadString('"')
					if err != nil {
						return err
					}

					ent[key[:len(key)-1]] = value[:len(value)-1]
					break
				} else {
					err = fmt.Errorf("bsp: unexpected char %q", c)
				}
			}
		} else if c != ' ' && c != '\t' && c != '\r' && c != '\n' && c != 0 {
			return fmt.Errorf("bsp: unexpected char %q", c)
		}
	}

	return nil
}

func q1BSPReadPlanes(r *io.SectionReader, m *Model) error {
	numPlanes := int(r.Size() / 20)
	planes := make([]struct {
		N [3]float32
		D float32
		T uint32
	}, numPlanes)

	if err := binary.Read(r, binary.LittleEndian, &planes); err != nil {
		return err
	}

	m.Planes = make([]Plane, numPlanes)
	for i, p := range planes {
		m.Planes[i] = Plane{
			Vector3{float64(p.N[0]), float64(p.N[1]), float64(p.N[2])},
			float64(p.D),
			PlaneType(p.T),
		}
	}

	return nil
}

func q1BSPReadTextures(r *io.SectionReader, m *Model) error {
	var numTex uint32
	var headerOffsets []uint32

	if err := binary.Read(r, binary.LittleEndian, &numTex); err != nil {
		return err
	} else {
		headerOffsets = make([]uint32, int(numTex))
		if err := binary.Read(r, binary.LittleEndian, &headerOffsets); err != nil {
			return err
		}
		m.Textures = make([]Texture, 0, int(numTex))
	}

	h := make([]byte, 28)

	for i := 0; i < cap(m.Textures); i++ {
		if _, err := r.ReadAt(h, int64(headerOffsets[i])); err != nil {
			return err
		}

		nameLen := bytes.IndexByte(h[0:16], 0)
		if nameLen < 0 || nameLen > 16 {
			nameLen = 16
		}

		m.Textures = append(m.Textures, Texture{
			Name:       string(bytes.ToLower(h[0:nameLen])),
			Width:      int(binary.LittleEndian.Uint32(h[16:])),
			Height:     int(binary.LittleEndian.Uint32(h[20:])),
			dataOffset: int64(headerOffsets[i]) + int64(binary.LittleEndian.Uint32(h[24:])),
			r:          r,
		})
	}

	return nil
}

func q1BSPReadVertices(r *io.SectionReader, m *Model) error {
	numVerts := int(r.Size() / 12)
	m.Verts = make([]Vector3, numVerts)
	verts := make([][3]float32, numVerts)

	if err := binary.Read(r, binary.LittleEndian, &verts); err != nil {
		return err
	}

	for i := 0; i < numVerts; i++ {
		v := verts[i]
		m.Verts[i] = Vector3{float64(v[0]), float64(v[1]), float64(v[2])}
	}

	return nil
}

func q1BSPReadVisibility(r *io.SectionReader, m *Model) error {
	return nil
}

func q1BSPReadNodes(r *io.SectionReader, m *Model) error {
	return nil
}

func q1BSPReadTextureInformation(r *io.SectionReader, m *Model) error {
	numInfo := int(r.Size() / 40)
	m.TexInfos = make([]TexInfo, 0, numInfo)
	infos := make([]struct {
		S     [3]float32
		Ds    float32
		T     [3]float32
		Dt    float32
		TexID uint32
		Anim  uint32
	}, numInfo)

	if err := binary.Read(r, binary.LittleEndian, &infos); err != nil {
		return err
	}

	for i := 0; i < numInfo; i++ {
		var flags TexFlags
		info := infos[i]

		if info.Anim != 0 {
			flags |= TexAnimated
		}

		m.TexInfos = append(m.TexInfos, TexInfo{
			S:       Vector3{float64(info.S[0]), float64(info.S[1]), float64(info.S[2])},
			T:       Vector3{float64(info.T[0]), float64(info.T[1]), float64(info.T[2])},
			Ds:      float64(info.Ds),
			Dt:      float64(info.Dt),
			Texture: &m.Textures[int(info.TexID)],
			Flags:   flags,
		})
	}

	return nil
}

func q1BSPReadFaces(r *io.SectionReader, m *Model) error {
	numFaces := int(r.Size() / 20)
	m.Faces = make([]Face, 0, numFaces)
	fs := make([]struct {
		Plane     uint16
		Side      uint16
		Edge      uint32
		NumEdges  uint16
		TexInfoID uint16
		LightType uint8
		LightBase uint8
		Light     [2]uint8
		LightMap  uint32
	}, numFaces)

	if err := binary.Read(r, binary.LittleEndian, &fs); err != nil {
		return err
	}

	for i := 0; i < numFaces; i++ {
		f := fs[i]
		if int(f.Edge)+int(f.NumEdges) >= len(m.Edges) {
			return fmt.Errorf("face's edges out of range (%d:%d)", f.Edge, f.Edge+uint32(f.NumEdges))
		}

		m.Faces = append(m.Faces, Face{
			Edges:   m.Edges[int(f.Edge) : int(f.Edge)+int(f.NumEdges)],
			Front:   f.Side == 0,
			Plane:   &m.Planes[int(f.Plane)],
			TexInfo: &m.TexInfos[int(f.TexInfoID)],
		})
	}

	return nil
}

func q1BSPReadLightmaps(r *io.SectionReader, m *Model) error {
	return nil
}

func q1BSPReadClipNodes(r *io.SectionReader, m *Model) error {
	return nil
}

func q1BSPReadLeaves(r *io.SectionReader, m *Model) error {
	return nil
}

func q1BSPReadMarkSurfaces(r *io.SectionReader, m *Model) error {
	return nil
}

func q1BSPReadEdges(r *io.SectionReader, m *Model) error {
	numEdges := int(r.Size() / 4)
	m.EdgeVIndices = make([]EdgeVIndex, numEdges)

	for i := 0; i < numEdges; i++ {
		var v uint16
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return err
		}
		m.EdgeVIndices[i].Ai = int(v)

		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return err
		}
		m.EdgeVIndices[i].Bi = int(v)

		if m.EdgeVIndices[i].Ai >= len(m.Verts) || m.EdgeVIndices[i].Bi >= len(m.Verts) {
			return fmt.Errorf("edge %d verts out of range: %v", i, m.EdgeVIndices[i])
		}
	}

	return nil
}

func q1BSPReadFaceEdgeTables(r *io.SectionReader, m *Model) error {
	numEdges := int(r.Size() / 2)
	m.Edges = make([]Edge, numEdges)
	edges := make([]int16, numEdges)

	if err := binary.Read(r, binary.LittleEndian, &edges); err != nil {
		return err
	}

	for i := 0; i < numEdges; i++ {
		ei := int(edges[i])

		if ei >= len(m.EdgeVIndices) || -ei >= len(m.EdgeVIndices) {
			return fmt.Errorf("edge index out of range (%d)", ei)
		}

		if ei < 0 {
			m.Edges[i].A = m.Verts[m.EdgeVIndices[-ei].Bi]
			m.Edges[i].B = m.Verts[m.EdgeVIndices[-ei].Ai]
		} else {
			m.Edges[i].A = m.Verts[m.EdgeVIndices[ei].Ai]
			m.Edges[i].B = m.Verts[m.EdgeVIndices[ei].Bi]
		}
	}

	return nil
}

func q1BSPReadModels(r *io.SectionReader, m *Model) error {
	return nil
}
