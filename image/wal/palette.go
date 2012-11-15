package wal

import (
	c "image/color"
)

// DefaultPalette is used as a default palette for WAL images.
var DefaultPalette = c.Palette{
	c.NRGBA{0x00, 0x00, 0x00, 0xff}, c.NRGBA{0x1f, 0x1f, 0x1f, 0xff},
	c.NRGBA{0x3f, 0x3f, 0x3f, 0xff}, c.NRGBA{0x5b, 0x5b, 0x5b, 0xff},
	c.NRGBA{0x7b, 0x7b, 0x7b, 0xff}, c.NRGBA{0x9b, 0x9b, 0x9b, 0xff},
	c.NRGBA{0xbb, 0xbb, 0xbb, 0xff}, c.NRGBA{0xdb, 0xdb, 0xdb, 0xff},
	c.NRGBA{0xff, 0xff, 0xff, 0xff}, c.NRGBA{0xff, 0xff, 0xff, 0xff},
	c.NRGBA{0xff, 0xff, 0xff, 0xff}, c.NRGBA{0xff, 0xff, 0xff, 0xff},
	c.NRGBA{0xff, 0xff, 0xff, 0xff}, c.NRGBA{0xff, 0xff, 0xff, 0xff},
	c.NRGBA{0xff, 0xff, 0xff, 0xff}, c.NRGBA{0xff, 0xff, 0xff, 0xff},
	c.NRGBA{0xcf, 0x97, 0x4b, 0xff}, c.NRGBA{0xa7, 0x7b, 0x3b, 0xff},
	c.NRGBA{0xa7, 0x7b, 0x3b, 0xff}, c.NRGBA{0xa7, 0x7b, 0x3b, 0xff},
	c.NRGBA{0x8b, 0x67, 0x2f, 0xff}, c.NRGBA{0x8b, 0x67, 0x2f, 0xff},
	c.NRGBA{0x6f, 0x53, 0x27, 0xff}, c.NRGBA{0x63, 0x4b, 0x23, 0xff},
	c.NRGBA{0x63, 0x4b, 0x23, 0xff}, c.NRGBA{0x53, 0x3f, 0x1f, 0xff},
	c.NRGBA{0x4f, 0x3b, 0x1b, 0xff}, c.NRGBA{0x43, 0x2b, 0x17, 0xff},
	c.NRGBA{0x33, 0x27, 0x13, 0xff}, c.NRGBA{0x2b, 0x1f, 0x13, 0xff},
	c.NRGBA{0x27, 0x1b, 0x0f, 0xff}, c.NRGBA{0x1f, 0x17, 0x0f, 0xff},
	c.NRGBA{0xb3, 0xc7, 0xd3, 0xff}, c.NRGBA{0xb3, 0xc7, 0xd3, 0xff},
	c.NRGBA{0xbb, 0xbb, 0xbb, 0xff}, c.NRGBA{0xab, 0xab, 0xab, 0xff},
	c.NRGBA{0x9b, 0x9b, 0x9b, 0xff}, c.NRGBA{0x9b, 0x9b, 0x9b, 0xff},
	c.NRGBA{0x8b, 0x8b, 0x8b, 0xff}, c.NRGBA{0x7b, 0x7b, 0x7b, 0xff},
	c.NRGBA{0x6b, 0x6b, 0x6b, 0xff}, c.NRGBA{0x5b, 0x5b, 0x5b, 0xff},
	c.NRGBA{0x5b, 0x5b, 0x5b, 0xff}, c.NRGBA{0x4b, 0x4b, 0x4b, 0xff},
	c.NRGBA{0x47, 0x3f, 0x43, 0xff}, c.NRGBA{0x3b, 0x37, 0x37, 0xff},
	c.NRGBA{0x2f, 0x2f, 0x2f, 0xff}, c.NRGBA{0x27, 0x27, 0x27, 0xff},
	c.NRGBA{0xff, 0xff, 0xa7, 0xff}, c.NRGBA{0xeb, 0x97, 0x7f, 0xff},
	c.NRGBA{0xeb, 0x97, 0x7f, 0xff}, c.NRGBA{0xcf, 0x97, 0x4b, 0xff},
	c.NRGBA{0xff, 0xff, 0xa7, 0xff}, c.NRGBA{0xff, 0xff, 0x7f, 0xff},
	c.NRGBA{0xff, 0xff, 0x53, 0xff}, c.NRGBA{0xcf, 0x97, 0x4b, 0xff},
	c.NRGBA{0xff, 0xff, 0x53, 0xff}, c.NRGBA{0xff, 0xff, 0x53, 0xff},
	c.NRGBA{0xff, 0xff, 0x53, 0xff}, c.NRGBA{0xff, 0xd7, 0x17, 0xff},
	c.NRGBA{0xeb, 0x9f, 0x27, 0xff}, c.NRGBA{0xaf, 0x77, 0x1f, 0xff},
	c.NRGBA{0x77, 0x4f, 0x17, 0xff}, c.NRGBA{0x43, 0x2b, 0x17, 0xff},
	c.NRGBA{0xeb, 0x97, 0x7f, 0xff}, c.NRGBA{0xff, 0x93, 0x00, 0xff},
	c.NRGBA{0xef, 0x7f, 0x00, 0xff}, c.NRGBA{0xe3, 0x6b, 0x00, 0xff},
	c.NRGBA{0xd3, 0x57, 0x00, 0xff}, c.NRGBA{0xc7, 0x47, 0x00, 0xff},
	c.NRGBA{0xc7, 0x47, 0x00, 0xff}, c.NRGBA{0xab, 0x2b, 0x00, 0xff},
	c.NRGBA{0x9b, 0x1f, 0x00, 0xff}, c.NRGBA{0x8f, 0x17, 0x00, 0xff},
	c.NRGBA{0x73, 0x17, 0x0b, 0xff}, c.NRGBA{0x67, 0x17, 0x07, 0xff},
	c.NRGBA{0x57, 0x13, 0x00, 0xff}, c.NRGBA{0x43, 0x0f, 0x00, 0xff},
	c.NRGBA{0x33, 0x0b, 0x00, 0xff}, c.NRGBA{0x23, 0x0b, 0x00, 0xff},
	c.NRGBA{0xd7, 0xbb, 0xb7, 0xff}, c.NRGBA{0xeb, 0x97, 0x7f, 0xff},
	c.NRGBA{0xeb, 0x97, 0x7f, 0xff}, c.NRGBA{0xcb, 0x9b, 0x93, 0xff},
	c.NRGBA{0xbf, 0x7b, 0x6f, 0xff}, c.NRGBA{0xa7, 0x8b, 0x77, 0xff},
	c.NRGBA{0x8f, 0x77, 0x53, 0xff}, c.NRGBA{0x8f, 0x77, 0x53, 0xff},
	c.NRGBA{0x87, 0x6b, 0x57, 0xff}, c.NRGBA{0x7b, 0x5f, 0x4b, 0xff},
	c.NRGBA{0x67, 0x4f, 0x3b, 0xff}, c.NRGBA{0x5f, 0x47, 0x37, 0xff},
	c.NRGBA{0x4b, 0x37, 0x2b, 0xff}, c.NRGBA{0x3f, 0x2f, 0x23, 0xff},
	c.NRGBA{0x2b, 0x1f, 0x13, 0xff}, c.NRGBA{0x1f, 0x17, 0x0f, 0xff},
	c.NRGBA{0xcb, 0x8b, 0x23, 0xff}, c.NRGBA{0xaf, 0x77, 0x1f, 0xff},
	c.NRGBA{0x9f, 0x57, 0x33, 0xff}, c.NRGBA{0x8b, 0x67, 0x2f, 0xff},
	c.NRGBA{0x63, 0x4b, 0x23, 0xff}, c.NRGBA{0x4f, 0x3b, 0x1b, 0xff},
	c.NRGBA{0x33, 0x27, 0x13, 0xff}, c.NRGBA{0x1f, 0x17, 0x0f, 0xff},
	c.NRGBA{0xff, 0xff, 0xa7, 0xff}, c.NRGBA{0xff, 0xff, 0xd3, 0xff},
	c.NRGBA{0xff, 0xff, 0xff, 0xff}, c.NRGBA{0xff, 0xff, 0xff, 0xff},
	c.NRGBA{0xff, 0xff, 0xff, 0xff}, c.NRGBA{0xff, 0xff, 0xff, 0xff},
	c.NRGBA{0xff, 0xff, 0xff, 0xff}, c.NRGBA{0xff, 0xff, 0xff, 0xff},
	c.NRGBA{0xff, 0xff, 0xff, 0xff}, c.NRGBA{0xff, 0xff, 0xff, 0xff},
	c.NRGBA{0xcb, 0xd7, 0xdf, 0xff}, c.NRGBA{0x9f, 0xb7, 0xc3, 0xff},
	c.NRGBA{0x77, 0x7b, 0xcf, 0xff}, c.NRGBA{0x5b, 0x87, 0x9b, 0xff},
	c.NRGBA{0x47, 0x77, 0x8b, 0xff}, c.NRGBA{0x2f, 0x67, 0x7f, 0xff},
	c.NRGBA{0x2f, 0x67, 0x7f, 0xff}, c.NRGBA{0x17, 0x53, 0x6f, 0xff},
	c.NRGBA{0x13, 0x4b, 0x67, 0xff}, c.NRGBA{0x0b, 0x3f, 0x53, 0xff},
	c.NRGBA{0x07, 0x2f, 0x3f, 0xff}, c.NRGBA{0x00, 0x1f, 0x2b, 0xff},
	c.NRGBA{0x00, 0x0f, 0x13, 0xff}, c.NRGBA{0x00, 0x00, 0x00, 0xff},
	c.NRGBA{0xeb, 0xd3, 0xc7, 0xff}, c.NRGBA{0xeb, 0x97, 0x7f, 0xff},
	c.NRGBA{0xeb, 0x97, 0x7f, 0xff}, c.NRGBA{0xeb, 0x97, 0x7f, 0xff},
	c.NRGBA{0xbf, 0x7b, 0x6f, 0xff}, c.NRGBA{0xc3, 0x73, 0x53, 0xff},
	c.NRGBA{0xb3, 0x5b, 0x4f, 0xff}, c.NRGBA{0xb3, 0x5b, 0x4f, 0xff},
	c.NRGBA{0x9f, 0x4b, 0x3f, 0xff}, c.NRGBA{0x7b, 0x47, 0x47, 0xff},
	c.NRGBA{0x63, 0x33, 0x33, 0xff}, c.NRGBA{0x57, 0x2b, 0x2b, 0xff},
	c.NRGBA{0x3f, 0x1f, 0x1f, 0xff}, c.NRGBA{0x27, 0x1b, 0x13, 0xff},
	c.NRGBA{0x17, 0x0f, 0x0b, 0xff}, c.NRGBA{0x00, 0x00, 0x00, 0xff},
	c.NRGBA{0xff, 0xff, 0xff, 0xff}, c.NRGBA{0xff, 0xff, 0xd3, 0xff},
	c.NRGBA{0xff, 0xff, 0xd3, 0xff}, c.NRGBA{0xff, 0xff, 0xd3, 0xff},
	c.NRGBA{0xff, 0xff, 0xd3, 0xff}, c.NRGBA{0xeb, 0xd3, 0xc7, 0xff},
	c.NRGBA{0xd7, 0xbb, 0xb7, 0xff}, c.NRGBA{0xc7, 0xab, 0x9b, 0xff},
	c.NRGBA{0xc7, 0xab, 0x9b, 0xff}, c.NRGBA{0x97, 0x9f, 0x7b, 0xff},
	c.NRGBA{0x87, 0x8b, 0x6b, 0xff}, c.NRGBA{0x73, 0x73, 0x57, 0xff},
	c.NRGBA{0x5b, 0x5b, 0x43, 0xff}, c.NRGBA{0x43, 0x43, 0x33, 0xff},
	c.NRGBA{0x2f, 0x2f, 0x23, 0xff}, c.NRGBA{0x1b, 0x1b, 0x17, 0xff},
	c.NRGBA{0xeb, 0x97, 0x7f, 0xff}, c.NRGBA{0xeb, 0x97, 0x7f, 0xff},
	c.NRGBA{0xeb, 0x97, 0x7f, 0xff}, c.NRGBA{0xc3, 0x73, 0x53, 0xff},
	c.NRGBA{0xc3, 0x73, 0x53, 0xff}, c.NRGBA{0xb3, 0x5b, 0x4f, 0xff},
	c.NRGBA{0xa7, 0x3b, 0x2b, 0xff}, c.NRGBA{0xa7, 0x3b, 0x2b, 0xff},
	c.NRGBA{0x9f, 0x2f, 0x23, 0xff}, c.NRGBA{0x8b, 0x27, 0x13, 0xff},
	c.NRGBA{0x6b, 0x2b, 0x1b, 0xff}, c.NRGBA{0x57, 0x1f, 0x13, 0xff},
	c.NRGBA{0x43, 0x17, 0x0b, 0xff}, c.NRGBA{0x2b, 0x0b, 0x00, 0xff},
	c.NRGBA{0x1b, 0x00, 0x00, 0xff}, c.NRGBA{0x00, 0x00, 0x00, 0xff},
	c.NRGBA{0xff, 0xff, 0xff, 0xff}, c.NRGBA{0xff, 0xff, 0xff, 0xff},
	c.NRGBA{0xff, 0xff, 0xff, 0xff}, c.NRGBA{0xeb, 0xeb, 0xeb, 0xff},
	c.NRGBA{0xcb, 0xd7, 0xdf, 0xff}, c.NRGBA{0xb3, 0xc7, 0xd3, 0xff},
	c.NRGBA{0x9f, 0xb7, 0xc3, 0xff}, c.NRGBA{0x77, 0x7b, 0xcf, 0xff},
	c.NRGBA{0x77, 0x7b, 0xcf, 0xff}, c.NRGBA{0x67, 0x6b, 0xb7, 0xff},
	c.NRGBA{0x5b, 0x5b, 0x9b, 0xff}, c.NRGBA{0x4b, 0x4f, 0x7f, 0xff},
	c.NRGBA{0x3f, 0x3f, 0x67, 0xff}, c.NRGBA{0x2f, 0x2f, 0x4b, 0xff},
	c.NRGBA{0x23, 0x1f, 0x2f, 0xff}, c.NRGBA{0x17, 0x0f, 0x0b, 0xff},
	c.NRGBA{0xff, 0xff, 0xff, 0xff}, c.NRGBA{0xff, 0xff, 0xd3, 0xff},
	c.NRGBA{0xff, 0xff, 0xd3, 0xff}, c.NRGBA{0xff, 0xff, 0xa7, 0xff},
	c.NRGBA{0xff, 0xff, 0xa7, 0xff}, c.NRGBA{0xff, 0xff, 0x7f, 0xff},
	c.NRGBA{0x9b, 0xab, 0x7b, 0xff}, c.NRGBA{0x9b, 0xab, 0x7b, 0xff},
	c.NRGBA{0x87, 0x97, 0x63, 0xff}, c.NRGBA{0x5f, 0xa7, 0x2f, 0xff},
	c.NRGBA{0x5f, 0x8f, 0x33, 0xff}, c.NRGBA{0x5f, 0x7b, 0x33, 0xff},
	c.NRGBA{0x3f, 0x4f, 0x1b, 0xff}, c.NRGBA{0x2f, 0x3b, 0x0b, 0xff},
	c.NRGBA{0x23, 0x2f, 0x07, 0xff}, c.NRGBA{0x1b, 0x23, 0x00, 0xff},
	c.NRGBA{0x00, 0xff, 0x00, 0xff}, c.NRGBA{0x00, 0xff, 0x00, 0xff},
	c.NRGBA{0xff, 0xff, 0x27, 0xff}, c.NRGBA{0xff, 0xff, 0x53, 0xff},
	c.NRGBA{0xff, 0xff, 0x53, 0xff}, c.NRGBA{0xff, 0xff, 0x53, 0xff},
	c.NRGBA{0xff, 0xff, 0x53, 0xff}, c.NRGBA{0xff, 0xff, 0xff, 0xff},
	c.NRGBA{0xff, 0xff, 0xff, 0xff}, c.NRGBA{0xff, 0xff, 0xff, 0xff},
	c.NRGBA{0xff, 0xff, 0xff, 0xff}, c.NRGBA{0xff, 0xff, 0xa7, 0xff},
	c.NRGBA{0xff, 0xff, 0x53, 0xff}, c.NRGBA{0xff, 0xff, 0x53, 0xff},
	c.NRGBA{0xff, 0xff, 0x27, 0xff}, c.NRGBA{0xff, 0xff, 0x27, 0xff},
	c.NRGBA{0xff, 0xff, 0x27, 0xff}, c.NRGBA{0xff, 0xff, 0x27, 0xff},
	c.NRGBA{0xff, 0xeb, 0x1f, 0xff}, c.NRGBA{0xff, 0xd7, 0x17, 0xff},
	c.NRGBA{0xff, 0xab, 0x07, 0xff}, c.NRGBA{0xff, 0x93, 0x00, 0xff},
	c.NRGBA{0xff, 0x93, 0x00, 0xff}, c.NRGBA{0xff, 0x93, 0x00, 0xff},
	c.NRGBA{0xff, 0x00, 0x00, 0xff}, c.NRGBA{0xff, 0x00, 0x00, 0xff},
	c.NRGBA{0xff, 0x00, 0x00, 0xff}, c.NRGBA{0xef, 0x00, 0x00, 0xff},
	c.NRGBA{0x9b, 0x1f, 0x00, 0xff}, c.NRGBA{0x7f, 0x0f, 0x00, 0xff},
	c.NRGBA{0x5f, 0x00, 0x00, 0xff}, c.NRGBA{0x2f, 0x00, 0x00, 0xff},
	c.NRGBA{0xff, 0x00, 0x00, 0xff}, c.NRGBA{0x37, 0x37, 0xff, 0xff},
	c.NRGBA{0xff, 0x00, 0x00, 0xff}, c.NRGBA{0x00, 0x00, 0xff, 0xff},
	c.NRGBA{0x5b, 0x5b, 0x43, 0xff}, c.NRGBA{0x37, 0x37, 0x2b, 0xff},
	c.NRGBA{0x23, 0x23, 0x1b, 0xff}, c.NRGBA{0xff, 0xff, 0xff, 0xff},
	c.NRGBA{0xff, 0xff, 0xa7, 0xff}, c.NRGBA{0xeb, 0x97, 0x7f, 0xff},
	c.NRGBA{0xeb, 0x9f, 0x27, 0xff}, c.NRGBA{0xff, 0xff, 0xff, 0xff},
	c.NRGBA{0xff, 0xff, 0xff, 0xff}, c.NRGBA{0xff, 0xff, 0xff, 0xff},
	c.NRGBA{0xeb, 0xd3, 0xc7, 0xff}, c.NRGBA{0x9f, 0x5b, 0x53, 0x00},
}
