package single_test

import (
	"fmt"
	"testing"
	bp "github.com/sw965/bippa"
	"github.com/sw965/bippa/battle/single"
	omwrand "github.com/sw965/omw/math/rand"
)

func TestTargetLeadPokemonsIndicesSingle(t *testing.T) {
	r := omwrand.NewMt19937()
	context := single.NewContext(r)

	battle := single.Battle{
		SelfLeadPokemons:bp.Pokemons{
			bp.NewKusanagi2009Toxicroak(),
		},

		OpponentLeadPokemons:bp.Pokemons{
			bp.NewRomanStan2009Metagross(),
			bp.NewRomanStan2009Latios(),
		},
	}

	action := single.SoloAction{
		MoveName:bp.SURF,
		SrcIndex:1,
		MoveTarget:bp.OTHERS_TARGET,
	}

	fmt.Println(battle.TargetPokemonsIndices(&action, &context))
	
	battle.OpponentLeadPokemons[0].CurrentHP = 0
	
	action = single.SoloAction{
		MoveName:bp.HYDRO_PUMP,
		SrcIndex:1,
		MoveTarget:bp.NORMAL_TARGET,
	}

	fmt.Println(battle.TargetPokemonsIndices(&action, &context))
}