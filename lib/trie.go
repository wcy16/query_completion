package lib

// trie to index queries
type Trie struct {
	Low      int // the lowest index, inclusive
	High     int // the highest index, exclusive
	Children map[rune]*Trie
}

func NewTrie(l int, h int) *Trie {
	return &Trie{
		Low:      l,
		High:     h,
		Children: make(map[rune]*Trie),
	}
}

func (t *Trie) Insert(r []rune, id int) {
	node := t
	for _, ele := range r {
		if child, ok := node.Children[ele]; ok { // node exist
			child.High = id + 1
		} else { // node does not exist
			node.Children[ele] = NewTrie(id, id+1)
		}
		node = node.Children[ele]
	}
}

// find the index range of a given string.
// return -1, -1 if not found
func (t *Trie) Find(r []rune) (int, int) {
	node := t
	for _, ele := range r {
		if child, ok := node.Children[ele]; ok { // node exist
			node = child
		} else { // node does not exist
			return -1, -1
		}
	}

	return node.Low, node.High
}
