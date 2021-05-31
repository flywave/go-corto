package corto

/*
#include "corto_api.h"
#cgo CFLAGS: -I ./
*/
import "C"
import "unsafe"

type Decoder struct {
	m *C.struct__corto_decoder_t
}

func NewDecoder(input []byte) *Decoder {
	return &Decoder{m: C.corto_new_decoder(C.int(len(input)), (*C.uchar)((unsafe.Pointer)(&input[0])))}
}

func (e *Decoder) Free() {
	if e.m != nil {
		C.corto_decoder_free(e.m)
	}
}

func (d *Decoder) HasAttr(name string) bool {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return bool(C.corto_encoder_has_attr(d.m, cname))
}

func (d *Decoder) SetPositions(buffer []float32) bool {
	return bool(C.corto_encoder_set_positions(d.m, (*C.float)(unsafe.Pointer(&buffer[0]))))
}

func (d *Decoder) SetNormals(buffer []float32) bool {
	return bool(C.corto_encoder_set_normals_float(d.m, (*C.float)(unsafe.Pointer(&buffer[0]))))
}

func (d *Decoder) SetNormalsInt16(buffer []int16) bool {
	return bool(C.corto_encoder_set_normals_short(d.m, (*C.short)(unsafe.Pointer(&buffer[0]))))
}

func (d *Decoder) SetUvs(buffer []float32) bool {
	return bool(C.corto_encoder_set_uvs(d.m, (*C.float)(unsafe.Pointer(&buffer[0]))))
}

func (d *Decoder) SetColors(buffer []byte, components int) bool {
	return bool(C.corto_encoder_set_colors(d.m, (*C.uchar)(unsafe.Pointer(&buffer[0])), C.int(components)))
}

func (d *Decoder) SetAttributeFormat(name string, buffer []byte, format FormatType) bool {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return bool(C.corto_encoder_set_attribute(d.m, cname, (*C.char)(unsafe.Pointer(&buffer[0])), C.uint(format)))
}

func (d *Decoder) SetIndexInt32(buffer []uint32) {
	C.corto_encoder_set_index32(d.m, (*C.uint)(unsafe.Pointer(&buffer[0])))
}

func (d *Decoder) SetIndexInt16(buffer []uint16) {
	C.corto_encoder_set_index16(d.m, (*C.ushort)(unsafe.Pointer(&buffer[0])))
}

func (d *Decoder) Decode() {
	C.corto_encoder_decode(d.m)
}
