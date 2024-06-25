package mcts

import (
	"github.com/sw965/crow/mcts/duct"
	"github.com/sw965/bippa/battle/single/game"
	"github.com/sw965/bippa/battle/single"
	"github.com/sw965/crow/ucb"
	"github.com/sw965/bippa/feature"
	"github.com/sw965/crow/model/1d"
)

func NewLeafNodeJointEvalFunc(model model1d.Sequential, f feature.SingleBattleFunc) duct.LeafNodeJointEvalFunc[single.Battle] {
	return func(battle *single.Battle) (duct.LeafNodeJointEvalY, error) {
		if isEnd, gameRetJointVal := game.IsEnd(battle); isEnd {
			y := make(duct.LeafNodeJointEvalY, len(gameRetJointVal))
			for i, v := range gameRetJointVal {
				y[i] = v
			}
			return y, nil
		} else {
			x := f(battle)
			y, err := model.Predict(x)
			v := y[0]
			return duct.LeafNodeJointEvalY{v, 1.0-v}, err
		}
	}
}

func New(context *single.Context) duct.MCTS[single.Battle, single.ActionSlices, single.Actions, single.Action] {
	mcts := duct.MCTS[single.Battle, single.ActionSlices, single.Actions, single.Action]{
		Game:game.New(context),
		UCBFunc:ucb.NewAlphaGoFunc(5),
		NextNodesCap:64,
		LastJointActionsCap:1,
	}
	mcts.SetUniformSeparateActionPolicyFunc()
	return mcts
}