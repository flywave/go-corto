package corto

import (
	"bytes"
	"encoding/binary"
	"math"
	"sort"
)

const rle_limit int = 255

var byteorder = binary.LittleEndian

type Symbol struct {
	symbol      byte
	probability byte
}

func (s *Symbol) less(t *Symbol) bool {
	return s.probability > t.probability
}

type TSymbol struct {
	offset      int
	length      int
	probability uint32
}

type Symbols []Symbol

func (p Symbols) Len() int { return len(p) }

func (p Symbols) Less(i, j int) bool {
	return p[i].probability > p[j].probability
}

func (p Symbols) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p Symbols) Data() []byte {
	ret := make([]byte, len(p)*2)
	writer := bytes.NewBuffer(ret)
	for i := range p {
		binary.Write(writer, byteorder, p[i])
	}
	return writer.Bytes()
}

type Tunstall struct {
	wordsize       int
	dictionarysize int
	probabilities  Symbols
	index          []int
	lengths        []int
	table          []byte
	lookup_size    int
	offsets        []int
	remap          []byte
}

func NewTunstall(wordsize, lookup int) *Tunstall {
	return &Tunstall{wordsize: wordsize, lookup_size: lookup}
}

func (t *Tunstall) getProbabilities(data []byte, size int) {
	t.probabilities = t.probabilities[:0]
	probs := [256]int{0}

	for i := 0; i < size; i++ {
		probs[data[i]]++
	}

	for i := 0; i < len(probs); i++ {
		if probs[i] > 0 {
			t.probabilities = append(t.probabilities, Symbol{symbol: byte(i), probability: byte(probs[i] * 255 / size)})
		}
	}

	sort.Sort(t.probabilities)
}

func (t *Tunstall) setProbabilities(probs []float32, n_symbols int) {
	t.probabilities = t.probabilities[:0]
	for i := 0; i < n_symbols; i++ {
		if probs[i] <= 0 {
			continue
		}
		t.probabilities = append(t.probabilities, Symbol{symbol: byte(i), probability: byte(probs[i] * 255)})
	}
}

func resize(buf []byte, si int) []byte {
	if si < len(buf) {
		return buf[:si]
	}
	new := make([]byte, si)
	copy(new, buf)
	return new
}

func (t *Tunstall) createDecodingTables() {
	n_symbols := len(t.probabilities)
	if n_symbols < 1 {
		return
	}

	queues := make([]*Deque, n_symbols)
	var buffer []byte

	for i := 0; i < n_symbols; i++ {
		var s TSymbol
		s.probability = uint32(t.probabilities[i].probability) << 8
		s.offset = len(buffer)
		s.length = 1

		if queues[i] == nil {
			queues[i] = NewDeque()
		}

		queues[i].PushBack(&s)
		buffer = append(buffer, t.probabilities[i].symbol)
	}

	dictionary_size := 1 << t.wordsize
	n_words := n_symbols
	table_length := n_symbols
	for n_words < dictionary_size-n_symbols+1 {
		best := int(0)
		max_prob := uint32(0)

		for i := 0; i < n_symbols; i++ {
			p := queues[i].Front().(*TSymbol).probability
			if p > max_prob {
				best = i
				max_prob = p
			}
		}

		symbol := queues[best].Front().(*TSymbol)
		pos := len(buffer)

		buffer = resize(buffer, pos+n_symbols*(symbol.length+1))
		for i := 0; i < n_symbols; i++ {
			p := t.probabilities[i].probability
			var s TSymbol
			s.probability = ((symbol.probability * uint32(p) << 8) >> 16)
			s.offset = pos
			s.length = symbol.length + 1

			copy(buffer[pos:], buffer[symbol.offset:symbol.offset+symbol.length])

			pos += symbol.length
			buffer[pos] = t.probabilities[i].symbol
			pos++
			queues[i].PushBack(&s)
		}
		table_length += (n_symbols-1)*(symbol.length+1) + 1
		n_words += n_symbols - 1
		queues[best].PopFront()
	}
	t.index = make([]int, n_words)
	t.lengths = make([]int, n_words)
	t.table = make([]byte, table_length)

	word := 0
	pos := 0
	for i := 0; i < len(queues); i++ {
		queue := queues[i]
		for k := 0; k < queue.Len(); k++ {
			s := queue.At(k).(*TSymbol)
			t.index[word] = pos
			t.lengths[word] = s.length
			word++
			copy(t.table[pos:], buffer[s.offset:s.offset+s.length])
			pos += s.length
		}
	}
	if len(t.index) > dictionary_size {
		panic("error")
	}
}

func (t *Tunstall) createDecodingTables2() {
	n_symbols := len(t.probabilities)
	if n_symbols <= 1 {
		return
	}

	dictionary_size := 1 << t.wordsize

	queues := make([]uint32, 2*dictionary_size)

	t.index = make([]int, 2*dictionary_size)
	t.lengths = make([]int, 2*dictionary_size)

	end := 0
	t.table = resize(t.table, 8192)
	if t.wordsize != 8 {
		panic("error")
	}

	buffer := t.table[:]
	pos := uint32(0)
	starts := make([]uint32, n_symbols)

	n_words := uint32(0)

	count := uint32(2)
	p0 := uint32(t.probabilities[0].probability) << 8
	p1 := uint32(t.probabilities[1].probability) << 8
	prob := (p0 * p0) >> 16
	max_count := uint32((dictionary_size - 1) / (n_symbols - 1))
	for prob > p1 && count < max_count {
		prob = (prob * p0) >> 16
		count++
	}

	if count >= 16 {
		buffer[pos] = t.probabilities[0].symbol
		pos++
		for k := 1; k < n_symbols; k++ {
			for i := 0; i < int(count-1); i++ {
				buffer[pos] = t.probabilities[0].symbol
				pos++
			}
			buffer[pos] = t.probabilities[k].symbol
			pos++
		}
		starts[0] = uint32((int(count) - 1) * n_symbols)
		for k := 1; k < n_symbols; k++ {
			starts[k] = uint32(k)
		}

		for col := 0; col < int(count); col++ {
			for row := 1; row < n_symbols; row++ {
				dest := row + col*n_symbols
				probability := &queues[dest]
				if col == 0 {
					*probability = uint32(t.probabilities[row].probability) << 8
				} else {
					*probability = prob * (uint32(t.probabilities[row].probability) << 8) >> 16
				}
				t.index[dest] = row*int(count) - col
				t.lengths[dest] = col + 1
			}
			if col == 0 {
				prob = p0
			} else {
				prob = (prob * p0) >> 16
			}
		}

		first := (int(count) - 1) * n_symbols
		queues[first] = uint32(prob)
		t.index[first] = 0
		t.lengths[first] = int(count)
		n_words = uint32(1 + int(count)*(n_symbols-1))
		end = int(count) * n_symbols
		if n_words != pos {
			panic("error")
		}
	} else {
		n_words = uint32(n_symbols)
		for i := 0; i < n_symbols; i++ {
			starts[i] = uint32(i)
			queues[end] = uint32(t.probabilities[i].probability) << 8
			t.index[end] = int(pos)
			t.lengths[end] = 1
			end++
			buffer[pos] = t.probabilities[i].symbol
			pos++
		}
	}

	for int(n_words) < dictionary_size {
		best := uint32(0)
		max_prob := uint32(0)
		for i := 0; i < n_symbols; i++ {
			p := queues[starts[i]]
			if p > max_prob {
				best = uint32(i)
				max_prob = p
			}
		}
		symbol := starts[best]
		probability := queues[symbol]
		offset := t.index[symbol]
		length := t.lengths[symbol]
		r := uint32(0)
		for ; int(r) < n_symbols; r++ {
			p := t.probabilities[r].probability
			queues[end] = ((probability * uint32(p) << 8) >> 16)
			t.index[end] = int(pos)
			t.lengths[end] = length + 1
			end++

			if (int(pos) + length) >= len(buffer) {
				panic("error")
			}
			copy(buffer[pos:], buffer[offset:offset+length])
			pos += uint32(length)
			buffer[pos] = t.probabilities[r].symbol
			pos++
			if int(n_words+r) == dictionary_size-1 {
				break
			}
		}
		if int(r) == n_symbols {
			starts[best] += uint32(n_symbols)
		}
		n_words += uint32(n_symbols) - 1
	}

	word := 0
	row := 0
	for i := 0; i < end; i++ {
		if row >= n_symbols {
			row = 0
		}
		if int(starts[row]) > i {
			continue
		}

		t.index[word] = t.index[i]
		t.lengths[word] = t.lengths[i]
		word++
		row++
	}

	t.index = t.index[:dictionary_size]
	t.lengths = t.lengths[:dictionary_size]
}

func (t *Tunstall) createEncodingTables() {
	n_symbols := len(t.probabilities)

	if n_symbols <= 1 {
		return
	}

	lookup_table_size := 1
	for i := 0; i < t.lookup_size; i++ {
		lookup_table_size *= n_symbols
	}

	t.remap = make([]byte, 256)

	for i := 0; i < n_symbols; i++ {
		s := t.probabilities[i]
		t.remap[s.symbol] = byte(i)
	}

	if int(t.probabilities[0].probability) > rle_limit {
		return
	}

	t.offsets = make([]int, lookup_table_size)
	for i := 0; i < lookup_table_size; i++ {
		t.offsets[i] = 0xffffff
	}

	for i := 0; i < len(t.index); i++ {
		var low, high int
		offset := 0
		table_offset := 0
		for {
			low, high = t.wordCode(t.table[t.index[i]+offset:], t.lengths[i]-offset)
			if t.lengths[i]-offset <= t.lookup_size {
				for k := low; k < high; k++ {
					t.offsets[table_offset+k] = i
				}
				break
			}

			w := t.offsets[table_offset+low]

			if w >= 0 {
				t.offsets[table_offset+low] = -len(t.offsets)

				newoffsets := make([]int, len(t.offsets)+lookup_table_size)
				copy(newoffsets, t.offsets)
				for i := len(t.offsets); i < len(newoffsets); i++ {
					newoffsets[i] = w
				}
			}
			table_offset = -t.offsets[table_offset+low]
			offset += t.lookup_size
		}
	}
}

func (t *Tunstall) compress(data []byte, input_size int) (output []byte, output_size int) {
	if len(t.probabilities) == 1 {
		output_size = 0
		return nil, 0
	}

	output = make([]byte, input_size*2)

	if t.wordsize > 16 {
		panic("error")
	}

	output_size = 0
	input_offset := 0
	word_offset := 0
	offset := 0

	for input_offset < input_size {
		d := input_size - input_offset
		if d > t.lookup_size {
			d = t.lookup_size
		}
		low, _ := t.wordCode(data[input_offset:], d)
		offset = t.offsets[-offset+low]
		if offset == 0xffffff {
			panic("error")
		}

		if offset >= 0 {
			output[output_size] = byte(offset)
			output_size++
			input_offset += t.lengths[offset] - word_offset
			offset = 0
			word_offset = 0
		} else {
			word_offset += t.lookup_size
			input_offset += t.lookup_size
		}
	}

	if offset < 0 {
		for offset < 0 {
			offset = t.offsets[-offset]
		}
		output[output_size] = byte(offset)
		output_size++
	}
	if output_size > input_size*2 {
		panic("error")
	}
	return
}

func (t *Tunstall) decompressWithSize(data []byte, output []byte, output_size int) []byte {
	end_output := output_size
	end_data := len(data) - 1

	if len(t.probabilities) == 1 {
		for i := 0; i < output_size; i++ {
			output[i] = t.probabilities[0].symbol
		}
		return output
	}

	startoutput := 0
	startdata := 0
	for ; startdata < end_data; startdata++ {
		symbol := data[startdata]
		start := t.index[symbol]
		length := t.lengths[symbol]
		copy(output[startoutput:], t.table[start:start+length])
		startoutput += length
	}

	symbol := data[startdata]
	start := t.index[symbol]
	length := (end_output - startdata)
	copy(output[startoutput:], t.table[start:start+length])
	return output
}

func (t *Tunstall) entropy() float32 {
	e := float32(0)
	for i := 0; i < len(t.probabilities); i++ {
		p := float64(t.probabilities[i].probability / 255.0)
		e += float32(p * math.Log(p) / math.Log(2))
	}
	return -e
}

func toUint(i int) uint32 {
	i *= 2
	if i < 0 {
		i = -i - 1
	}
	return uint32(i)
}

func toInt(i uint32) int {
	k := int(i)
	if k&0x1 > 0 {
		k = (-k - 1) / 2
	} else {
		k /= 2
	}
	return k
}

func (t *Tunstall) wordCode(w []byte, length int) (low, high int) {
	n_symbols := len(t.probabilities)

	c := 0
	for i := 0; i < length && i < t.lookup_size; i++ {
		c *= n_symbols
		c += int(t.remap[w[i]])
	}
	low = c
	high = c
	high++
	for i := length; i < t.lookup_size; i++ {
		low *= n_symbols
		high *= n_symbols
	}
	return
}

func roundUp(v int) int {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
