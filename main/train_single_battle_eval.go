package main

import (
	"fmt"
	"github.com/sw965/bippa/battle/single/game"
	"github.com/sw965/crow/mcts/duct"
	omwrand "github.com/sw965/omw/math/rand"
	"github.com/sw965/bippa/feature"
	"github.com/sw965/crow/model/1d"
	"github.com/sw965/bippa/battle/single"
	"github.com/sw965/crow/ucb"
	"github.com/sw965/bippa/battle/dmgtools"
	bp "github.com/sw965/bippa"
	"github.com/sw965/crow/tensor"
)

type RawTeacher struct {
	Battle single.Battle
	Q float64
	GameResult float64
}

type RawTeachers []*RawTeacher

func (ts RawTeachers) TrainXY(f feature.SingleBattleFunc, qRatio float64) tensor.D2 {
	x := make(tensor.D1, 0, len(ts))
	y := make(tensor.D2, 0, len(ts))
	for i, t := range ts {
		x[i] = f(&t.Battle)
		y[i] = (qRatio * t.Q) + ((1.0 - qRatio) * t.GameResult) 
	}
	return x, y
}

type CaitlinModel[single.Battle, single.Actionss, single.Actions] struct {
	Game simultaneous.Game
	MCTS duct.MCTS[single.Battle, single.Actionss, single.Actions]
	FeatureFunc feature.SingleBattleFunc
	Affine model1d.Sequential
	NewInitBattleFunc func() Battle 
	RawTeachers RawTeachers
	RandDmgBonuses dmgtools.RandBonuses
	Rand *rand.Rand
	rawTeacherStart int
}

func (lady *Caitlin) SetMCTS() {
	leafNodeEvalsFunc := func(battle *single.Battle) (duct.LeafNodeEvalYs, error) {
		gameResult := battle.GameResult()
		if gameResult == -1.0 {
			x := lady.FeatureFunc(battle)
			y, err := affine.Predict(x)
			if err != nil {
				return duct.LeafNodeEvalYs{}, err
			}
			v := y[0]
			return duct.LeafNodeEvalYs{
				duct.LeafNodeEvalY(v),
				duct.LeafNodeEvalY(1.0-v),
			}, nil 
		} else {
			return duct.LeafNodeEvalYs{gameResult, 1.0-gameResult}, nil
		}
	}

	gm := game.New(dmgtools.RandBonuses{1.0}, lady.Rand)
	gm.Player = func(battle *Battle) single.Action {}

	mcts := duct.MCTS[single.Battle, single.Actionss, single.Actions, single.Action]{
		UCBFunc:ucb.NewAlphaGoFunc(5),
		Game:gm,
		LeafNodeEvalsFunc:leafNodeEvalsFunc,
		NextNodesCap:64,
		LastJointActionsCap:1,
	}

	mcts.SetUniformActionPoliciesFunc()
	lady.MCTS = mcts
}

func (lady *Caitlin) SetSelfBattlePlayer(simulation int) {
	lady.Game.Player = func(battle *single.Battle) (single.Actions, error) {
		rootNode := lady.MCTS.NewNode(battle) 
		err := lady.MCTS.Run(simulation, rootNode, lady.Rand)
		if err != nil {
			return single.Actions{}, err
		}
		
		jointAvg := rootNode.UCBManagers.JointAverageValue()
		teacher := RawTeacher{Battle:*battle, Q:jointAvg[0]}
		return rootNode.UCBManagers.JointActionByMaxTrial(lady.Rand), nil
	}
}

func (lady *Caitlin) SelfBattle(gameNum, maxRandActionNum int) error {
	for i := 0; i < n; i++ {
		initBattle := lady.NewInitBattleFunc()

		lady.Game.SetRandActionPlayer(lady.Rand)
		lady.rawTeacherStart := len(lady.RawTeachers)
		battle, err := lady.Game.Play(initBattle, func(_ *single.Battle, i int) bool { return i == lady.Rand.Intn(maxRandActionNum+1)} )
		if err != nil {
			return err
		}
		if lady.IsEnd(&battle) {
			continue
		}

		lady.SetSelfBattlePlayer()
		endBattle, err := lady.Game.Playout(battle)
		if err != nil {
			return err
		}

		gameResult := endBattle.GameResult()
		for i := lady.rawTeacherStart; i < len(lady.RawTeachers); i++ {
			lady.RawTeachers[i].GameResult = gameResult
		}
	}
}

func (lady *Caitlin) Train() {

}

func main() {
	r := omwrand.NewMt19937()
	xn := 90
	fmt.Println("xn", xn)
	u1n := 64
	u2n := 16
	yn := 1
	affine, variable := model1d.NewStandardAffine(xn, u1n, u2n, yn, 0.0001, 64.0, r)

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

	selfBattleNum := 1280
	for i := 0; i < selfBattleNum; i++ {
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
	variable.Param.WriteJSON("c:/Users/kuroko/Desktop/test.json")
}