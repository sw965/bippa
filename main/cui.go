package main

import (
	"github.com/sw965/crow/model/1d"
)

func main() {
	r := omwrand.NewMt19937()
	xn := 90
	u1n := 64
	u2n := 16
	yn := 1
	affine, variable := model1d.NewStandardAffine(xn, u1n, u2n, yn, 0.0001, 64.0, r)
	param := model1d.LoadParamJSON("C:/Users/kuroko/Desktop/test.json")
	variable.SetParam(param)

	leafNodeEvalsFunc := func(battle *single.Battle) (duct.LeafNodeEvalYs, error) {
		isP1AllFaint := battle.P1Fighters.IsAllFaint()
		isP2AllFaint := battle.P2Fighters.IsAllFaint()
		if isP1AllFaint && isP2AllFaint {
			return duct.LeafNodeEvalYs{0.5, 0.5}, nil
		} else if isP1AllFaint {
			return duct.LeafNodeEvalYs{0.0, 1.0}, nil
		} else if isP2AllFaint {
			return duct.LeafNodeEvalYs{1.0, 0.0}, nil
		} else {
			x := feature.SingleBattleClosure(2, feature.ExpectedDamageRatioToCurrentHP, feature.DPSRatioToCurrentHP)(battle)
			y, err := affine.Predict(x)
			if err != nil {
				return duct.LeafNodeEvalYs{}, err
			}
			return duct.LeafNodeEvalYs{
				duct.LeafNodeEvalY(y[0]),
				duct.LeafNodeEvalY(1.0-y[0]),
			}, nil
		}
	}

	mcts := duct.MCTS[single.Battle, single.Actionss, single.Actions, single.Action]{
		UCBFunc:ucb.NewAlphaGoFunc(5),
		Game:game.New(dmgtools.RandBonuses{1.0}, r),
		LeafNodeEvalsFunc:leafNodeEvalsFunc,
		NextNodesCap:64,
		LastJointActionsCap:1,
	}
	mcts.SetUniformActionPoliciesFunc()

	batchSize := 1280
	trainX := make(tensor.D2, 0, batchSize)
	trainY := make(tensor.D2, 0, batchSize)

	gm := game.New(dmgtools.RandBonuses{1.0}, r)
	mctsPlayer := func(battle *single.Battle) (single.Actions, error) {
		rootNode := mcts.NewNode(battle)
		err := mcts.Run(128, rootNode, r)
		jointAction := rootNode.UCBManagers.JointActionByMaxTrial(r)
		avgs := rootNode.UCBManagers.JointAverageValue()
		swaped := battle.SwapPlayers()
		trainX = append(trainX, feature.SingleBattleClosure(2, feature.ExpectedDamageRatioToCurrentHP, feature.DPSRatioToCurrentHP)(battle))
		trainX = append(trainX, feature.SingleBattleClosure(2, feature.ExpectedDamageRatioToCurrentHP, feature.DPSRatioToCurrentHP)(&swaped))
		trainY = append(trainY, tensor.D1{avgs[0]})
		trainY = append(trainY, tensor.D1{avgs[1]})
		return jointAction, err
	}

	initBattle := single.Battle{
		P1Fighters:single.Fighters{bp.NewTemplateBulbasaur(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
		P2Fighters:single.Fighters{bp.NewTemplateBulbasaur(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
	}

	gm.SetRandActionPlayer(r)
	battle, err := gm.Play(initBattle, func(_ *single.Battle, i int) bool { return i == r.Intn(8) })
	if err != nil {
		panic(err)
	}
	if gm.IsEnd(&battle) {
		continue
	}

	gm.Player = mctsPlayer
		_, err = gm.Playout(battle)
		if err != nil {
			panic(err)
		}
		fmt.Println("i=", i)
		if len(trainX) >= batchSize - 100 {
			for j := 0; j < batchSize; j++ {
				idx := r.Intn(len(trainX))
				affine.SGD(trainX[idx], trainY[idx], 0.01)
			}
			trainX = make(tensor.D2, 0, batchSize)
			trainY = make(tensor.D2, 0, batchSize)
		}
	}
}