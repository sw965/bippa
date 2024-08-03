package game_test

import (
	"testing"
	"fmt"
	"github.com/sw965/bippa/battle/single"
	"github.com/sw965/bippa/battle/single/game"
	bp "github.com/sw965/bippa"
)

func Test(t *testing.T) {
	battle := single.Battle{
		SelfLeadPokemons:bp.Pokemons{
			bp.NewKusanagi2009Toxicroak(),
			bp.NewKusanagi2009Empoleon(),
		},

		OpponentLeadPokemons:bp.Pokemons{
			bp.NewMoruhu2008Metagross(),
			bp.NewRomanStan2009Latios(),
		},
	}

	fmt.Println(battle.SelfLeadPokemons[0].UsableMoveNames())
	// result := battle.SelfLegalMoveSoloActions()
	// for _, action := range result {
	// 	fmt.Println(action.SrcPokemon(&battle).Name.ToString() + " の " + action.MoveName.ToString() + " : " + action.TargetPokemon(&battle).Name.ToString())
	// }

	result := game.LegalSeparateActionsSlice(&battle)
	fmt.Println(result)
	for _, action := range result[0] {
		fmt.Println(action[0].SrcPokemon(&battle).Name.ToString() + " の " + action[0].MoveName.ToString() + " : " + action[0].TargetPokemon(&battle).Name.ToString())
		fmt.Println(action[1].SrcPokemon(&battle).Name.ToString() + " の " + action[1].MoveName.ToString() + " : " + action[1].TargetPokemon(&battle).Name.ToString())
		fmt.Println("")
	}
}