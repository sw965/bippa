package main

import (
	"github.com/sw965/bippa"
	"fmt"
	"math/rand"
	"github.com/seehuhn/mt19937"
	"time"
)

func main() {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())

	pbks := map[bippa.PokeName]*bippa.PokemonBuildCommonKnowledge{}
	for _, pokeName := range bippa.ALL_POKE_NAMES {
		pbk, err := bippa.LoadJsonPokemonBuildCommonKnowledge(pokeName)
		if err != nil {
			fmt.Println(pokeName, "は pbk を 作れなかった")
		}
		pbks[pokeName] = &pbk
	}
	for _, pokeName := range bippa.ALL_POKE_NAMES {
		pscs := bippa.NewMoveNameCombinations(pokeName)
		pscms := bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		pscms.WriteJson(pokeName, "move_name.json")

		pscs = bippa.NewGenderCombinations(pokeName)
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		pscms.WriteJson(pokeName, "gender.json")

		pscs = bippa.NewAbilityCombinations(pokeName)
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		pscms.WriteJson(pokeName, "ability.json")

		pscs = bippa.NewItemCombinations()
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		pscms.WriteJson(pokeName, "item.json")

		pscs = bippa.NewNatureCombinations()
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		pscms.WriteJson(pokeName, "nature.json")

		for _, key := range []string{"HP", "Atk", "Def", "SpAtk", "SpDef", "Speed"} {
			pscs = bippa.NewIndividualCombinations(key)
			pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
			pscms.WriteJson(pokeName, key + "_individual.json")
	
			pscs = bippa.NewEffortCombinations(key)
			pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
			pscms.WriteJson(pokeName, key + "_effort.json")
		}
	}
}