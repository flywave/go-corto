package corto

import (
	"math"
	"reflect"
	"unsafe"

	"github.com/flywave/go3d/vec2"
	"github.com/flywave/go3d/vec3"
)

type PredictionType uint32

const (
	PREDICTION_DIFF      PredictionType = 0x0
	PREDICTION_ESTIMATED PredictionType = 0x1
	PREDICTION_BORDER    PredictionType = 0x2
)

type NormalAttr struct {
	VertexAttribute
	buffer   []byte
	N        int
	q        float32
	strategy StrategyType

	format FormatType
	size   uint32
	bits   int

	prediction PredictionType
	boundary   []int32
	values     []int32
	diffs      []int32
}

func NewNormalAttr(bits *int) *NormalAttr {
	cbits := 10
	if bits != nil {
		cbits = *bits
	}
	q := float32(math.Pow(2.0, float64(cbits-1)))
	return &NormalAttr{N: 3, q: q, prediction: PREDICTION_DIFF, strategy: 0 | CORRELATED}
}

func (a *NormalAttr) Codec() CodecType {
	return NORMAL_CODEC
}

func markBoundary(nvert uint32, nface uint32, index interface{}) (boundary []int32) {
	boundary = make([]int32, nvert)

	endpos := nface * 3
	for f := 0; f < int(endpos); f += 3 {
		f0 := GenericGetInt(index, f)
		f1 := GenericGetInt(index, f+1)
		f2 := GenericGetInt(index, f+2)

		boundary[f0] ^= int32(f1)
		boundary[f0] ^= int32(f2)
		boundary[f1] ^= int32(f2)
		boundary[f1] ^= int32(f0)
		boundary[f2] ^= int32(f0)
		boundary[f2] ^= int32(f1)
	}
	return
}

func estimateNormals(nvert uint32, coords []Point3i, nface uint32, index interface{}) (estimated []vec3.T) {
	estimated = make([]vec3.T, nvert)

	endpos := nface * 3
	for f := 0; f < int(endpos); f += 3 {
		f0 := GenericGetInt(index, f)
		f1 := GenericGetInt(index, f+1)
		f2 := GenericGetInt(index, f+2)

		p0 := coords[f0]
		p1 := coords[f1]
		p2 := coords[f2]

		v0 := vec3.T{float32(p0[0]), float32(p0[1]), float32(p0[2])}
		v1 := vec3.T{float32(p1[0]), float32(p1[1]), float32(p1[2])}
		v2 := vec3.T{float32(p2[0]), float32(p2[1]), float32(p2[2])}

		a := vec3.Sub(&v1, &v0)
		b := vec3.Sub(&v2, &v0)
		n := vec3.Cross(&a, &b)

		estimated[f0] = vec3.Add(&estimated[f0], &n)
		estimated[f1] = vec3.Add(&estimated[f1], &n)
		estimated[f2] = vec3.Add(&estimated[f2], &n)
	}
	return
}

func toOctaFloat(v vec3.T, unit int) Point2i {
	p := vec2.T{v[0], v[1]}
	p = p.Scaled(float32(1 / (math.Abs(float64(v[0])) + math.Abs(float64(v[1])) + math.Abs(float64(v[2])))))

	if v[2] < 0 {
		p = vec2.T{float32(1.0 - math.Abs(float64(p[1]))), float32(1.0 - math.Abs(float64(p[0])))}
		if v[0] < 0 {
			p[0] = -p[0]
		}
		if v[1] < 0 {
			p[1] = -p[1]
		}
	}
	return Point2i{int32(p[0] * float32(unit)), int32(p[1] * float32(unit))}
}

func toOctaInt(v Point3i, unit int) Point2i {
	len := (math.Abs(float64(v[0])) + math.Abs(float64(v[1])) + math.Abs(float64(v[2])))
	if len == 0 {
		return Point2i{0, 0}
	}

	p := Point2i{v[0] * int32(unit), v[1] * int32(unit)}
	p = Point2i{int32(float64(p[0]) / len), int32(float64(p[1]) / len)}

	if v[2] < 0 {
		p = Point2i{int32(float64(unit) - math.Abs(float64(p[1]))), int32(float64(unit) - math.Abs(float64(p[0])))}
		if v[0] < 0 {
			p[0] = -p[0]
		}
		if v[1] < 0 {
			p[1] = -p[1]
		}
	}
	return p
}

func toSphereFloat(v Point2i, unit int) vec3.T {
	n := vec3.T{float32(v[0]), float32(v[1]), float32(float64(unit) - math.Abs(float64(v[0])) - math.Abs(float64(v[1])))}
	if n[2] < 0 {
		if v[0] > 0 {
			n[0] = float32(1 * (float64(unit) - math.Abs(float64(v[1]))))
		} else {
			n[0] = float32(-1 * (float64(unit) - math.Abs(float64(v[1]))))
		}
		if v[1] > 0 {
			n[1] = float32(1 * (float64(unit) - math.Abs(float64(v[0]))))
		} else {
			n[1] = float32(-1 * (float64(unit) - math.Abs(float64(v[0]))))
		}
	}
	n.Normalize()
	return n
}

func toSphereInt(v Point2s, unit int) Point3s {
	n := vec3.T{float32(v[0]), float32(v[1]), float32(float64(unit) - math.Abs(float64(v[0])) - math.Abs(float64(v[1])))}
	if n[2] < 0 {
		if v[0] > 0 {
			n[0] = float32(1 * (float64(unit) - math.Abs(float64(v[1]))))
		} else {
			n[0] = float32(-1 * (float64(unit) - math.Abs(float64(v[1]))))
		}
		if v[1] > 0 {
			n[1] = float32(1 * (float64(unit) - math.Abs(float64(v[0]))))
		} else {
			n[1] = float32(-1 * (float64(unit) - math.Abs(float64(v[0]))))
		}
	}
	n.Normalize()
	return Point3s{int16(n[0] * 32767), int16(n[1] * 32767), int16(n[2] * 32767)}
}

func (a *NormalAttr) Quantize(nvert uint32, buffer []byte) {
	n := int(uint32(a.N) * nvert)

	values := make([]int32, n)
	diffs := make([]int32, n)

	a.values = values[:]
	a.diffs = diffs[:]

	var normals []Point2i
	pointsHeader := (*reflect.SliceHeader)((unsafe.Pointer(&normals)))
	pointsHeader.Cap = int(n)
	pointsHeader.Len = int(n)
	pointsHeader.Data = uintptr(unsafe.Pointer(&values[0]))

	switch a.format {
	case FORMAT_FLOAT:
		var points []vec3.T
		cpointsHeader := (*reflect.SliceHeader)((unsafe.Pointer(&points)))
		cpointsHeader.Cap = int(n)
		cpointsHeader.Len = int(n)
		cpointsHeader.Data = uintptr(unsafe.Pointer(&buffer[0]))

		for i := 0; i < int(nvert); i++ {
			normals[i] = toOctaFloat(points[i], int(a.q))
		}
		break
	case FORMAT_INT32:
		var points []Point3i
		cpointsHeader := (*reflect.SliceHeader)((unsafe.Pointer(&points)))
		cpointsHeader.Cap = int(n)
		cpointsHeader.Len = int(n)
		cpointsHeader.Data = uintptr(unsafe.Pointer(&buffer[0]))

		for i := 0; i < int(nvert); i++ {
			normals[i] = toOctaInt(points[i], int(a.q))
		}
		break
	case FORMAT_INT16:
		var points []Point3s
		cpointsHeader := (*reflect.SliceHeader)((unsafe.Pointer(&points)))
		cpointsHeader.Cap = int(n)
		cpointsHeader.Len = int(n)
		cpointsHeader.Data = uintptr(unsafe.Pointer(&buffer[0]))

		for i := 0; i < int(nvert); i++ {
			normals[i] = toOctaInt(Point3i{int32(points[i][0]), int32(points[i][1]), int32(points[i][2])}, int(a.q))
		}
		break
	case FORMAT_INT8:
		var points []Point3b
		cpointsHeader := (*reflect.SliceHeader)((unsafe.Pointer(&points)))
		cpointsHeader.Cap = int(n)
		cpointsHeader.Len = int(n)
		cpointsHeader.Data = uintptr(unsafe.Pointer(&buffer[0]))

		for i := 0; i < int(nvert); i++ {
			normals[i] = toOctaInt(Point3i{int32(points[i][0]), int32(points[i][1]), int32(points[i][2])}, int(a.q))
		}
		break
	}
	min := Point2i{values[0], values[1]}
	max := min

	for i := 1; i < int(nvert); i++ {
		min.setMin(normals[i])
		max.setMax(normals[i])
	}
	max = Point2i{max[0] - min[0], max[1] - min[1]}
	a.bits = imax(ilog2(int(max[0])-1), ilog2(int(max[1])-1)) + 1
}

func (attr *NormalAttr) PreDelta(nvert uint32, nface uint32, attrs map[string]VertexAttribute, index *IndexAttribute) {
	if attr.prediction == PREDICTION_DIFF {
		return
	}

	var va VertexAttribute
	ok := false

	if va, ok = attrs["position"]; !ok {
		panic("No position attribute found. Use DIFF normal strategy instead.")
	}

	coord, ok := va.(*GenericAttr)
	if !ok {
		panic("Position attr has been overloaded, Use DIFF normal strategy instead.")
	}

	var points []Point3i
	cpointsHeader := (*reflect.SliceHeader)((unsafe.Pointer(&points)))
	cpointsHeader.Cap = int(nface)
	cpointsHeader.Len = int(nface)
	cpointsHeader.Data = GenericGetPtr(coord.values)

	estimated := estimateNormals(nvert, points, nface, index.faces)

	if attr.prediction == PREDICTION_BORDER {
		attr.boundary = markBoundary(nvert, nface, index.faces)
	}

	var v []Point2i
	vHeader := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
	vHeader.Cap = int(nvert)
	vHeader.Len = int(nvert)
	vHeader.Data = uintptr(unsafe.Pointer(&attr.values[0]))

	for i := 0; i < int(nvert); i++ {
		n := toOctaFloat(estimated[i], int(attr.q))
		v[i] = Point2i{v[i][0] - n[0], v[i][1] - n[1]}
	}
}

func (attr *NormalAttr) DeltaEncode(context []Quad) {
	if attr.prediction == PREDICTION_DIFF {
		attr.diffs[0] = attr.values[context[0].t*2]
		attr.diffs[1] = attr.values[context[0].t*2+1]

		for i := 1; i < len(context); i++ {
			quad := context[i]
			attr.diffs[i*2+0] = attr.diffs[quad.t*2+0] - attr.diffs[quad.a*2+0]
			attr.diffs[i*2+1] = attr.diffs[quad.t*2+1] - attr.diffs[quad.a*2+1]
		}
		GenericResize(attr.diffs, len(context)*2)
	} else {
		count := 0

		for i := 1; i < len(context); i++ {
			quad := context[i]
			if attr.prediction != PREDICTION_BORDER || attr.boundary[quad.t] != 0 {
				attr.diffs[count*2+0] = attr.values[quad.t*2+0]
				attr.diffs[count*2+1] = attr.values[quad.t*2+1]
				count++
			}
		}
		GenericResize(attr.diffs, count*2)
	}
}

func (a *NormalAttr) Encode(nvert uint32, stream *OutStream) {
	stream.write(a.prediction)
	stream.restart()
	stream.encodeArray(uint32(len(a.diffs)/2), a.diffs, 2)
	a.size = stream.elapsed()
}

func (a *NormalAttr) Decode(nvert uint32, stream *InStream) {
	a.prediction = PredictionType(stream.readUint8())
	a.diffs = make([]int32, nvert*2)
	readed := stream.decodeArray(a.diffs, 2)

	if a.prediction == PREDICTION_BORDER {
		newdiffs := make([]int32, readed*2)
		copy(newdiffs, a.diffs[:readed])
		a.diffs = newdiffs
	}
}

func (attr *NormalAttr) DeltaDecode(nvert uint32, context []Face) {
	if attr.buffer == nil {
		return
	}

	if attr.prediction != PREDICTION_DIFF {
		return
	}

	if len(context) > 0 {
		for i := 1; i < len(context); i++ {
			f := context[i]
			for c := 0; c < 2; c++ {
				d := &attr.diffs[i*2+c]
				*d += attr.diffs[int(f.a)*2+c]
			}

		}
	} else {
		for i := 2; i < int(nvert*2); i++ {
			d := &attr.diffs[i]
			*d += attr.diffs[i-2]
		}
	}
}

func (attr *NormalAttr) PostDelta(nvert uint32, nface uint32, attrs map[string]VertexAttribute, index *IndexAttribute) {
	if attr.buffer == nil {
		return
	}

	if attr.prediction == PREDICTION_DIFF {
		return
	}

	var va VertexAttribute
	ok := false

	if va, ok = attrs["position"]; !ok {
		panic("No position attribute found. Use DIFF normal strategy instead.")
	}

	coord, ok := va.(*GenericAttr)
	if !ok {
		panic("Position attr has been overloaded, Use DIFF normal strategy instead.")
	}

	var points []Point3i
	cpointsHeader := (*reflect.SliceHeader)((unsafe.Pointer(&points)))
	cpointsHeader.Cap = int(nface)
	cpointsHeader.Len = int(nface)
	cpointsHeader.Data = GenericGetPtr(coord.values)

	var estimated []vec3.T

	if index.faces32 != nil {
		estimated = estimateNormals(nvert, points, nface, index.faces32)
	} else {
		estimated = estimateNormals(nvert, points, nface, index.faces16)
	}

	if attr.prediction == PREDICTION_BORDER {
		if index.faces32 != nil {
			attr.boundary = markBoundary(nvert, nface, index.faces32)
		} else {
			attr.boundary = markBoundary(nvert, nface, index.faces16)
		}
	}

	switch attr.format {
	case FORMAT_FLOAT:
		var points []vec3.T
		cpointsHeader := (*reflect.SliceHeader)((unsafe.Pointer(&points)))
		cpointsHeader.Cap = int(nvert)
		cpointsHeader.Len = int(nvert)
		cpointsHeader.Data = uintptr(unsafe.Pointer(&attr.buffer[0]))

		attr.computeNormalsFloat(points, estimated)
		break
	case FORMAT_INT16:
		var points []Point3s
		cpointsHeader := (*reflect.SliceHeader)((unsafe.Pointer(&points)))
		cpointsHeader.Cap = int(nvert)
		cpointsHeader.Len = int(nvert)
		cpointsHeader.Data = uintptr(unsafe.Pointer(&attr.buffer[0]))

		attr.computeNormalsInt16(points, estimated)
		break
	}
}

func (attr *NormalAttr) Dequantize(nvert uint32) {
	if attr.buffer == nil {
		return
	}

	if attr.prediction != PREDICTION_DIFF {
		return
	}

	switch attr.format {
	case FORMAT_FLOAT:
		var points []vec3.T
		cpointsHeader := (*reflect.SliceHeader)((unsafe.Pointer(&points)))
		cpointsHeader.Cap = int(nvert)
		cpointsHeader.Len = int(nvert)
		cpointsHeader.Data = uintptr(unsafe.Pointer(&attr.buffer[0]))

		for i := 0; i < int(nvert); i++ {
			points[i] = toSphereFloat(Point2i{attr.diffs[i*2], attr.diffs[i*2+1]}, int(attr.q))
		}
		break
	case FORMAT_INT16:
		var points []Point3s
		cpointsHeader := (*reflect.SliceHeader)((unsafe.Pointer(&points)))
		cpointsHeader.Cap = int(nvert)
		cpointsHeader.Len = int(nvert)
		cpointsHeader.Data = uintptr(unsafe.Pointer(&attr.buffer[0]))

		for i := 0; i < int(nvert); i++ {
			points[i] = toSphereInt(Point2s{int16(attr.diffs[i*2]), int16(attr.diffs[i*2+1])}, int(attr.q))
		}
		break
	}
}

func (attr *NormalAttr) computeNormalsInt16(normals []Point3s, estimated []vec3.T) {
	nvert := len(estimated)

	var diffp []Point2i
	diffpHeader := (*reflect.SliceHeader)((unsafe.Pointer(&diffp)))
	diffpHeader.Cap = int(nvert)
	diffpHeader.Len = int(nvert)
	diffpHeader.Data = uintptr(unsafe.Pointer(&attr.diffs[0]))

	count := 0
	for i := 0; i < nvert; i++ {
		e := &estimated[i]
		n := &normals[i]

		if attr.prediction == PREDICTION_ESTIMATED || attr.boundary[i] > 0 {
			d := &diffp[count]
			count++
			qn := toOctaFloat(*e, int(attr.q))
			*n = toSphereInt(Point2s{int16(qn[0] + d[0]), int16(qn[1] + d[1])}, int(attr.q))
		} else {
			len := e.Length()
			if len < 0.00001 {
				*e = vec3.T{0, 0, 1}
			} else {
				len = 32767.0 / len
				for k := 0; k < 3; k++ {
					n[k] = int16(e[k] * len)
				}
			}
		}
	}
}

func (attr *NormalAttr) computeNormalsFloat(normals []vec3.T, estimated []vec3.T) {
	nvert := len(estimated)

	var diffp []Point2i
	diffpHeader := (*reflect.SliceHeader)((unsafe.Pointer(&diffp)))
	diffpHeader.Cap = int(nvert)
	diffpHeader.Len = int(nvert)
	diffpHeader.Data = uintptr(unsafe.Pointer(&attr.diffs[0]))

	count := 0
	for i := 0; i < nvert; i++ {
		e := &estimated[i]
		n := &normals[i]

		if attr.prediction == PREDICTION_ESTIMATED || attr.boundary[i] > 0 {
			qn := toOctaFloat(*e, int(attr.q))
			d := &diffp[count]
			count++
			*n = toSphereFloat(Point2i{int32(qn[0] + d[0]), int32(qn[1] + d[1])}, int(attr.q))
		} else {
			*n = vec3.T{e[0], e[1], e[2]}
			n.Normalize()
		}
	}
}
