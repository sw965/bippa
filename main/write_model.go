package main

import (
	"github.com/sw965/bippa"
)

func main() {
	pbks := map[PokeName]*PokemonBuildCommonKnowledge{}

	for pokeName := range ALL_POKE_NAMES {
		pbk, err := LoadJsonPokemonBuildCommonKnowledge(pokeName)
		if err != nil {
			fmt.Println(pokeName, " は pbk を 作れなかった")
		}
		pbks[pokeName] = LoadJsonPokemonBuildCommonKnowledge()
	}
}