package battle_test

import (
	"testing"
	"fmt"
	bp "github.com/sw965/bippa"
	"github.com/sw965/bippa/battle"
)

func soloActionPrintHelper(a *battle.SoloAction, m *battle.Manager) {
	src := m.CurrentSelfLeadPokemons[a.SrcIndex]
	if a.IsMove() {
		msg := src.Name.ToString()
		msg += "の " + a.MoveName.ToString()
		if a.TargetIndex != -1 {
			msg += " 対象 →【"
			if a.IsSelfLeadTarget {
				target := m.CurrentSelfLeadPokemons[a.TargetIndex]
				msg += "自分：" + target.Name.ToString() + "】"
			} else {
				target := m.CurrentOpponentLeadPokemons[a.TargetIndex]
				msg += "相手：" + target.Name.ToString() + "】"
			}
		}
		fmt.Println(msg)
		return
	}

	target := m.CurrentSelfBenchPokemons[a.TargetIndex]
	fmt.Println(src.Name.ToString() + " から " + target.Name.ToString() + " へ 交代")
}

func TestNewLegalActions(t *testing.T) {
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

	result := battle.NewLegalActions(&m)
	for i, a := range result {
		fmt.Println("i = ", i)
		for _, soloA := range a {
			soloActionPrintHelper(&soloA, &m)
		}
		fmt.Println("")
	}
}