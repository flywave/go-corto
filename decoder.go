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

func (d *Decoder) hasAttr(name string) bool {
	return false
}

func (d *Decoder) setPositions(buffer []float32) bool {
	return false
}

func (d *Decoder) setNormals(buffer []float32) bool {
	return false
}

func (d *Decoder) setNormalsInt16(buffer []int16) bool {
	return false
}

func (d *Decoder) setUvs(buffer []float32) bool {
	return false
}

func (d *Decoder) setColors(buffer []byte, components int) bool {
	return false
}

func (d *Decoder) setAttributeFormat(name string, buffer interface{}, format FormatType) bool {
	return false
}

func (d *Decoder) setIndexInt32(buffer []uint32) {
}

func (d *Decoder) setIndexInt16(buffer []uint16) {
}

func (d *Decoder) decode() {

}

func (d *Decoder) decodePointCloud() {
}

func (d *Decoder) decodeMesh() {

}

func (d *Decoder) decodeFaces(start, end uint32) uint32 {
	return 0
}
