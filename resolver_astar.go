package main

import (
	"container/heap"
	"fmt"
)

const (
	MaxDepth = 50
)

var (
	signs map[string]bool
)

type Nodes []*Node

func (ns Nodes) Len() int {
	return len(ns)
}

func (ns Nodes) Less(i, j int) bool {
	return ns[i].priority < ns[j].priority
}

func (ns Nodes) Swap(i, j int) {
	fmt.Println("Swap", i, j)
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
	// current switch
	s      int
	depth  int
	parent *Node
	// lvl represents a copy of the level
	// at the current node of the tree
	lvl      Level
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
		s = n.lvl.switches[n.s].name + s
		n = n.parent
	}
	if n.s >= 0 {
		s = n.lvl.switches[n.s].name + s
	}
	return s
}

func Resolve(lvl Level) *Node {
	//f, err := os.Create("resolver.prof")
	//if err != nil {
	//	panic(err)
	//}
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()

	ns := make(Nodes, 0)
	heap.Init(&ns)

	init := &Node{
		s:        -1,
		depth:    0,
		lvl:      lvl.Copy(),
		priority: lvl.HowFar(),
	}
	heap.Push(&ns, init)
	fmt.Println("INIT NODE", init)
	signs = make(map[string]bool)
	signs[init.lvl.blockSignature()] = true

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
	fmt.Println("Processing node", n)
	if n.lvl.Win() {
		return n
	}
	if n.depth > MaxDepth {
		return nil
	}
	for i, _ := range n.lvl.switches {
		nn := &Node{
			s:        i,
			depth:    n.depth + 1,
			lvl:      n.lvl.Copy(),
			parent:   n,
			priority: n.lvl.HowFar() + n.depth,
		}
		nn.lvl.RotateSwitch(n.lvl.switches[i])
		sign := nn.lvl.blockSignature()
		if _, ok := signs[sign]; ok {
			// Already processed skip
			continue
		}
		signs[sign] = true
		heap.Push(ns, nn)
		//fmt.Println("Added node", nn)
	}
	return nil
}
