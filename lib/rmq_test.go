package lib

import (
	"math/rand"
	"sort"
	"testing"
)

func TestRMQ_Top1(t *testing.T) {
	for i := 0; i != 1000; i++ {
		scores := rand.Perm(i)
		r := CreateRMQ(scores)
		got := r.Top1(0, i)
		if i == 0 {
			if got != -1 {
				t.Errorf("TestRMQ_Top1; want -1; got %d", got)
			}
		} else {
			if scores[got] != i-1 {
				t.Errorf("TestRMQ_Top1; want %d; got %d", i-1, scores[got])
			}
		}

	}
}

func TestRMQ_TopK(t *testing.T) {
	const n = 100
	const maxTopK = 50

	scores := rand.Perm(n)
	//fmt.Printf("%v\n", scores)
	r := CreateRMQ(scores)
	for i := 0; i != n; i++ {
		for j := i + 1; j != n; j++ {
			c := make([]int, j-i)
			copy(c, scores[i:j])
			sort.Ints(c)
			// revert c
			for x := len(c)/2 - 1; x >= 0; x-- {
				opp := len(c) - 1 - x
				c[x], c[opp] = c[opp], c[x]
			}

			// from top 2
			for k := 2; k <= maxTopK; k++ {
				topk := r.TopK(i, j, k)
				if !compareSlices(topk, c, scores) {
					t.Errorf("TestRMQ_TopK; i=%d; j=%d; k=%d; got %v; fail", i, j, k, topk)
					//t.Fail()
				}
			}
		}
	}
}

// 6332 ns/op
func BenchmarkRMQ_TopK(b *testing.B) {
	queries := 1000000
	topk := 10
	// initialize
	scores := rand.Perm(queries)
	r := CreateRMQ(scores)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		begin := rand.Intn(queries)
		end := rand.Intn(queries)
		if begin < end {
			r.TopK(begin, end, topk)
		} else if begin > end {
			r.TopK(end, begin, topk)
		} else {
			r.TopK(begin, end+1, topk)
		}
	}
}

// 381600148800 ns/op
func BenchmarkCreateRMQ(b *testing.B) {
	queries := 1000000
	// initialize
	scores := generateScores(queries)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateRMQ(scores)
	}
}

// helper to compare
// s1: index of topk scores
// s2: topk scores
func compareSlices(s1, s2, scores []int) bool {
	l1 := len(s1)
	l2 := len(s2)
	if l1 > l2 {
		return false
	}

	for i := 0; i != l1; i++ {
		if !(scores[s1[i]] == s2[i]) {
			return false
		}
	}
	return true
}

func generateScores(length int) []int {
	scores := make([]int, length)

	for i := 0; i != length; i++ {
		scores[i] = rand.Int()
	}

	return scores
}
