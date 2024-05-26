package main

import (
	"fmt"
	sb "github.com/sw965/bippa/battle/single"
	//"github.com/sw965/crow/game/simultaneous"
	"github.com/sw965/bippa/dmgtools"
	orand "github.com/sw965/omw/rand"
	omath "github.com/sw965/omw/math"
	bp "github.com/sw965/bippa"
	"github.com/sw965/crow/ucb"
	"github.com/sw965/crow/model"
	"github.com/sw965/crow/tensor"
)

func main() {
	p1Fighters := sb.Fighters{
		bp.NewTemplateBulbasaur(),
		bp.NewTemplateCharmander(),
		bp.NewTemplateSquirtle(),
	}

	p2Fighters := sb.Fighters{
		bp.NewTemplateSquirtle(),
		bp.NewTemplateCharmander(),
		bp.NewTemplateBulbasaur(),
	}

	initBattle := sb.Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}

	makeFeature := func(attacker, defender *bp.Pokemon) tensor.D1 {
		calculator := dmgtools.Calculator{
			dmgtools.Attacker{
				PokeName:attacker.Name,
				Level:attacker.Level,
				Atk:attacker.Atk,
				SpAtk:attacker.SpAtk,
			},
			dmgtools.Defender{
				PokeName:defender.Name,
				Level:defender.Level,
				Def:defender.Def,
				SpDef:defender.SpDef,
			},
		}

		maxExpectationDmgRatio := 0.0
		maxKORatio := 0.0

		for moveName, pp := range attacker.Moveset {
			calculator.Attacker.MoveName = moveName
			moveData := bp.MOVEDEX[moveName]
			if pp.Current > 0 {
				dmg := calculator.Calculation(omath.Mean(dmgtools.RAND_BONUSES...))
				accuracy := moveData.Accuracy

				expectationDmgRatio := omath.Min(1.0, (float64(dmg) * float64(accuracy) / 100.0) / float64(defender.CurrentHP))
				if expectationDmgRatio > maxExpectationDmgRatio {
					maxExpectationDmgRatio = expectationDmgRatio
				}

				koRatio := maxExpectationDmgRatio * float64(accuracy) / 100.0
				if koRatio > maxKORatio {
					maxKORatio = koRatio
				}
			}
		}
		return tensor.D1{maxExpectationDmgRatio, maxKORatio}
	}

	f := sb.MakeEvalFeatureFunc(makeFeature, 2)
	fmt.Println(len(f(&initBattle)))
	r := orand.NewMt19937()
	game := sb.NewGame(dmgtools.RandBonuses{1.0}, r)
	game.SetRandActionPlayer(r)
	mcts := sb.NewMCTSRandPlayout(ucb.New1Func(5), dmgtools.RandBonuses{1.0}, r)
	nn, _ := model.NewThreeLayerAffineParamReLUInput1DOutputSigmoid1D(90, 32, 8, 1, 0.0001, 64.0, r)
	mctsNNEval := sb.NewMCTSNNEval(ucb.NewAlphaGoFunc(5), &nn, dmgtools.RandBonuses{1.0}, f, r)
	mctsGame := mcts.Game.Clone()

	batchCap := 5120
	qs := make(tensor.D2, 0, batchCap)
	mctsGame.Player = func(battle *sb.Battle) (sb.Actions, error) {
		rootNode := mctsNNEval.NewNode(battle)
		err := mcts.Run(512, rootNode, r)
		vs := rootNode.UCBManagers.AverageValues()
		qs = append(qs, tensor.D1{vs[0]})
		qs = append(qs, tensor.D1{vs[1]})
		return rootNode.UCBManagers.JointActionByMaxTrial(r), err
	}

	trainX := make(tensor.D2, 0, batchCap)
	trainY := make(tensor.D2, 0, batchCap)

	rn := mctsNNEval.NewNode(&initBattle)
	err := mctsNNEval.Run(5120, rn, r)
	if err != nil {
		panic(err)
	}

	predictedJointActions, predictedJointValues := rn.Predict(r, 8)
	for i, jointAction := range predictedJointActions {
		fmt.Println(jointAction.ToStrings(), predictedJointValues[i])
	}
	fmt.Println("")

	for i := 0; i < 12800; i++ {
		tmpBattle, err := game.Play(initBattle, r.Intn(8))
		if err != nil {
			panic(err)
		}
		if !sb.IsEnd(&tmpBattle) {
			endBattle, battleHistory, _, err := mctsGame.PlayoutWithHistory(tmpBattle, 128)
			if err != nil {
				panic(err)
			}
			ys, err := endBattle.EndLeafNodeEvalYs()
			if err != nil {
				panic(err)
			}

			for _, b := range battleHistory {
				swaped := b.SwapPlayers()
				trainX = append(trainX, f(&b))
				trainX = append(trainX, f(&swaped))
				trainY = append(trainY, tensor.D1{float64(ys[0])})
				trainY = append(trainY, tensor.D1{float64(ys[1])})
			}

			if len(trainX) > batchCap {
				m := len(trainX)
				for i := 0; i < m*2; i++ {
					idx := r.Intn(m)
					nn.SGD(trainX[idx], tensor.D1{(trainY[idx][0] * 0.5) + (qs[idx][0] * 0.5)}, 0.01)
				}

				rn = mctsNNEval.NewNode(&initBattle)
				err = mctsNNEval.Run(12800, rn, r)
				if err != nil {
					panic(err)
				}
			
				predictedJointActions, predictedJointValues = rn.Predict(r, 8)
				for i, jointAction := range predictedJointActions {
					fmt.Println(jointAction.ToStrings(), predictedJointValues[i])
				}
				fmt.Println("")

				clone := initBattle
				clone.P1Fighters[1].CurrentHP = 0
				clone.P1Fighters[2].CurrentHP = 0
				clone.P2Fighters[1].CurrentHP = 0
				clone.P2Fighters[2].CurrentHP = 0
				y, err := nn.Predict(f(&clone))
				if err != nil {
					panic(err)
				}
				fmt.Println("最後1",  y)
				
				clone.P1Fighters[0].CurrentHP = 1
				clone.P2Fighters[0].CurrentHP = 1
				y, err = nn.Predict(f(&clone))
				if err != nil {
					panic(err)
				}
				fmt.Println("最後2", y)

				clone.P1Fighters[0].Speed = 1
				y, err = nn.Predict(f(&clone))
				if err != nil {
					panic(err)
				}
				fmt.Println("最後3", y)

				clone.P2Fighters[0].Speed = 1
				y, err = nn.Predict(f(&clone))
				if err != nil {
					panic(err)
				}
				fmt.Println("最後4", y)

				clone = initBattle
				clone.P2Fighters[0] = bp.NewTemplateSuicune()
				y, err = nn.Predict(f(&clone))
				if err != nil {
					panic(err)
				}
				fmt.Println("最後5", y)

				trainX = make(tensor.D2, 0, batchCap)
				trainY = make(tensor.D2, 0, batchCap)
				qs = make(tensor.D2, 0, batchCap)
			}
		}
	}

	rn = mctsNNEval.NewNode(&initBattle)
	err = mctsNNEval.Run(12800, rn, r)
	if err != nil {
		panic(err)
	}

	predictedJointActions, predictedJointValues = rn.Predict(r, 8)
	for i, jointAction := range predictedJointActions {
		fmt.Println(jointAction.ToStrings(), predictedJointValues[i])
	}
	fmt.Println("")
}