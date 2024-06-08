package mcts

import (
	"github.com/sw965/bippa/team/game"
	"github.com/sw965/crow/mcts/duct"
)

func New() uct.MCTS[team.Team, team.Action] {
	g := game.New()
	mcts := ucb.MCTS[team.Team, team.Action]{
		UCBFunc:ucb.NewAlphaGoFunc(5),
		Game:game,
		
	}
}