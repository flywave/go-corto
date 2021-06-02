package corto

import (
	"testing"

	"github.com/flywave/go3d/vec3"
)

var (
	testMesh = Geom{
		Vertices: []vec3.T{
			{0., 0., 0.},
			{1., 0., 0.},
			{0., 1., 0.},
			{0., 1., 0.},
			{1., 0., 0.},
			{1., 1., 0.},
			{0., 1., 1.},
			{1., 0., 1.},
			{0., 0., 1.},
			{1., 1., 1.},
			{1., 0., 1.},
			{0., 1., 1.},
			{0., 1., 0.},
			{1., 1., 0.},
			{0., 1., 1.},
			{0., 1., 1.},
			{1., 1., 0.},
			{1., 1., 1.},
			{0., 0., 1.},
			{1., 0., 0.},
			{0., 0., 0.},
			{1., 0., 1.},
			{1., 0., 0.},
			{0., 0., 1.},
			{1., 0., 0.},
			{1., 0., 1.},
			{1., 1., 0.},
			{1., 1., 0.},
			{1., 0., 1.},
			{1., 1., 1.},
			{0., 1., 0.},
			{0., 0., 1.},
			{0., 0., 0.},
			{0., 1., 1.},
			{0., 0., 1.},
			{0., 1., 0.}},
		Indices16: []Face16{
			{0, 1, 2},
			{3, 4, 5},
			{6, 7, 8},
			{9, 10, 11},
			{12, 13, 14},
			{15, 16, 17},
			{18, 19, 20},
			{21, 22, 23},
			{24, 25, 26},
			{27, 28, 29},
			{30, 31, 32},
			{33, 34, 35},
		},
	}
)

func TestCorto(t *testing.T) {
	ctx := NewEncoderContext(0.1)
	ctx.VertexBits = 1
	buf := EncodeGeom(ctx, &testMesh)

	if len(buf) == 0 {
		t.FailNow()
	}

	dctx := NewDecoderContext(len(testMesh.Indices16), len(testMesh.Vertices), 3, true, false)

	geom := DecodeGeom(dctx, buf)

	if len(geom.Indices16) != len(testMesh.Indices16) {
		t.FailNow()
	}
}
