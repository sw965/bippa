package main

import (
	"fmt"
	"github.com/sw965/bippa/battle/single/game"
	omwrand "github.com/sw965/omw/math/rand"
	"github.com/sw965/bippa/feature"
	"github.com/sw965/crow/model/1d"
	"github.com/sw965/bippa/battle/single"
	bp "github.com/sw965/bippa"
	"github.com/sw965/crow/tensor"
	omwslices "github.com/sw965/omw/slices"
	"github.com/sw965/bippa/battle/single/mcts"
)

type RawTeacher struct {
	Battle single.Battle
	JointQ []float64
	GameResultJointValue []float64
}

type RawTeachers []*RawTeacher

func NewRawTeachers(battleHistory []single.Battle, jointQs [][]float64, gameRetJointVal []float64) RawTeachers {
	ts := make(RawTeachers, len(battleHistory))
	for i := range ts {
		ts[i] = &RawTeacher{
			Battle:battleHistory[i],
			JointQ:jointQs[i],
			GameResultJointValue:gameRetJointVal,
		}
	}
	return ts
}

func (ts RawTeachers) DataAugmentation() RawTeachers {
	data := make(RawTeachers, 0, len(ts) * 3)
	for _, t := range ts {
		data = append(data, &RawTeacher{
			Battle:t.Battle.SwapPlayers(),
			JointQ:omwslices.Reverse(t.JointQ),
			GameResultJointValue:omwslices.Reverse(t.GameResultJointValue),
		})

		data = append(data, &RawTeacher{
			Battle:single.Battle{SelfFighters:t.Battle.SelfFighters, OpponentFighters:t.Battle.SelfFighters},
			JointQ:[]float64{0.5, 0.5},
			GameResultJointValue:[]float64{0.5, 0.5},
		})

		data = append(data, &RawTeacher{
			Battle:single.Battle{SelfFighters:t.Battle.OpponentFighters, OpponentFighters:t.Battle.OpponentFighters},
			JointQ:[]float64{0.5, 0.5},
			GameResultJointValue:[]float64{0.5, 0.5},
		})
	}
	return omwslices.Concat(ts, data)
}

func (ts RawTeachers) ToTrainXY(f feature.SingleBattleFunc, qRatio float64) (tensor.D2, tensor.D2) {
	n := len(ts)
	x := make(tensor.D2, n)
	y := make(tensor.D2, n)
	for i, t := range ts {
		x[i] = f(&t.Battle)
		y[i] = tensor.D1{(qRatio * t.JointQ[0]) + (1.0-qRatio) * t.GameResultJointValue[0]}
	}
	return x, y
}

func main() {
	rg := omwrand.NewMt19937()
	battleContext := single.NewContext(rg)
	gameManager := game.New(&battleContext)

	xn := 90
	u1n := 64
	u2n := 16
	yn := 1
	affine, variable := model1d.NewStandardAffine(xn, u1n, u2n, yn, 0.0001, 64.0, rg)
	featureFunc := feature.NewSingleBattleFunc(feature.ExpectedDamageRatioToCurrentHP, feature.DPSRatioToCurrentHP)

	mctSearch := mcts.New(&battleContext)
	mctSearch.LeafNodeJointEvalFunc = mcts.NewLeafNodeJointEvalFunc(affine, featureFunc)

	batchSize := 2560
	trainX := make(tensor.D2, 0, batchSize)
	trainY := make(tensor.D2, 0, batchSize)

	randActionPlayer := gameManager.NewRandActionPlayer(rg)
	mctsPlayer := mctSearch.NewPlayer(512, rg)
	selfBattleNum := 512

	for i := 0; i < selfBattleNum; i++ {
		initBattle := single.Battle{
			SelfFighters:bp.Pokemons{bp.NewTemplateBulbasaur(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
			OpponentFighters:bp.Pokemons{bp.NewTemplateBulbasaur(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
		}

		battle, err := gameManager.Play(randActionPlayer, initBattle, func(_ *single.Battle, i int) bool { return i == rg.Intn(8) })
		if err != nil {
			panic(err)
		}

		if isEnd, _ := game.IsEnd(&battle); isEnd {
			continue
		}

		endBattle, battleHistory, _, qHistory, err := gameManager.PlayoutWithHistory(mctsPlayer, battle, 128)
		if err != nil {
			panic(err)
		}

		_, gameRetJointVal := game.IsEnd(&endBattle)
		oneGameTrainX, oneGameTrainY := NewRawTeachers(battleHistory, qHistory, gameRetJointVal).DataAugmentation().ToTrainXY(featureFunc, 0.5)
		trainX = append(trainX, oneGameTrainX...)
		trainY = append(trainY, oneGameTrainY...)

		if len(trainX) >= batchSize - 196 {
			fmt.Println("minibatch")
			for j := 0; j < batchSize; j++ {
				idx := rg.Intn(len(trainX))
				affine.SGD(trainX[idx], trainY[idx], 0.01)
			}
			trainX = make(tensor.D2, 0, batchSize)
			trainY = make(tensor.D2, 0, batchSize)
		}
		fmt.Println("i = ", i)
	}
	variable.Param.WriteJSON("C:/Go/project/bippa/main/single_battle/nn_param.json")
}