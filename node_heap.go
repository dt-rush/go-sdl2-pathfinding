package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/disiqueira/gotree"
)

type NodeHeap struct {
	PC  *DstarPathComputer
	Arr []Position
}

func NewNodeHeap(pc *DstarPathComputer) *NodeHeap {
	// 0 is a nil element not considered, because it makes the
	// array shifting math cleaner
	h := NodeHeap{PC: pc, Arr: []Position{NOWHERE}}
	return &h
}

func (h *NodeHeap) bubbleUp(ix int) int {
	// if ix > 1 and K(ix) < K(ix>>1)
	// (if we're not the top node and our K value is less than the parent)
	for ix > 1 &&
		h.PC.K[h.Arr[ix].X][h.Arr[ix].Y] <
			h.PC.K[h.Arr[ix>>1].X][h.Arr[ix>>1].Y] {
		// swap Nodes in heap
		h.Arr[ix], h.Arr[ix>>1] = h.Arr[ix>>1], h.Arr[ix]
		// swap HeapIX for nodes
		h.PC.HeapIX[h.Arr[ix].X][h.Arr[ix].Y],
			h.PC.HeapIX[h.Arr[ix>>1].X][h.Arr[ix>>1].Y] =
			h.PC.HeapIX[h.Arr[ix>>1].X][h.Arr[ix>>1].Y],
			h.PC.HeapIX[h.Arr[ix].X][h.Arr[ix].Y]
		// set ix for next iter to be the ix we just swapped to
		ix = ix >> 1
	}
	return ix
}

func (h *NodeHeap) Add(p Position) (ix int) {

	// append to array end (will bubble up below)
	h.Arr = append(h.Arr, p)
	ix = len(h.Arr) - 1
	// calculate the new K value (min of P and H)
	var K int
	P := h.PC.P[h.Arr[ix].X][h.Arr[ix].Y]
	H := h.PC.H[h.Arr[ix].X][h.Arr[ix].Y]
	if P < H {
		K = P
	} else {
		K = H
	}
	// assign the new K value
	h.PC.K[h.Arr[ix].X][h.Arr[ix].Y] = K
	// bubble up (needs HeapIX set)
	h.PC.HeapIX[p.X][p.Y] = ix
	ix = h.bubbleUp(ix)
	// return ix to user
	return ix
}

func (h *NodeHeap) bubbleDown(ix int) int {
	for {
		// beginning state: greater node is the one we're on
		greater := ix
		// l index, r index
		lix := (ix << 1)
		rix := (ix << 1) + 1
		// if we've reached the bottom, return
		if !(lix < len(h.Arr) || rix < len(h.Arr)) {
			return ix
		}
		// if left child exists and is greater, set greater to lix
		if lix < len(h.Arr) &&
			h.PC.K[h.Arr[lix].X][h.Arr[lix].Y] <
				h.PC.K[h.Arr[greater].X][h.Arr[greater].Y] {
			greater = lix
		}
		// if left child exists and is greater, set greater to rix
		if rix < len(h.Arr) &&
			h.PC.K[h.Arr[rix].X][h.Arr[rix].Y] <
				h.PC.K[h.Arr[greater].X][h.Arr[greater].Y] {
			greater = rix
		}
		// if one of children was greater, swap and continue bubble down
		// from that node
		if greater != ix {
			// swap Nodes in heap
			h.Arr[ix], h.Arr[greater] = h.Arr[greater], h.Arr[ix]
			// swap HeapIX for nodes
			h.PC.HeapIX[h.Arr[ix].X][h.Arr[ix].Y],
				h.PC.HeapIX[h.Arr[greater].X][h.Arr[greater].Y] =
				h.PC.HeapIX[h.Arr[greater].X][h.Arr[greater].Y],
				h.PC.HeapIX[h.Arr[ix].X][h.Arr[ix].Y]
			// continue to bubble down
			ix = greater
			continue
		} else {
			// else, no child was greater. we can stop here
			return ix
		}
	}
}

func (h *NodeHeap) Pop() (Position, error) {
	if h.Len() == 0 {
		return Position{}, errors.New("heap empty")
	}
	// get root elem and replace with last element (shrink slice)
	p := h.Arr[1]
	last_ix := len(h.Arr) - 1
	h.Arr[1] = h.Arr[last_ix]
	h.PC.HeapIX[h.Arr[1].X][h.Arr[1].Y] = 1
	h.Arr = h.Arr[:last_ix]
	// bubble element down to its place
	h.bubbleDown(1)
	return p, nil
}

func (h *NodeHeap) Modify(ix int, H int) {
	// get the old K value
	oldVal := h.PC.K[h.Arr[ix].X][h.Arr[ix].Y]
	// set the new G
	h.PC.H[h.Arr[ix].X][h.Arr[ix].Y] = H
	// calculate the new K value (min of P and H)
	var K int
	P := h.PC.P[h.Arr[ix].X][h.Arr[ix].Y]
	if P < H {
		K = P
	} else {
		K = H
	}
	// assign the new K value
	h.PC.K[h.Arr[ix].X][h.Arr[ix].Y] = K
	// bubble up if needed (setting HeapIX)
	if K < oldVal {
		h.bubbleUp(ix)
	}
	// bubble down if needed
	if K > oldVal {
		h.bubbleDown(ix)
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
