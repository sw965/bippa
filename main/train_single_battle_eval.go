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

func main() {
	r := omwrand.NewMt19937()
	xn := 90
	fmt.Println("xn", xn)
	u1n := 64
	u2n := 16
	yn := 1
	affine, variable := model1d.NewStandardAffine(xn, u1n, u2n, yn, 0.0001, 64.0, r)

	leafNodeEvalsFunc := func(battle *single.Battle) (duct.LeafNodeJointEvalY, error) {
		if isEnd, gameRetJointVal := gm.IsEnd(battle); isEnd {
			y := make(duct.LeafNodeJointEvalY, len(gameRetVal))
			for i, v := range gameRetJointVal {
				y[i] = v
			}
			return y
		} else {
			x := feature.NewSingleBattleFunc(2, feature.ExpectedDamageRatioToCurrentHP, feature.DPSRatioToCurrentHP)(battle)
			y, err := affine.Predict(x)
			v := y[0]
			return duct.LeafNodeEvalYs{v, 1.0-v}, nil
		}
	}

	batchSize := 1280
	trainX := make(tensor.D2, 0, batchSize)
	trainY := make(tensor.D2, 0, batchSize)

	gm := game.New(dmgtools.RandBonuses{1.0}, r)
	mcts := mcts.New(dmgtools.RandBonuses{1.0}, r)

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

		if isEnd, _ := gm.IsEnd(&battle); isEnd {
			continue
		}

		gm.Player = mcts.NewPlayer(512, r)
		_, battleHistory, _, qHistory, err = gm.PlayoutWithHistory(battle)
		if err != nil {
			panic(err)
		}

		swappedBattleHistory := make(single.Battle, len(battleHistory))
		for i, b := range battleHistory {
			swappedBattleHistory[i] = b.SwapPlayers()
		}

		for 
		
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