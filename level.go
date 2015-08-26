package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/sprite/clock"
)

type Level struct {
	sync.Mutex
	blocks       [][]*Block
	switches     []*Switch
	winSignature [][]Color
	// rotated represents the historics of rotations
	rotated []int
	// rotating represents a rotate which
	// is currently rotating
	rotating *Switch
	solution string
	maxMoves int
	moves    int
}

type Color rune

const (
	Empty       = '-'
	Red         = '0'
	Yellow      = '1'
	Blue        = '2'
	Green       = '3'
	Pink        = '4'
	Orange      = '5'
	LightBlue   = '6'
	Purple      = '7'
	Brown       = '8'
	LightGreen  = '9'
	Cyan        = 'A'
	LightPink   = 'B'
	White       = 'C'
	LightPurple = 'D'
	LightBrown  = 'E'
	OtherWhite  = 'F'
)

type Block struct {
	Object
	Color Color
}

type Switch struct {
	Object
	line, col int
	name      string
}

// Blocks returns the block arround the switch in parameter.
func (l *Level) Blocks(sw *Switch) []*Block {
	topLeft := l.blocks[sw.line][sw.col]
	topRight := l.blocks[sw.line][sw.col+1]
	bottomLeft := l.blocks[sw.line+1][sw.col]
	bottomRight := l.blocks[sw.line+1][sw.col+1]
	return []*Block{topLeft, topRight, bottomLeft, bottomRight}
}

func (s *Switch) String() string {
	return fmt.Sprintf("sw{line:%d, col:%d}", s.line, s.col)
}

func (l *Level) Copy() Level {
	lcp := new(Level)
	lcp.blocks = make([][]*Block, len(l.blocks))
	for i := range l.blocks {
		lcp.blocks[i] = make([]*Block, len(l.blocks[i]))
		for j := range l.blocks[i] {
			lcp.blocks[i][j] = &Block{Color: l.blocks[i][j].Color}
		}
	}

	lcp.switches = make([]*Switch, len(l.switches))
	for i := range l.switches {
		sw := l.switches[i]
		lcp.switches[i] = &Switch{col: sw.col, line: sw.line, name: sw.name}
	}
	lcp.winSignature = l.winSignature
	return *lcp
}

// Win returns true if player has win
func (l *Level) Win() bool {
	for i := range l.winSignature {
		for j := range l.winSignature[i] {
			if l.winSignature[i][j] != l.blocks[i][j].Color {
				return false
			}
		}
	}
	return true
}

// UndoLastMove cancels the last player move
func (l *Level) UndoLastMove() {
	if l.rotating != nil {
		return
	}
	sw := l.PopLastRotated()
	if sw != nil {
		g.level.rotating = sw
		blocks := l.Blocks(sw)
		for i := range blocks {
			b := blocks[i]
			v := switchSize / 2
			b.Time = 0
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

func newBlock(color Color, line, col int, size, padding float32) *Block {
	colf, linef := float32(col), float32(line)
	b := &Block{Color: color}
	b.Object = Object{
		X:      colf * (size + padding),
		Y:      linef * (size + padding),
		Width:  size,
		Height: size,
		Data:   b,
	}
	return b
}

func (l *Level) addBlock(color Color, line, col int) {
	b := newBlock(color, line, col, blockSize, blockPadding)
	b.X += xMin
	b.Y += yMin
	b.Action = wait{until: clock.Time(line*10 + col*5), next: ActionFunc(blockPopIn)}
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
	log.Println("Switch added", s.X, s.Y)
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
			l.triggerSwitch(i)
		}
	}
}

func (l *Level) triggerSwitchName(name string) {
	for i := 0; i < len(l.switches); i++ {
		if l.switches[i].name == name {
			l.triggerSwitch(i)
			return
		}
	}
}

func (l *Level) triggerSwitch(i int) {
	sw := l.switches[i]
	blocks := l.Blocks(sw)
	for i := range blocks {
		b := blocks[i]
		v := switchSize / 2
		// Prepare a rotation around the center of the switch
		b.Rx, b.Ry = sw.X+v, sw.Y+v
		b.Sx, b.Sy = b.Rx, b.Ry
		b.Time = 0
		b.Action = ActionFunc(blockRotate)
	}
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
			signature.WriteRune(rune(l.blocks[i][j].Color))
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

func ctoa(c Color) string {
	return fmt.Sprintf("%d", c)
}

// LoadLevel loads the level number in parameter
func LoadLevel(level int) Level {
	f, err := asset.Open(fmt.Sprintf("levels/%d", level))
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	l := ParseLevel(string(b))
	log.Printf("Level loaded %d\n", level)
	return l
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
				l.addBlock(Color(c), i, j)
			}
		case 1:
			// read switch locations
			tokens := strings.Split(lines[i], ",")
			l.addSwitch(atoi(tokens[0]), atoi(tokens[1]))
		case 2:
			//read win
			wline := make([]Color, len(lines[i]))
			for j, c := range lines[i] {
				wline[j] = Color(c)
			}
			l.winSignature = append(l.winSignature, wline)
		case 3:
			// read the max move count
			l.maxMoves = atoi(lines[i])
		case 4:
			// read the solution
			l.solution = lines[i]
		}
	}
	return l
}

// RotateSwitch swaps bocks according to the 90d rotation
func (l *Level) RotateSwitch(s *Switch) {
	li, co := s.line, s.col
	log.Println("Swap from", s.name, li, co)
	color := l.blocks[li][co].Color
	l.blocks[li][co].Color = l.blocks[li+1][co].Color
	l.blocks[li+1][co].Color = l.blocks[li+1][co+1].Color
	l.blocks[li+1][co+1].Color = l.blocks[li][co+1].Color
	l.blocks[li][co+1].Color = color
	l.moves++
}

// RotateSwitchInverse swaps bocks according to the -90d rotation
func (l *Level) RotateSwitchInverse(s *Switch) {
	li, co := s.line, s.col
	color := l.blocks[li][co].Color
	l.blocks[li][co].Color = l.blocks[li][co+1].Color
	l.blocks[li][co+1].Color = l.blocks[li+1][co+1].Color
	l.blocks[li+1][co+1].Color = l.blocks[li+1][co].Color
	l.blocks[li+1][co].Color = color
	l.moves--
}
