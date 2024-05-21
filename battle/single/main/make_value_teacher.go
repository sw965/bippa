package main

import (
	"fmt"
	sb "github.com/sw965/bippa/battle/single"
	//"github.com/sw965/crow/game/simultaneous"
	"github.com/sw965/bippa/dmgtools"
	orand "github.com/sw965/omw/rand"
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

	r := orand.NewMt19937()
	game := sb.NewGame(dmgtools.RandBonuses{1.0}, r)
	game.SetRandActionPlayer(r)
	mcts := sb.NewMCTSRandPlayout(ucb.New1Func(5), dmgtools.RandBonuses{1.0}, r)
	nn, _ := model.NewThreeLayerAffineLeakyReLUInput1DOutputSigmoid1D(len(bp.ALL_TYPES) * 4 * sb.FIGHTER_NUM * 2, 64, 16, 1, 0.0001, 64.0, r)
	mctsNNEval := sb.NewMCTSNNEval(ucb.NewAlphaGoFunc(5), &nn, dmgtools.RandBonuses{1.0}, r)
	mctsGame := mcts.Game.Clone()
	mctsGame.Player = func(battle *sb.Battle) (sb.Actions, error) {
		rootNode := mctsNNEval.NewNode(battle)
		err := mcts.Run(512, rootNode, r)
		return rootNode.UCBManagers.JointActionByMaxTrial(r), err
	}

	batchCap := 5120
	trainX := make(tensor.D2, 0, batchCap)
	trainY := make(tensor.D2, 0, batchCap)

	rn := mctsNNEval.NewNode(&initBattle)
	err := mctsNNEval.Run(5120, rn, r)
	if err != nil {
		panic(err)
	}

	for _, jointAction := range rn.MaxTrialJointActionPath(r, 8) {
		fmt.Println(jointAction.ToStrings())
	}
	fmt.Println("")

	for i := 0; i < 1280; i++ {
		tmpBattle, err := game.Play(initBattle, r.Intn(16))
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
				trainX = append(trainX, b.ToIndexFeature(15000))
				trainX = append(trainX, swaped.ToIndexFeature(15000))
				trainY = append(trainY, tensor.D1{float64(ys[0])})
				trainY = append(trainY, tensor.D1{float64(ys[1])})
			}

			if len(trainX) > batchCap {
				m := len(trainX)
				for i := 0; i < m; i++ {
					idx := r.Intn(m)
					nn.SGD(trainX[idx], trainY[idx], 0.01)
				}

				rn = mctsNNEval.NewNode(&initBattle)
				err = mctsNNEval.Run(5120, rn, r)
				if err != nil {
					panic(err)
				}
			
				for _, jointAction := range rn.MaxTrialJointActionPath(r, 8) {
					fmt.Println(jointAction.ToStrings())
				}
				fmt.Println("")

				trainX = make(tensor.D2, 0, batchCap)
				trainY = make(tensor.D2, 0, batchCap)
			}
		}
	}

	rn = mctsNNEval.NewNode(&initBattle)
	err = mctsNNEval.Run(5120, rn, r)
	if err != nil {
		panic(err)
	}

	for _, jointAction := range rn.MaxTrialJointActionPath(r, 8) {
		fmt.Println(jointAction.ToStrings())
	}
}