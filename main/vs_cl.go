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
	battlemsg "github.com/sw965/bippa/battle/msg"
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

	gm := game.New(r)

	mctSearch := mcts.New(r)
	mctSearch.Game.SetRandActionPlayer(r)

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

	mctSearch.LeafNodeJointEvalFunc = leafNodeJointEvalFunc

	var lastMoveset bp.Moveset
	var lastUsedMoveName bp.MoveName
	var lastLeadPokeName bp.PokeName
	var lastCurrentHP int

	cui := func(battle *single.Battle, step single.Step) {
		switch step {
			case single.BEFORE_SWITCH_STEP:
				lastLeadPokeName = battle.SelfFighters[0].Name
			case single.AFTER_SWITCH_STEP:
				for _, m := range battlemsg.NewBack("", lastLeadPokeName, battle.IsRealSelf).Accumulate() {
					fmt.Println(m)
				}

				for _, m := range battlemsg.NewGo("", battle.SelfFighters[0].Name, battle.IsRealSelf).Accumulate() {
					fmt.Println(m)
				}
			case single.BEFORE_MOVE_USE_STEP:
				lastMoveset = battle.SelfFighters[0].Moveset.Clone()
			case single.AFTER_MOVE_USE_STEP:
				currentMoveset := battle.SelfFighters[0].Moveset
				for moveName, pp := range currentMoveset {
					if pp.Current < lastMoveset[moveName].Current {
						lastUsedMoveName = moveName
						break
					}
				}
				for _, m := range battlemsg.NewMoveUse(battle.SelfFighters[0].Name, lastUsedMoveName, battle.IsRealSelf).Accumulate() {
					fmt.Println(m)
				}
			case single.BEFORE_MOVE_DAMAGE_STEP:
				lastCurrentHP = battle.OpponentFighters[0].CurrentHP
			case single.AFTER_MOVE_DAMAGE_STEP:
				dmg := lastCurrentHP - battle.OpponentFighters[0].CurrentHP
				fmt.Println("dmg", dmg)
			case single.SELF_FAINT_STEP:
				fmt.Println(battle.IsRealSelf, "real koko")
				for _, m := range battlemsg.NewFaint(battle.SelfFighters[0].Name, battle.IsRealSelf).Accumulate() {
					fmt.Println(m)
				}
			case single.OPPONENT_FAINT_STEP:
				fmt.Println(battle.IsRealSelf, "real")
				for _, m := range battlemsg.NewFaint(battle.OpponentFighters[0].Name, battle.IsRealSelf).Accumulate() {
					fmt.Println(m)
				}
		}
	}

	initBattle := single.Battle{
		SelfFighters:single.Fighters{bp.NewTemplateBulbasaur(), bp.NewTemplateSuicune(), bp.NewTemplateSquirtle()},
		OpponentFighters:single.Fighters{bp.NewTemplateGarchomp(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
		IsRealSelf:true,
		RandDmgBonuses:dmgtools.RandBonuses{1.0},
		Observer:cui,
	}

	gm.Player = func(battle *single.Battle) (single.Actions, []float64, error) {
		fmt.Println(
			"p1", battle.SelfFighters.Names().ToStrings(),
			battle.SelfFighters[0].CurrentHP,
			battle.SelfFighters[1].CurrentHP,
			battle.SelfFighters[2].CurrentHP,
		)
		fmt.Println(
			"p2", battle.OpponentFighters.Names().ToStrings(),
			battle.OpponentFighters[0].CurrentHP,
			battle.OpponentFighters[1].CurrentHP,
			battle.OpponentFighters[2].CurrentHP,
		)

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

		battle.Observer = func(_ *single.Battle, _ single.Step) {}
		jointAction, jointQ, err := mctSearch.NewPlayer(51200, r)(battle)
		if err != nil {
			return single.Actions{}, []float64{}, err
		}
		battle.Observer = cui

		fmt.Println(jointAction[0].CmdMoveName.ToString(), jointAction[0].SwitchPokeName.ToString(), jointAction[0].IsPlayer1)
		fmt.Println(jointAction[1].CmdMoveName.ToString(), jointAction[1].SwitchPokeName.ToString(), jointAction[1].IsPlayer1)
		fmt.Println("jointQ", jointQ)
		fmt.Println("")

		return single.Actions{jointAction[0], action}, []float64{}, nil
	}

	endBattle, _, _, _, err := gm.PlayoutWithHistory(initBattle, 128)
	if err != nil {
		panic(err)
	}
	isEnd, gameRetJointVal := single.IsEnd(&endBattle)
	fmt.Println("isEnd", isEnd, gameRetJointVal)
}