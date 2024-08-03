package battle_test

import (
	"testing"
	"github.com/sw965/bippa/battle/game"
	omwrand "github.com/sw965/omw/math/rand"
	"github.com/sw965/bippa/battle"
	bp "github.com/sw965/bippa"
)

func Test(t *testing.T) {
	m := battle.Manager{
		SelfLeadPokemons:bp.Pokemons{
			bp.NewKusanagi2009Toxicroak(),
			bp.NewKusanagi2009Empoleon(),
		},

		SelfBenchPokemons:bp.Pokemons{
			bp.NewKusanagiSalamence2009(),
			bp.NewKusanagi2009Snorlax(),
		},

		OpponentLeadPokemons:bp.Pokemons{
			bp.NewMoruhu2007Bronzong(),
			bp.NewMoruhu2007Smeargle(),
		},

		OpponentBenchPokemons:bp.Pokemons{
			bp.NewMoruhu2007Metagross(),
			bp.NewMoruhu2007Snorlax(),
		},
		IsPlayer1View:true,
	}
	p2View := m.Clone()
	p2View.SwapView()

	r := omwrand.NewMt19937()
	context := battle.NewContext(r)
	ui := battle.ObserverUI{
		LastP1ViewManager:m,
		LastP2ViewManager:p2View,
	}
	context.Observer = ui.Observer
	push := game.NewPushFunc(&context)

	actions := battle.Actions{
		battle.Action{
			battle.SoloAction{
				MoveName:bp.FOLLOW_ME,
				SrcIndex:1,
				TargetIndex:1,
				IsSelfLeadTarget:true,
				Speed:m.SelfLeadPokemons[1].Stat.Speed,
				IsSelfView:false,
			},
			
			battle.SoloAction{
				MoveName:bp.TRICK_ROOM,
				SrcIndex:0,
				Speed:m.SelfLeadPokemons[0].Stat.Speed,
				IsSelfView:false,
			},
		},

		battle.Action{
			battle.SoloAction{
				MoveName:bp.FAKE_OUT,
				SrcIndex:0,
				TargetIndex:0,
				Speed:m.OpponentLeadPokemons[0].Stat.Speed,
				IsSelfView:true,
			},

			battle.SoloAction{
				MoveName:bp.SURF,
				SrcIndex:1,
				Speed:m.OpponentLeadPokemons[1].Stat.Speed,
				IsSelfView:true,
			},
		},
	}

	_, err := push(m, actions)
	if err != nil {
		panic(err)
	}
}