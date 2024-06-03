package mcts

import (
	"fmt"
	"github.com/sw965/crow/ucb"
	"github.com/sw965/crow/model"
	"github.com/sw965/bippa/battle/single"
	"github.com/sw965/crow/mcts/duct"
	omwslices "github.com/sw965/omw/slices"
	"math/rand"
	"github.com/sw965/bippa/battle/dmgtools"
	"github.com/sw965/bippa/battle/single/game"
	"github.com/sw965/crow/tensor"
)

var (
	P1_WIN_LEAF_NODE_EVAL_YS = duct.LeafNodeEvalYs{1.0, 0.5}
	P2_WIN_LEAF_NODE_EVAL_YS = omwslices.Reverse(P1_WIN_LEAF_NODE_EVAL_YS)
	DRAW_LEAF_NODE_EVAL_YS = duct.LeafNodeEvalYs{0.5, 0.5}
)

func EndLeafNodeEvalsFunc(battle *single.Battle) (duct.LeafNodeEvalYs, error) {
	isP1AllFaint := battle.P1Fighters.IsAllFaint()
	isP2AllFaint := battle.P2Fighters.IsAllFaint()
	if isP1AllFaint && isP2AllFaint {
		return DRAW_LEAF_NODE_EVAL_YS, nil
	} else if isP1AllFaint {
		return P2_WIN_LEAF_NODE_EVAL_YS, nil
	} else if isP2AllFaint {
		return P1_WIN_LEAF_NODE_EVAL_YS, nil
	} else {
		return duct.LeafNodeEvalYs{}, fmt.Errorf("ゲームが終わっていないので、EndLeafNodeEvalsFuncを計算できません。")
	}
}

func NewInputOutput1DModelLeafNodeEvalsFunc(io1dModel model.SequentialInputOutput1D, f func(*single.Battle) tensor.D1) duct.LeafNodeEvalsFunc[single.Battle] {
	return func(battle *single.Battle) (duct.LeafNodeEvalYs, error) {
		if ys, err := EndLeafNodeEvalsFunc(battle); err == nil {
			return ys, err
		} else {
			x := f(battle)
			v, err := io1dModel.Predict(x)
			p1v := duct.LeafNodeEvalY(v[0])
			p2v := duct.LeafNodeEvalY(1.0-p1v)
			return duct.LeafNodeEvalYs{p1v, p2v}, err
		}
	}
}

func New(randDmgBonuses dmgtools.RandBonuses, r *rand.Rand) duct.MCTS[single.Battle, single.Actionss, single.Actions, single.Action] {
	gm := game.New(randDmgBonuses, r)

	leafNodeEvalsFunc := func(battle *single.Battle) (duct.LeafNodeEvalYs, error) {
		endBattle, err := gm.Playout(*battle)
		if err != nil {
			return duct.LeafNodeEvalYs{}, err
		}
		ys, err := EndLeafNodeEvalsFunc(&endBattle)
		return ys, err
	}

	mcts := duct.MCTS[single.Battle, single.Actionss, single.Actions, single.Action]{
		UCBFunc:ucb.NewAlphaGoFunc(5),
		Game:gm,
		LeafNodeEvalsFunc:leafNodeEvalsFunc,
		NextNodesCap:64,
	}
	mcts.SetUniformActionPoliciesFunc()
	return mcts
}