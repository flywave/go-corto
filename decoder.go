package corto

type Decoder struct {
}

func NewDecoder(input []byte) *Decoder {
	return nil
}

func (d *Decoder) hasAttr(name string) bool {
	return ok
}

func (d *Decoder) setPositions(buffer []float32) bool {
	return d.setAttributeFormat("position", buffer, FORMAT_FLOAT)
}

func (d *Decoder) setNormals(buffer []float32) bool {
	return d.setAttributeFormat("normal", buffer, FORMAT_FLOAT)
}

func (d *Decoder) setNormalsInt16(buffer []int16) bool {
	return d.setAttributeFormat("normal", buffer, FORMAT_INT16)
}

func (d *Decoder) setUvs(buffer []float32) bool {
	return d.setAttributeFormat("uv", buffer, FORMAT_FLOAT)
}

func (d *Decoder) setColors(buffer []byte, components int) bool {
	return false
}

func (d *Decoder) setAttributeFormat(name string, buffer interface{}, format FormatType) bool {
	return false
}

func (d *Decoder) setAttribute(name string, buffer interface{}, attr VertexAttribute) bool {
	return false
}

func (d *Decoder) setIndexInt32(buffer []uint32) {
	d.index.faces32 = buffer
}

func (d *Decoder) setIndexInt16(buffer []uint16) {
	d.index.faces16 = buffer
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
