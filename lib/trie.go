package lib

// trie to index queries
type Trie struct {
	low      int // the lowest index, inclusive
	high     int // the highest index, exclusive
	children map[rune]*Trie
}

func NewTrie(l int, h int) *Trie {
	return &Trie{
		low:      l,
		high:     h,
		children: make(map[rune]*Trie),
	}
}

func (t *Trie) Insert(r []rune, id int) {
	node := t
	for _, ele := range r {
		if child, ok := node.children[ele]; ok { // node exist
			child.high = id + 1
		} else { // node does not exist
			node.children[ele] = NewTrie(id, id+1)
		}
		node = node.children[ele]
	}
}

// find the index range of a given string.
// return -1, -1 if not found
func (t *Trie) Find(r []rune) (int, int) {
	node := t
	for _, ele := range r {
		if child, ok := node.children[ele]; ok { // node exist
			node = child
		} else { // node does not exist
			return -1, -1
		}
	}

	return node.low, node.high
}
