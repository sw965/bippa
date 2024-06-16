package mcts

import (
	"math/rand"
	"github.com/sw965/crow/mcts/duct"
	"github.com/sw965/bippa/battle/single/game"
	"github.com/sw965/bippa/battle/single"
	"github.com/sw965/crow/ucb"
)

func New(r *rand.Rand) duct.MCTS[single.Battle, single.ActionSlices, single.Actions, single.Action] {
	mcts := duct.MCTS[single.Battle, single.ActionSlices, single.Actions, single.Action]{
		Game:game.New(r),
		UCBFunc:ucb.NewAlphaGoFunc(5),
		NextNodesCap:64,
		LastJointActionsCap:1,
	}
	mcts.SetUniformSeparateActionPolicyFunc()
	return mcts
}