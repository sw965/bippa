package mcts

import (
	"github.com/sw965/crow/mcts/duct"
	"github.com/sw965/crow/ucb"
	"github.com/sw965/bippa/battle"
	"github.com/sw965/bippa/battle/game"
)

func New(c float64) duct.MCTS[battle.Manager, battle.ActionsSlice, battle.Actions, battle.Action] {
	mcts := duct.MCTS[battle.Manager, battle.ActionsSlice, battle.Actions, battle.Action]{
		Game:game.New(),
		UCBFunc:ucb.NewAlphaGoFunc(c),
		NextNodesCap:128,
	}
	mcts.SetSeparateUniformActionPolicyFunc()
	mcts.SetRandomPlayoutLeafNodeJointEvalFunc(battle.GlobalContext.Rand)
	return mcts
}