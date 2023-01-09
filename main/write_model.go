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
		err := pscms.WriteJson(pokeName, "move_name.json")
		if err != nil {
			panic(err)
		}

		pscs = bippa.NewGenderCombinations(pokeName)
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err = pscms.WriteJson(pokeName, "gender.json")
		if err != nil {
			panic(err)
		}

		pscs = bippa.NewAbilityCombinations(pokeName)
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err = pscms.WriteJson(pokeName, "ability.json")
		if err != nil {
			panic(err)
		}

		pscs = bippa.NewItemCombinations()
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err = pscms.WriteJson(pokeName, "item.json")
		if err != nil {
			panic(err)
		}

		pscs = bippa.NewNatureCombinations()
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err = pscms.WriteJson(pokeName, "nature.json")
		if err != nil {
			panic(err)
		}

		for _, key := range []string{"HP", "Atk", "Def", "SpAtk", "SpDef", "Speed"} {
			pscs = bippa.NewIndividualCombinations(key)
			pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
			err = pscms.WriteJson(pokeName, key + "_individual.json")
			if err != nil {
				panic(err)
			}

			pscs = bippa.NewEffortCombinations(key)
			pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
			err = pscms.WriteJson(pokeName, key + "_effort.json")
			if err != nil {
				panic(err)
			}
		}
	}
}