package main

const (
	WindowWidth          = 576
	WindowHeight         = 704
	BlockSize            = 128
	BlockRadius          = 10
	SwitchSize           = 48
	DashboardHeight      = 128
	XMin                 = 32
	YMin                 = 32
	XMax                 = WindowHeight - 32
	YMax                 = WindowWidth - 32 - DashboardHeight
	SignatureBlockSize   = 32
	SignatureBlockRadius = 6
	LineWidth            = 2
	SignatureLineWidth   = 1
)

type World struct {
	background *Background
	switches   []*SwitchModel
	blocks     map[*Block]*BlockModel
}

func (w *World) Reset() {
	// Clean
	if len(w.switches) > 0 {
		for _, s := range w.switches {
			s.Destroy()
		}
	}
	w.switches = nil
	if len(w.blocks) > 0 {
		for k, m := range w.blocks {
			m.Destroy()
			delete(w.blocks, k)

		}
	}
	// Create the blocks
	w.blocks = make(map[*Block]*BlockModel)
	for i := 0; i < len(g.level.blocks); i++ {
		for j := 0; j < len(g.level.blocks[i]); j++ {
			b := g.level.blocks[i][j]
			if b != nil {
				w.blocks[b] = NewBlockModel(b, BlockSize, BlockRadius)
			}
		}
	}
	for _, sw := range g.level.switches {
		w.switches = append(w.switches, NewSwitchModel(sw))
	}
}

type Game struct {
	currentLevel int
	level        Level
	listen       bool
	world        *World
}

func NewGame() *Game {
	return &Game{currentLevel: 1, listen: true, world: &World{}}
}

func (g *Game) Start() {
	// Load first level
	g.level = LoadLevel(g.currentLevel)
	g.world.background = NewBackground()
	g.world.Reset()
}

func (g *Game) Stop() {
}

func (g *Game) Click(x, y int) {
	if g.listen {
		g.level.PressSwitch(x, y)
	}
}

func (g *Game) Listen() bool {
	return g.listen && g.level.rotating == nil
}

func (g *Game) Update() {
	for _, s := range g.level.switches {
		s.state.Update(g, s)
	}
}

func (g *Game) Continue() {
	if g.level.Win() {
		g.Warp()
	}
}

func (g *Game) Warp() {
	if g.listen {
		// Next level
		g.currentLevel++
		g.level = LoadLevel(g.currentLevel)
		g.world.Reset()
	}
}

func (g *Game) UndoLastMove() {
	if g.listen {
		g.level.UndoLastMove()
	}
}

func (g *Game) Reset() {
	sw := g.level.PopLastRotated()
	if sw != nil {
		g.listen = false
		sw.ChangeState(NewResetState())
	}
}
