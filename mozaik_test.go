package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() {
	g = NewGame()
}

func fill() {
	g.level = ParseLevel(
		`02
14

0,0`)
}

func TestFindSwitch(t *testing.T) {
	setup()
	fill()

	_, s := g.level.findSwitch(XMin+BlockSize, YMin+BlockSize)

	assert.NotNil(t, s, "Should found a switch")
}

func TestParseLevel(t *testing.T) {
	lvl := `01
24

0,0

24`

	l := ParseLevel(lvl)

	assert.Equal(t, 1, len(l.switches))
	assert.Equal(t, 0, l.switches[0].line)
	assert.Equal(t, 0, l.switches[0].col)
	assert.Equal(t, 2, len(l.blocks))
	assert.Equal(t, 2, len(l.blocks[0]))
	assert.Equal(t, 2, len(l.blocks[1]))
	assert.Equal(t, 0, l.blocks[0][0].Color)
	assert.Equal(t, 1, l.blocks[0][1].Color)
	assert.Equal(t, 2, l.blocks[1][0].Color)
	assert.Equal(t, 4, l.blocks[1][1].Color)
	assert.Equal(t, "24\n", l.winSignature)
}

func TestRotateState(t *testing.T) {
	setup()
	fill()

	g.level.switches[0].ChangeState(NewRotateState())
	g.level.switches[0].ChangeState(NewIdleState())

	l := g.level
	assert.Equal(t, 2, len(l.blocks))
	assert.Equal(t, 2, len(l.blocks[0]))
	assert.Equal(t, 2, len(l.blocks[1]))
	assert.Equal(t, Yellow, l.blocks[0][0].Color)
	assert.Equal(t, Red, l.blocks[0][1].Color)
	assert.Equal(t, Blue, l.blocks[1][1].Color)
	assert.Equal(t, Pink, l.blocks[1][0].Color)
}

func TestBlockSignature(t *testing.T) {
	setup()

	l := ParseLevel(`01
24

0,0`)

	signature := `01
24
`
	assert.Equal(t, signature, l.blockSignature())
}
