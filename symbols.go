package main

import "huffman_coding/heap"

const (
	newChar     rune                = 1<<31 - 1 - iota // value representing a new character
	eof                                                // value representing end of data
	customChars = iota                                 // number of custom characters
	maxChars    = 256 + customChars                    // number of possible bytes + custom characters
)

type symbols struct {
	root   *Node
	leaves NodeHeap
	chars  map[rune]*Node
}

func newSymbols() *symbols {
	s := new(symbols)

	leaves := make(NodeHeap, customChars, maxChars)
	leaves[0] = &Node{Freq: 1, Char: newChar, index: 1}
	leaves[1] = &Node{Freq: 1, Char: eof, index: 1}
	heap.Init(&leaves)
	s.leaves = leaves

	chars := make(map[rune]*Node, cap(leaves))
	for _, node := range leaves {
		chars[node.Char] = node
	}
	s.chars = chars

	s.buildTree()

	return s
}

func (s *symbols) insert(char rune) {
	node := &Node{Freq: 1, Char: char, index: len(s.leaves)}
	heap.Push(&s.leaves, node)
	s.chars[char] = node
	s.buildTree()
}

func (s *symbols) update(node *Node) {
	node.Freq++
	heap.Fix(&s.leaves, node.index)
	s.buildTree()
}

func (s *symbols) buildTree() {
	h := make(NodeHeap, len(s.leaves))
	copy(h, s.leaves)
	for h.Len() > 1 {
		left := heap.Pop(&h).(*Node)
		right := heap.Pop(&h).(*Node)
		parent := &Node{Freq: left.Freq + right.Freq}
		left.Parent = parent
		right.Parent = parent
		parent.Left = left
		parent.Right = right
		heap.Push(&h, parent)
	}
	s.root = heap.Pop(&h).(*Node)
}
