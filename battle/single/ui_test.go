package single_test

import (
	"testing"
	omwrand "github.com/sw965/omw/math/rand"
	"github.com/sw965/bippa/battle/single"
	bp "github.com/sw965/bippa"
)

func Test(t *testing.T) {
	battle := single.Battle{
		SelfLeadPokemons:bp.Pokemons{
			bp.NewRomanStan2009Metagross(),
			bp.NewRomanStan2009Latios(),
		},

		OpponentLeadPokemons:bp.Pokemons{
			bp.NewKusanagi2009Toxicroak(),
			bp.NewKusanagi2009Empoleon(),
		},
	}

	action := single.SoloAction{
		MoveName:bp.COMET_PUNCH,
		SrcIndex:0,
		TargetIndex:0,
		IsOpponentLeadTarget:true,
	}

	r := omwrand.NewMt19937()
	context := single.NewContext(r)
	battle.MoveUse(&action, &context)
}