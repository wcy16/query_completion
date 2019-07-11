package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type RMQ struct {
	Data  []int   // data
	Index [][]int // store range max index
}

// max of a range of positive scores
func max(scores []int, begin int, end int) (int, int) {
	m := -1
	id := -1
	for i := begin; i != end; i++ {
		if m < scores[i] {
			m = scores[i]
			id = i
		}
	}

	return id, m
}

// create range maximum queries
func CreateRMQ(scores []int) *RMQ {
	r := new(RMQ)
	l := len(scores)
	r.Index = make([][]int, 0, l)

	// not a very efficient one, maybe can cache and reuse some values
	for i := 0; i != l; i++ {
		tmp := make([]int, 0, 1)

		for pos, gap, m, id := i, 1, -1, -1; i+gap <= l; gap = gap * 2 {
			if id2, m2 := max(scores, pos, i+gap); m2 > m {
				m = m2
				id = id2
			}
			tmp = append(tmp, id)
			pos = i + gap
		}

		r.Index = append(r.Index, tmp)
	}

	r.Data = scores
	return r
}

// save the rmq data to file
func (r *RMQ) Save(path string) {
	if data, err := json.Marshal(r); err == nil {
		if fp, err := os.Create(path); err == nil {
			_, _ = fp.Write(data)
		} else {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
}

// load the rmq
func (r *RMQ) Load(path string) {
	if data, err := ioutil.ReadFile(path); err == nil {
		err = json.Unmarshal(data, r)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
}

// return top 1 index from begin (inclusive) to end (exclusive)
func (r *RMQ) Top1(begin int, end int) int {
	if begin >= end {
		return -1
	}

	// find gaps
	gap := 1
	counter := 0
	for ; begin+gap < end; gap = gap << 1 {
		counter++
	}

	if begin+gap == end {
		return r.Index[begin][counter]
	} else {
		counter--
		gap = gap >> 1
		id1 := r.Index[begin][counter]
		id2 := r.Index[end-gap][counter]
		m1 := r.Data[id1]
		m2 := r.Data[id2]
		if m1 < m2 {
			return id2
		} else {
			return id1
		}
	}

}

type indexRange struct {
	begin    int
	end      int
	maxIndex int
}

// find the top-k results from begin (inclusive) to end (exclusive)
func (r *RMQ) TopK(begin int, end int, k int) []int {
	if begin == end {
		return make([]int, 0)
	}

	topk := make([]int, 0, k)
	buf := make(map[indexRange]int)

	id := r.Top1(begin, end)
	initialRange := indexRange{
		begin,
		end,
		id,
	}
	buf[initialRange] = r.Data[id]

	// begin iteration
	var split indexRange
	for ; k != 0 && len(buf) != 0; k-- {
		m := -1
		for k, v := range buf {
			if v > m {
				m = v
				split = k
			}
		}

		topk = append(topk, split.maxIndex)

		// split range
		if split.begin < split.maxIndex {
			id1 := r.Top1(split.begin, split.maxIndex)
			r1 := indexRange{split.begin, split.maxIndex, id1}
			buf[r1] = r.Data[id1]
		}
		if tmp := split.maxIndex + 1; tmp < split.end {
			id2 := r.Top1(tmp, split.end)
			r2 := indexRange{tmp, split.end, id2}
			buf[r2] = r.Data[id2]
		}
		delete(buf, split)

	}

	return topk
}
