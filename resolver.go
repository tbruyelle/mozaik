package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"runtime/pprof"
)

const (
	MaxDepth = 50
)

var (
	signs map[string]bool
	lvl   Level
)

type Board [4][4]Color

// IsPlain returns true if all the blocks of the switch
// have the same color
func (b Board) isPlain(li, col int) bool {
	return b[li][col] == b[li+1][col] && b[li+1][col] == b[li][col+1] && b[li][col+1] == b[li+1][col+1]
}

func (b *Board) cp(board Board) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			b[i][j] = board[i][j]
		}
	}
}

func (b Board) win() bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if lvl.winSignature[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func (b *Board) rotate(li, co int) {
	color := b[li][co]
	b[li][co] = b[li+1][co]
	b[li+1][co] = b[li+1][co+1]
	b[li+1][co+1] = b[li][co+1]
	b[li][co+1] = color
}

func (b Board) signature() string {
	var signature bytes.Buffer
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			signature.WriteRune(rune(b[i][j]))
		}
		signature.WriteString("\n")
	}
	return signature.String()
}

func (b Board) findManhattan(x, y int) int {
	c := b[x][y]
	max := 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if c == lvl.winSignature[i][j] {
				m := manhattan(x, y, i, j)
				if m > max {
					max = m
				}
			}
		}
	}
	return max
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func manhattan(x1, y1, x2, y2 int) int {
	return abs(x1-x2) + abs(y1-y2)
}

func (b Board) howFar() int {
	howfar := 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if lvl.winSignature[i][j] != b[i][j] {
				howfar += b.findManhattan(i, j)
			}
		}
	}
	return howfar
}

type Nodes []*Node

func (ns Nodes) Len() int {
	return len(ns)
}

func (ns Nodes) Less(i, j int) bool {
	return ns[i].priority < ns[j].priority
}

func (ns Nodes) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

func (ns *Nodes) Push(x interface{}) {
	node := x.(*Node)
	*ns = append(*ns, node)
}

func (ns *Nodes) Pop() interface{} {
	old := *ns
	n := len(old)
	node := old[n-1]
	*ns = old[0 : n-1]
	return node
}

type Node struct {
	board Board
	depth int
	// current switch
	s        int
	parent   *Node
	priority int
}

func (n *Node) String() string {
	//return fmt.Sprintf("s%d, d=%d, childs=%+v", n.s, n.depth, n.childs)
	//return fmt.Sprintf("s%d, d=%d, parent=[%+v] win=%t", n.s, n.depth, n.parent, n.lvl.Win())
	depth := n.depth
	return fmt.Sprintf("d=%d, p=%d, sws=%s", depth, n.priority, n.road())
}

// Returns the switch combination used so far
func (n *Node) road() string {
	var s string
	for n.parent != nil && n.s >= 0 {
		s = lvl.switches[n.s].name + s
		n = n.parent
	}
	if n.s >= 0 {
		s = lvl.switches[n.s].name + s
	}
	return s
}

func Resolve(l Level) *Node {
	f, err := os.Create("resolver.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	lvl = l
	ns := make(Nodes, 0)
	heap.Init(&ns)

	init := &Node{s: -1}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			init.board[i][j] = lvl.blocks[i][j].Color
		}
	}
	init.priority = init.board.howFar()
	heap.Push(&ns, init)
	//fmt.Println("INIT NODE", init)
	signs = make(map[string]bool)
	signs[init.board.signature()] = true

	loop := 0
	for {
		n := process(&ns)
		if n != nil {
			return n
		}
		loop++
	}
	return nil
}

func process(ns *Nodes) *Node {
	n := heap.Pop(ns).(*Node)
	if n.depth > MaxDepth {
		return nil
	}
	if n.board.win() {
		return n
	}
	for i, sw := range lvl.switches {
		if n.board.isPlain(sw.line, sw.col) {
			// Useless to rotate a plain switch
			continue
		}
		if n.s == i && n.parent != nil && n.parent.s == i && n.parent.parent != nil && n.parent.parent.s == i {
			// Useless to rotate 4 times in a row the same switch
			continue
		}

		nn := &Node{
			s:      i,
			depth:  n.depth + 1,
			parent: n,
		}
		nn.board.cp(n.board)
		nn.board.rotate(sw.line, sw.col)
		sign := nn.board.signature()
		if _, ok := signs[sign]; ok {
			// Already processed skip
			continue
		}
		signs[sign] = true
		nn.priority = nn.board.howFar() + nn.depth

		heap.Push(ns, nn)
	}
	return nil
}
