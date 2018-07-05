package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/disiqueira/gotree"
)

type NodeHeap struct {
	Arr []*Node
}

func NewNodeHeap() *NodeHeap {
	// 0 is a nil element not considered, because it makes the
	// array shifting math cleaner
	return &NodeHeap{Arr: []*Node{nil}}
}

func (h *NodeHeap) Add(n *Node) (ix int) {
	// set tag of node to OPEN
	n.T = OPEN
	// append to array end (will bubble up below)
	h.Arr = append(h.Arr, n)
	ix = len(h.Arr) - 1
	// bubble up (needs HeapIX set)
	n.heapIX = ix
	ix = h.bubbleUp(ix)
	// return ix to user
	return ix
}

func (h *NodeHeap) bubbleUp(ix int) int {
	// if ix > 1 and K(ix) < K(ix>>1)
	// (if we're not the top node and our K value is less than the parent)
	for ix > 1 && h.Arr[ix].K < h.Arr[ix>>1].K {
		// swap Nodes in heap
		h.Arr[ix], h.Arr[ix>>1] = h.Arr[ix>>1], h.Arr[ix]
		// swap HeapIX for nodes
		h.Arr[ix].heapIX, h.Arr[ix>>1].heapIX =
			h.Arr[ix>>1].heapIX, h.Arr[ix].heapIX
		// set ix for next iter to be the ix we just swapped to
		ix = ix >> 1
	}
	return ix
}

func (h *NodeHeap) bubbleDown(ix int) int {
	for {
		// beginning state: lesser node is the one we're on
		lesser := ix
		// l index, r index
		lix := (ix << 1)
		rix := (ix << 1) + 1
		// if we've reached the bottom, return
		if !(lix < len(h.Arr) || rix < len(h.Arr)) {
			return ix
		}
		// if left child exists and is lesser, set lesser to lix
		if lix < len(h.Arr) && h.Arr[lix].K < h.Arr[lesser].K {
			lesser = lix
		}
		// if left child exists and is lesser, set lesser to rix
		if rix < len(h.Arr) && h.Arr[rix].K < h.Arr[lesser].K {
			lesser = rix
		}
		// if one of children was lesser, swap and continue bubble down
		// from that node
		if lesser != ix {
			// swap Nodes in heap
			h.Arr[ix], h.Arr[lesser] = h.Arr[lesser], h.Arr[ix]
			// swap HeapIX for nodes
			h.Arr[ix].heapIX, h.Arr[lesser].heapIX =
				h.Arr[lesser].heapIX, h.Arr[ix].heapIX
			// continue to bubble down
			ix = lesser
			continue
		} else {
			// else, no child was lesser. we can stop here
			return ix
		}
	}
}

func (h *NodeHeap) Pop() (*Node, error) {
	if h.Len() == 0 {
		return nil, errors.New("heap empty")
	}
	// get root elem and replace with last element
	n := h.Arr[1]
	last_ix := len(h.Arr) - 1
	// bring last element to root, bubble new root node down if
	// needed, before return
	h.Arr[1] = h.Arr[last_ix]
	h.Arr[1].heapIX = 1
	h.Arr = h.Arr[:last_ix]
	h.bubbleDown(1)
	return n, nil
}

func (h *NodeHeap) Modified(n *Node) {
	// if less than parent, bubble up
	parentIX := n.heapIX >> 1
	if parentIX > 0 {
		parent := h.Arr[parentIX]
		if n.K < parent.K {
			h.bubbleUp(n.heapIX)
			return
		}
	}
	// if greater than either child, bubble down
	lix := (n.heapIX << 1)
	rix := (n.heapIX << 1) + 1
	for _, ix := range []int{lix, rix} {
		if ix < len(n.Arr) {
			child := h.Arr[ix]
			if n.K > child.K {
				h.bubbleDown(n.heapIX)
				return
			}
		}
	}
}

func (h *NodeHeap) Len() int {
	return len(h.Arr) - 1
}

func (h *NodeHeap) Clear() {
	h.Arr = h.Arr[:1]
}

func (h *NodeHeap) String() string {
	// for building string
	var buffer bytes.Buffer
	// print array
	// if elements, print tree
	if len(h.Arr) > 1 {
		buffer.WriteString("\n")
		// build tree using gotree package by descending recursively
		var addChildren func(node gotree.Tree, ix int)
		addChildren = func(node gotree.Tree, ix int) {
			rix := (ix << 1) + 1
			if rix < len(h.Arr) {
				r := node.Add(fmt.Sprintf("[%d]%s",
					rix, h.Arr[rix].String()))
				addChildren(r, rix)
			}
			lix := (ix << 1)
			if lix < len(h.Arr) {
				l := node.Add(fmt.Sprintf("[%d]%s",
					lix, h.Arr[lix].String()))
				addChildren(l, lix)
			}
		}
		tree := gotree.New(fmt.Sprintf("[%d]%s",
			1, h.Arr[1].String()))
		addChildren(tree, 1)
		buffer.WriteString(tree.Print())
	}
	return buffer.String()
}
