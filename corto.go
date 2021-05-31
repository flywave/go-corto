package corto

import (
	"math"
	"reflect"
	"unsafe"

	"github.com/flywave/go3d/vec2"
	"github.com/flywave/go3d/vec3"
)

type DecoderContext struct {
	NFace            uint32
	NVert            uint32
	ColorsComponents uint32
	Index16          bool
	Normal16         bool
}

func NewDecoderContext(nface, nvert, colors int, index16, normal16 bool) *DecoderContext {
	return &DecoderContext{NFace: uint32(nface), NVert: uint32(nvert), ColorsComponents: uint32(colors), Index16: index16, Normal16: normal16}
}

type Face [3]uint32
type Face16 [3]uint16
type Normal16 [3]int16

type Color [4]byte
type Color3 [3]byte

type EncoderContext struct {
	Entropy    EntropyType
	VertexQ    float32
	VertexBits int
	NormBits   int
	ColorBits  [4]int
	UvBits     float32
}

func NewEncoderContext(coordQ int) *EncoderContext {
	return &EncoderContext{Entropy: ENTROPY_TUNSTALL, VertexQ: float32(math.Pow(2, float64(coordQ))), VertexBits: 0, NormBits: 10, ColorBits: [4]int{6, 7, 6, 5}, UvBits: 12}
}

type Geom struct {
	Groups    []int
	Vertices  []vec3.T
	Indices   []Face
	Indices16 []Face16
	Normals   []vec3.T
	Normals16 []Normal16
	Colors    []Color
	Colors3   []Color3
	TexCoord  []vec2.T
}

func (m *Geom) MakeVerticesPlane(nvert int) []float32 {
	m.Vertices = make([]vec3.T, nvert)

	var floatSlice []float32
	floatHeader := (*reflect.SliceHeader)((unsafe.Pointer(&floatSlice)))
	floatHeader.Cap = int(nvert * 3)
	floatHeader.Len = int(nvert * 3)
	floatHeader.Data = uintptr(unsafe.Pointer(&m.Vertices[0]))

	return floatSlice
}

func (m *Geom) GetVerticesPlane() []float32 {
	nvert := len(m.Vertices)

	var floatSlice []float32
	floatHeader := (*reflect.SliceHeader)((unsafe.Pointer(&floatSlice)))
	floatHeader.Cap = int(nvert * 3)
	floatHeader.Len = int(nvert * 3)
	floatHeader.Data = uintptr(unsafe.Pointer(&m.Vertices[0]))

	return floatSlice
}

func (m *Geom) MakeNormalsPlane(nvert int) []float32 {
	m.Vertices = make([]vec3.T, nvert)

	var floatSlice []float32
	floatHeader := (*reflect.SliceHeader)((unsafe.Pointer(&floatSlice)))
	floatHeader.Cap = int(nvert * 3)
	floatHeader.Len = int(nvert * 3)
	floatHeader.Data = uintptr(unsafe.Pointer(&m.Normals[0]))

	return floatSlice
}

func (m *Geom) GetNormalsPlane() []float32 {
	nvert := len(m.Normals)

	var floatSlice []float32
	floatHeader := (*reflect.SliceHeader)((unsafe.Pointer(&floatSlice)))
	floatHeader.Cap = int(nvert * 3)
	floatHeader.Len = int(nvert * 3)
	floatHeader.Data = uintptr(unsafe.Pointer(&m.Normals[0]))

	return floatSlice
}

func (m *Geom) MakeNormals16Plane(nvert int) []int16 {
	m.Normals16 = make([]Normal16, nvert)

	var int16Slice []int16
	int16Header := (*reflect.SliceHeader)((unsafe.Pointer(&int16Slice)))
	int16Header.Cap = int(nvert * 3)
	int16Header.Len = int(nvert * 3)
	int16Header.Data = uintptr(unsafe.Pointer(&m.Normals16[0]))

	return int16Slice
}

func (m *Geom) GetNormals16Plane() []int16 {
	nvert := len(m.Normals16)

	var int16Slice []int16
	int16Header := (*reflect.SliceHeader)((unsafe.Pointer(&int16Slice)))
	int16Header.Cap = int(nvert * 3)
	int16Header.Len = int(nvert * 3)
	int16Header.Data = uintptr(unsafe.Pointer(&m.Normals16[0]))

	return int16Slice
}

func (m *Geom) MakeColors3Plane(nvert int) []byte {
	m.Colors3 = make([]Color3, nvert)

	var byteSlice []byte
	byteHeader := (*reflect.SliceHeader)((unsafe.Pointer(&byteSlice)))
	byteHeader.Cap = int(nvert * 3)
	byteHeader.Len = int(nvert * 3)
	byteHeader.Data = uintptr(unsafe.Pointer(&m.Colors3[0]))

	return byteSlice
}

func (m *Geom) GetColors3Plane() []byte {
	nvert := len(m.Colors3)

	var byteSlice []byte
	byteHeader := (*reflect.SliceHeader)((unsafe.Pointer(&byteSlice)))
	byteHeader.Cap = int(nvert * 3)
	byteHeader.Len = int(nvert * 3)
	byteHeader.Data = uintptr(unsafe.Pointer(&m.Colors3[0]))

	return byteSlice
}

func (m *Geom) MakeColorsPlane(nvert int) []byte {
	m.Colors = make([]Color, nvert)

	var byteSlice []byte
	byteHeader := (*reflect.SliceHeader)((unsafe.Pointer(&byteSlice)))
	byteHeader.Cap = int(nvert * 4)
	byteHeader.Len = int(nvert * 4)
	byteHeader.Data = uintptr(unsafe.Pointer(&m.Colors3[0]))

	return byteSlice
}

func (m *Geom) GetColorsPlane() []byte {
	nvert := len(m.Colors)

	var byteSlice []byte
	byteHeader := (*reflect.SliceHeader)((unsafe.Pointer(&byteSlice)))
	byteHeader.Cap = int(nvert * 4)
	byteHeader.Len = int(nvert * 4)
	byteHeader.Data = uintptr(unsafe.Pointer(&m.Colors3[0]))

	return byteSlice
}

func (m *Geom) MakeTexCoordPlane(nvert int) []float32 {
	m.TexCoord = make([]vec2.T, nvert)

	var floatSlice []float32
	floatHeader := (*reflect.SliceHeader)((unsafe.Pointer(&floatSlice)))
	floatHeader.Cap = int(nvert * 2)
	floatHeader.Len = int(nvert * 2)
	floatHeader.Data = uintptr(unsafe.Pointer(&m.TexCoord[0]))

	return floatSlice
}

func (m *Geom) GetTexCoordPlane() []float32 {
	nvert := len(m.TexCoord)

	var floatSlice []float32
	floatHeader := (*reflect.SliceHeader)((unsafe.Pointer(&floatSlice)))
	floatHeader.Cap = int(nvert * 2)
	floatHeader.Len = int(nvert * 2)
	floatHeader.Data = uintptr(unsafe.Pointer(&m.Colors3[0]))

	return floatSlice
}

func (m *Geom) MakeIndex16Plane(nface int) []uint16 {
	m.Indices16 = make([]Face16, nface)

	var uint16Slice []uint16
	uint16Header := (*reflect.SliceHeader)((unsafe.Pointer(&uint16Slice)))
	uint16Header.Cap = int(nface * 3)
	uint16Header.Len = int(nface * 3)
	uint16Header.Data = uintptr(unsafe.Pointer(&m.Indices16[0]))

	return uint16Slice
}

func (m *Geom) GetIndex16Plane() []uint16 {
	nface := len(m.Indices16)

	var uint16Slice []uint16
	uint16Header := (*reflect.SliceHeader)((unsafe.Pointer(&uint16Slice)))
	uint16Header.Cap = int(nface * 3)
	uint16Header.Len = int(nface * 3)
	uint16Header.Data = uintptr(unsafe.Pointer(&m.Indices16[0]))

	return uint16Slice
}

func (m *Geom) MakeIndex32Plane(nface int) []uint32 {
	m.Indices = make([]Face, nface)

	var uint32Slice []uint32
	uint32Header := (*reflect.SliceHeader)((unsafe.Pointer(&uint32Slice)))
	uint32Header.Cap = int(nface * 3)
	uint32Header.Len = int(nface * 3)
	uint32Header.Data = uintptr(unsafe.Pointer(&m.Indices[0]))

	return uint32Slice
}

func (m *Geom) GetIndex32Plane() []uint32 {
	nface := len(m.Indices)

	var uint32Slice []uint32
	uint32Header := (*reflect.SliceHeader)((unsafe.Pointer(&uint32Slice)))
	uint32Header.Cap = int(nface * 3)
	uint32Header.Len = int(nface * 3)
	uint32Header.Data = uintptr(unsafe.Pointer(&m.Indices16[0]))

	return uint32Slice
}

func (m *Geom) HasGroup() bool {
	return len(m.Groups) > 0
}

func (m *Geom) HasFace() bool {
	return len(m.Indices) > 0
}

func (m *Geom) HasFace16() bool {
	return len(m.Indices16) > 0
}

func (m *Geom) HasNormal() bool {
	return len(m.Normals) > 0
}

func (m *Geom) HasNormal16() bool {
	return len(m.Normals16) > 0
}

func (m *Geom) HasColor() bool {
	return len(m.Colors) > 0
}

func (m *Geom) HasColor3() bool {
	return len(m.Colors3) > 0
}

func (m *Geom) HasTexCoord() bool {
	return len(m.TexCoord) > 0
}

func DecodeGeom(ctx *DecoderContext, input []byte) *Geom {
	nvert := uint32(ctx.NVert)
	nface := uint32(ctx.NFace)

	geom := &Geom{}

	dec := NewDecoder(input)
	defer dec.Free()

	dec.SetPositions(geom.MakeVerticesPlane(int(nvert)))

	if dec.HasAttr("normal") {
		if ctx.Normal16 {
			dec.SetNormalsInt16(geom.MakeNormals16Plane(int(nvert)))
		} else {
			dec.SetNormals(geom.MakeNormalsPlane(int(nvert)))
		}
	}

	if dec.HasAttr("color") {
		if ctx.ColorsComponents == 3 {
			dec.SetColors(geom.MakeColors3Plane(int(nvert)), int(ctx.ColorsComponents))
		} else {
			dec.SetColors(geom.MakeColorsPlane(int(nvert)), int(ctx.ColorsComponents))
		}
	}

	if dec.HasAttr("uv") {
		dec.SetUvs(geom.MakeTexCoordPlane(int(nvert)))
	}

	if nface > 0 {
		if ctx.Index16 {
			dec.SetIndexInt16(geom.MakeIndex16Plane(int(nface)))
		} else {
			dec.SetIndexInt32(geom.MakeIndex32Plane(int(nface)))
		}
	}
	dec.Decode()
	return geom
}

func EncodeGeom(ctx *EncoderContext, geom *Geom) []byte {
	nvert := uint32(len(geom.Vertices))
	nface := uint32(len(geom.Indices))

	enc := NewEncoder(nvert, nface, ctx.Entropy)
	defer enc.Free()

	if geom.HasGroup() {
		for i := range geom.Groups {
			enc.AddGroup(geom.Groups[i])
		}
	}

	if geom.HasFace() {
		if geom.HasFace16() {
			if ctx.VertexBits > 0 {
				enc.AddPositionsBitsInt16(geom.GetVerticesPlane(), geom.GetIndex16Plane(), ctx.VertexBits)
			} else {
				enc.AddPositionsInt16(geom.GetVerticesPlane(), geom.GetIndex16Plane(), ctx.VertexQ, vec3.Zero)
			}
		} else {
			if ctx.VertexBits > 0 {
				enc.AddPositionsBitsInt32(geom.GetVerticesPlane(), geom.GetIndex32Plane(), ctx.VertexBits)
			} else {
				enc.AddPositionsInt32(geom.GetVerticesPlane(), geom.GetIndex32Plane(), ctx.VertexQ, vec3.Zero)
			}
		}
	} else {
		if ctx.VertexBits > 0 {
			enc.AddPositionsBits(geom.GetVerticesPlane(), ctx.VertexBits)
		} else {
			enc.AddPositions(geom.GetVerticesPlane(), ctx.VertexQ, vec3.Zero)
		}
	}

	if geom.HasNormal() {
		if geom.HasFace() {
			enc.AddNormals(geom.GetNormalsPlane(), ctx.NormBits, PREDICTION_ESTIMATED)
		} else {
			enc.AddNormals(geom.GetNormalsPlane(), ctx.NormBits, PREDICTION_DIFF)
		}
	} else if geom.HasNormal16() {
		if geom.HasFace() {
			enc.AddNormalsInt16(geom.GetNormals16Plane(), ctx.NormBits, PREDICTION_ESTIMATED)
		} else {
			enc.AddNormalsInt16(geom.GetNormals16Plane(), ctx.NormBits, PREDICTION_DIFF)
		}
	}

	if geom.HasColor() {
		enc.AddColors(geom.GetColorsPlane(), ctx.ColorBits[:])
	} else if geom.HasColor3() {
		enc.AddColors3(geom.GetColors3Plane(), ctx.ColorBits[:])
	}

	if geom.HasTexCoord() {
		enc.AddUvs(geom.GetTexCoordPlane(), ctx.UvBits)
	}

	return enc.Encode()
}
