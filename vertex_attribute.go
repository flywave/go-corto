package corto

import "math"

type FormatType uint32
type StrategyType uint32
type CodecType uint32

const (
	FORMAT_UINT32 FormatType = 0
	FORMAT_INT32  FormatType = 1
	FORMAT_UINT16 FormatType = 2
	FORMAT_INT16  FormatType = 3
	FORMAT_UINT8  FormatType = 4
	FORMAT_INT8   FormatType = 5
	FORMAT_FLOAT  FormatType = 6
	FORMAT_DOUBLE FormatType = 7
)

const (
	PARALLEL   StrategyType = 0x1
	CORRELATED StrategyType = 0x2
)

const (
	GENERIC_CODEC CodecType = 1
	NORMAL_CODEC  CodecType = 2
	COLOR_CODEC   CodecType = 3
	CUSTOM_CODEC  CodecType = 100
)

type VertexAttribute interface {
	Codec() CodecType
	Quantize(nvert uint32, buffer []byte)
	PreDelta(nvert uint32, nface uint32, attrs map[string]VertexAttribute, index *IndexAttribute)
	DeltaEncode(context []Quad)
	Encode(nvert uint32, stream *OutStream)
	Decode(nvert uint32, stream *InStream)
	DeltaDecode(nvert uint32, faces []Face)
	PostDelta(nvert uint32, nface uint32, attrs map[string]VertexAttribute, index *IndexAttribute)
	Dequantize(nvert uint32)
	Q() float32
	Dim() byte
	Format() byte
	Strategy() byte
}

type GenericAttr struct {
	VertexAttribute
	buffer   []byte
	N        int
	q        float32
	strategy StrategyType

	format FormatType
	size   uint32
	bits   int

	values interface{}
	diffs  interface{}
}

func NewGenericAttr(dim int, values interface{}, diffs interface{}) *GenericAttr {
	return &GenericAttr{N: dim, values: values, diffs: diffs}
}

func (a *GenericAttr) Q() float32 {
	return a.q
}

func (a *GenericAttr) Dim() byte {
	return byte(a.N)
}

func (a *GenericAttr) Format() byte {
	return byte(a.format)
}

func (a *GenericAttr) Strategy() byte {
	return byte(a.strategy)
}

func (a *GenericAttr) Codec() CodecType {
	return GENERIC_CODEC
}

func (a *GenericAttr) Quantize(nvert uint32, buffer []byte) {
	n := int(uint32(a.N) * nvert)

	a.values = GenericResize(a.values, n)
	a.diffs = GenericResize(a.diffs, n)

	switch a.format {
	case FORMAT_INT32:
		for i := 0; i < n; i++ {
			GenericSet(a.values, i, int32(float32(byteorder.Uint32(buffer[i:]))/a.q))
		}
		break
	case FORMAT_INT16:
		for i := 0; i < n; i++ {
			GenericSet(a.values, i, int16(float32(byteorder.Uint16(buffer[i:]))/a.q))
		}
		break
	case FORMAT_INT8:
		for i := 0; i < n; i++ {
			GenericSet(a.values, i, int8(float32(buffer[i])/a.q))
		}
		break
	case FORMAT_FLOAT:
		for i := 0; i < n; i++ {
			GenericSet(a.values, i, math.Float32frombits(byteorder.Uint32(buffer[i:]))/a.q)
		}
		break
	case FORMAT_DOUBLE:
		for i := 0; i < n; i++ {
			GenericSet(a.values, i, math.Float64frombits(byteorder.Uint64(buffer[i:]))/float64(a.q))
		}
		break
	}
	a.bits = 0
	for k := 0; k < a.N; k++ {
		min := GenericGet(a.values, k)
		max := min
		for i := k; i < n; i += a.N {
			b := GenericGet(a.values, i)
			if GenericGreater(min, b) {
				min = b
			}
			if GenericLess(max, b) {
				max = b
			}
		}
		max = genericSub(max, min)
		a.bits = imax(a.bits, ilog2(GenericInt(max)-1)+1)
	}
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func genericSub(max, min interface{}) interface{} {
	switch data := max.(type) {
	case int8:
		return data - min.(int8)
	case uint8:
		return data - min.(uint8)
	case int16:
		return data - min.(int16)
	case uint16:
		return data - min.(uint16)
	case int32:
		return data - min.(int32)
	case uint32:
		return data - min.(uint32)
	case float32:
		return data - min.(float32)
	case float64:
		return data - min.(float64)
	}
	return 0
}

func genericAdd(max, min interface{}) interface{} {
	switch data := max.(type) {
	case int8:
		return data + min.(int8)
	case uint8:
		return data + min.(uint8)
	case int16:
		return data + min.(int16)
	case uint16:
		return data + min.(uint16)
	case int32:
		return data + min.(int32)
	case uint32:
		return data + min.(uint32)
	case float32:
		return data + min.(float32)
	case float64:
		return data + min.(float64)
	}
	return 0
}

func (attr *GenericAttr) DeltaEncode(context []Quad) {
	for c := 0; c < attr.N; c++ {
		t := GenericGet(attr.values, int(context[0].t)*attr.N+c)
		GenericSet(attr.values, c, t)
	}
	for i := 1; i < len(context); i++ {
		q := &context[i]
		if (q.a != q.b) && (attr.strategy&PARALLEL) > 0 {
			for c := 0; c < attr.N; c++ {
				t := GenericGet(attr.values, int(q.t)*attr.N+c)
				a := GenericGet(attr.values, int(q.a)*attr.N+c)
				b := GenericGet(attr.values, int(q.b)*attr.N+c)
				cc := GenericGet(attr.values, int(q.c)*attr.N+c)

				r := genericSub(genericAdd(genericSub(t, a), b), cc)

				GenericSet(attr.diffs, i*attr.N+c, r)
			}
		} else {
			for c := 0; c < attr.N; c++ {
				t := GenericGet(attr.values, int(q.t)*attr.N+c)
				a := GenericGet(attr.values, int(q.a)*attr.N+c)

				r := genericSub(t, a)

				GenericSet(attr.diffs, i*attr.N+c, r)
			}
		}
	}
	GenericResize(attr.diffs, len(context)*attr.N)
}

func (a *GenericAttr) Encode(nvert uint32, stream *OutStream) {
	stream.restart()
	if (a.strategy & CORRELATED) > 0 {
		stream.encodeArray(nvert, a.diffs, a.N)
	} else {
		stream.encodeValues(nvert, a.diffs, a.N)
	}

	a.size = stream.elapsed()
}

func (a *GenericAttr) Decode(nvert uint32, stream *InStream) {
	if (a.strategy & CORRELATED) > 0 {
		stream.decodeArray(a.buffer, a.N)
	} else {
		stream.decodeValues(a.buffer, a.N)
	}
}

func (attr *GenericAttr) DeltaDecode(nvert uint32, context []Face) {
	if attr.buffer == nil {
		return
	}

	values := attr.buffer[:]

	if (attr.strategy & PARALLEL) > 0 {
		for i := 1; i < len(context); i++ {
			f := context[i]
			for c := 0; c < attr.N; c++ {
				a := GenericGet(values, int(f.a)*attr.N+c)
				b := GenericGet(values, int(f.b)*attr.N+c)
				cc := GenericGet(values, int(f.c)*attr.N+c)
				z := GenericGet(values, i*attr.N+c)

				r := genericAdd(genericSub(genericAdd(a, b), cc), z)

				GenericSet(values, i*attr.N+c, r)
			}
		}
	} else if len(context) > 0 {
		for i := 1; i < len(context); i++ {
			f := context[i]
			for c := 0; c < attr.N; c++ {
				a := GenericGet(values, int(f.a)*attr.N+c)
				z := GenericGet(values, i*attr.N+c)
				r := genericAdd(a, z)
				GenericSet(values, i*attr.N+c, r)
			}
		}
	} else {
		for i := attr.N; i < int(nvert)*attr.N; i++ {
			a := GenericGet(values, i-attr.N)
			z := GenericGet(values, i)
			r := genericAdd(a, z)
			GenericSet(values, i, r)
		}
	}
}

func (attr *GenericAttr) Dequantize(nvert uint32) {
	if attr.buffer == nil {
		return
	}

	buffer := attr.buffer[:]

	n := attr.N * int(nvert)
	switch attr.format {
	case FORMAT_FLOAT:
		for i := 0; i < n; i++ {
			v := math.Float32frombits(byteorder.Uint32(buffer[i:])) * attr.q
			byteorder.PutUint32(buffer[i:], math.Float32bits(v))
		}
		break

	case FORMAT_INT16:
		for i := 0; i < n; i++ {
			v := float32(int16(byteorder.Uint16(buffer[i:]))) * attr.q
			byteorder.PutUint16(buffer[i:], uint16(v))
		}
		break

	case FORMAT_INT32:
		for i := 0; i < n; i++ {
			v := float32(int32(byteorder.Uint32(buffer[i:]))) * attr.q
			byteorder.PutUint32(buffer[i:], uint32(v))
		}
		break

	case FORMAT_INT8:
		for i := 0; i < n; i++ {
			v := float32(buffer[i]) * attr.q
			buffer[i] = byte(v)
		}
		break

	case FORMAT_DOUBLE:
		for i := 0; i < n; i++ {
			v := math.Float64frombits(byteorder.Uint64(buffer[i:])) * float64(attr.q)
			byteorder.PutUint64(buffer[i:], math.Float64bits(v))
		}
		break

	case FORMAT_UINT16:
		for i := 0; i < n; i++ {
			v := float32(byteorder.Uint16(buffer[i:])) * attr.q
			byteorder.PutUint16(buffer[i:], uint16(v))
		}
		break

	case FORMAT_UINT32:
		for i := 0; i < n; i++ {
			v := float32(byteorder.Uint32(buffer[i:])) * attr.q
			byteorder.PutUint32(buffer[i:], uint32(v))
		}
		break

	case FORMAT_UINT8:
		for i := 0; i < n; i++ {
			buffer[i] = byte(float32(buffer[i]) * attr.q)
		}
		break
	}
}
