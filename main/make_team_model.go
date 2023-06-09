package main

import (
	bp "github.com/sw965/bippa"
	"fmt"
)

func main() {
	teamPokemonModel := bp.NewMoveNamesAndAbilityTeamPokemonModel(bp.POKEDEX["フシギバナ"].Learnset, 2, bp.POKEDEX["フシギバナ"].AllAbilities)
	for _, part := range teamPokemonModel {
		fmt.Println(part)
	}
}