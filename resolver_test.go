package main

import (
	"fmt"
	"testing"
	"time"
)

func TestPaths_Level1_(t *testing.T) {
	lvl := LoadLevel(1)
	t0 := time.Now()

	n := Resolve(lvl)

	d := time.Now().Sub(t0)
	fmt.Printf("Level1 (%s) %+v\n", d, n)
}

func TestPaths_Level2(t *testing.T) {
	lvl := LoadLevel(2)
	t0 := time.Now()

	n := Resolve(lvl)

	d := time.Now().Sub(t0)
	fmt.Printf("Level2 (%s) %+v\n", d, n)
}

func TestPaths_Level3(t *testing.T) {
	lvl := LoadLevel(3)
	t0 := time.Now()

	n := Resolve(lvl)

	d := time.Now().Sub(t0)
	fmt.Printf("Level3 (%s) %+v\n", d, n)
}

func TestPaths_Level4(t *testing.T) {
	lvl := LoadLevel(4)
	t0 := time.Now()

	n := Resolve(lvl)

	d := time.Now().Sub(t0)
	fmt.Printf("Level4 (%s) %+v\n", d, n)
}

func TestPaths_Level5(t *testing.T) {
	lvl := LoadLevel(5)
	t0 := time.Now()

	n := Resolve(lvl)

	d := time.Now().Sub(t0)
	fmt.Printf("Level5 (%s) %+v\n", d, n)
}

func TestPaths_Level6(t *testing.T) {
	lvl := LoadLevel(6)
	t0 := time.Now()

	n := Resolve(lvl)

	d := time.Now().Sub(t0)
	fmt.Printf("Level6 (%s) %+v\n", d, n)
}

func TestPaths_Level7(t *testing.T) {
	lvl := LoadLevel(7)
	t0 := time.Now()

	n := Resolve(lvl)

	d := time.Now().Sub(t0)
	fmt.Printf("Level8 (%s) %+v\n", d, n)
}

func TestPaths_Level9(t *testing.T) {
	lvl := LoadLevel(9)
	t0 := time.Now()

	n := Resolve(lvl)

	d := time.Now().Sub(t0)
	fmt.Printf("Level9 (%s) %+v\n", d, n)
}

func TestPaths_Level10(t *testing.T) {
	lvl := LoadLevel(10)
	t0 := time.Now()

	n := Resolve(lvl)

	d := time.Now().Sub(t0)
	fmt.Printf("Level10 (%s) %+v\n", d, n)
}

func TestPaths_Level11(t *testing.T) {
	lvl := LoadLevel(11)
	t0 := time.Now()

	n := Resolve(lvl)

	d := time.Now().Sub(t0)
	fmt.Printf("Level11 (%s) %+v\n", d, n)
}

func TestPaths_Level12(t *testing.T) {
	lvl := LoadLevel(12)
	t0 := time.Now()

	n := Resolve(lvl)

	d := time.Now().Sub(t0)
	fmt.Printf("Level12 (%s) %+v\n", d, n)
}

func TestPaths_Level13(t *testing.T) {
	lvl := LoadLevel(13)
	t0 := time.Now()

	n := Resolve(lvl)

	d := time.Now().Sub(t0)
	fmt.Printf("Level13 (%s) %+v\n", d, n)
}

func TestPaths_Level14(t *testing.T) {
	lvl := LoadLevel(14)
	t0 := time.Now()

	n := Resolve(lvl)

	d := time.Now().Sub(t0)
	fmt.Printf("Level14 (%s) %+v\n", d, n)
}

func TestPaths_Level15(t *testing.T) {
	lvl := LoadLevel(15)

	t0 := time.Now()
	n := Resolve(lvl)

	d := time.Now().Sub(t0)
	fmt.Printf("Level15 (%s) %+v\n", d, n)
}
