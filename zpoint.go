package corto

import "github.com/flywave/go3d/vec3"

const ONE64 uint64 = 1

type ZPoint struct {
	Bits uint64
	Pos  uint32
}

func NewZPoint(b uint64) *ZPoint {
	return &ZPoint{Bits: b, Pos: 0}
}

func NewZPointFromXYZ(x, y, z uint64, levels int, i int) *ZPoint {
	l := uint64(1)
	r := &ZPoint{}
	for i := 0; i < levels; i++ {
		r.Bits |= (x&l<<i)<<(2*i) | (y&l<<i)<<(2*i+1) | (z&l<<i)<<(2*i+2)
	}
	return r
}

func morton2(x uint64) uint64 {
	x = x & 0x55555555
	x = (x | (x >> 1)) & 0x33333333
	x = (x | (x >> 2)) & 0x0F0F0F0F
	x = (x | (x >> 4)) & 0x00FF00FF
	x = (x | (x >> 8)) & 0x0000FFFF
	return x
}

func morton3(x uint64) uint64 {
	//1001 0010  0100 1001  0010 0100  1001 0010  0100 1001  0010 0100  1001 0010  0100 1001 => 249
	//1001001001001001001001001001001001001001001001001001001001001001
	//a  b  c  d  e  f  g  h  i  l  m  n  o  p  q  r  s  t  u  v  x  y
	// a  b  c  d  e  f  g  h  i  l  m  n  o  p  q  r  s  t  u  v  x
	//  a  b  c  d  e  f  g  h  i  l  m  n  o  p  q  r  s  t  u  v  x
	//
	//0011000011000011000011000011000011000011000011000011000011000011
	x = x & 0x9249249249249249
	x = (x | (x >> 2)) & 0x30c30c30c30c30c3
	x = (x | (x >> 4)) & 0xf00f00f00f00f00f
	x = (x | (x >> 8)) & 0x00ff0000ff0000ff
	x = (x | (x >> 16)) & 0xffff00000000ffff
	x = (x | (x >> 32)) & 0x00000000ffffffff
	return x
}

func (p *ZPoint) ToPoint2(min Point3i, step float32) vec3.T {
	x := int(morton3(p.Bits))
	y := int(morton3(p.Bits >> 1))
	z := int(morton3(p.Bits >> 2))

	pr := vec3.T{float32(int32(x) + min[0]), float32(int32(y) + min[1]), float32(int32(z) + min[2])}
	pr.Scale(step)
	return pr
}

func (p *ZPoint) ToPoint(step float32) vec3.T {
	x := int(morton3(p.Bits))
	y := int(morton3(p.Bits >> 1))
	z := int(morton3(p.Bits >> 2))

	return vec3.T{float32(x) * step, float32(y) * step, float32(z) * step}
}

func (p *ZPoint) ClearBit(i int) {
	p.Bits &= ^(ONE64 << i)
}

func (p *ZPoint) SetBit(i int) {
	p.Bits |= (ONE64 << i)
}

func (p *ZPoint) SetBit2(i int, val uint64) {
	p.Bits &= ^(ONE64 << i)
	p.Bits |= (val << i)
}

func (p *ZPoint) TestBit(i int) bool {
	return (p.Bits & (ONE64 << i)) != 0
}

func (p *ZPoint) Eq(t *ZPoint) bool {
	return p.Bits == t.Bits
}

func (p *ZPoint) NotEq(t *ZPoint) bool {
	return p.Bits != t.Bits
}

func (p *ZPoint) Less(t *ZPoint) bool {
	return p.Bits > t.Bits
}

func log2(p uint64) int {
	k := 0
	for {
		if p = p >> 1; p > 0 {
			k++
		} else {
			break
		}
	}
	return k
}

func (p *ZPoint) Difference(t *ZPoint) int {
	diff := p.Bits ^ t.Bits
	return log2(diff)
}
