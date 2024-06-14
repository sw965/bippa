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
	"github.com/sw965/bippa/battle/single/mcts"
	"bufio"
	"os"
)

func main() {
	r := omwrand.NewMt19937()
	xn := 90
	u1n := 64
	u2n := 16
	yn := 1
	affine, variable := model1d.NewStandardAffine(xn, u1n, u2n, yn, 0.0001, 64.0, r)
	param, err := model1d.LoadParamJSON("C:/Users/kuroko/Desktop/test.json")
	if err != nil {
		panic(err)
	}
	variable.SetParam(param)

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

	initBattle := single.Battle{
		P1Fighters:single.Fighters{bp.NewTemplateBulbasaur(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
		P2Fighters:single.Fighters{bp.NewTemplateBulbasaur(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
	}

	gm.Player = func(battle *single.Battle) (single.Actions, []float64, error) {
		fmt.Println("p1", battle.P1Fighters.Names().ToStrings(), battle.P1Fighters[0].CurrentHP)
		fmt.Println("p2", battle.P2Fighters.Names().ToStrings(), battle.P2Fighters[0].CurrentHP)

		var action single.Action
		for {
			s := bufio.NewScanner(os.Stdin)
			s.Scan()
			text := s.Text()
			if moveName, ok := bp.STRING_TO_MOVE_NAME[text]; ok {
				action = single.Action{CmdMoveName:moveName, IsPlayer1:false}
				break
			} else if pokeName, ok := bp.STRING_TO_POKE_NAME[text]; ok {
				action = single.Action{SwitchPokeName:pokeName, IsPlayer1:false}
				break
			}
		}

		jointAction, jointQ, err := mctSearch.NewPlayer(128, r)(battle)
		if err != nil {
			return single.Actions{}, []float64{}, err
		}

		fmt.Println(jointAction[0].CmdMoveName.ToString(), jointAction[0].SwitchPokeName.ToString(), jointAction[0].IsPlayer1)
		fmt.Println(jointAction[1].CmdMoveName.ToString(), jointAction[1].SwitchPokeName.ToString(), jointAction[1].IsPlayer1)
		fmt.Println("jointQ", jointQ)
		fmt.Println("")

		return single.Actions{jointAction[0], action}, []float64{}, nil
	}

	endBattle, _, jointActionHistory, _, err := gm.PlayoutWithHistory(initBattle, 25600)
	if err != nil {
		panic(err)
	}
	isEnd, _ := single.IsEnd(&endBattle)
	fmt.Println("isEnd", isEnd)

	for i, jointAction := range jointActionHistory {
		fmt.Println(i, jointAction[0].CmdMoveName.ToString(), jointAction[0].SwitchPokeName.ToString())
		fmt.Println(i, jointAction[1].CmdMoveName.ToString(), jointAction[1].SwitchPokeName.ToString())
		fmt.Println("")
	}
}