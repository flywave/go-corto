package corto

import (
	"bytes"
	"encoding/binary"
)

type EntropyType uint32

const (
	ENTROPY_NONE     EntropyType = 0
	ENTROPY_TUNSTALL EntropyType = 1
)

type Stream struct {
	Entropy EntropyType
}

type OutStream struct {
	Stream
	buffer    []byte
	stopwatch int
}

func (s *OutStream) size() int {
	return len(s.buffer)
}

func (s *OutStream) data() []byte {
	return s.buffer[:]
}

func (s *OutStream) reserve(r int) {
	len := len(s.buffer)

	if r <= len {
		return
	}

	buf := make([]byte, r)
	copy(buf, s.buffer)
	s.buffer = buf
}

func (s *OutStream) restart() {
	s.stopwatch = len(s.buffer)
}

func (s *OutStream) elapsed() uint32 {
	e := s.size() - s.stopwatch
	s.stopwatch = s.size()
	return uint32(e)
}

func (s *OutStream) compress(size uint32, data []byte) int {
	switch s.Entropy {
	case ENTROPY_NONE:
		s.write(size)
		s.write(data[:size])
		return 4 + int(size)

	case ENTROPY_TUNSTALL:
		return s.tunstall_compress(data, size)
	}
	return 0
}

func (s *OutStream) tunstall_compress(data []byte, size uint32) int {
	var t Tunstall
	t.getProbabilities(data, int(size))

	t.createDecodingTables2()
	t.createEncodingTables()

	compressed_data, compressed_size := t.compress(data, int(size))

	s.write(uint32(len(t.probabilities)))
	s.write(t.probabilities.Data())

	s.write(size)
	s.write(uint32(compressed_size))
	s.write(compressed_data[:compressed_size])
	return 1 + len(t.probabilities)*2 + 4 + 4 + compressed_size
}

func (s *OutStream) writeString(str string) {
	bytes := uint16(len(str) + 1)
	s.write(bytes)
	s.write([]byte(str))
}

func (s *OutStream) write(v interface{}) {
	n := binary.Size(v)
	pos := s.grow(n)
	buf := make([]byte, n)
	writer := bytes.NewBuffer(buf)
	binary.Write(writer, byteorder, v)
	copy(s.buffer[pos:], writer.Bytes())
}

func (s *OutStream) writeBitStream(stream *BitStream) {
	stream.flush()
	s.write(int32(stream.size))

	pad := s.size() & 0x3
	if pad != 0 {
		pad = 4 - pad
	}

	s.grow(pad)
	s.write(stream.buffer[:stream.size])
}

func (s *OutStream) grow(l int) int {
	len := len(s.buffer)

	buf := make([]byte, len+l)
	copy(buf, s.buffer)
	s.buffer = buf

	return len
}

func needed(a int) int {
	if a == 0 {
		return 0
	}
	if a == -1 {
		return 1
	}
	if a < 0 {
		a = -a - 1
	}
	n := 2

	for ; a > 0; a >>= 1 {
		n++
	}
	return n
}

func (s *OutStream) encodeValues(size uint32, values interface{}) {

}

func (s *OutStream) encodeArray(size uint32, values interface{}) {

}

func (s *OutStream) encodeDiffs(size uint32, values interface{}) {

}

func (s *OutStream) encodeIndices(size uint32, values interface{}) {

}

type InStream struct {
	Stream
	buffer []byte
	pos    int
}

func (s *InStream) decompress(data []byte) {}

func (s *InStream) tunstall_decompress(data []byte) []byte {
	return nil
}

func (s *InStream) rewind() {

}

func (s *InStream) readArray(t interface{}) interface{} {
	return nil
}

func (s *InStream) readUint8() uint8 {
	return 0
}

func (s *InStream) readUint16() uint16 {
	return 0
}

func (s *InStream) readUint32() uint32 {
	return 0
}

func (s *InStream) readFloat() float32 {
	return 0
}

func (s *InStream) readString() string {
	return ""
}

func (s *InStream) read(stream *BitStream) {

}

func (s *InStream) decodeValues(values interface{}) int {
	return 0
}

func (s *InStream) decodeArray(values interface{}) int {
	return 0
}

func (s *InStream) decodeIndices(values interface{}) int {
	return 0
}

func (s *InStream) decodeDiffs(values interface{}) int {
	return 0
}
