package main

import (
	"fmt"
	"testing"
)

func TestPaths_Level1(t *testing.T) {
	lvl := LoadLevel(1)

	//paths := FindShortestPaths(lvl)
	n := Resolve(lvl)

	if n != nil {
		fmt.Printf("test result %s\n", n.road())
	}
}

func TestPaths_Level2(t *testing.T) {
	lvl := LoadLevel(2)

	//paths := FindShortestPaths(lvl)
	n := Resolve(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestPaths_Level3(t *testing.T) {
	lvl := LoadLevel(3)

	//paths := FindShortestPaths(lvl)
	n := Resolve(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestPaths_Level4(t *testing.T) {
	lvl := LoadLevel(4)

	//paths := FindShortestPaths(lvl)
	n := Resolve(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestPaths_Level5(t *testing.T) {
	lvl := LoadLevel(5)

	//paths := FindShortestPaths(lvl)
	n := Resolve(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestPaths_Level6(t *testing.T) {
	lvl := LoadLevel(6)

	//paths := FindShortestPaths(lvl)
	n := Resolve(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestPaths_Level7(t *testing.T) {
	lvl := LoadLevel(7)

	//paths := FindShortestPaths(lvl)
	n := Resolve(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestPaths_Level8(t *testing.T) {
	lvl := LoadLevel(8)

	//paths := FindShortestPaths(lvl)
	n := Resolve(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestPaths_Level9(t *testing.T) {
	lvl := LoadLevel(9)

	n := Resolve(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestPaths_Level10(t *testing.T) {
	lvl := LoadLevel(10)

	n := Resolve(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestPaths_Level11(t *testing.T) {
	lvl := LoadLevel(11)

	n := Resolve(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestPaths_Level12(t *testing.T) {
	lvl := LoadLevel(12)

	n := Resolve(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestPaths_Level13(t *testing.T) {
	lvl := LoadLevel(13)

	n := Resolve(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestPaths_Level14(t *testing.T) {
	lvl := LoadLevel(14)

	n := Resolve(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestPaths_Level15(t *testing.T) {
	lvl := LoadLevel(15)

	n := Resolve(lvl)

	fmt.Printf("test result %+v\n", n)
}
