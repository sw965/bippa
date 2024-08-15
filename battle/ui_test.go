package battle_test

import (
	"testing"
	"github.com/sw965/bippa/battle/game"
	"github.com/sw965/bippa/battle"
	bp "github.com/sw965/bippa"
)

func Test(t *testing.T) {
	m := battle.Manager{
		GuestHumanTitle:"ポケモントレーナー",
		GuestHumanName:"メイ",

		CurrentSelfLeadPokemons:bp.Pokemons{
			bp.NewKusanagi2009Toxicroak(),
			bp.NewKusanagi2009Empoleon(),
		},

		CurrentSelfBenchPokemons:bp.Pokemons{
			bp.NewKusanagiSalamence2009(),
			bp.NewKusanagi2009Snorlax(),
		},

		CurrentOpponentLeadPokemons:bp.Pokemons{
			bp.NewMoruhu2007Bronzong(),
			bp.NewMoruhu2007Smeargle(),
		},

		CurrentOpponentBenchPokemons:bp.Pokemons{
			bp.NewMoruhu2007Metagross(),
			bp.NewMoruhu2007Snorlax(),
		},
	}

	m.Init()
	p2View := m.Clone()
	p2View.SwapView()

	ui := battle.ObserverUI{
		LastP1ViewManager:m,
		LastP2ViewManager:p2View,
	}

	battle.GlobalContext.Observer = ui.Observer

	actions := battle.Actions{
		battle.Action{
			battle.SoloAction{
				MoveName:bp.FOLLOW_ME,
				SrcIndex:1,
				TargetIndex:1,
				IsSelfLeadTarget:true,
				Speed:m.CurrentOpponentLeadPokemons[1].Stat.Speed,
				IsSelf:false,
			},
			
			battle.SoloAction{
				MoveName:bp.TRICK_ROOM,
				SrcIndex:0,
				Speed:m.CurrentOpponentLeadPokemons[0].Stat.Speed,
				IsSelf:false,
			},
		},

		battle.Action{
			battle.SoloAction{
				MoveName:bp.FAKE_OUT,
				SrcIndex:0,
				TargetIndex:0,
				Speed:m.CurrentOpponentLeadPokemons[0].Stat.Speed,
				IsSelf:true,
			},

			// battle.SoloAction{
			// 	SrcIndex:0,
			// 	TargetIndex:0,
			// 	Speed:m.CurrentSelfLeadPokemons[0].Stat.Speed,
			// 	IsSelf:true,
			// },

			battle.SoloAction{
				MoveName:bp.SURF,
				SrcIndex:1,
				Speed:m.CurrentSelfLeadPokemons[1].Stat.Speed,
				IsSelf:true,
			},
		},
	}

	_, err := game.Push(m, actions)
	if err != nil {
		panic(err)
	}
}