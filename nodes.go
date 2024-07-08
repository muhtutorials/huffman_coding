package main

import (
	"fmt"
	"strconv"
)

type Node struct {
	Parent, Left, Right *Node
	Freq                int
	Char                rune
	// index of the item in the heap.
	// It's needed by update method and is maintained by the heap.Interface methods.
	index int
}

// Code returns the Huffman code of the node and number of bits set.
// Left children get bit 0, Right children get bit 1.
// Implementation uses Node.Parent to "walk up" the tree.
func (n *Node) Code() (r uint64, count uint8) {
	for parent := n.Parent; parent != nil; n, parent = parent, parent.Parent {
		if parent.Right == n {
			// count = 3
			// ... 0111 | 1 << 3
			// ... 0111 | 1000
			// ... 1111
			r |= 1 << count
		} // else bit 0 => nothing to do with r
		count++
	}
	return r, count
}

// Print traverses the Huffman tree and prints the values with their code in binary representation.
// Function is used for debugging purposes.
func Print(root *Node) {
	// traverse traverses a subtree from the given node,
	// using the prefix code leading to this node, having the number of set bits specified.
	var traverse func(n *Node, code uint64, count uint8)

	traverse = func(n *Node, code uint64, count uint8) {
		if n.Left == nil {
			// it's a leaf
			fmt.Printf("'%c' (%d): %0"+strconv.Itoa(int(count))+"b, freq: %d\n", n.Char, n.Char, code, n.Freq)
			return
		}
		count++
		traverse(n.Left, code<<1, count)
		traverse(n.Right, code<<1+1, count)
	}

	traverse(root, 0, 0)
}

type NodeHeap []*Node

func (h *NodeHeap) Len() int {
	return len(*h)
}

func (h *NodeHeap) Less(i, j int) bool {
	return (*h)[i].Freq < (*h)[j].Freq
}

func (h *NodeHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
	(*h)[i].index = i
	(*h)[j].index = j
}

func (h *NodeHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents
	value, ok := x.(*Node)
	if !ok {
		fmt.Println("value must be of *Node type")
		return
	}
	length := len(*h)
	value.index = length
	*h = append(*h, value)
}

func (h *NodeHeap) Pop() any {
	last := len(*h) - 1
	x := (*h)[last]
	(*h)[last] = nil // avoid memory leak
	*h = (*h)[:last]
	return x
}
