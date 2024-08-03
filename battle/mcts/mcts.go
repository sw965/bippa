package mcts

import (
	"github.com/sw965/crow/mcts/duct"
	"github.com/sw965/bippa/battle/single/game"
	"github.com/sw965/bippa/battle/single"
	"github.com/sw965/crow/ucb"
)

func New(context *single.Context) duct.MCTS[single.Battle, single.ActionsSlice, single.Actions, single.Action] {
	mcts := duct.MCTS[single.Battle, single.ActionsSlice, single.Actions, single.Action]{
		Game:game.New(context),
		UCBFunc:ucb.NewAlphaGoFunc(5),
		NextNodesCap:64,
		LastJointActionsCap:1,
	}
	mcts.SetUniformSeparateActionPolicyFunc()
	return mcts
}