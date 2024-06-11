package mcts

import (
	"github.com/sw965/bippa/team/game"
	"github.com/sw965/crow/mcts/uct"
	"github.com/sw965/crow/ucb"
	"github.com/sw965/crow/model/1d"
	"github.com/sw965/bippa/team"
	"github.com/sw965/bippa/feature"

)

func NewStandard(affine model1d.Sequential) uct.MCTS[team.Team, team.Actions, team.Action] {
	g := game.New()

	leafNodeEvalFunc := func(party *team.Team) uct.LeafNodeEvalY {
		x := feature.Team(
			2,
			feature.FirePowerIndex,
			feature.DefenseIndex,
		)(*party)
		v, err := affine.Predict(x)
		if err != nil {
			panic(err)
		}
		return uct.LeafNodeEvalY(v[0])
	}

	eachPlayerEvalFunc := func(y uct.LeafNodeEvalY, _ *team.Team) uct.EachPlayerEvalY {
		return uct.EachPlayerEvalY(y)
	}

	mcts := uct.MCTS[team.Team, team.Actions, team.Action]{
		UCBFunc:ucb.NewAlphaGoFunc(5),
		Game:g,
		EvalFunc:uct.EvalFunc[team.Team]{LeafNode:leafNodeEvalFunc, EachPlayer:eachPlayerEvalFunc},
		NextNodesCap:64,
	}
	mcts.SetUniformActionPolicyFunc()
	return mcts
}