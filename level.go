package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/sprite/clock"
)

type Level struct {
	sync.Mutex
	blocks       [][]*Block
	switches     []*Switch
	winSignature [][]ColorDef
	// rotated represents the historics of rotations
	rotated []int
	// rotating represents a rotate which
	// is currently rotating
	rotating *Switch
	solution string
}

type ColorDef int

const (
	Red         ColorDef = iota //0
	Yellow                      //1
	Blue                        //2
	Green                       //3
	Pink                        //4
	Orange                      //5
	LightBlue                   //6
	Purple                      //7
	Brown                       //8
	LightGreen                  //9
	Cyan                        //A
	LightPink                   //B
	White                       //C
	LightPurple                 //D
	LightBrown                  //E
	OtherWhite                  //F
	Empty       ColorDef = -1
)

type Block struct {
	Object
	Color ColorDef
}

type Switch struct {
	Object
	line, col int
	name      string
}

func (sw *Switch) Rotate() {
	blocks := sw.Blocks()
	for i := range blocks {
		b := blocks[i]
		v := switchSize / 2
		// Prepare a rotation around the center of the switch
		b.Rx, b.Ry = sw.X+v, sw.Y+v
		b.Sx, b.Sy = b.Rx, b.Ry
		b.Time = 0
		b.Action = ActionFunc(blockRotate)
	}
}

// Blocks returns the block arround the switch in parameter.
func (sw *Switch) Blocks() []*Block {
	topLeft := g.level.blocks[sw.line][sw.col]
	topRight := g.level.blocks[sw.line][sw.col+1]
	bottomLeft := g.level.blocks[sw.line+1][sw.col]
	bottomRight := g.level.blocks[sw.line+1][sw.col+1]
	return []*Block{topLeft, topRight, bottomLeft, bottomRight}
}

func (s *Switch) String() string {
	return fmt.Sprintf("sw{line:%d, col:%d}", s.line, s.col)
}

func (l *Level) Copy() Level {
	lvl := new(Level)
	lvl.blocks = make([][]*Block, len(l.blocks))
	for i := range l.blocks {
		lvl.blocks[i] = make([]*Block, len(l.blocks[i]))
		copy(lvl.blocks[i], l.blocks[i])
	}
	lvl.switches = make([]*Switch, len(l.switches))
	copy(lvl.switches, l.switches)
	lvl.winSignature = l.winSignature
	return *lvl
}

// IsPlain returns true if all the blocks of the switch
// have the same color
func (l *Level) IsPlain(sw int) bool {
	x, y := l.switches[sw].line, l.switches[sw].col
	b1 := l.blocks[x][y]
	b2 := l.blocks[x+1][y]
	b3 := l.blocks[x][y+1]
	b4 := l.blocks[x+1][y+1]

	return b1.Color == b2.Color && b2.Color == b3.Color && b3.Color == b4.Color
}

// Win returns true if player has win
func (l *Level) Win() bool {
	for i := range l.winSignature {
		for j := range l.winSignature[i] {
			if l.blocks[i][j] != nil && l.winSignature[i][j] != l.blocks[i][j].Color {
				return false
			}
		}
	}
	return true
}

func (l *Level) findManhattan(x, y int) int {
	c := l.blocks[x][y].Color
	maxM := 0
	for i := range l.winSignature {
		for j := range l.winSignature[i] {
			if c == l.winSignature[i][j] {
				m := manhattan(x, y, i, j)
				if m > maxM {
					maxM = m
				}
			}
		}
	}
	return maxM
}

func manhattan(x1, y1, x2, y2 int) int {
	return int(math.Abs(float64(x1)-float64(x2))) + int(math.Abs(float64(y1)-float64(y2)))
}

func (l *Level) HowFar() int {
	howfar := 0
	for i := range l.blocks {
		for j := range l.blocks[i] {
			// FIXME populate always blocks in parse
			if l.blocks[i][j] != nil && l.winSignature[i][j] != l.blocks[i][j].Color {
				howfar += l.findManhattan(i, j)
			}
		}
	}
	return howfar
}

// UndoLastMove cancels the last player move
func (l *Level) UndoLastMove() {
	if l.rotating != nil {
		return
	}
	sw := l.PopLastRotated()
	if sw != nil {
		g.level.rotating = sw
		blocks := sw.Blocks()
		for i := range blocks {
			b := blocks[i]
			v := switchSize / 2
			b.Rx, b.Ry = sw.X+v, sw.Y+v
			b.Sx, b.Sy = b.Rx, b.Ry
			b.Action = ActionFunc(blockRotateInverse)
		}
	}
}

func (l *Level) PopLastRotated() *Switch {
	if len(l.rotated) == 0 {
		return nil
	}
	i := len(l.rotated) - 1
	res := l.rotated[i]
	l.rotated = l.rotated[:i]
	return l.switches[res]
}

func (l *Level) addBlock(color ColorDef, line, col int) {
	colf, linef := float32(col), float32(line)
	b := &Block{Color: color}
	b.Object = Object{
		X:      xMin + colf*(blockSize+blockPadding),
		Y:      yMin + linef*(blockSize+blockPadding),
		Width:  blockSize,
		Height: blockSize,
		Data:   b,
		Action: wait{until: clock.Time(line*10 + col*5), next: ActionFunc(blockPopIn)},
	}
	l.blocks[line][col] = b
}

// addSwitch appends a new switch at the bottom right
// of the coordinates in parameters.
func (l *Level) addSwitch(line, col int) {
	v := switchSize / 2
	colf, linef := float32(col), float32(line)
	s := &Switch{
		line: line, col: col,
		name: determineName(line, col),
	}
	s.Object = Object{
		X:      xMin + (colf+1)*blockSize + colf*blockPadding*2 - v,
		Y:      yMin + (linef+1)*blockSize + linef*blockPadding*2 - v,
		Width:  switchSize,
		Height: switchSize,
		Action: wait{until: 70, next: ActionFunc(switchPopIn)},
		Data:   s,
	}
	l.switches = append(l.switches, s)
	//fmt.Println("Switch added", s.X, s.Y)
}

func determineName(line, col int) string {
	switch line {
	case 0:
		switch col {
		case 0:
			return "7"
		case 1:
			return "8"
		case 2:
			return "9"
		}
	case 1:
		switch col {
		case 0:
			return "4"
		case 1:
			return "5"
		case 2:
			return "6"
		}
	case 2:
		switch col {
		case 0:
			return "1"
		case 1:
			return "2"
		case 2:
			return "3"
		}
	}
	return "x"
}

// PressSwitch tries to find a swicth from the coordinates
// and activate it.
func (l *Level) PressSwitch(x, y float32) {
	// Handle click only when no switch are rotating
	if l.rotating == nil {
		if i, s := l.findSwitch(x, y); s != nil {
			l.rotating = s
			l.TriggerSwitch(i)
		}
	}
}

func (l *Level) TriggerSwitchName(name string) {
	for i := 0; i < len(l.switches); i++ {
		if l.switches[i].name == name {
			l.TriggerSwitch(i)
			return
		}
	}
}

func (l *Level) TriggerSwitch(i int) {
	l.switches[i].Rotate()
	l.rotated = append(l.rotated, i)
}

func (l *Level) findSwitch(x, y float32) (int, *Switch) {
	for i, s := range l.switches {
		if x >= s.X && x <= s.X+switchSize && y >= s.Y && y <= s.Y+switchSize {
			return i, s
		}
	}
	return -1, nil
}

func (l *Level) blockSignature() string {
	var signature bytes.Buffer
	for i := 0; i < len(l.blocks); i++ {
		for j := 0; j < len(l.blocks[i]); j++ {
			if l.blocks[i][j] == nil {
				signature.WriteString("-")
			} else {
				signature.WriteString(ctoa(l.blocks[i][j].Color))
			}
		}
		signature.WriteString("\n")
	}
	return signature.String()
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func atoc(s string) ColorDef {
	if s == "-" {
		return Empty
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		switch s {
		case "A":
			return Cyan
		case "B":
			return LightPink
		case "C":
			return White
		case "D":
			return LightPurple
		case "E":
			return LightBrown
		case "F":
			return OtherWhite
		}
	}
	return ColorDef(i)
}

func ctoa(c ColorDef) string {
	return fmt.Sprintf("%d", c)
}

// LoadLevel loads the level number in parameter
func LoadLevel(level int) Level {
	f, err := app.Open(fmt.Sprintf("levels/%d", level))
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return ParseLevel(string(b))
}

// ParseLevel reads level information
func ParseLevel(str string) Level {
	lines := strings.Split(str, "\n")
	step := 0
	l := Level{}

	for i := 0; i < len(lines); i++ {
		if len(lines[i]) == 0 {
			step++
			continue
		}
		switch step {
		case 0:
			// read block colors
			bline := make([]*Block, len(lines[i]))
			l.blocks = append(l.blocks, bline)
			for j, c := range lines[i] {
				if c != '-' {
					l.addBlock(atoc(string(c)), i, j)
				} else {
					l.addBlock(Empty, i, j)
				}
			}
		case 1:
			// read switch locations
			tokens := strings.Split(lines[i], ",")
			l.addSwitch(atoi(tokens[0]), atoi(tokens[1]))
		case 2:
			//read win
			wline := make([]ColorDef, len(lines[i]))
			for j, c := range lines[i] {
				wline[j] = atoc(string(c))
			}
			l.winSignature = append(l.winSignature, wline)
		case 3:
			// read the solution
			l.solution = lines[i]
		}
	}
	//fmt.Printf("Level loaded blocks=%d, swicthes=%d\n", len(l.blocks), len(l.switches))

	//for i := 0; i < len(l.blocks); i++ {
	//	fmt.Printf("line %d blocks %d\n", i, len(l.blocks[i]))
	//}
	//fmt.Printf("winSignature\n%s\n---\n", l.winSignature)
	return l
}

// RotateSwitch swaps bocks according to the 90d rotation
func (lvl *Level) RotateSwitch(s *Switch) {
	l, c := s.line, s.col
	//fmt.Println("Swap from", l, c)
	color := lvl.blocks[l][c].Color
	lvl.blocks[l][c].Color = lvl.blocks[l+1][c].Color
	lvl.blocks[l+1][c].Color = lvl.blocks[l+1][c+1].Color
	lvl.blocks[l+1][c+1].Color = lvl.blocks[l][c+1].Color
	lvl.blocks[l][c+1].Color = color
}

// RotateSwitchInverse swaps bocks according to the -90d rotation
func (lvl *Level) RotateSwitchInverse(s *Switch) {
	l, c := s.line, s.col
	color := lvl.blocks[l][c].Color
	lvl.blocks[l][c].Color = lvl.blocks[l][c+1].Color
	lvl.blocks[l][c+1].Color = lvl.blocks[l+1][c+1].Color
	lvl.blocks[l+1][c+1].Color = lvl.blocks[l+1][c].Color
	lvl.blocks[l+1][c].Color = color
}
