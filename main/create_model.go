package main

import (
	"github.com/sw965/bippa"
	"fmt"
	"math/rand"
	"github.com/seehuhn/mt19937"
	"time"
)

var KEY_DATA = map[string]string{"HP":"hp", "Atk":"atk", "Def":"def", "SpAtk":"sp_atk", "SpDef":"sp_def", "Speed":"speed"}

func main() {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())

	pbkList := make(bippa.PokemonBuildCommonKnowledgeList, 0, len(bippa.ALL_POKE_NAMES))
	pbkPokeNames := make(bippa.PokeNames, 0, len(bippa.ALL_POKE_NAMES))

	for _, pokeName := range bippa.ALL_POKE_NAMES {
		pbk, err := bippa.LoadJsonPokemonBuildCommonKnowledge(pokeName)
		if err != nil {
			fmt.Println(pokeName, "の PokemonBuildCommonKnowledge を 作れなかった")
			continue
		}
		pbkList = append(pbkList, pbk)
		pbkPokeNames = append(pbkPokeNames, pokeName)
	}

	permutation2PBKList, err := pbkList.Permutation(2)
	if err != nil {
		panic(err)
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

		for key, lowerKey := range KEY_DATA {
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

	for _, pbk := range pbkList {
		pscs := bippa.NewMoveNameAndAbilityCombinations(&pbk)
		pscms := bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err := pscms.WriteJson(pbk.PokeName, "move_name_and_ability.json")
		if err != nil {
			panic(err)
		}

		pscs = bippa.NewMoveNameAndItemCombinations(&pbk)
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err = pscms.WriteJson(pbk.PokeName, "move_name_and_item.json")
		if err != nil {
			panic(err)
		}

		pscs = bippa.NewMoveNameAndNatureCombinations(&pbk)
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err = pscms.WriteJson(pbk.PokeName, "move_name_and_nature.json")
		if err != nil {
			panic(err)
		}

		for key, lowerKey := range KEY_DATA {
			pscs = bippa.NewMoveNameAndIndividualCombinations(&pbk, key)
			pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
			err = pscms.WriteJson(pbk.PokeName, "move_name_and_" + lowerKey + "_individual.json")
			if err != nil {
				panic(err)
			}
	
			pscs = bippa.NewMoveNameAndEffortCombinations(&pbk, key)
			pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
			err = pscms.WriteJson(pbk.PokeName, "move_name_and_" + lowerKey + "_effort.json")
			if err != nil {
				panic(err)
			}
		}
		
		pscs, err = bippa.NewMoveNamesCombinations(&pbk, 2)
		if err != nil {
			panic(err)
		}
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err = pscms.WriteJson(pbk.PokeName, "move_names2.json")
		if err != nil {
			panic(err)
		}

		pscs, err = bippa.NewMoveNames2AndAbilityCombinations(&pbk)
		if err != nil {
			panic(err)
		}
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err = pscms.WriteJson(pbk.PokeName, "move_names2_and_ability.json")
		if err != nil {
			panic(err)
		}

		pscs, err = bippa.NewMoveNames2AndItemCombinations(&pbk)
		if err != nil {
			panic(err)
		}
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err = pscms.WriteJson(pbk.PokeName, "move_names2_and_item.json")
		if err != nil {
			panic(err)
		}

		pscs, err = bippa.NewMoveNames2AndNatureCombinations(&pbk)
		if err != nil {
			panic(err)
		}
		pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
		err = pscms.WriteJson(pbk.PokeName, "move_names2_and_nature.json")
		if err != nil {
			panic(err)
		}

		for key, lowerKey := range KEY_DATA {
			pscs, err = bippa.NewMoveNames2AndIndividualCombinations(&pbk, key)
			if err != nil {
				panic(err)
			}
			pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
			err = pscms.WriteJson(pbk.PokeName, "move_names2_and_" + lowerKey + "_individual.json")
			if err != nil {
				panic(err)
			}
	
			pscs, err = bippa.NewMoveNames2AndEffortCombinations(&pbk, key)
			if err != nil {
				panic(err)
			}
			pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
			err = pscms.WriteJson(pbk.PokeName, "move_names2_and_" + lowerKey + "_effort.json")
			if err != nil {
				panic(err)
			}

			pscs = bippa.NewMoveNameAndNatureAndIndividualCombinations(&pbk, key)
			pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
			err = pscms.WriteJson(pbk.PokeName, "move_name_and_nature_and_" + lowerKey + "_individual.json")
			if err != nil {
				panic(err)
			}

			pscs = bippa.NewMoveNameAndNatureAndEffortCombinations(&pbk, key)
			pscms = bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
			err = pscms.WriteJson(pbk.PokeName, "move_name_and_nature_and_" + lowerKey + ".json")
			if err != nil {
				panic(err)
			}
		}

		for individualKey, individualLowerKey := range KEY_DATA {
			for effortKey, effortLowerKey := range KEY_DATA {
				pscs := bippa.NewMoveNameAndIndividualAndEffortCombinations(&pbk, individualKey, effortKey)
				pscms := bippa.NewPokemonStateCombinationModels(pscs, mtRandom)
				err := pscms.WriteJson(pbk.PokeName, "move_name_and_" + individualLowerKey + "_individual_and_" + effortLowerKey + "_effort.json")
				if err != nil {
					panic(err)
				}
			}
		}
	}

	for _, pbk := range pbkList {
		for _, pokeName := range bippa.ALL_POKE_NAMES {
			mpscs := bippa.NewPokemon1MoveNameAndPokemon2NameCombinations(&pbk, pokeName)
			mpscms := bippa.NewMultiplePokemonStateCombinationModels(mpscs, mtRandom)
			err := mpscms.WriteJson(pbk.PokeName, pokeName)
			if err != nil {
				panic(err)
			}
		}
	}

	for _, pbkList := range permutation2PBKList {
		mpscs := bippa.NewPokemon1MoveNameAndPokemon2MoveNameCombinations(&pbkList[0], &pbkList[1])
		mpscms := bippa.NewMultiplePokemonStateCombinationModels(mpscs, mtRandom)
		err := mpscms.WriteJson(pbkList[0].PokeName, pbkList[1].PokeName)
		if err != nil {
			panic(err)
		}
	}
}