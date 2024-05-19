package single_test

import (
	"testing"
	"fmt"
	"github.com/sw965/crow/ucb"
	bp "github.com/sw965/bippa"
	sb "github.com/sw965/bippa/battle/single"
	orand "github.com/sw965/omw/rand"
	"github.com/sw965/bippa/dmgtools"
)

func Test(t *testing.T) {
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

	r := orand.NewMt19937()
	battle := sb.Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters, RandDamageBonuses:dmgtools.RandBonuses{1.0}}
	push := sb.Push(r)
	battle, err := push(
		battle,
		sb.Actions{
			sb.Action{CmdMoveName:bp.TACKLE, IsPlayer1:true},
			sb.Action{CmdMoveName:bp.WATER_GUN, IsPlayer1:false},
		},
	)
	if err != nil {
		panic(err)
	}
	battle.P1Fighters[0].CurrentHP = 0
	fmt.Println(battle.P1Fighters[0].CurrentHP, battle.P2Fighters[0].CurrentHP)
}

func Test2(t *testing.T) {
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

	p1Fighters[0].CurrentHP = 0
	battle := sb.Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters, RandDamageBonuses:dmgtools.RandBonuses{1.0}}
	legalActionss := sb.LegalActionss(&battle)
	fmt.Println(legalActionss)
}

// func Test2(t *testing.T) {
// 	p1Fighters := sb.Fighters{
// 		bp.NewTemplateBulbasaur(),
// 		bp.NewTemplateCharmander(),
// 		bp.NewTemplateSquirtle(),
// 	}

// 	p2Fighters := sb.Fighters{
// 		bp.NewTemplateSquirtle(),
// 		bp.NewTemplateCharmander(),
// 		bp.NewTemplateBulbasaur(),
// 	}
// 	initBattle := sb.Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}
// 	r := omw.NewMt19937()
// 	mcts := sb.NewMCTS(r)
// 	testNum := 196
// 	for i, pokemon := range initBattle.P1Fighters {
// 		fmt.Println(bp.POKE_NAME_TO_STRING[pokemon.Name])
// 		fmt.Println(pokemon.CurrentHP)
// 		fmt.Println(bp.POKE_NAME_TO_STRING[initBattle.P2Fighters[i].Name])
// 		fmt.Println(initBattle.P2Fighters[i].CurrentHP)
// 	}

// 	for i := 0; i < testNum; i++ {
// 		_, err := mcts.Game.Playout(initBattle)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 	fmt.Println("ok")
// 	for i, pokemon := range initBattle.P1Fighters {
// 		fmt.Println(bp.POKE_NAME_TO_STRING[pokemon.Name])
// 		fmt.Println(pokemon.CurrentHP)
// 		fmt.Println(bp.POKE_NAME_TO_STRING[initBattle.P2Fighters[i].Name])
// 		fmt.Println(initBattle.P2Fighters[i].CurrentHP)
// 	}
// }

//turnが正常かテスト
//訪問回数が正常化テスト
//参照透過のテスト


func TestMCTS(t *testing.T) {
	r := orand.NewMt19937()
	mcts := sb.NewMCTS(r)
	mcts.UCBFunc = ucb.NewAlphaGoFunc(5)
	//mcts.UCBFunc = ucb.New1Func(math.Sqrt(25))

	p1Fighters := sb.Fighters{
		bp.NewTemplateBulbasaur(),
		bp.NewTemplateCharmander(),
		bp.NewTemplateSquirtle(),
	}

	// for i := range p1Fighters {
	// 	p1Fighters[i].CurrentHP = 1
	// }
	//p1Fighters[1].CurrentHP = 0
	//p1Fighters[2].CurrentHP = 0

	p2Fighters := sb.Fighters{
		bp.NewTemplateSquirtle(),
		bp.NewTemplateCharmander(),
		bp.NewTemplateBulbasaur(),
	}

	battle := sb.Battle{
		P1Fighters:p1Fighters,
		P2Fighters:p2Fighters,
		Actions:sb.Actions{sb.Action{IsPlayer1:true}, sb.Action{IsPlayer1:false}},
		RandDamageBonuses:dmgtools.RandBonuses{1.0},
	}

	rootNode := mcts.NewNode(&battle)
	simulation := 1600

	for i := 0; i < simulation; i++ {
		err := mcts.Run(1, rootNode, r)
		if err != nil {
			panic(err)
		}
	}

	for _, jointAction := range rootNode.MaxTrialJointActionPath(r, 16) {
		fmt.Println(
			bp.MOVE_NAME_TO_STRING[jointAction[0].CmdMoveName],
			bp.POKE_NAME_TO_STRING[jointAction[0].SwitchPokeName],
			bp.MOVE_NAME_TO_STRING[jointAction[1].CmdMoveName],
			bp.POKE_NAME_TO_STRING[jointAction[1].SwitchPokeName],
		)
	}

	for _, nextNode := range rootNode.NextNodes {
		if len(nextNode.LastJointActions) != 1 {
			fmt.Println("無理っす")
			return
		}
		stateJointAction := nextNode.State.Actions
		jointAction := nextNode.LastJointActions[0]

		if jointAction[0].CmdMoveName != stateJointAction[0].CmdMoveName {
			fmt.Println("無理")
			return
		}

		if jointAction[0].SwitchPokeName != stateJointAction[0].SwitchPokeName {
			fmt.Println("無理")
			return
		}

		if jointAction[1].CmdMoveName != stateJointAction[1].CmdMoveName {
			fmt.Println("無理")
			return
		}

		if jointAction[1].SwitchPokeName != stateJointAction[1].SwitchPokeName {
			fmt.Println("無理")
			return
		}


		fmt.Println(
			"state",
			bp.POKE_NAME_TO_STRING[nextNode.State.P1Fighters[0].Name],
			bp.POKE_NAME_TO_STRING[nextNode.State.P2Fighters[0].Name],
			"action",
			bp.MOVE_NAME_TO_STRING[jointAction[0].CmdMoveName],
			bp.MOVE_NAME_TO_STRING[stateJointAction[0].CmdMoveName],
			bp.POKE_NAME_TO_STRING[jointAction[0].SwitchPokeName],
			bp.POKE_NAME_TO_STRING[stateJointAction[0].SwitchPokeName],

			bp.MOVE_NAME_TO_STRING[jointAction[1].CmdMoveName],
			bp.MOVE_NAME_TO_STRING[stateJointAction[1].CmdMoveName],
			bp.POKE_NAME_TO_STRING[jointAction[1].SwitchPokeName],
			bp.POKE_NAME_TO_STRING[stateJointAction[1].SwitchPokeName],
		)
	}
}

func TestActionData(t *testing.T) {
	r := orand.NewMt19937()
	p1Fighters := sb.Fighters{
		bp.NewTemplateBulbasaur(),
		bp.NewTemplateCharmander(),
		bp.NewTemplateSquirtle(),
	}

	p2Fighters := sb.Fighters{
		bp.NewTemplateBulbasaur(),
		bp.NewTemplateSquirtle(),
		bp.NewTemplateCharmander(),
	}

	// p1Fighters[0].CurrentHP = 1
	// p1Fighters[1].CurrentHP = 0
	// p1Fighters[2].CurrentHP = 0

	// p2Fighters[0].CurrentHP = 1
	// p2Fighters[1].CurrentHP = 0
	// p2Fighters[2].CurrentHP = 0

	battle := sb.Battle{
		P1Fighters:p1Fighters,
		P2Fighters:p2Fighters,
		Actions:sb.Actions{sb.Action{IsPlayer1:true}, sb.Action{IsPlayer1:false}},
		RandDamageBonuses:dmgtools.RandBonuses{1.0},
	}

	n := 64
	outputs := make([]map[sb.Action]float64, n) 
	battles := make([]sb.Battle, n)
	values := make([]float64, n)
	maxvs := make([]float64, n)
	avgvs := make([]float64, n)
	game := sb.NewMCTSPlayerGame(battles, outputs, values, maxvs, avgvs, 51200, r)
	_, err := game.Playout(battle)
	if err != nil {
		panic(err)
	}
	for i := range battles {
		b := battles[i]
		fmt.Println(
			bp.POKE_NAME_TO_STRING[b.P1Fighters[0].Name],
			b.P1Fighters[0].CurrentHP,
			bp.POKE_NAME_TO_STRING[b.P2Fighters[0].Name],
			b.P2Fighters[0].CurrentHP,
		)
		fmt.Println(outputs[i])
		fmt.Println("v=", values[i])
		fmt.Println("maxv=", maxvs[i])
		fmt.Println("avgv=", avgvs[i])
	}
}