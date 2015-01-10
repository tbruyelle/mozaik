// +build ignore

package main

import (
	"fmt"
	"math"
	"sync"
)

const (
	MaxDepth = 30
)

type Node struct {
	// current switch
	s      int
	depth  int
	parent *Node
	childs []*Node
	// lvl represents a copy of the level
	// at the current node of the tree
	lvl Level
}

var signs map[string]bool
var mutex = new(sync.Mutex)

func addSign(sign string) {
	mutex.Lock()
	defer mutex.Unlock()
	signs[sign] = true
}

func hasSign(sign string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	return signs[sign]
}

func (n *Node) String() string {
	//return fmt.Sprintf("s%d, d=%d, childs=%+v", n.s, n.depth, n.childs)
	//return fmt.Sprintf("s%d, d=%d, parent=[%+v] win=%t", n.s, n.depth, n.parent, n.lvl.Win())
	depth := n.depth
	return fmt.Sprintf("d=%d, sws=%s", depth, n.road())
}

// Returns the switch combination used so far
func (n *Node) road() string {
	var s string
	for n.parent != nil {
		s = n.lvl.switches[n.s].name + s
		n = n.parent
	}
	s = n.lvl.switches[n.s].name + s
	return s
}

func FindShortestPaths(lvl Level) []*Node {
	ns := FindPaths(lvl)
	fmt.Printf("Found %d path\n", len(ns))
	shortest := MaxDepth
	for _, n := range ns {
		if n.depth < shortest {
			shortest = n.depth
		}
	}
	// keep only the shortest
	var shortestPaths []*Node
	for _, n := range ns {
		if n.depth == shortest {
			shortestPaths = append(shortestPaths, n)
		}
	}
	fmt.Printf("Shortest is %d with %d paths", shortest, len(shortestPaths))
	return shortestPaths
}

func FindPaths(lvl Level) []*Node {
	var paths []*Node
	for i := range lvl.switches {
		fmt.Printf("find path starting switch %d %t\n", i, lvl.IsPlain(i))
		n := &Node{s: i, depth: 1, lvl: lvl.Copy()}
		if n.lvl.IsPlain(n.s) {
			continue
		}
		check(n, &paths)
		//fmt.Printf("switch %d path=%+v\n", i, n)
	}
	return paths
}

func check(n *Node, paths *[]*Node) {
	//fmt.Printf("check %+v\n", n)
	if n.depth > MaxDepth {
		return
	}
	n.lvl.RotateSwitch(n.lvl.switches[n.s])
	if n.lvl.Win() {
		//fmt.Printf("WIN %+v\n", n)
		*paths = append(*paths, n)
		return
	}

	// Add childs
	//fmt.Printf("n%d %t %t %s\n", n.s, n.lvl.IsPlain(n.s), n.hasRotatedTwice(), n.lvl.blockSignature())
	for i := range n.lvl.switches {
		if n.s == i {
			// For current child, add only if not plain
			// and has not rotated twice successively
			if !n.lvl.IsPlain(n.s) && !n.hasRotatedTwice() {
				n.addChild(n.s)
			}
		} else {
			n.addChild(i)
		}
	}
	//fmt.Printf("Check %d childs depth=%d\n", len(n.childs), n.depth)
	for _, c := range n.childs {
		check(c, paths)
	}
}

var origHowFar int

func FindPathAsync(lvl Level) *Node {
	origHowFar = lvl.HowFar()
	fmt.Printf("orig howfar=%d\n", origHowFar)

	signs = make(map[string]bool)
	quit := make(chan bool)
	nodes := make(chan *Node)
	for i := range lvl.switches {
		if lvl.IsPlain(i) {
			continue
		}
		fmt.Printf("find path switch %d\n", i)
		n := &Node{s: i, depth: 1, lvl: lvl.Copy()}
		go checkAsync(n, nodes, quit)
		//fmt.Printf("switch %d path=%+v\n", i, n)
	}
	win := <-nodes
	close(quit)
	return win
}

func checkAsync(n *Node, nodes chan *Node, quit chan bool) {
	select {
	case <-quit:
		return
	default:
		//howfar := n.lvl.HowFar()
		//fmt.Printf("check \n%+v%d\n", n.lvl.blockSignature(), howfar)
		if n.depth > MaxDepth {
			fmt.Println("max depth reached")
			return
		}
		n.lvl.RotateSwitch(n.lvl.switches[n.s])
		//if n.lvl.HowFar() >= howfar {
		//	//fmt.Println("too far")
		//	return
		//}
		if n.lvl.Win() {
			//fmt.Printf("WIN %+v\n", n)
			nodes <- n
			return
		}

		// Check if the new signature has already been reached
		sign := n.lvl.blockSignature()
		if hasSign(sign) {
			// Signature already reached, skip
			//fmt.Println("already reached")
			return
		}
		addSign(sign)

		// Add childs
		//fmt.Printf("n%d %t %t %s\n", n.s, n.lvl.IsPlain(n.s), n.hasRotatedTwice(), n.lvl.blockSignature())
		for i := range n.lvl.switches {
			if n.s == i {
				// For current child, add only if not plain
				// and has not rotated twice successively
				if !n.lvl.IsPlain(n.s) && !n.hasRotatedTwice() {
					n.addChild(n.s)
				}
			} else {
				n.addChild(i)
			}
		}
		//fmt.Printf("Check %d childs depth=%d\n", len(n.childs), n.depth)
		for _, c := range n.childs {
			go checkAsync(c, nodes, quit)
		}
	}
}

func (n *Node) addChild(s int) *Node {
	c := &Node{s: s, depth: n.depth + 1, parent: n, lvl: n.lvl.Copy()}
	n.childs = append(n.childs, c)
	return c
}

func (n *Node) hasRotatedTwice() bool {
	if n.parent != nil && n.parent.s == n.s && n.parent.parent != nil && n.parent.parent.s == n.s {
		return true
	}
	return false
}

// DetermineNearSwitches returns a map of nearest switches
// for each switches in the lvl in parameter
func DetermineNearSwitches(lvl Level) map[int][]int {
	res := make(map[int][]int, len(lvl.switches))
	for i := range lvl.switches {
		for j := range lvl.switches {
			if i == j {
				continue
			}
			lineDiff := math.Abs(float64(lvl.switches[i].line - lvl.switches[j].line))
			colDiff := math.Abs(float64(lvl.switches[i].col - lvl.switches[j].col))
			if lineDiff < 2 && colDiff < 2 {
				res[i] = append(res[i], j)
			}
		}

	}
	return res
}
