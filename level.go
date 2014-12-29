package main

import (
	"fmt"
	"golang.org/x/mobile/app"
	"io/ioutil"
	"strconv"
	"strings"
)

type Level struct {
	blocks       [][]*Block
	switches     []*Switch
	winSignature string
	// rotated represents the historics of rotations
	rotated []int
	// rotating represents a rotate which
	// is currently rotating
	rotating *Switch
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
	return l.winSignature == l.blockSignature()
}

func (l *Level) HowFar() int {
	howfar := 0
	signature := l.blockSignature()
	for i := range l.winSignature {
		if l.winSignature[i] != signature[i] {
			howfar++
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
		sw.ChangeState(NewRotateStateReverse())
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

// addSwitch appends a new switch at the bottom right
// of the coordinates in parameters.
func (l *Level) addSwitch(line, col int) {

	s := &Switch{
		line: line, col: col,
		X:    XMin + (col+1)*BlockSize + col*BlockPadding*2 - SwitchSize/2,
		Y:    YMin + (line+1)*BlockSize + line*BlockPadding*2 - SwitchSize/2,
		name: determineName(line, col),
	}
	s.ChangeState(NewIdleState())
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
func (l *Level) PressSwitch(x, y int) {
	// Handle click only when no switch are rotating
	if l.rotating == nil {
		if i, s := l.findSwitch(x, y); s != nil {
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

func (l *Level) findSwitch(x, y int) (int, *Switch) {
	for i, s := range l.switches {
		if x >= s.X && x <= s.X+SwitchSize && y >= s.Y && y <= s.Y+SwitchSize {
			return i, s
		}
	}
	return -1, nil
}

func (l *Level) blockSignature() string {
	var signature string
	for i := 0; i < len(l.blocks); i++ {
		for j := 0; j < len(l.blocks[i]); j++ {
			if l.blocks[i][j] == nil {
				signature += "-"
			} else {
				signature += ctoa(l.blocks[i][j].Color)
			}
		}
		signature += "\n"
	}
	return signature
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func atoc(s string) ColorDef {
	return ColorDef(atoi(s))
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
					bline[j] = &Block{Color: atoc(string(c))}
				}
			}
		case 1:
			// read switch locations
			tokens := strings.Split(lines[i], ",")
			l.addSwitch(atoi(tokens[0]), atoi(tokens[1]))
		case 2:
			//read win
			l.winSignature += lines[i] + "\n"
		}
	}
	//fmt.Printf("Level loaded blocks=%d, swicthes=%d\n", len(l.blocks), len(l.switches))

	//for i := 0; i < len(l.blocks); i++ {
	//	fmt.Printf("line %d blocks %d\n", i, len(l.blocks[i]))
	//}
	//fmt.Printf("winSignature\n%s\n---\n", l.winSignature)
	return l
}

func (lvl *Level) RotateSwitch(s *Switch) {
	// Swap bocks according to the 90d rotation
	l, c := s.line, s.col
	//fmt.Println("Swap from", l, c)
	b := lvl.blocks[l][c]
	lvl.blocks[l][c] = lvl.blocks[l+1][c]
	lvl.blocks[l+1][c] = lvl.blocks[l+1][c+1]
	lvl.blocks[l+1][c+1] = lvl.blocks[l][c+1]
	lvl.blocks[l][c+1] = b
}

func (lvl *Level) RotateSwitchInverse(s *Switch) {
	// Swap bocks according to the -90d rotation
	l, c := s.line, s.col
	b := lvl.blocks[l][c]
	lvl.blocks[l][c] = lvl.blocks[l][c+1]
	lvl.blocks[l][c+1] = lvl.blocks[l+1][c+1]
	lvl.blocks[l+1][c+1] = lvl.blocks[l+1][c]
	lvl.blocks[l+1][c] = b
}
