package main

import (
	"github.com/sw965/bippa"
	"fmt"
	"math/rand"
	"github.com/seehuhn/mt19937"
	"time"
	"strings"
)

func main() {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())

	pbks := map[bippa.PokeName]*bippa.PokemonBuildCommonKnowledge{}
	pbkPokeNames := make(bippa.PokeNames, 0, len(bippa.ALL_POKE_NAMES))

	for _, pokeName := range bippa.ALL_POKE_NAMES {
		pbk, err := bippa.LoadJsonPokemonBuildCommonKnowledge(pokeName)
		if err != nil {
			fmt.Println(pokeName, "は pbk を 作れなかった")
			continue
		}
		pbks[pokeName] = &pbk
		pbkPokeNames = append(pbkPokeNames, pokeName)
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
			lowerKey := strings.ToLower(key)

			pscs = bippa.NewIndividualCombinations(key)
			pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
			err = pscms.WriteJson(pokeName, lowerKey + "_individual.json")
			if err != nil {
				panic(err)
			}

			pscs = bippa.NewEffortCombinations(key)
			pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
			err = pscms.WriteJson(pokeName, lowerKey + "_effort.json")
			if err != nil {
				panic(err)
			}
		}
	}

	for pokeName, pbk := range pbks {
		pscs := bippa.NewMoveNameAndAbilityCombinations(pbk)
		pscms := bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err := pscms.WriteJson(pokeName, "move_name_and_ability.json")
		if err != nil {
			panic(err)
		}

		pscs = bippa.NewMoveNameAndItemCombinations(pbk)
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err = pscms.WriteJson(pokeName, "move_name_and_item.json")
		if err != nil {
			panic(err)
		}

		pscs = bippa.NewMoveNameAndNatureCombinations(pbk)
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err = pscms.WriteJson(pokeName, "move_name_and_nature.json")
		if err != nil {
			panic(err)
		}

		for _, key := range []string{"HP", "Atk", "Def", "SpAtk", "SpDef", "Speed"} {
			lowerKey := strings.ToLower(key)

			pscs = bippa.NewMoveNameAndIndividualCombinations(pbk, key)
			pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
			err = pscms.WriteJson(pokeName, "move_name_and_" + lowerKey + "_individual.json")
			if err != nil {
				panic(err)
			}
	
			pscs = bippa.NewMoveNameAndEffortCombinations(pbk, key)
			pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
			err = pscms.WriteJson(pokeName, "move_name_and_" + lowerKey + "_effort.json")
			if err != nil {
				panic(err)
			}
		}
		
		pscs, err = bippa.NewMoveNamesCombinations(pbk, 2)
		if err != nil {
			panic(err)
		}
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err = pscms.WriteJson(pokeName, "move_names2.json")
		if err != nil {
			panic(err)
		}

		pscs, err = bippa.NewMoveNames2AndAbilityCombinations(pbk)
		if err != nil {
			panic(err)
		}
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err = pscms.WriteJson(pokeName, "move_names2_and_ability.json")
		if err != nil {
			panic(err)
		}

		pscs, err = bippa.NewMoveNames2AndItemCombinations(pbk)
		if err != nil {
			panic(err)
		}
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err = pscms.WriteJson(pokeName, "move_names2_and_item.json")
		if err != nil {
			panic(err)
		}

		pscs, err = bippa.NewMoveNames2AndNatureCombinations(pbk)
		if err != nil {
			panic(err)
		}
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err = pscms.WriteJson(pokeName, "move_names2_and_nature.json")
		if err != nil {
			panic(err)
		}
	}
}