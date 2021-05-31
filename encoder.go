package corto

/*
#include "corto_api.h"
#cgo CFLAGS: -I ./
*/
import "C"

import (
	"unsafe"

	"github.com/flywave/go3d/vec3"
)

type Encoder struct {
	m *C.struct__corto_encoder_t
}

func NewEncoder(_nvert uint32, _nface uint32, entropy EntropyType) *Encoder {
	return &Encoder{m: C.corto_new_encoder(C.uint(_nvert), C.uint(_nface), C.uint(entropy))}
}

func (e *Encoder) Free() {
	if e.m != nil {
		C.corto_encoder_free(e.m)
	}
}

func (e *Encoder) AddPositions(buffer []float32, q float32, o vec3.T) bool {
	if o.IsZero() {
		return bool(C.corto_encoder_add_positions(e.m, (*C.float)(unsafe.Pointer(&buffer[0])), C.float(q), nil))
	}
	return bool(C.corto_encoder_add_positions(e.m, (*C.float)(unsafe.Pointer(&buffer[0])), C.float(q), (*C.float)(unsafe.Pointer(&o[0]))))
}

func (e *Encoder) AddPositionsInt32(buffer []float32, _index []uint32, q float32, o vec3.T) bool {
	if o.IsZero() {
		return bool(C.corto_encoder_add_positions_index32(e.m, (*C.float)(unsafe.Pointer(&buffer[0])), (*C.uint)(unsafe.Pointer(&_index[0])), C.float(q), nil))
	}
	return bool(C.corto_encoder_add_positions_index32(e.m, (*C.float)(unsafe.Pointer(&buffer[0])), (*C.uint)(unsafe.Pointer(&_index[0])), C.float(q), (*C.float)(unsafe.Pointer(&o[0]))))
}

func (e *Encoder) AddPositionsInt16(buffer []float32, _index []uint16, q float32, o vec3.T) bool {
	if o.IsZero() {
		return bool(C.corto_encoder_add_positions_index16(e.m, (*C.float)(unsafe.Pointer(&buffer[0])), (*C.ushort)(unsafe.Pointer(&_index[0])), C.float(q), nil))
	}
	return bool(C.corto_encoder_add_positions_index16(e.m, (*C.float)(unsafe.Pointer(&buffer[0])), (*C.ushort)(unsafe.Pointer(&_index[0])), C.float(q), (*C.float)(unsafe.Pointer(&o[0]))))
}

func (e *Encoder) AddPositionsBits(buffer []float32, bits int) bool {
	return bool(C.corto_encoder_add_positions_bits(e.m, (*C.float)(unsafe.Pointer(&buffer[0])), C.int(bits)))
}

func (e *Encoder) AddPositionsBitsInt32(buffer []float32, index []uint32, bits int) bool {
	return bool(C.corto_encoder_add_positions_bits_index32(e.m, (*C.float)(unsafe.Pointer(&buffer[0])), (*C.uint)(unsafe.Pointer(&index[0])), C.int(bits)))
}

func (e *Encoder) AddPositionsBitsInt16(buffer []float32, index []uint16, bits int) bool {
	return bool(C.corto_encoder_add_positions_bits_index16(e.m, (*C.float)(unsafe.Pointer(&buffer[0])), (*C.ushort)(unsafe.Pointer(&index[0])), C.int(bits)))
}

func (e *Encoder) AddNormals(buffer []float32, bits int, no PredictionType) bool {
	return bool(C.corto_encoder_add_normals_float(e.m, (*C.float)(unsafe.Pointer(&buffer[0])), C.int(bits), C.uint(no)))
}

func (e *Encoder) AddNormalsInt16(buffer []int16, bits int, no PredictionType) bool {
	return bool(C.corto_encoder_add_normals_short(e.m, (*C.short)(unsafe.Pointer(&buffer[0])), C.int(bits), C.uint(no)))
}

func (e *Encoder) AddColors(buffer []byte, cbits []int) bool {
	return bool(C.corto_encoder_add_colors(e.m, (*C.uchar)(unsafe.Pointer(&buffer[0])), C.int(cbits[0]), C.int(cbits[1]), C.int(cbits[2]), C.int(cbits[3])))
}

func (e *Encoder) AddColors3(buffer []byte, cbits []int) bool {
	return bool(C.corto_encoder_add_colors3(e.m, (*C.uchar)(unsafe.Pointer(&buffer[0])), C.int(cbits[0]), C.int(cbits[1]), C.int(cbits[2])))
}

func (e *Encoder) AddUvs(buffer []float32, q float32) bool {
	return bool(C.corto_encoder_add_uvs(e.m, (*C.float)(unsafe.Pointer(&buffer[0])), C.float(q)))
}

func (e *Encoder) AddAttributeFormat(name string, buffer []byte, format FormatType, components int, q float32, strategy StrategyType) bool {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return bool(C.corto_encoder_add_attribute(e.m, cname, (*C.char)(unsafe.Pointer(&buffer[0])), C.uint(format), C.int(components), C.float(q), C.uint(strategy)))
}

func (e *Encoder) AddGroup(end_triangle int) {
	C.corto_encoder_add_group(e.m, C.int(end_triangle))
}

func (e *Encoder) AddGroupProps(end_triangle int, props map[string]string) {
	keys := make([]*C.char, len(props))
	values := make([]*C.char, len(props))
	pos := 0
	for k, v := range props {
		keys[pos] = C.CString(k)
		values[pos] = C.CString(v)
		pos++
	}
	defer func() {
		for i := 0; i < pos; i++ {
			C.free(unsafe.Pointer(keys[i]))
			C.free(unsafe.Pointer(values[i]))
		}
	}()
	C.corto_encoder_add_group_props(e.m, C.int(end_triangle), &keys[0], &values[0], C.int(pos))
}

func (e *Encoder) Encode() []byte {
	si := int(C.corto_encoder_encode(e.m))
	ret := make([]byte, si)
	C.corto_encoder_get_data(e.m, (*C.char)(unsafe.Pointer(&ret[0])), C.size_t(si))
	return ret
}
