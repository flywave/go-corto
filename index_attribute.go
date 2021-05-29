package corto

type ClersType uint8

const (
	CLERS_VERTEX   = 0
	CLERS_LEFT     = 1
	CLERS_RIGHT    = 2
	CLERS_END      = 3
	CLERS_BOUNDARY = 4
	CLERS_DELAY    = 5
	CLERS_SPLIT    = 6
)

type Quad struct {
	t uint32
	a uint32
	b uint32
	c uint32
}

type Face struct {
	a uint32
	b uint32
	c uint32
}

type Group struct {
	end        uint32
	properties map[string]string
}

type IndexAttribute struct {
	faces32    []uint32
	faces16    []uint16
	faces      []uint32
	prediction []Face
	groups     []Group
	clers      []byte
	bitstream  BitStream
	max_front  uint32
	size       uint32
}

func NewIndexAttribute() *IndexAttribute {
	return &IndexAttribute{faces32: nil, faces16: nil, max_front: 0}
}

func (a *IndexAttribute) encode(stream *OutStream) {
	stream.write(a.max_front)
	stream.restart()
	stream.compress(uint32(len(a.clers)), a.clers)
	stream.write(a.bitstream)
	a.size = uint32(stream.elapsed())
}

func (a *IndexAttribute) encodeGroups(stream *OutStream) {
	stream.write(uint32(len(a.groups)))
	for _, g := range a.groups {
		stream.write(g.end)
		stream.write(byte(len(g.properties)))
		for f, s := range g.properties {
			stream.writeString(f)
			stream.writeString(s)
		}
	}
}

func (a *IndexAttribute) decode(stream *InStream) {
	a.max_front = stream.readUint32()
	a.clers = stream.decompress()
	stream.read(&a.bitstream)
}

func (a *IndexAttribute) decodeGroups(stream *InStream) {
	a.groups = make([]Group, stream.readUint32())
	for _, g := range a.groups {
		g.end = stream.readUint32()
		size := stream.readUint8()
		for i := 0; i < int(size); i++ {
			key := stream.readString()
			g.properties[key] = stream.readString()
		}
	}
}
