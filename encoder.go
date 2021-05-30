package corto

/*
#include "corto_api.h"
#cgo CFLAGS: -I ./
*/
import "C"

import (
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

func (e *Encoder) addPositions(buffer []float32, q float32, o vec3.T) bool {
	return false
}

func (e *Encoder) addPositionsInt32(buffer []float32, _index []uint32, q float32, o vec3.T) bool {
	return false
}

func (e *Encoder) addPositionsInt16(buffer []float32, _index []uint16, q float32, o vec3.T) bool {
	return false
}

func (e *Encoder) addPositionsBits(buffer []float32, bits int) bool {
	return false
}

func (e *Encoder) addPositionsBitsInt32(buffer []float32, index []uint32, bits int) bool {
	return false
}

func (e *Encoder) addPositionsBitsInt16(buffer []float32, index []uint16, bits int) bool {
	return false
}

func (e *Encoder) addNormals(buffer []float32, bits int, no PredictionType) bool {
	return false
}

func (e *Encoder) addNormalsInt16(buffer []uint16, bits int, no PredictionType) bool {
	return false

}

func (e *Encoder) addColors(buffer []byte, rbits int, gbits, bbits int, abits int) bool {
	return false

}

func (e *Encoder) addColors3(buffer []byte, rbits int, gbits int, bbits int) bool {
	return false

}

func (e *Encoder) addUvs(buffer []float32, q float32) bool {
	return false

}

func (e *Encoder) addAttributeFormat(name string, buffer []byte, format FormatType, components int, q float32, strategy StrategyType) bool {
	return false

}

func (e *Encoder) addGroup(end_triangle int) {
}

func (e *Encoder) addGroupProps(end_triangle int, props map[string]string) {
}

func (e *Encoder) encode() {

}

func (e *Encoder) encodePointCloud() {

}

func (e *Encoder) encodeMesh() {

}
