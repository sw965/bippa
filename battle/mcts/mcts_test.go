package mcts_test

import (
	"testing"
	bp "github.com/sw965/bippa"
	omwrand "github.com/sw965/omw/math/rand"
	"github.com/sw965/bippa/battle"
	"github.com/sw965/bippa/battle/mcts"
)

func Test(t *testing.T) {
	m := battle.Manager{
		CurrentSelfLeadPokemons:bp.Pokemons{
			bp.NewKusanagi2009Toxicroak(),
			bp.NewKusanagi2009Empoleon(),
		},

		CurrentSelfBenchPokemons:bp.Pokemons{
			bp.NewKusanagi2009Snorlax(),
			bp.NewKusanagiSalamence2009(),
		},

		CurrentOpponentLeadPokemons:bp.Pokemons{
			bp.NewMoruhu2007Bronzong(),
			bp.NewMoruhu2007Smeargle(),
		},

		CurrentOpponentBenchPokemons:bp.Pokemons{
			bp.NewMoruhu2007Snorlax(),
			bp.NewMoruhu2007Metagross(),
		},
	}
	m.Init("四天王", "カトレア")
	m.CurrentSelfLeadPokemons[0].Stat.CurrentHP = 0

	r := omwrand.NewMt19937()

	mctSearch := mcts.New()
	mctSearch.SetRandomPlayoutLeafNodeJointEvalFunc(r)
	//NewNodePointerという命名に変更すべき？
	rootNode := mctSearch.NewNode(&m)
	err := mctSearch.Run(128, rootNode, r)
	if err != nil {
		panic(err)
	}
}