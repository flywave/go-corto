package corto

type Decoder struct {
	nvert uint32
	nface uint32
	exif  map[string]string
	data  map[string]VertexAttribute
	index IndexAttribute

	stream       *InStream
	vertex_count uint32
}

func NewDecoder(input []byte) *Decoder {
	/**if((uintptr_t)input & 0x3)
		throw "Memory must be alignegned on 4 bytes.";

	stream.init(len, input);
	uint32_t magic = stream.readUint32();
	if(magic != 0x787A6300)
		throw "Not a crt file.";
	uint32_t version = stream.readUint32();
	stream.entropy = (Stream::Entropy)stream.readUint8();

	uint32_t size = stream.readUint32();
	for(uint32_t i = 0; i < size; i++) {
		const char *key = stream.readString();
		exif[key] = stream.readString();
	}

	int nattr = stream.readUint32();

	for(int i = 0; i < nattr; i++) {
		std::string name =  stream.readString();
		int codec = stream.readUint32();
		float q = stream.readFloat();
		uint32_t components = stream.readUint8();
		uint32_t format = stream.readUint8();
		uint32_t strategy = stream.readUint8();

		VertexAttribute *attr = nullptr;
		switch(codec) {
		case VertexAttribute::NORMAL_CODEC: attr = new NormalAttr(); break;
		case VertexAttribute::COLOR_CODEC: attr = new ColorAttr(components); break;

		case VertexAttribute::GENERIC_CODEC:
		default:  //
			attr = new GenericAttr<int>(components);
		}

		attr->q = q;
		attr->format = (VertexAttribute::Format)format;
		attr->strategy = strategy;
		data[name] = attr;
	}
	nvert = stream.readUint32();
	nface = stream.readUint32();**/

	return nil
}

func (d *Decoder) hasAttr(name string) bool {
	_, ok := d.data[name]
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
	/**if(data.find("color") == data.end()) return false;
	ColorAttr *attr = dynamic_cast<crt::ColorAttr *>(data["color"]);
	attr->format = VertexAttribute::UINT8;
	attr->buffer = (char *)buffer;
	attr->out_components = components;
	return true;**/
	return false
}

func (d *Decoder) setAttributeFormat(name string, buffer interface{}, format FormatType) bool {
	/**if(data.find(name) == data.end()) return false;
	VertexAttribute *attr = data[name];
	attr->format = format;
	attr->buffer = buffer;
	return true;**/
	return false
}

func (d *Decoder) setAttribute(name string, buffer interface{}, attr VertexAttribute) bool {
	/**	if(data.find(name) == data.end()) return false;
	VertexAttribute *found = data[name];
	attr->q = found->q;
	attr->strategy = found->strategy;
	attr->N = found->N;
	attr->buffer = buffer;
	delete data[name];
	data[name] = attr;
	return true;**/
	return false
}

func (d *Decoder) setIndexInt32(buffer []uint32) {
	d.index.faces32 = buffer
}

func (d *Decoder) setIndexInt16(buffer []uint16) {
	d.index.faces16 = buffer
}

func (d *Decoder) decode() {
	/**if(nface > 0)
		decodeMesh();
	else
		decodePointCloud();**/
}

func (d *Decoder) decodePointCloud() {
	/**std::vector<crt::Face> dummy;

	index.decodeGroups(stream);
	for(auto it: data)
		it.second->decode(nvert, stream);
	for(auto it: data)
		it.second->deltaDecode(nvert, dummy);
	for(auto it: data)
		it.second->dequantize(nvert);**/
}

func (d *Decoder) decodeMesh() {
	/**index.decodeGroups(stream);
		index.decode(stream);

		for(auto it: data)
			it.second->decode(nvert, stream);

		index.prediction.resize(nvert);

		uint32_t start = 0;
		uint32_t cler = 0; //keeps track of current cler
		for(Group &g: index.groups) {
			decodeFaces(start*3, g.end*3, cler);
			start = g.end;
		}

	#ifdef PRESERVED_UNREFERENCED
		//decode unreferenced vertices
		while(vertex_count < nvert) {
			int last = vertex_count-1;
			prediction[vertex_count++] = Face(last, last, last);
		}
	#endif

		for(auto it: data)
			it.second->deltaDecode(nvert, index.prediction);

		for(auto it: data)
			it.second->postDelta(nvert, nface, data, index);

		for(auto it: data)
			it.second->dequantize(nvert);**/
}

func (d *Decoder) decodeFaces(start, end uint32) uint32 {
	/**
	  vector<DEdge2> front;
	  front.reserve(index.max_front);

	  //faceorder is used to minimize split occourrence positioning in front and in back the new edges to be processed.
	  vector<int>faceorder;
	  faceorder.reserve((end - start)/2);
	  uint32_t order = 0;

	  //delayed again minimize split by further delay problematic splits
	  vector<int> delayed;

	  //TODO test if recording number of bits needed for splits improves anything. (very small but cost is zero.
	  int splitbits = ilog2(nvert) + 1;

	  int new_edge = -1; //last edge added which sohuld be the first to be processed, no need to store it in faceorder.

	  while(start < end) {
	  	if(new_edge == -1 && order >= faceorder.size() && !delayed.size()) {

	  		uint32_t last_index = vertex_count-1;
	  		int vindex[3];

	  		int split =  0; //bitmask for vertex already decoded/
	  		int c = index.clers[cler++];
	  		if(c == SPLIT) { //lookahead
	  			split = index.bitstream.read(3);
	  		} else
	  			assert(c == VERTEX);

	  		for(int k = 0; k < 3; k++) {
	  			int v; //TODO just use last_index.
	  			if(split & (1<<k)) {
	  				v = index.bitstream.read(splitbits);
	  			} else {
	  				assert(vertex_count < index.prediction.size());
	  				index.prediction[vertex_count] = Face(last_index, last_index, last_index);
	  				last_index = v = vertex_count++;
	  			}
	  			vindex[k] = v;
	  			if(index.faces16)
	  				index.faces16[start++] = v;
	  			else
	  				index.faces32[start++] = v;
	  		}
	  		int current_edge = front.size();
	  		faceorder.push_back(front.size());
	  		front.emplace_back(vindex[1], vindex[2], vindex[0], current_edge + 2, current_edge + 1);
	  		faceorder.push_back(front.size());
	  		front.emplace_back(vindex[2], vindex[0], vindex[1], current_edge + 0, current_edge + 2);
	  		faceorder.push_back(front.size());
	  		front.emplace_back(vindex[0], vindex[1], vindex[2], current_edge + 1, current_edge + 0);
	  		continue;
	  	}

	  	int f;
	  	if(new_edge != -1) {
	  		f = new_edge;
	  		new_edge = -1;

	  	} else if(order < faceorder.size()) {
	  		f = faceorder[order++];
	  	} else if (delayed.size()){
	  		f = delayed.back();
	  		delayed.pop_back(); //or popfront?

	  	} else {
	  		throw "Decoding topology failed";
	  	}

	  	const DEdge2 e = front[f];
	  	if(e.deleted) continue;
	  	//e.deleted = true; //each edge is processed once at most.

	  	int c = index.clers[cler++];
	  	if(c == BOUNDARY) continue;

	  	int v0 = e.v0;
	  	int v1 = e.v1;

	  	const DEdge2 previous_edge = front[e.prev];
	  	const DEdge2 next_edge = front[e.next];

	  	new_edge = front.size(); //index of the next edge to be added.
	  	int opposite = -1;

	  	if(c == VERTEX || c == SPLIT) {
	  		if(c == SPLIT) {
	  			opposite = index.bitstream.read(splitbits);
	  		} else {
	  			//Edge is inverted respect to encoding hence v1-v0 inverted.
	  			index.prediction[vertex_count] = Face(v1, v0, e.v2);
	  			opposite = vertex_count++;
	  		}
	  		assert(opposite < nvert);

	  		front[e.prev].next = new_edge;
	  		front[e.next].prev = new_edge + 1;

	  		front.emplace_back(v0, opposite, v1, e.prev, new_edge + 1);
	  		faceorder.push_back(front.size());
	  		front.emplace_back(opposite, v1, v0, new_edge, e.next);

	  	} else if(c == LEFT) {
	  		front[e.prev].deleted = true;
	  		front[previous_edge.prev].next = new_edge;
	  		front[e.next].prev = new_edge;
	  		opposite = previous_edge.v0;

	  		front.emplace_back(opposite, v1, v0, previous_edge.prev, e.next);

	  	} else if(c == RIGHT) {
	  		front[e.next].deleted = true;
	  		front[next_edge.next].prev = new_edge;
	  		front[e.prev].next = new_edge;
	  		opposite = next_edge.v1;

	  		front.emplace_back(v0, opposite, v1, e.prev, next_edge.next);

	  	} else if(c == DELAY) {
	  		//e.deleted = false;
	  		delayed.push_back(f);
	  		new_edge = -1;
	  		continue;

	  	} else if(c == END) {
	  		front[e.prev].deleted = true;
	  		front[e.next].deleted = true;
	  		front[previous_edge.prev].next = next_edge.next;
	  		front[next_edge.next].prev = previous_edge.prev;
	  		opposite = previous_edge.v0;
	  		new_edge = -1;

	  	} else {
	  		assert(0);
	  	}
	  	assert(v0 != v1);
	  	assert(v1 != opposite);
	  	assert(v0 != opposite);

	  	if(index.faces16) {
	  		index.faces16[start++] = v1;
	  		index.faces16[start++] = v0;
	  		index.faces16[start++] = opposite;
	  	} else {
	  		index.faces32[start++] = v1;
	  		index.faces32[start++] = v0;
	  		index.faces32[start++] = opposite;
	  	}
	  }**/
	return 0
}
