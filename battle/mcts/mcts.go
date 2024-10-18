package mcts

import (
	"github.com/sw965/bippa/battle"
	"github.com/sw965/bippa/battle/game"
	"github.com/sw965/crow/mcts/duct"
	"github.com/sw965/crow/ucb"
)

func New(c float64) duct.MCTS[battle.Manager, battle.ActionsSlice, battle.Actions, battle.Action] {
	mcts := duct.MCTS[battle.Manager, battle.ActionsSlice, battle.Actions, battle.Action]{
		GameLogic:         game.NewLogic(),
		UCBFunc:      ucb.NewAlphaGoFunc(c),
		NextNodesCap: 128,
	}
	mcts.SetSeparateUniformActionPolicyProvider()
	mcts.SetRandPlayout(game.ResultJointEvaluator, battle.GlobalContext.Rand)
	return mcts
}
