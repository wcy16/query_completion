package lib

import (
	"math/rand"
	"strings"
	"testing"
)

var sortedQueries = []string{
	"ada",
	"apple",
	"artificial",
	"artificial intelligence",

	"basic",

	"c plus plus",
	"c programming language",

	"data",
	"database",
	"data warehousing",
	"delete",

	"mab",
	"mac",
	"mac pro",
	"machine learning",
	"memory",

	"nearest neighbour",
	"net",
	"network",
	"networking",
	"new",
	"not",
}

func TestTrie(t *testing.T) {
	trie := NewTrie(-1, -1)
	// insert
	for id, ele := range sortedQueries {
		trie.Insert([]rune(ele), id)
	}

	// iterate
	for _, ele := range sortedQueries {
		query := []rune(ele)
		for i := 1; i <= len(query); i++ {
			low, high := trie.Find(query[0:i])
			prefix := string(query[0:i])

			// check if exclude valid string
			if low > 0 {
				if strings.HasPrefix(sortedQueries[low-1], prefix) {
					t.Errorf("TestTrie; query=%s; exclude=%s", prefix, sortedQueries[low-1])
				}
			}
			if high < len(sortedQueries) {
				if strings.HasPrefix(sortedQueries[high], prefix) {
					t.Errorf("TestTrie; query=%s; exclude=%s", prefix, sortedQueries[high])
				}
			}

			// check if include valid strings
			for j := low; j != high; j++ {
				if !strings.HasPrefix(sortedQueries[j], prefix) {
					t.Errorf("TestTrie; query=%s; exclude=%s", prefix, sortedQueries[high+1])
				}
			}
		}
	}

}

func BenchmarkTrie_Insert(b *testing.B) {
	length := 10
	iter := 100000

	for i := 0; i < b.N; i++ {
		trie := NewTrie(-1, -1)
		for j := 0; j != iter; j++ {
			trie.Insert(generator(length), j)
		}
	}
}

func BenchmarkTrie_Find(b *testing.B) {
	length := 10
	iter := 100000

	history := make([][]rune, iter)

	trie := NewTrie(-1, -1)
	for j := 0; j != iter; j++ {
		query := generator(length)
		trie.Insert(query, j)
		history[j] = query
	}

	for i := 0; i < b.N; i++ {
		// history find
		for _, q := range history {
			trie.Find(q)
		}

		// random find
		for j := 0; j != iter; j++ {
			trie.Find(generator(length))
		}
	}
}

var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// generate fake queries of length [1, n]
func generator(n int) []rune {
	l := rand.Intn(n) + 1
	b := make([]rune, l)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return b
}
