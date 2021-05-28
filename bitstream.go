package corto

var (
	bmask = []uint32{
		0x00, 0x01, 0x03, 0x07, 0x0f, 0x01f, 0x03f, 0x07f,
		0xff, 0x01ff, 0x03ff, 0x07ff, 0x0fff, 0x01fff, 0x03fff, 0x07fff,
		0xffff, 0x01ffff, 0x03ffff, 0x07ffff, 0x0fffff, 0x01fffff, 0x03fffff, 0x07fffff,
		0xffffff, 0x01ffffff, 0x03ffffff, 0x07ffffff, 0x0fffffff, 0x01fffffff, 0x03fffffff, 0x7fffffff}
)

const (
	BITS_PER_WORD = 32
)

type BitStream struct {
	size   int
	buffer []uint32

	allocated int
	pos       int
	buff      uint32
	bits      int
}

func NewBitStream(reserved *int) *BitStream {
	s := &BitStream{size: 0, buffer: nil, pos: 0, buff: 0, bits: 0}
	if reserved == nil {
		return s
	}
	s.reserve(*reserved)
	return s
}

func NewBitStreamWithBuf(size int, buffer []uint32) *BitStream {
	s := &BitStream{size: size, buffer: buffer, pos: 0, buff: 0, bits: 0}
	return s
}

func (s *BitStream) init(size int, buffer []uint32) {
	s.size = size
	s.buffer = buffer
}

func (s *BitStream) reserve(reserved int) {
	s.allocated = reserved
	s.buffer = make([]uint32, s.allocated)
	s.size = 0
	s.buff = 0
	s.bits = BITS_PER_WORD
	s.pos = 0
}

func (s *BitStream) write(value uint32, numbits int) {
	if s.allocated == 0 {
		s.reserve(256)
	}
	if numbits >= s.bits {
		s.buff = (s.buff << s.bits) | (value >> (numbits - s.bits))
		s.push_back(s.buff)
		value &= bmask[numbits-s.bits]
		numbits -= s.bits
		s.bits = BITS_PER_WORD
		s.buff = 0
	}

	if numbits > 0 {
		s.buff = (s.buff << numbits) | value
		s.bits -= numbits
	}
}

func (s *BitStream) read(numbits int) uint32 {
	if numbits > s.bits {
		s.bits = numbits - s.bits
		result := (s.buff << s.bits)
		s.bits = 32 - s.bits

		s.buff = s.buffer[s.pos]
		s.pos++
		result |= (s.buff >> s.bits)
		s.buff = (s.buff & ((1 << s.bits) - 1))
		return result

	} else {
		s.bits -= numbits
		result := (s.buff >> s.bits)
		s.buff = (s.buff & ((1 << s.bits) - 1))
		return result
	}
}

func (s *BitStream) writtenBits() uint32 {
	return uint32((s.size+1)*BITS_PER_WORD - s.bits)
}

func (s *BitStream) flush() {
	if s.bits != BITS_PER_WORD {
		s.push_back((s.buff << s.bits))
		s.buff = 0
		s.bits = BITS_PER_WORD
	}
}

func (s *BitStream) push_back(w uint32) {
	if s.size >= s.allocated {
		b := make([]uint32, s.allocated*2)
		copy(b, s.buffer[:s.allocated])
		s.buffer = b
		s.allocated *= 2
	}
	s.buffer[s.size] = w
	s.size++
}
