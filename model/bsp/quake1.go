package bsp

import (
	"io"
	"io/ioutil"
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
	if entities, err := ioutil.ReadAll(r); err != nil {
		return err
	} else {
		m.Entities = string(entities)
	}

	return nil
}

func q1BSPReadPlanes(r *io.SectionReader, m *Model) error {
	var planesstruct {
		Normal   [3]float
		Distance float
		Type     uint32
	}

	return nil
}

func q1BSPReadTextures(r *io.SectionReader, m *Model) error {
	return nil
}

func q1BSPReadVertices(r *io.SectionReader, m *Model) error {
	return nil
}

func q1BSPReadVisibility(r *io.SectionReader, m *Model) error {
	return nil
}

func q1BSPReadNodes(r *io.SectionReader, m *Model) error {
	return nil
}

func q1BSPReadTextureInformation(r *io.SectionReader, m *Model) error {
	return nil
}

func q1BSPReadFaces(r *io.SectionReader, m *Model) error {
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
	return nil
}

func q1BSPReadFaceEdgeTables(r *io.SectionReader, m *Model) error {
	return nil
}

func q1BSPReadModels(r *io.SectionReader, m *Model) error {
	return nil
}
