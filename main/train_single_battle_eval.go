package main

import (
	"fmt"
	"github.com/sw965/bippa/battle/single/game"
	"github.com/sw965/crow/mcts/duct"
	omwrand "github.com/sw965/omw/math/rand"
	"github.com/sw965/bippa/feature"
	"github.com/sw965/crow/model/1d"
	"github.com/sw965/bippa/battle/single"
	"github.com/sw965/bippa/battle/dmgtools"
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
			Battle:single.Battle{P1Fighters:t.Battle.P1Fighters, P2Fighters:t.Battle.P1Fighters},
			JointQ:[]float64{0.5, 0.5},
			GameResultJointValue:[]float64{0.5, 0.5},
		})

		data = append(data, &RawTeacher{
			Battle:single.Battle{P1Fighters:t.Battle.P2Fighters, P2Fighters:t.Battle.P2Fighters},
			JointQ:[]float64{0.5, 0.5},
			GameResultJointValue:[]float64{0.5, 0.5},
		})
	}
	return omwslices.Concat(ts, data)
}

func (ts RawTeachers) TrainXY(f feature.SingleBattleFunc, qRatio float64) (tensor.D2, tensor.D2) {
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
	r := omwrand.NewMt19937()
	xn := 90
	u1n := 64
	u2n := 16
	yn := 1
	affine, variable := model1d.NewStandardAffine(xn, u1n, u2n, yn, 0.0001, 64.0, r)
	gm := game.New(dmgtools.RandBonuses{1.0}, r)

	leafNodeJointEvalFunc := func(battle *single.Battle) (duct.LeafNodeJointEvalY, error) {
		if isEnd, gameRetJointVal := single.IsEnd(battle); isEnd {
			y := make(duct.LeafNodeJointEvalY, len(gameRetJointVal))
			for i, v := range gameRetJointVal {
				y[i] = v
			}
			return y, nil
		} else {
			x := feature.NewSingleBattleFunc(2, feature.ExpectedDamageRatioToCurrentHP, feature.DPSRatioToCurrentHP)(battle)
			y, err := affine.Predict(x)
			v := y[0]
			return duct.LeafNodeJointEvalY{v, 1.0-v}, err
		}
	}

	mctSearch := mcts.New(dmgtools.RandBonuses{1.0}, r)
	mctSearch.LeafNodeJointEvalFunc = leafNodeJointEvalFunc

	batchSize := 2560
	trainX := make(tensor.D2, 0, batchSize)
	trainY := make(tensor.D2, 0, batchSize)

	selfBattleNum := 1960
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

		if isEnd, _ := single.IsEnd(&battle); isEnd {
			continue
		}

		gm.Player = mctSearch.NewPlayer(512, r)
		endBattle, battleHistory, _, qHistory, err := gm.PlayoutWithHistory(battle, 128)
		if err != nil {
			panic(err)
		}

		_, gameRetJointVal := single.IsEnd(&endBattle)
		oneGameTrainX, oneGameTrainY := NewRawTeachers(battleHistory, qHistory, gameRetJointVal).DataAugmentation().TrainXY(feature.NewSingleBattleFunc(2, feature.ExpectedDamageRatioToCurrentHP, feature.DPSRatioToCurrentHP), 0.5)
		trainX = append(trainX, oneGameTrainX...)
		trainY = append(trainY, oneGameTrainY...)

		if len(trainX) >= batchSize - 100 {
			fmt.Println("minibatch")
			for j := 0; j < batchSize; j++ {
				idx := r.Intn(len(trainX))
				affine.SGD(trainX[idx], trainY[idx], 0.01)
			}
			trainX = make(tensor.D2, 0, batchSize)
			trainY = make(tensor.D2, 0, batchSize)
		}
		fmt.Println("i = ", i)
	}
	variable.Param.WriteJSON("c:/Users/kuroko/Desktop/test.json")
}