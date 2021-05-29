package corto

import (
	"math"
	"reflect"
	"unsafe"
)

type ColorAttr struct {
	GenericAttr
	qc             [4]int
	out_components int
}

func NewColorAttr(components *int) *ColorAttr {
	c := 4
	if components != nil {
		c = *components
	}
	return &ColorAttr{qc: [4]int{4, 4, 4, 8}, GenericAttr: GenericAttr{N: c, values: make([]byte, 0), diffs: make([]byte, 0)}}
}

func (a *ColorAttr) Codec() CodecType {
	return COLOR_CODEC
}

func (a *ColorAttr) Quantize(nvert uint32, buffer []byte) {
	n := a.N * int(nvert)

	values := make([]byte, n)
	diffs := make([]byte, n)

	switch a.format {
	case FORMAT_UINT8:
		bpos := 0
		vpos := 0
		var y Color4b
		for i := 0; i < int(nvert); i++ {
			for k := 0; k < a.N; k++ {
				y[k] = byte(int(buffer[bpos:][k]) / a.qc[k])
			}
			y = y.toYCC()
			for k := 0; k < a.N; k++ {
				values[vpos:][k] = y[k]
			}
			bpos += a.N
			vpos += a.N
		}
		break
	case FORMAT_FLOAT:
		var y Color4b
		y[3] = 255
		bpos := 0
		vpos := 0
		for i := 0; i < int(nvert); i++ {
			for k := 0; k < a.N; k++ {
				f := math.Float32frombits(byteorder.Uint32(buffer[bpos:]))
				y[k] = byte((f * 255.0) / float32(a.qc[k]))
			}
			y = y.toYCC()
			for k := 0; k < a.N; k++ {
				values[vpos:][k] = y[k]
			}
			bpos += a.N * 4
			vpos += a.N
		}
		break
	}
	a.bits = 0

	a.values = values
	a.diffs = diffs
}

func (a *ColorAttr) Dequantize(nvert uint32) {
	if a.buffer == nil {
		return
	}
	n := a.N * int(nvert)

	switch a.format {
	case FORMAT_UINT8:
		bpos := n
		tpos := a.out_components * int(nvert)

		var color Color4b
		color[3] = 255

		for bpos > 0 {
			bpos -= a.N
			tpos -= a.out_components

			for k := 0; k < a.N; k++ {
				color[k] = a.buffer[bpos:][k]
			}
			color = color.toRGB()
			for k := 0; k < a.out_components; k++ {
				a.buffer[tpos:][k] = byte(int(color[k]) * a.qc[k])
			}
		}
		break
	case FORMAT_FLOAT:
		var colorsSlice []Color4b
		colorsHeader := (*reflect.SliceHeader)((unsafe.Pointer(&colorsSlice)))
		colorsHeader.Cap = int(nvert)
		colorsHeader.Len = int(nvert)
		colorsHeader.Data = uintptr(unsafe.Pointer(&a.buffer[0]))

		var floatSlice []float32
		floatHeader := (*reflect.SliceHeader)((unsafe.Pointer(&floatSlice)))
		floatHeader.Cap = int(nvert)
		floatHeader.Len = int(nvert)
		floatHeader.Data = uintptr(unsafe.Pointer(&a.buffer[0]))

		for i := 0; i < int(nvert); i++ {
			rgb := &colorsSlice[i]
			*rgb = rgb.toRGB()
			for k := 0; k < a.out_components; k++ {
				floatSlice[i:][k] = (floatSlice[i:][k] * float32(a.qc[k])) / 255.0
			}
		}
		break
	}
}

func (a *ColorAttr) setQ(r_bits, g_bits, b_bits, a_bits int) {
	a.qc[0] = 1 << (8 - r_bits)
	a.qc[1] = 1 << (8 - g_bits)
	a.qc[2] = 1 << (8 - b_bits)
	a.qc[3] = 1 << (8 - a_bits)
}

func (a *ColorAttr) Encode(nvert uint32, stream *OutStream) {
	stream.restart()
	for c := 0; c < a.N; c++ {
		stream.write(byte(a.qc[c]))
	}

	stream.encodeValues(nvert, a.diffs.([]byte), a.N)
	a.size = stream.elapsed()
}

func (a *ColorAttr) Decode(nvert uint32, stream *InStream) {
	for c := 0; c < a.N; c++ {
		a.qc[c] = int(stream.readUint8())
	}
	stream.decodeValues(a.buffer, a.N)
}
