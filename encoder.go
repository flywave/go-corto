package corto

import (
	"math"
	"reflect"
	"unsafe"

	"github.com/flywave/go3d/vec3"
)

type Encoder struct {
	nvert          uint32
	nface          uint32
	exif           map[string]string
	index          *IndexAttribute
	groups         []Group
	data           map[string]VertexAttribute
	header_size    int
	stream         *OutStream
	current_vertex uint32
	last_index     uint32
	boundary       []bool
	encoded        []int
	prediction     []Quad
}

func NewEncoder(_nvert uint32, _nface uint32, entropy EntropyType) *Encoder {
	e := &Encoder{}
	e.nvert = _nvert
	e.nface = _nface
	e.header_size = 0
	e.current_vertex = 0
	e.last_index = 0
	e.stream.Entropy = entropy
	e.index.faces = make([]uint32, _nface*3)
	return e
}

func quantizationStep(nvert int, buffer []float32, bits int) float32 {
	var input []vec3.T
	cpointsHeader := (*reflect.SliceHeader)((unsafe.Pointer(&input)))
	cpointsHeader.Cap = int(nvert)
	cpointsHeader.Len = int(nvert)
	cpointsHeader.Data = uintptr(unsafe.Pointer(&buffer[0]))

	min := input[0]
	max := input[0]
	for i := 0; i < nvert; i++ {
		min.SetMin(input[i])
		max.SetMin(input[i])
	}

	max = vec3.Sub(&max, &min)
	intervals := math.Pow(2.0, float64(bits))
	max.Scale(1.0 / float32(intervals))
	q := math.Max(math.Max(float64(max[0]), float64(max[1])), float64(max[2]))
	return float32(q)
}

func (e *Encoder) addPositions(buffer []float32, q float32, o vec3.T) bool {
	coordsbuf := make([]byte, e.nvert*3*4)

	var coords []vec3.T
	coordsHeader := (*reflect.SliceHeader)((unsafe.Pointer(&coords)))
	coordsHeader.Cap = int(e.nvert)
	coordsHeader.Len = int(e.nvert)
	coordsHeader.Data = uintptr(unsafe.Pointer(&coordsbuf[0]))

	var input []vec3.T
	cpointsHeader := (*reflect.SliceHeader)((unsafe.Pointer(&input)))
	cpointsHeader.Cap = int(e.nvert)
	cpointsHeader.Len = int(e.nvert)
	cpointsHeader.Data = uintptr(unsafe.Pointer(&buffer[0]))

	for i := 0; i < int(e.nvert); i++ {
		coords[i] = vec3.Sub(&input[i], &o)
	}

	if q == 0 {
		max := vec3.MinVal
		min := vec3.MaxVal
		for i := 0; i < int(e.nvert); i++ {
			min.SetMin(coords[i])
			max.SetMax(coords[i])
		}
		max = vec3.Sub(&max, &min)
		q = float32(0.02 * math.Pow(float64(max[0]*max[1]*max[2]), 2.0/3.0) / float64(e.nvert))
	}
	strategy := CORRELATED
	if e.nface > 0 {
		strategy |= PARALLEL
	}

	return e.addAttributeFormat("position", coordsbuf, FORMAT_FLOAT, 3, q, strategy)
}

func (e *Encoder) addPositionsInt32(buffer []float32, _index []uint32, q float32, o vec3.T) bool {
	copy(e.index.faces, _index)

	var coords []vec3.T
	coordsHeader := (*reflect.SliceHeader)((unsafe.Pointer(&coords)))
	coordsHeader.Cap = int(e.nvert)
	coordsHeader.Len = int(e.nvert)
	coordsHeader.Data = uintptr(unsafe.Pointer(&buffer[0]))

	if q == 0 {
		average := float32(0)
		for f := 0; f < int(e.nface*3); f += 3 {
			sub := vec3.Sub(&coords[_index[f]], &coords[_index[f+1]])
			average += sub.Length()
		}
		q = (average / float32(e.nface)) / 20.0
	}
	return e.addPositions(buffer, q, o)
}

func (e *Encoder) addPositionsInt16(buffer []float32, _index []uint16, q float32, o vec3.T) bool {
	tmp := make([]uint32, e.nface*3)
	for i := 0; i < int(e.nface*3); i++ {
		tmp[i] = uint32(_index[i])
	}
	return e.addPositionsInt32(buffer, tmp[:], q, o)
}

func (e *Encoder) addPositionsBits(buffer []float32, bits int) bool {
	return e.addPositions(buffer, quantizationStep(int(e.nvert), buffer, bits), vec3.Zero)
}

func (e *Encoder) addPositionsBitsInt32(buffer []float32, index []uint32, bits int) bool {
	return e.addPositionsInt32(buffer, index, quantizationStep(int(e.nvert), buffer, bits), vec3.Zero)
}

func (e *Encoder) addPositionsBitsInt16(buffer []float32, index []uint16, bits int) bool {
	return e.addPositionsInt16(buffer, index, quantizationStep(int(e.nvert), buffer, bits), vec3.Zero)
}

func (e *Encoder) addNormals(buffer []float32, bits int, no PredictionType) bool {
	normal := NewNormalAttr(&bits)

	normal.format = FORMAT_FLOAT
	normal.prediction = no

	var bbuffer []byte
	bbufferHeader := (*reflect.SliceHeader)((unsafe.Pointer(&bbuffer)))
	bbufferHeader.Cap = int(len(buffer) * 4)
	bbufferHeader.Len = int(len(buffer) * 4)
	bbufferHeader.Data = uintptr(unsafe.Pointer(&buffer[0]))

	ok := e.addAttribute("normal", bbuffer, normal)
	return ok
}

func (e *Encoder) addNormalsInt16(buffer []uint16, bits int, no PredictionType) bool {
	tmpbuf := make([]float32, e.nvert*3)

	var tmp []vec3.T
	tmpHeader := (*reflect.SliceHeader)((unsafe.Pointer(&tmp)))
	tmpHeader.Cap = int(e.nvert)
	tmpHeader.Len = int(e.nvert)
	tmpHeader.Data = uintptr(unsafe.Pointer(&tmpbuf[0]))

	for i := 0; i < int(e.nvert); i++ {
		for k := 0; k < 3; k++ {
			tmp[i][k] = float32(buffer[i*3+k] / 32767.0)
		}
	}

	return e.addNormals(tmpbuf, bits, no)
}

func (e *Encoder) addColors(buffer []byte, rbits int, gbits, bbits int, abits int) bool {
	color := NewColorAttr(nil)
	color.setQ(rbits, gbits, bbits, abits)
	color.format = FORMAT_UINT8
	ok := e.addAttribute("color", buffer, color)
	return ok
}

func (e *Encoder) addColors3(buffer []byte, rbits int, gbits int, bbits int) bool {
	color := NewColorAttr(NewInt(3))
	color.setQ(rbits, gbits, bbits, 8)
	color.format = FORMAT_UINT8
	ok := e.addAttribute("color", buffer, color)
	return ok
}

func (e *Encoder) addUvs(buffer []float32, q float32) bool {
	uv := NewGenericAttr(2, make([]int32, 0), make([]int32, 0))
	uv.q = q
	uv.format = FORMAT_FLOAT

	var tmp []byte
	tmpHeader := (*reflect.SliceHeader)((unsafe.Pointer(&tmp)))
	tmpHeader.Cap = int(len(buffer) * 4)
	tmpHeader.Len = int(len(buffer) * 4)
	tmpHeader.Data = uintptr(unsafe.Pointer(&buffer[0]))

	ok := e.addAttribute("uv", tmp, uv)
	return ok
}

func (e *Encoder) addAttributeFormat(name string, buffer []byte, format FormatType, components int, q float32, strategy StrategyType) bool {
	if _, ok := e.data[name]; ok {
		return false
	}
	attr := NewGenericAttr(components, make([]int32, 0), make([]int32, 0))

	attr.q = q
	attr.strategy = strategy
	attr.format = format
	attr.Quantize(e.nvert, buffer)
	e.data[name] = attr
	return true
}

func (e *Encoder) addAttribute(name string, buffer []byte, attr VertexAttribute) bool {
	if _, ok := e.data[name]; ok {
		return true
	}
	attr.Quantize(e.nvert, buffer)
	e.data[name] = attr
	return true
}

func (e *Encoder) addGroup(end_triangle int) {
	e.index.groups = append(e.index.groups, Group{end: uint32(end_triangle)})
}

func (e *Encoder) addGroupAttr(end_triangle int, props map[string]string) {
	g := Group{end: uint32(end_triangle)}
	g.properties = props
	e.index.groups = append(e.index.groups, g)
}

func (e *Encoder) encode() {
	e.stream.reserve(int(e.nvert))

	e.stream.write(uint32(0x787A6300))
	e.stream.write(uint32(0x1)) //version
	e.stream.write(byte(e.stream.Entropy))

	e.stream.write(uint32(len(e.exif)))
	for first, second := range e.exif {
		e.stream.writeString(first)
		e.stream.writeString(second)
	}

	e.stream.write(int32(len(e.data)))
	for first, second := range e.data {
		e.stream.writeString(first) //name
		e.stream.write(int32(second.Codec()))
		e.stream.write(float32(second.Q()))
		e.stream.write(second.Dim())
		e.stream.write(second.Format())
		e.stream.write(second.Strategy())
	}

	if e.nface > 0 {
		e.encodeMesh()
	} else {
		e.encodePointCloud()
	}
}

func (e *Encoder) encodePointCloud() {
	/**
	if(data.find("position") == data.end())
		throw "No position attribute found. Use DIFF normal strategy instead.";

	GenericAttr<int> *coord = dynamic_cast<GenericAttr<int> *>(data["position"]);
	if(!coord)
		throw "Position attr has been overloaded, Use DIFF normal strategy instead.";

	Point3i *coords = (Point3i *)coord->values.data();

	std::vector<ZPoint> zpoints(nvert);

	Point3i min(0, 0, 0);
	for(uint32_t i = 0; i < nvert; i++) {
		Point3i &q = coords[i];
		min.setMin(q);
	}
	for(uint32_t i = 0; i < nvert; i++) {
		Point3i q = coords[i] - min;
		zpoints[i] = ZPoint(q[0], q[1], q[2], 21, i);
	}
	sort(zpoints.rbegin(), zpoints.rend());//, greater<ZPoint>());

	int count = 0;
	for(unsigned int i = 1; i < nvert; i++) {
		if(zpoints[i] == zpoints[count])
			continue;
		count++;
		zpoints[count] = zpoints[i];
	}
	count++;
	nvert = count;
	zpoints.resize(nvert);

	header_size = stream.elapsed();

	stream.write<uint32_t>(nvert);
	stream.write<uint32_t>(0); //nface

	index.encodeGroups(stream);

	prediction.resize(nvert);
	prediction[0] = Quad(zpoints[0].pos, -1, -1, -1);
	for(uint32_t i = 1; i < nvert; i++)
		prediction[i] = Quad(zpoints[i].pos, zpoints[i-1].pos, zpoints[i-1].pos, zpoints[i-1].pos);

	for(auto it: data)
		it.second->preDelta(nvert, nface, data, index);

	for(auto it: data)
		it.second->deltaEncode(prediction);

	for(auto it: data)
		it.second->encode(nvert, stream);**/
}

func (e *Encoder) encodeMesh() {
	/**encoded.resize(nvert, -1);

		if(!index.groups.size()) index.groups.push_back(nface);
		//remove degenerate faces
		uint32_t start =  0;
		uint32_t count = 0;
		for(Group &g: index.groups) {
			for(uint32_t i = start; i < g.end; i++) {
				uint32_t *f = &index.faces[i*3];

				if(f[0] == f[1] || f[0] == f[2] || f[1] == f[2])
					continue;

				if(count != i) {
					uint32_t *dest = &index.faces[count*3];
					dest[0] = f[0];
					dest[1] = f[1];
					dest[2] = f[2];
				}
				count++;
			}
			start = g.end;
			g.end = count;
		}
		index.faces.resize(count*3);
		nface = count;

		//BitStream bitstream(nvert/4);
		index.bitstream.reserve(nvert/4);
		prediction.resize(nvert);

		start =  0;
		for(Group &g: index.groups) {
			encodeFaces(start, g.end);
			start = g.end;
		}
	#ifdef PRESERVED_UNREFERENCED
		//encoding unreferenced vertices
		for(uint32_t i = 0; i < nvert; i++) {
			if(encoded[i] != -1)
				continue;
			int last = current_vertex-1;
			prediction.emplace_back(Quad(i, last, last, last));
			current_vertex++;
		}
	#endif

		//predelta works using the original indexes, we will deal with unreferenced vertices later (due to prediction resize)
		for(auto it: data)
			it.second->preDelta(nvert, nface, data, index);

		//cout << "Unreference vertices: " << nvert - current_vertex << " remaining: " << current_vertex << endl;
		nvert = current_vertex;
		prediction.resize(nvert);



		for(auto it: data)
			it.second->deltaEncode(prediction);

		stream.write<int>(nvert);
		stream.write<int>(nface);
		header_size = stream.elapsed();
		index.encodeGroups(stream);
		index.encode(stream);

		for(auto it: data)
			it.second->encode(nvert, stream);**/
}

/**
class McFace {
	public:
		uint32_t f[3];
		uint32_t t[3]; //topology: opposite face
		uint32_t i[3]; //index in the opposite face of this face: faces[f.t[k]].t[f.i[k]] = f;
		McFace(uint32_t v0 = 0, uint32_t v1 = 0, uint32_t v2 = 0) {
			f[0] = v0; f[1] = v1; f[2] = v2;
			t[0] = t[1] = t[2] = 0xffffffff;
		}
		bool operator<(const McFace &face) const {
			if(f[0] < face.f[0]) return true;
			if(f[0] > face.f[0]) return false;
			if(f[1] < face.f[1]) return true;
			if(f[1] > face.f[1]) return false;
			return f[2] < face.f[2];
		}
		bool operator>(const McFace &face) const {
			if(f[0] > face.f[0]) return true;
			if(f[0] < face.f[0]) return false;
			if(f[1] > face.f[1]) return true;
			if(f[1] < face.f[1]) return false;
			return f[2] > face.f[2];
		}
	};

	class CEdge { //compression edges
	public:
		uint32_t face;
		uint32_t side; //opposited to side vertex of face (edge 2 opposite to vertex 2)
		uint32_t prev, next;
		bool deleted;
		CEdge(uint32_t f = 0, uint32_t s = 0, uint32_t p = 0, uint32_t n = 0):
			face(f), side(s), prev(p), next(n), deleted(false) {}
	};


	class McEdge { //topology edges
	public:
		uint32_t face, side;
		uint32_t v0, v1;
		bool inverted;
		McEdge() {}
		McEdge(uint32_t _face, uint32_t _side, uint32_t _v0, uint32_t _v1): face(_face), side(_side), inverted(false) {
			if(_v0 < _v1) {
				v0 = _v0; v1 = _v1;
				inverted = false;
			} else {
				v1 = _v0; v0 = _v1;
				inverted = true;
			}
		}
		bool operator==(const McEdge &e) const {
			return face == e.face && v0 == e.v0 && v1 == e.v1;
		}

		bool operator<(const McEdge &edge) const {
			if(v0 < edge.v0) return true;
			if(v0 > edge.v0) return false;
			return v1 < edge.v1;
		}
		bool match(const McEdge &edge) const {
			if(inverted && edge.inverted) return false;
			if(!inverted && !edge.inverted) return false;
			return v0 == edge.v0 && v1 == edge.v1;
		}
	};

	static void buildTopology(vector<McFace> &faces, uint32_t nvert) {
		//compute buckets size for edges with lower vertex in common
		vector<uint32_t> count(nvert, 0);
		for(McFace &face: faces) {
			count[min(face.f[0], face.f[1])]++;
			count[min(face.f[1], face.f[2])]++;
			count[min(face.f[2], face.f[0])]++;
		}

		//find start of every bucket.
		uint32_t partial = 0;
		for(uint32_t &c: count) {
			uint32_t tmp = c;
			c = partial;
			partial += tmp;
		}

		//write edges in the buckets.
		vector<McEdge> edges(faces.size()*3);
		for(size_t i = 0; i < faces.size(); i++) {
			McFace &face = faces[i];
			uint32_t v0 = min(face.f[1], face.f[2]);
			edges[count[v0]++] = McEdge(i, 0, face.f[1], face.f[2]);

			uint32_t v1 = min(face.f[2], face.f[0]);
			edges[count[v1]++] = McEdge(i, 1, face.f[2], face.f[0]);

			uint32_t v2 = min(face.f[0], face.f[1]);
			edges[count[v2]++] = McEdge(i, 2, face.f[0], face.f[1]);
		}

		//sort the buckets
		sort(edges.begin(), edges.begin() + count[0]);
		for(uint32_t i = 0; i < count.size()-1; i++) {
			if(count[i] == 0) continue;
			sort(edges.begin() + count[i], edges.begin() + count[i+1]);
		}

		McEdge previous(0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff);

		//identify opposite edges
		for(const McEdge &edge: edges) {
			if(edge.match(previous)) {
				uint32_t &edge_side_face = faces[edge.face].t[edge.side];
				uint32_t &previous_side_face = faces[previous.face].t[previous.side];
				if(edge_side_face == 0xffffffff && previous_side_face == 0xffffffff) {
					edge_side_face = previous.face;
					faces[edge.face].i[edge.side] = previous.side;
					previous_side_face = edge.face;
					faces[previous.face].i[previous.side] = edge.side;
				}
			} else
				previous = edge;
		}
	}

	static int next_(int t) {
		t++;
		if(t == 3) t = 0;
		return t;
	}

	static uint32_t countReferenced(vector<uint32_t> &faces, uint32_t nvert) {
		vector<bool> referenced(nvert, false);
		for(auto &i: faces)
			referenced[i] = true;
		uint32_t count = 0;
		for(bool b: referenced)
			if(b) count++;
		return count;
	}
	**/
func (e *Encoder) encodeFaces(start, end int) {
	/**
		vector<McFace> faces(end - start);
		for(int i = start; i < end; i++) {
			uint32_t * f = &index.faces[i*3];
			faces[i - start] = McFace(f[0], f[1], f[2]);
			assert(f[0] != f[1] && f[1] != f[2] && f[2] != f[0]);
		}

		buildTopology(faces, nvert);

		unsigned int current = 0;          //keep track of connected component start

		vector<int> delayed;
		//TODO move to vector + order
		vector<int> faceorder;
		faceorder.reserve(end - start);
		uint32_t order = 0;
		vector<CEdge> front;
		front.reserve(end - start);

		vector<bool> visited(faces.size(), false);
		unsigned int totfaces = faces.size();

		//unreferenced vertices will not be saved, we need to know the number of referenced vertices before computing splitbits
		uint32_t nreferenced = countReferenced(index.faces, nvert);
		int splitbits = ilog2(nreferenced) + 1;

		int new_edge = -1;
		int counting = 0;
		while(totfaces > 0) {
			if(new_edge == -1 && order >= faceorder.size() && !delayed.size()) {

				while(current != faces.size()) {   //find first triangle non visited
					if(!visited[current]) break;
					current++;
				}
				if(current == faces.size()) break; //no more faces to encode exiting

				//encode first face: 3 vertices indexes, and add edges
				unsigned int current_edge = front.size();
				McFace &face = faces[current];

				int split = 0;
				for(int k = 0; k < 3; k++) {
					int vindex = face.f[k];
					if(encoded[vindex] != -1)
						split |= 1<<k;
				}

				if(split) {
					index.clers.push_back(SPLIT);
					index.bitstream.write(split, 3);
				} else
					index.clers.push_back(VERTEX);

				for(int k = 0; k < 3; k++) {
					uint32_t vindex = face.f[k];
					assert(vindex < nvert);
					int &enc = encoded[vindex];

					if(enc != -1) {
						index.bitstream.write(enc, splitbits);
					} else {
						//quad uses presorting indexing. (diff in attribute are sorted, values are not).
						assert(last_index < nvert);
						prediction[current_vertex] = Quad(vindex, last_index, last_index, last_index);
						enc = current_vertex++;
						last_index = vindex;
					}
				}

				faceorder.push_back(front.size());
				front.emplace_back(current, 0, current_edge + 2, current_edge + 1);
				faceorder.push_back(front.size());
				front.emplace_back(current, 1, current_edge + 0, current_edge + 2);
				faceorder.push_back(front.size());
				front.emplace_back(current, 2, current_edge + 1, current_edge + 0);


				counting++;
				visited[current] = true;
				current++;
				totfaces--;
				continue;
			}
			int c;
			if(new_edge != -1) {
				c = new_edge;
				new_edge = -1;

			} else if(order < faceorder.size()) {
				c =  faceorder[order++];

			} else if(delayed.size()) {
				c = delayed.back();
				delayed.pop_back();


			} else {
				throw "Decoding topology failed";
			}
			CEdge &e = front[c];
			if(e.deleted) continue;
			//e.deleted = true;

			//opposite face is the triangle we are encoding
			uint32_t opposite_face = faces[e.face].t[e.side];
			int opposite_side = faces[e.face].i[e.side];

			if(opposite_face == 0xffffffff || visited[opposite_face]) { //boundary edge or glue
				index.clers.push_back(BOUNDARY);
				continue;
			}

			assert(opposite_face < faces.size());
			McFace &face = faces[opposite_face];

			int k2 = opposite_side;
			int k0 = next_(k2);
			int k1 = next_(k0);

			//check for closure on previous or next edge
			int eprev = e.prev;
			int enext = e.next;
			assert(eprev >= 0);
			assert(enext < (int)front.size());
			assert(eprev < (int)front.size());
			const CEdge previous_edge = front[eprev];
			const CEdge next_edge = front[enext];

			bool close_left = (faces[previous_edge.face].t[previous_edge.side] == opposite_face);
			bool close_right = (faces[next_edge.face].t[next_edge.side] == opposite_face);

			new_edge = front.size(); //index of the next edge to be added.

			if(close_left && close_right) {
				index.clers.push_back(END);
				front[eprev].deleted = true;
				front[enext].deleted = true;
				front[previous_edge.prev].next = next_edge.next;
				front[next_edge.next].prev = previous_edge.prev;
				new_edge = -1;

			} else if(close_left) {
				index.clers.push_back(LEFT);
				front[eprev].deleted = true;
				front[previous_edge.prev].next = new_edge;
				front[enext].prev = new_edge;

				front.emplace_back(opposite_face, k1, previous_edge.prev, enext);

			} else if(close_right) {
				index.clers.push_back(RIGHT);
				front[enext].deleted = true;
				front[next_edge.next].prev = new_edge;
				front[eprev].next = new_edge;

				front.emplace_back(opposite_face, k0, eprev, next_edge.next);

			} else {
				int v0 = face.f[k0];
				int v1 = face.f[k1];
				int opposite = face.f[k2];

				if(encoded[opposite] != -1 && order < faceorder.size()) { //split, but we can still delay it.
					e.deleted = false; //undelete it.
					delayed.push_back(c);
					index.clers.push_back(DELAY);
					new_edge = -1;
					continue;
				}
				if(encoded[opposite] != -1) {
					index.clers.push_back(SPLIT);
					index.bitstream.write(encoded[opposite], splitbits);

				} else {
					index.clers.push_back(VERTEX);
					//vertex needed for parallelogram prediction
					int v2 = faces[e.face].f[e.side];
					prediction[current_vertex] = Quad(opposite, v0, v1, v2);
					encoded[opposite] = current_vertex++;
					last_index = opposite;
				}

				front[eprev].next = new_edge;
				front[enext].prev = new_edge + 1;

				front.emplace_back(opposite_face, k0, eprev, new_edge+1);
				faceorder.push_back(front.size());
				front.emplace_back(opposite_face, k1, new_edge, enext);
			}

			counting++;
			assert(!visited[opposite_face]);
			visited[opposite_face] = true;
			totfaces--;
		}
		index.max_front = std::max(index.max_front, (uint32_t)front.size());
	**/
}
