package corto

type Point2i [2]int32
type Point3i [3]int32
type Point2s [2]int16
type Point3s [3]int16
type Color4b [4]byte
type Point3b [3]byte

func (p *Point2i) setMin(c Point2i) {
	if c[0] < p[0] {
		p[0] = c[0]
	}
	if c[1] < p[1] {
		p[1] = c[1]
	}
}
func (p *Point2i) setMax(c Point2i) {
	if c[0] > p[0] {
		p[0] = c[0]
	}
	if c[1] > p[1] {
		p[1] = c[1]
	}
}

func (p *Point3i) setMin(c Point3i) {
	if c[0] < p[0] {
		p[0] = c[0]
	}
	if c[1] < p[1] {
		p[1] = c[1]
	}
	if c[2] < p[2] {
		p[2] = c[2]
	}
}
func (p *Point3i) setMax(c Point3i) {
	if c[0] > p[0] {
		p[0] = c[0]
	}
	if c[1] > p[1] {
		p[1] = c[1]
	}
	if c[2] > p[2] {
		p[2] = c[2]
	}
}

func (p *Point2s) setMin(c Point2s) {
	if c[0] < p[0] {
		p[0] = c[0]
	}
	if c[1] < p[1] {
		p[1] = c[1]
	}
}
func (p *Point2s) setMax(c Point2s) {
	if c[0] > p[0] {
		p[0] = c[0]
	}
	if c[1] > p[1] {
		p[1] = c[1]
	}
}

func (p *Point3s) setMin(c Point3s) {
	if c[0] < p[0] {
		p[0] = c[0]
	}
	if c[1] < p[1] {
		p[1] = c[1]
	}
	if c[2] < p[2] {
		p[2] = c[2]
	}
}
func (p *Point3s) setMax(c Point3s) {
	if c[0] > p[0] {
		p[0] = c[0]
	}
	if c[1] > p[1] {
		p[1] = c[1]
	}
	if c[2] > p[2] {
		p[2] = c[2]
	}
}

func (b *Color4b) toYCC() Color4b {
	return Color4b{b[1], byte(b[2] - b[1]), byte(b[0] - b[1]), b[3]}
}

func (b *Color4b) toRGB() Color4b {
	return Color4b{byte(b[2] + b[0]), b[0], byte(b[1] + b[0]), b[3]}
}
