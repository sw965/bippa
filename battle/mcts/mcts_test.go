package mcts_test

import (
	"fmt"
	bp "github.com/sw965/bippa"
	"github.com/sw965/bippa/battle"
	"github.com/sw965/bippa/battle/mcts"
	omwrand "github.com/sw965/omw/math/rand"
	"testing"
)

func Test(t *testing.T) {
	m := battle.Manager{
		CurrentSelfLeadPokemons: bp.Pokemons{
			bp.NewKusanagi2009Toxicroak(),
			bp.NewKusanagi2009Empoleon(),
		},

		CurrentSelfBenchPokemons: bp.Pokemons{
			bp.NewKusanagi2009Snorlax(),
			bp.NewKusanagiSalamence2009(),
		},

		CurrentOpponentLeadPokemons: bp.Pokemons{
			bp.NewMoruhu2007Bronzong(),
			bp.NewMoruhu2007Smeargle(),
		},

		CurrentOpponentBenchPokemons: bp.Pokemons{
			bp.NewMoruhu2007Snorlax(),
			bp.NewMoruhu2007Metagross(),
		},
	}
	m.Init("四天王", "カトレア")
	// p2View := m.Clone()
	// p2View.SwapView()

	// ui := battle.ObserverUI{
	// 	LastP1ViewManager:m,
	// 	LastP2ViewManager:p2View,
	// }
	// battle.GlobalContext.Observer = ui.Observer
	battle.GlobalContext.DamageRandBonuses = []float64{1.0}
	r := omwrand.NewMt19937()

	mctSearch := mcts.New(5.0)

	//NewNodePointerという命名に変更すべき？
	rootNode, err := mctSearch.NewNode(&m)
	if err != nil {
		panic(err)
	}

	err = mctSearch.Run(25600, rootNode, r)
	if err != nil {
		panic(err)
	}

	jointAction := rootNode.SeparateUCBManager.JointActionByMaxTrial(r)
	for _, action := range jointAction {
		for _, soloAction := range action {
			fmt.Println(soloAction.MoveName.ToString())
		}
	}
	fmt.Println(jointAction)
	fmt.Println(rootNode.SeparateUCBManager.JointAverageValue())
}
