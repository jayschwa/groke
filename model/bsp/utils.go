package bsp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"
)

type bspLump struct {
	Offset uint32
	Size   uint32
}

func (lump bspLump) String() string {
	return fmt.Sprintf("{0x%x, 0x%x}", lump.Offset, lump.Size)
}

func (lump bspLump) Data(headerLen uint32, b []byte) []byte {
	if lump.Offset == 0 || lump.Size == 0 {
		return []byte{}
	}
	offset := int(lump.Offset - headerLen)
	if len(b) < offset+int(lump.Size) {
		panic(fmt.Sprintf("lump data out of range (offset=%x, size=%x, requested [%x:%x])", lump.Offset, lump.Size, offset, offset+int(lump.Size)))
	}
	return b[offset : offset+int(lump.Size)]
}

func bspLumpsFrom(b []byte, numLumps int) []bspLump {
	h := sliceHeader(&b)
	h.Len = numLumps
	h.Cap = numLumps
	return *(*[]bspLump)(unsafe.Pointer(&h))
}

func bspReadEntities(b []byte) (ents []Entity, err error) {
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

					if valueIndex == 0 {
						ent[key] = ""
					} else {
						ent[key] = stringFrom(b[i : i+valueIndex])
					}

					i += valueIndex + 1
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

func qVector3(v [3]float32) Vector3 {
	return Vector3{float64(v[0]), float64(v[1]), float64(v[2])}
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

func Uint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}
