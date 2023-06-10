package main

import (
	bp "github.com/sw965/bippa"
	//"fmt"
	"github.com/sw965/omw"
)

func main() {
	r := omw.NewMt19937()
	pokeData := bp.POKEDEX["フシギバナ"]
	teamPokemonModel := bp.NewMoveNamesAndAbilityTeamPokemonModel(pokeData.Learnset, 2, pokeData.AllAbilities)
	teamPokemonModel.Init(r)
	err := teamPokemonModel.Write("フシギバナ", "move_name_and_ability.json", false)
	if err != nil {
		panic(err)
	}
}
