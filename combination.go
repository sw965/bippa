package bippa

import (
	"os"
	"math/rand"
	"encoding/json"
	"io/ioutil"
	"github.com/sw965/omw"
)

type PokemonBuildCommonKnowledge struct {
	PokeName PokeName
	MoveNames MoveNames
	Items     Items
	Natures   Natures
}

func LoadJsonPokemonBuildCommonKnowledge(pokeName PokeName) (PokemonBuildCommonKnowledge, error) {
	filePath := PBCK_PATH + string(pokeName) + ".json"

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return PokemonBuildCommonKnowledge{}, err
	}

	result := PokemonBuildCommonKnowledge{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return PokemonBuildCommonKnowledge{}, err
	}

	return result, nil
}

type PokemonBuildCommonKnowledgeList []PokemonBuildCommonKnowledge

func (pbkl PokemonBuildCommonKnowledgeList) Access(indices []int) PokemonBuildCommonKnowledgeList {
	result := make(PokemonBuildCommonKnowledgeList, len(indices))
	for i, index := range indices {
		result[i] = pbkl[index]
	}
	return result
}

func (pbkl PokemonBuildCommonKnowledgeList) Permutation(r int) ([]PokemonBuildCommonKnowledgeList, error) {
	n := len(pbkl)
	permutationTotalNum := omw.PermutationTotalNum(n, r)
	permutationNumbers, err := omw.PermutationNumbers(n, r)
	if err != nil {
		return []PokemonBuildCommonKnowledgeList{}, err
	}
	result := make([]PokemonBuildCommonKnowledgeList, permutationTotalNum)

	for i, indices := range permutationNumbers {
		result[i] = pbkl.Access(indices)
	}
	return result, nil
}

type PokemonStateCombination struct {
	MoveNames MoveNames
	Gender Gender
	Ability Ability
	Item Item
	Nature Nature

	LowerAndUpperLimitHPIndividuals Individuals
	LowerAndUpperLimitAtkIndividuals Individuals
	LowerAndUpperLimitDefIndividuals Individuals
	LowerAndUpperLimitSpAtkIndividuals Individuals
	LowerAndUpperLimitSpDefIndividuals Individuals
	LowerAndUpperLimitSpeedIndividuals Individuals

	LowerAndUpperLimitHPEfforts Efforts
	LowerAndUpperLimitAtkEfforts Efforts
	LowerAndUpperLimitDefEfforts Efforts
	LowerAndUpperLimitSpAtkEfforts Efforts
	LowerAndUpperLimitSpDefEfforts Efforts
	LowerAndUpperLimitSpeedEfforts Efforts
}

func (psc *PokemonStateCombination) OK(pokemon *Pokemon) bool {
	if psc.Ability != "" {
		if psc.Ability != pokemon.Ability {
			return false
		}
	}

	if psc.Item != "" {
		if psc.Item != pokemon.Item {
			return false
		}
	}

	if len(psc.MoveNames) != 0 {
		for _, moveName := range psc.MoveNames {
			_, ok := pokemon.Moveset[moveName]
			if !ok {
				return false
			}
		}
	}

	if psc.Nature != "" {
		if psc.Nature != pokemon.Nature {
			return false
		}
	}

	isIndividualOK := func(individual Individual, lowerAndUpperLimitIndividuals Individuals) bool {
		if len(lowerAndUpperLimitIndividuals) == 0 {
			return true
		}
		return (individual >= lowerAndUpperLimitIndividuals[0]) &&  (individual < lowerAndUpperLimitIndividuals[1])
	}

	if !isIndividualOK(pokemon.IndividualState.HP, psc.LowerAndUpperLimitHPIndividuals) {
		return false
	}

	if !isIndividualOK(pokemon.IndividualState.Atk, psc.LowerAndUpperLimitAtkIndividuals) {
		return false
	}

	if !isIndividualOK(pokemon.IndividualState.Def, psc.LowerAndUpperLimitDefIndividuals) {
		return false
	}

	if !isIndividualOK(pokemon.IndividualState.SpAtk, psc.LowerAndUpperLimitSpAtkIndividuals) {
		return false
	}

	if !isIndividualOK(pokemon.IndividualState.SpDef, psc.LowerAndUpperLimitSpDefIndividuals) {
		return false
	}

	if !isIndividualOK(pokemon.IndividualState.Speed, psc.LowerAndUpperLimitSpeedIndividuals) {
		return false
	}

	isEffortOK := func(effort Effort, lowerAndUpperLimitEfforts Efforts) bool {
		if len(lowerAndUpperLimitEfforts) == 0 {
			return true
		}
		return (effort >= lowerAndUpperLimitEfforts[0]) && (effort < lowerAndUpperLimitEfforts[1])
	}

	if !isEffortOK(pokemon.EffortState.HP, psc.LowerAndUpperLimitHPEfforts) {
		return false
	}

	if !isEffortOK(pokemon.EffortState.Atk, psc.LowerAndUpperLimitAtkEfforts) {
		return false
	}

	if !isEffortOK(pokemon.EffortState.Def, psc.LowerAndUpperLimitDefEfforts) {
		return false
	}

	if !isEffortOK(pokemon.EffortState.SpAtk, psc.LowerAndUpperLimitSpAtkEfforts) {
		return false
	}

	if !isEffortOK(pokemon.EffortState.SpDef, psc.LowerAndUpperLimitSpDefEfforts) {
		return false
	}

	if !isEffortOK(pokemon.EffortState.Speed, psc.LowerAndUpperLimitSpeedEfforts) {
		return false
	}
	return true
}

type PokemonStateCombinations []PokemonStateCombination

func NewMoveNameCombinations(pokeName PokeName) PokemonStateCombinations {
	pokeData := POKEDEX[pokeName]
	learnset := pokeData.Learnset
	result := make(PokemonStateCombinations, len(learnset))
	for i, moveName := range learnset {
		result[i] = PokemonStateCombination{MoveNames:MoveNames{moveName}}
	}
	return result
}

func NewGenderCombinations(pokeName PokeName) PokemonStateCombinations {
	validGenders := NewVaildGenders(pokeName)
	result := make(PokemonStateCombinations, len(validGenders))
	for i, gender := range validGenders {
		result[i] = PokemonStateCombination{Gender:gender}
	}
	return result
}

func NewAbilityCombinations(pokeName PokeName) PokemonStateCombinations {
	pokeData := POKEDEX[pokeName]
	allAbilities := pokeData.AllAbilities
	result := make(PokemonStateCombinations, len(allAbilities))
	for i, ability := range allAbilities {
		result[i] = PokemonStateCombination{Ability:ability}
	}
	return result
}

func NewItemCombinations() PokemonStateCombinations {
	result := make(PokemonStateCombinations, len(ALL_ITEMS))
	for i, item := range ALL_ITEMS {
		result[i] = PokemonStateCombination{Item:item}
	}
	return result
}

func NewNatureCombinations() PokemonStateCombinations {
	result := make(PokemonStateCombinations, len(ALL_NATURES))
	for i, nature := range ALL_NATURES {
		result[i] = PokemonStateCombination{Nature:nature}
	}
	return result
}

func NewIndividualCombinations(key string) PokemonStateCombinations {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_INDIVIDUALS[key]
	result := make(PokemonStateCombinations, len(ALL_LOWER_AND_UPPER_LIMIT_INDIVIDUALS))
	for i, lowerAndUpperLimit := range ALL_LOWER_AND_UPPER_LIMIT_INDIVIDUALS {
		psc := PokemonStateCombination{}
		setter(&psc, lowerAndUpperLimit)
		result[i] = psc
	}
	return result

}

func NewEffortCombinations(key string) PokemonStateCombinations {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[key]
	result := make(PokemonStateCombinations, len(ALL_LOWER_AND_UPPER_LIMIT_EFFORTS))

	for i, lowerAndUpperLimit := range ALL_LOWER_AND_UPPER_LIMIT_EFFORTS {
		psc := PokemonStateCombination{}
		setter(&psc, lowerAndUpperLimit)
		result[i] = psc
	}
	return result

}

func NewMoveNamesCombinations(pbk *PokemonBuildCommonKnowledge, r int) (PokemonStateCombinations, error) {
	combinationMoveNames, err := pbk.MoveNames.Combination(r)
	if err != nil {
		return PokemonStateCombinations{}, err
	}
	length := len(combinationMoveNames)
	result := make(PokemonStateCombinations, length)

	for i, moveNames := range combinationMoveNames {
		result[i] = PokemonStateCombination{MoveNames:moveNames}
	}
	return result, nil
}

func NewMoveNameAndAbilityCombinations(pbk *PokemonBuildCommonKnowledge) PokemonStateCombinations {
	allAbilities :=  POKEDEX[pbk.PokeName].AllAbilities
	length := len(pbk.MoveNames) * len(allAbilities)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for _, ability := range allAbilities {
			result = append(result, PokemonStateCombination{MoveNames:MoveNames{moveName}, Ability:ability})
		}
	}
	return result
}

func NewMoveNameAndItemCombinations(pbk *PokemonBuildCommonKnowledge) PokemonStateCombinations {
	length := len(pbk.MoveNames) * len(pbk.Items)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for _, item := range pbk.Items {
			result = append(result, PokemonStateCombination{MoveNames:MoveNames{moveName}, Item:item})
		}
	}
	return result
}

func NewMoveNameAndNatureCombinations(pbk *PokemonBuildCommonKnowledge) PokemonStateCombinations {
	length := len(pbk.MoveNames) * len(pbk.Natures)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for _, nature := range pbk.Natures {
			result = append(result, PokemonStateCombination{MoveNames:MoveNames{moveName}, Nature:nature})
		}
	}
	return result
}

func NewMoveNameAndIndividualCombinations(pbk *PokemonBuildCommonKnowledge, key string) PokemonStateCombinations {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_INDIVIDUALS[key]
	length := len(pbk.MoveNames) * len(SET_LOWER_AND_UPPER_LIMIT_INDIVIDUALS)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for _, lowerAndUpperLimit := range SET_LOWER_AND_UPPER_LIMIT_INDIVIDUALS {
			psc := PokemonStateCombination{}
			psc.MoveNames = MoveNames{moveName}
			setter(&psc, lowerAndUpperLimit)
			result = append(result, psc)
		}
	}
	return result
}

func NewMoveNameAndEffortCombinations(pbk *PokemonBuildCommonKnowledge, key string) PokemonStateCombinations {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[key]
	length := len(pbk.MoveNames) * len(SET_LOWER_AND_UPPER_LIMIT_EFFORTS)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for _, lowerAndUpperLimit := range SET_LOWER_AND_UPPER_LIMIT_EFFORTS {
			psc := PokemonStateCombination{}
			psc.MoveNames = MoveNames{moveName}
			setter(&psc, lowerAndUpperLimit)
			result = append(result, psc)
		}
	}
	return result
}

func NewMoveNames2AndAbilityCombinations(pbk *PokemonBuildCommonKnowledge) (PokemonStateCombinations, error) {
	combination2MoveNames, err := pbk.MoveNames.Combination(2)
	if err != nil {
		return PokemonStateCombinations{}, err
	}
	allAbilities := POKEDEX[pbk.PokeName].AllAbilities
	length := len(combination2MoveNames) * len(allAbilities)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveNames := range combination2MoveNames {
		for _, ability := range allAbilities {
			result = append(result, PokemonStateCombination{MoveNames:moveNames, Ability:ability})
		}
	}
	return result, nil
}

func NewMoveNames2AndItemCombinations(pbk *PokemonBuildCommonKnowledge) (PokemonStateCombinations, error) {
	combination2MoveNames, err := pbk.MoveNames.Combination(2)
	if err != nil {
		return PokemonStateCombinations{}, err
	}
	length := len(combination2MoveNames) * len(pbk.Items)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveNames := range combination2MoveNames {
		for _, item := range pbk.Items {
			result = append(result, PokemonStateCombination{MoveNames:moveNames, Item:item})
		}
	}
	return result, nil
}

func NewMoveNames2AndNatureCombinations(pbk *PokemonBuildCommonKnowledge) (PokemonStateCombinations, error) {
	combination2MoveNames, err := pbk.MoveNames.Combination(2)
	if err != nil {
		return PokemonStateCombinations{}, err
	}
	length := len(combination2MoveNames) * len(pbk.Natures)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveNames := range combination2MoveNames {
		for _, nature := range pbk.Natures {
			result = append(result, PokemonStateCombination{MoveNames:moveNames, Nature:nature})
		}
	}
	return result, nil
}

func NewMoveNames2AndIndividualCombinations(pbk *PokemonBuildCommonKnowledge, key string) (PokemonStateCombinations, error) {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_INDIVIDUALS[key]
	combination2MoveNames, err := pbk.MoveNames.Combination(2)
	if err != nil {
		return PokemonStateCombinations{}, err
	}
	length := len(combination2MoveNames) * len(SET_LOWER_AND_UPPER_LIMIT_INDIVIDUALS)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveNames := range combination2MoveNames {
		for _, lowerAndUpperLimit := range SET_LOWER_AND_UPPER_LIMIT_INDIVIDUALS {
			psc := PokemonStateCombination{MoveNames:moveNames}
			setter(&psc, lowerAndUpperLimit)
			result = append(result, psc)
		}
	}
	return result, nil
}

func NewMoveNames2AndEffortCombinations(pbk *PokemonBuildCommonKnowledge, key string) (PokemonStateCombinations, error) {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[key]
	combination2MoveNames, err := pbk.MoveNames.Combination(2)
	if err != nil {
		return PokemonStateCombinations{}, err
	}
	length := len(combination2MoveNames) * len(SET_LOWER_AND_UPPER_LIMIT_EFFORTS)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveNames := range combination2MoveNames {
		for _, lowerAndUpperLimit := range SET_LOWER_AND_UPPER_LIMIT_EFFORTS {
			psc := PokemonStateCombination{MoveNames:moveNames}
			setter(&psc, lowerAndUpperLimit)
			result = append(result, psc)
		}
	}
	return result, nil
}

func NewMoveNameAndNatureAndIndividualCombinations(pbk *PokemonBuildCommonKnowledge, key string) PokemonStateCombinations {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_INDIVIDUALS[key]
	length := len(pbk.MoveNames) * len(pbk.Natures) * len(SET_LOWER_AND_UPPER_LIMIT_INDIVIDUALS)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for _, nature := range pbk.Natures {
			for _, lowerAndUpperLimit := range SET_LOWER_AND_UPPER_LIMIT_INDIVIDUALS {
				psc := PokemonStateCombination{MoveNames:MoveNames{moveName}, Nature:nature}
				setter(&psc, lowerAndUpperLimit)
				result = append(result, psc)
			}
		}
	}
	return result
}

func NewMoveNameAndNatureAndEffortCombinations(pbk *PokemonBuildCommonKnowledge, key string) PokemonStateCombinations {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[key]
	length := len(pbk.MoveNames) * len(pbk.Natures) * len(SET_LOWER_AND_UPPER_LIMIT_EFFORTS)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for _, nature := range pbk.Natures {
			for _, lowerAndUpperLimit := range SET_LOWER_AND_UPPER_LIMIT_EFFORTS {
				psc := PokemonStateCombination{MoveNames:MoveNames{moveName}, Nature:nature}
				setter(&psc, lowerAndUpperLimit)
				result = append(result, psc)
			}
		}
	}
	return result
}

func NewMoveNameAndIndividualAndEffortCombinations(pbk *PokemonBuildCommonKnowledge, individualKey, effortKey string) PokemonStateCombinations {
	individualSetter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_INDIVIDUALS[individualKey]
	effortSetter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[effortKey]
	length := len(pbk.MoveNames) * len(SET_LOWER_AND_UPPER_LIMIT_INDIVIDUALS) * len(SET_LOWER_AND_UPPER_LIMIT_EFFORTS)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for _, lowerAndUpperLimitIndividuals := range SET_LOWER_AND_UPPER_LIMIT_INDIVIDUALS {
			for _, lowerAndUpperLimitEfforts := range SET_LOWER_AND_UPPER_LIMIT_EFFORTS {
				psc := PokemonStateCombination{MoveNames:MoveNames{moveName}}
				individualSetter(&psc, lowerAndUpperLimitIndividuals)
				effortSetter(&psc, lowerAndUpperLimitEfforts)
				result = append(result, psc)
			}
		}
	}
	return result
}

func SetPokemonStateCombinationLowerAndUpperLimitHPIndividuals(psc *PokemonStateCombination, individuals Individuals) {
	psc.LowerAndUpperLimitHPIndividuals = individuals
}

func SetPokemonStateCombinationLowerAndUpperLimitAtkIndividuals(psc *PokemonStateCombination, individuals Individuals) {
	psc.LowerAndUpperLimitAtkIndividuals = individuals
}

func SetPokemonStateCombinationLowerAndUpperLimitDefIndividuals(psc *PokemonStateCombination, individuals Individuals) {
	psc.LowerAndUpperLimitDefIndividuals = individuals
}

func SetPokemonStateCombinationLowerAndUpperLimitSpAtkIndividuals(psc *PokemonStateCombination, individuals Individuals) {
	psc.LowerAndUpperLimitSpAtkIndividuals = individuals
}

func SetPokemonStateCombinationLowerAndUpperLimitSpDefIndividuals(psc *PokemonStateCombination, individuals Individuals) {
	psc.LowerAndUpperLimitSpDefIndividuals = individuals
}

func SetPokemonStateCombinationLowerAndUpperLimitSpeedIndividuals(psc *PokemonStateCombination, individuals Individuals) {
	psc.LowerAndUpperLimitSpeedIndividuals = individuals
}

var SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_INDIVIDUALS = map[string]func(*PokemonStateCombination, Individuals) {
	"HP":SetPokemonStateCombinationLowerAndUpperLimitHPIndividuals, "Atk":SetPokemonStateCombinationLowerAndUpperLimitAtkIndividuals,
	"Def":SetPokemonStateCombinationLowerAndUpperLimitDefIndividuals, "SpAtk":SetPokemonStateCombinationLowerAndUpperLimitSpAtkIndividuals,
	"SpDef":SetPokemonStateCombinationLowerAndUpperLimitSpDefIndividuals, "Speed":SetPokemonStateCombinationLowerAndUpperLimitSpeedIndividuals,
}

func SetPokemonStateCombinationLowerAndUpperLimitHPEfforts(psc *PokemonStateCombination, efforts Efforts) {
	psc.LowerAndUpperLimitHPEfforts = efforts
}

func SetPokemonStateCombinationLowerAndUpperLimitAtkEfforts(psc *PokemonStateCombination, efforts Efforts) {
	psc.LowerAndUpperLimitAtkEfforts = efforts
}

func SetPokemonStateCombinationLowerAndUpperLimitDefEfforts(psc *PokemonStateCombination, efforts Efforts) {
	psc.LowerAndUpperLimitDefEfforts = efforts
}

func SetPokemonStateCombinationLowerAndUpperLimitSpAtkEfforts(psc *PokemonStateCombination, efforts Efforts) {
	psc.LowerAndUpperLimitSpAtkEfforts = efforts
}

func SetPokemonStateCombinationLowerAndUpperLimitSpDefEfforts(psc *PokemonStateCombination, efforts Efforts) {
	psc.LowerAndUpperLimitSpAtkEfforts = efforts
}

func SetPokemonStateCombinationLowerAndUpperLimitSpeedEfforts(psc *PokemonStateCombination, efforts Efforts) {
	psc.LowerAndUpperLimitSpeedEfforts = efforts
}

var SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS = map[string]func(*PokemonStateCombination, Efforts) {
	"HP":SetPokemonStateCombinationLowerAndUpperLimitHPEfforts, "Atk":SetPokemonStateCombinationLowerAndUpperLimitAtkEfforts,
	"Def":SetPokemonStateCombinationLowerAndUpperLimitDefEfforts, "SpAtk":SetPokemonStateCombinationLowerAndUpperLimitSpDefEfforts,
	"SpDef":SetPokemonStateCombinationLowerAndUpperLimitSpDefEfforts, "Speed":SetPokemonStateCombinationLowerAndUpperLimitSpeedEfforts,
}

type PokemonStateCombinationModel struct {
	X PokemonStateCombination
	Value float64
	Number int
}

type PokemonStateCombinationModels []*PokemonStateCombinationModel

func NewPokemonStateCombinationModels(pscs PokemonStateCombinations, random *rand.Rand) PokemonStateCombinationModels {
	result := make(PokemonStateCombinationModels, len(pscs))
	for i, psc := range pscs {
		value, err := omw.RandomFloat64(0, 16.0, random)
		if err != nil {
			panic(err)
		}
		pscm := PokemonStateCombinationModel{X:psc, Value:value}
		result[i] = &pscm
	}
	result.InitNumber()
	return result
}

func (pscms PokemonStateCombinationModels) InitNumber() {
	length := len(pscms)
	for i := 0; i < length; i++ {
		pscms[i].Number = i
	}
}

func (pscms PokemonStateCombinationModels) WriteJson(pokeName PokeName, fileName string) error {
	folderDirectory := PSCMS_PATH + string(pokeName) + "/"

	if _, err := os.Stat(folderDirectory); err != nil {
		if mkErr := os.Mkdir(folderDirectory, os.ModePerm); mkErr != nil {
			return mkErr
		}
	}

	fullPath := folderDirectory + fileName
	file, err := json.MarshalIndent(pscms, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullPath, file, 0644)
}

type MultiplePokemonStateCombination map[PokeName]*PokemonStateCombination

func (mpsc MultiplePokemonStateCombination) OK(pokemons ...*Pokemon) bool {
	if len(mpsc) > len(pokemons) {
		return false
	}

	for _, pokemon := range pokemons {
		psc, ok := mpsc[pokemon.Name]
		if !ok {
			return false
		}
		if !psc.OK(pokemon) {
			return false
		}
	}
	return true
}

type MultiplePokemonStateCombinations []MultiplePokemonStateCombination

func NewPokemon1MoveNameAndPokemon2NameCombinations(pbk1 *PokemonBuildCommonKnowledge, pokeName2 PokeName) MultiplePokemonStateCombinations {
	result := make(MultiplePokemonStateCombinations, 0, len(pbk1.MoveNames))
	for _, moveName := range pbk1.MoveNames {
		psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName}}
		psc2 := PokemonStateCombination{}
		result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pokeName2:&psc2})
	}
	return result
}

func NewPokemon1MoveNameAndPokemon2MoveNameCombinations(pbk1, pbk2 *PokemonBuildCommonKnowledge) MultiplePokemonStateCombinations {
	length := len(pbk1.MoveNames) * len(pbk2.MoveNames)
	result := make(MultiplePokemonStateCombinations, 0, length)

	for _, moveName1 := range pbk1.MoveNames {
		for _, moveName2 := range pbk2.MoveNames {
			psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}}
			psc2 := PokemonStateCombination{MoveNames:MoveNames{moveName2}}
			mpsc := MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pbk2.PokeName:&psc2}
			result = append(result, mpsc)
		}
	}
	return result
}

func NewPokemon1MoveNameAndPokemon2AbilityCombinations(pbk1, pbk2 *PokemonBuildCommonKnowledge) MultiplePokemonStateCombinations {
	allAbilities := POKEDEX[pbk2.PokeName].AllAbilities
	length := len(pbk1.MoveNames) * len(allAbilities)
	result := make(MultiplePokemonStateCombinations, 0, length)

	for _, moveName1 := range pbk1.MoveNames {
		for _, ability2 := range allAbilities {
			pbc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}}
			pbc2 := PokemonStateCombination{Ability:ability2}
			mpsc := MultiplePokemonStateCombination{pbk1.PokeName:&pbc1, pbk2.PokeName:&pbc2}
			result = append(result, mpsc)
		}
	}
	return result
}

func NewPokemon1MoveNameAndPokemon2ItemCombinations(pbk1, pbk2 *PokemonBuildCommonKnowledge) MultiplePokemonStateCombinations {
	length := len(pbk1.MoveNames) * len(pbk2.Items)
	result := make(MultiplePokemonStateCombinations, 0, length)

	for _, moveName1 := range pbk1.MoveNames {
		for _, item2 := range pbk2.Items {
			pbc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}}
			pbc2 := PokemonStateCombination{Item:item2}
			mpsc := MultiplePokemonStateCombination{pbk1.PokeName:&pbc1, pbk2.PokeName:&pbc2}
			result = append(result, mpsc)
		}
	}
	return result
}

func NewPokemon1MoveNameAndPokemon2NatureCombinations(pbk1, pbk2 *PokemonBuildCommonKnowledge) MultiplePokemonStateCombinations {
	length := len(pbk1.MoveNames) * len(pbk2.Natures)
	result := make(MultiplePokemonStateCombinations, 0, length)

	for _, moveName1 := range pbk1.MoveNames {
		for _, nature2 := range pbk2.Natures {
			pbc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}}
			pbc2 := PokemonStateCombination{Nature:nature2}
			mpsc := MultiplePokemonStateCombination{pbk1.PokeName:&pbc1, pbk2.PokeName:&pbc2}
			result = append(result, mpsc)
		}
	}
	return result
}

func NewPokemon1MoveNameAndPokemon2EffortCombinations(pbk1, pbk2 *PokemonBuildCommonKnowledge, key string) MultiplePokemonStateCombinations {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[key]
	length := len(pbk1.MoveNames) * len(SET_LOWER_AND_UPPER_LIMIT_EFFORTS)
	result := make(MultiplePokemonStateCombinations, 0, length)

	for _, moveName1 := range pbk1.MoveNames {
		for _, lowerAndUpperLimit := range SET_LOWER_AND_UPPER_LIMIT_EFFORTS {
			pbc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}}
			pbc2 := PokemonStateCombination{}
			setter(&pbc2, lowerAndUpperLimit)
			mpsc := MultiplePokemonStateCombination{pbk1.PokeName:&pbc1, pbk2.PokeName:&pbc2}
			result = append(result, mpsc)
		}
	}
	return result
}

func NewPokemon1MoveNames2AndPokemon2MoveNameCombinations(pbk1, pbk2 *PokemonBuildCommonKnowledge) (MultiplePokemonStateCombinations, error) {
	combination2MoveNames, err := pbk1.MoveNames.Combination(2)
	if err != nil {
		return MultiplePokemonStateCombinations{}, err
	}
	length := len(combination2MoveNames) * len(pbk2.MoveNames)
	result := make(MultiplePokemonStateCombinations, 0, length)

	for _, moveNames1 := range combination2MoveNames {
		for _, moveName2 := range pbk2.MoveNames {
			psc1 := PokemonStateCombination{MoveNames:moveNames1}
			psc2 := PokemonStateCombination{MoveNames:MoveNames{moveName2}}
			result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pbk2.PokeName:&psc2})
		}
	}
	return result, nil
}

func NewPokemon1MoveNames2AndPokemon2AbilityCombinations(pbk1, pbk2 *PokemonBuildCommonKnowledge) (MultiplePokemonStateCombinations, error) {
	combination2MoveNames, err := pbk1.MoveNames.Combination(2)
	if err != nil {
		return MultiplePokemonStateCombinations{}, err
	}
	allAbilities := POKEDEX[pbk2.PokeName].AllAbilities
	length := len(combination2MoveNames) * len(allAbilities)
	result := make(MultiplePokemonStateCombinations, 0, length)

	for _, moveNames1 := range combination2MoveNames {
		for _, ability2 := range allAbilities {
			psc1 := PokemonStateCombination{MoveNames:moveNames1}
			psc2 := PokemonStateCombination{Ability:ability2}
			result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pbk2.PokeName:&psc2})
		}
	}
	return result, nil

}

func NewPokemon1MoveNames2AndPokemon2ItemCombinations(pbk1, pbk2 *PokemonBuildCommonKnowledge) (MultiplePokemonStateCombinations, error) {
	combination2MoveNames, err := pbk1.MoveNames.Combination(2)
	if err != nil {
		return MultiplePokemonStateCombinations{}, err
	}
	length := len(combination2MoveNames) * len(pbk2.Items)
	result := make(MultiplePokemonStateCombinations, 0, length)

	for _, moveNames1 := range combination2MoveNames {
		for _, item2 := range pbk2.Items {
			psc1 := PokemonStateCombination{MoveNames:moveNames1}
			psc2 := PokemonStateCombination{Item:item2}
			result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pbk2.PokeName:&psc2})
		}
	}
	return result, nil
}

func NewPokemon1MoveNames2AndPokemon2NatureCombinations(pbk1, pbk2 *PokemonBuildCommonKnowledge) (MultiplePokemonStateCombinations, error) {
	combination2MoveNames, err := pbk1.MoveNames.Combination(2)
	if err != nil {
		return MultiplePokemonStateCombinations{}, err
	}
	length := len(combination2MoveNames) * len(pbk2.Natures)
	result := make(MultiplePokemonStateCombinations, 0, length)
	
	for _, moveNames1 := range combination2MoveNames {
		for _, nature2 := range pbk2.Natures {
			psc1 := PokemonStateCombination{MoveNames:moveNames1}
			psc2 := PokemonStateCombination{Nature:nature2}
			result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pbk2.PokeName:&psc2})
		}
	}
	return result, nil
}

func NewPokemon1MoveNames2AndPokemon2EffortCombinations(pbk1, pbk2 *PokemonBuildCommonKnowledge, key string) (MultiplePokemonStateCombinations, error) {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[key]
	combination2MoveNames, err := pbk1.MoveNames.Combination(2)
	if err != nil {
		return MultiplePokemonStateCombinations{}, err
	}
	length := len(combination2MoveNames) * len(SET_LOWER_AND_UPPER_LIMIT_EFFORTS)
	result := make(MultiplePokemonStateCombinations, 0, length)

	for _, moveNames1 := range combination2MoveNames {
		for _, lowerAndUpperLimit2 := range SET_LOWER_AND_UPPER_LIMIT_EFFORTS {
			psc1 := PokemonStateCombination{MoveNames:moveNames1}
			psc2 := PokemonStateCombination{}
			setter(&psc2, lowerAndUpperLimit2)
			result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pbk2.PokeName:&psc2})
		}
	}
	return result, nil
}

func NewPokemon1MoveNameAndAbilityAndPokemon2AbilityCombinations(pbk1, pbk2 *PokemonBuildCommonKnowledge) (MultiplePokemonStateCombinations, error) {
	allAbilities1 := POKEDEX[pbk1.PokeName].AllAbilities
	allAbilities2 := POKEDEX[pbk2.PokeName].AllAbilities
	length := len(pbk1.MoveNames) * len(allAbilities1) * len(allAbilities2)
	result := make(MultiplePokemonStateCombinations, 0, length)

	for _, moveName1 := range pbk1.MoveNames {
		for _, ability1 := range allAbilities1 {
			for _, ability2 := range allAbilities2 {
				psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}, Ability:ability1}
				psc2 := PokemonStateCombination{Ability:ability2}
				result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pbk2.PokeName:&psc2})
			}
		}
	}
	return result, nil
}

func NewPokemon1MoveNameAndNatureAndPokemon2NatureCombinations(pbk1, pbk2 *PokemonBuildCommonKnowledge) MultiplePokemonStateCombinations {
	length := len(pbk1.MoveNames) * len(pbk1.Natures) * len(pbk2.Natures)
	result := make(MultiplePokemonStateCombinations, 0, length)

	for _, moveName1 := range pbk1.MoveNames {
		for _, nature1 := range pbk1.Natures {
			for _, nature2 := range pbk2.Natures {
				psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}, Nature:nature1}
				psc2 := PokemonStateCombination{Nature:nature2}
				result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pbk2.PokeName:&psc2})
			}
		}
	}
	return result
}

func NewPokemon1MoveNameAndNatureAndPokemon2EffortCombinations(pbk1, pbk2 *PokemonBuildCommonKnowledge, key string) MultiplePokemonStateCombinations {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[key]
	length := len(pbk1.MoveNames) * len(pbk1.Natures) * len(SET_LOWER_AND_UPPER_LIMIT_EFFORTS)
	result := make(MultiplePokemonStateCombinations, 0, length)

	for _, moveName1 := range pbk1.MoveNames {
		for _, nature1 := range pbk1.Natures {
			for _, lowerAndUpperLimit2 := range SET_LOWER_AND_UPPER_LIMIT_EFFORTS {
				psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}, Nature:nature1}
				psc2 := PokemonStateCombination{}
				setter(&psc2, lowerAndUpperLimit2)
				result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pbk2.PokeName:&psc2})
			}
		}
	}
	return result
}

func NewPokemon1MoveNameAndEffortAndPokemon2NatureCombinations(pbk1, pbk2 *PokemonBuildCommonKnowledge, key string) MultiplePokemonStateCombinations {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[key]
	length := len(pbk1.MoveNames) * len(SET_LOWER_AND_UPPER_LIMIT_EFFORTS) * len(pbk2.Natures)
	result := make(MultiplePokemonStateCombinations, 0, length)

	for _, moveName1 := range pbk1.MoveNames {
		for _, lowerAndUpperLimit1 := range SET_LOWER_AND_UPPER_LIMIT_EFFORTS {
			for _, nature2 := range pbk2.Natures {
				psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}}
				setter(&psc1, lowerAndUpperLimit1)
				psc2 := PokemonStateCombination{Nature:nature2}
				result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pbk2.PokeName:&psc2})
			}
		}
	}
	return result
}

func NewPokemon1MoveNameAndEffortAndPokemon2EffortCombinations(pbk1, pbk2 *PokemonBuildCommonKnowledge, key1, key2 string) MultiplePokemonStateCombinations {
	setter1 := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[key1]
	setter2 := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[key2]
	length := len(pbk1.MoveNames) * len(SET_LOWER_AND_UPPER_LIMIT_EFFORTS) * len(SET_LOWER_AND_UPPER_LIMIT_EFFORTS)
	result := make(MultiplePokemonStateCombinations, 0, length)

	for _, moveName1 := range pbk1.MoveNames {
		for _, lowerAndUpperLimit1 := range SET_LOWER_AND_UPPER_LIMIT_EFFORTS {
			for _, lowerAndUpperLimit2 := range SET_LOWER_AND_UPPER_LIMIT_EFFORTS {
				psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}}
				setter1(&psc1, lowerAndUpperLimit1)
				psc2 := PokemonStateCombination{}
				setter2(&psc2, lowerAndUpperLimit2)
				result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pbk2.PokeName:&psc2})
			}
		}
	}
	return result
}

func NewPokemon1MoveNames2AndPokemon2NameCombinations(pbk1 *PokemonBuildCommonKnowledge, pokeName2 PokeName) (MultiplePokemonStateCombinations, error) {
	combination2MoveNames, err := pbk1.MoveNames.Combination(2)
	if err != nil {
		return MultiplePokemonStateCombinations{}, err
	}
	result := make(MultiplePokemonStateCombinations, 0, len(combination2MoveNames))

	for _, moveNames1 := range combination2MoveNames {
		psc1 := PokemonStateCombination{MoveNames:moveNames1}
		psc2 := PokemonStateCombination{}
		result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pokeName2:&psc2})
	}
	return result, nil
}

func NewPokemon1MoveNameAndGenderAndPokemon2NameCombinations(pbk1 *PokemonBuildCommonKnowledge, pokeName2 PokeName) MultiplePokemonStateCombinations {
	validGenders1 := NewVaildGenders(pbk1.PokeName)
	result := make(MultiplePokemonStateCombinations, 0, len(pbk1.MoveNames) * len(validGenders1))

	for _, moveName1 := range pbk1.MoveNames {
		for _, gender1 := range validGenders1 {
			psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}, Gender:gender1}
			psc2 := PokemonStateCombination{}
			result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pokeName2:&psc2})
		}
	}
	return result
}

func NewPokemon1MoveNameAndAbilityAndPokemon2NameCombinations(pbk1 *PokemonBuildCommonKnowledge, pokeName2 PokeName) MultiplePokemonStateCombinations {
	allAbilities1 := POKEDEX[pbk1.PokeName].AllAbilities
	result := make(MultiplePokemonStateCombinations, 0, len(pbk1.MoveNames) * len(allAbilities1))

	for _, moveName1 := range pbk1.MoveNames {
		for _, ability1 := range allAbilities1 {
			psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}, Ability:ability1}
			psc2 := PokemonStateCombination{}
			result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pokeName2:&psc2})
		}
	}
	return result
}

func NewPokemon1MoveNameAndItemAndPokemon2NameCombinations(pbk1 *PokemonBuildCommonKnowledge, pokeName2 PokeName) MultiplePokemonStateCombinations {
	result := make(MultiplePokemonStateCombinations, 0, len(pbk1.MoveNames) * len(pbk1.Items))
	for _, moveName1 := range pbk1.MoveNames {
		for _, item1 := range pbk1.Items {
			psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}, Item:item1}
			psc2 := PokemonStateCombination{}
			result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pokeName2:&psc2})
		}
	}
	return result
}

func NewPokemon1MoveNameAndNatureAndPokemon2NameCombinations(pbk1 *PokemonBuildCommonKnowledge, pokeName2 PokeName) MultiplePokemonStateCombinations {
	result := make(MultiplePokemonStateCombinations, 0, len(pbk1.MoveNames) * len(pbk1.Natures))
	for _, moveName1 := range pbk1.MoveNames {
		for _, nature1 := range pbk1.Natures {
			psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}, Nature:nature1}
			psc2 := PokemonStateCombination{}
			result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pokeName2:&psc2})
		}
	}
	return result
}

func NewPokemon1MoveNameAndEffortAndPokemon2NameCombinations(pbk1 *PokemonBuildCommonKnowledge, pokeName2 PokeName, key string) MultiplePokemonStateCombinations {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[key]
	result := make(MultiplePokemonStateCombinations, 0, len(pbk1.MoveNames) * len(SET_LOWER_AND_UPPER_LIMIT_EFFORTS))

	for _, moveName1 := range pbk1.MoveNames {
		for _, lowerAndUpperLimit1 := range SET_LOWER_AND_UPPER_LIMIT_EFFORTS {
			psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}}
			setter(&psc1, lowerAndUpperLimit1)
			psc2 := PokemonStateCombination{}
			result = append(result, MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pokeName2:&psc2})
		}
	}
	return result
}

func NewPokemon1MoveNameAndPokemon2NameAndPokemon3Name(pbk1 *PokemonBuildCommonKnowledge, pokeName2, pokeName3 PokeName) MultiplePokemonStateCombinations {
	result := make(MultiplePokemonStateCombinations, len(pbk1.MoveNames))
	for i, moveName1 := range pbk1.MoveNames {
		psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}}
		psc2 := PokemonStateCombination{}
		psc3 := PokemonStateCombination{}
		result[i] = MultiplePokemonStateCombination{pbk1.PokeName:&psc1, pokeName2:&psc2, pokeName3:&psc3}
	}
	return result
}

type MultiplePokemonStateCombinationModel struct {
	X MultiplePokemonStateCombination
	Value float64
	Number int
}

type MultiplePokemonStateCombinationModels []MultiplePokemonStateCombinationModel

func NewMultiplePokemonStateCombinationModels(mpscs MultiplePokemonStateCombinations, random *rand.Rand) MultiplePokemonStateCombinationModels {
	length := len(mpscs)
	result := make(MultiplePokemonStateCombinationModels, 0, length)
	for i, mpsc := range mpscs {
		value, err := omw.RandomFloat64(0.01, 16.0, random)
		if err != nil {
			panic(err)
		}
		result = append(result, MultiplePokemonStateCombinationModel{X:mpsc, Value:value, Number:i})
	}
	return result
}

func (mpscms MultiplePokemonStateCombinationModels) WriteJson(folderName string, pokeNames ...PokeName) error {
	lastIndex := len(pokeNames) - 1
	folderDirectory := MPSCMS_PATH + folderName

	for _, pokeName := range pokeNames[:lastIndex] {
		folderDirectory += string(pokeName) + "/"
	}

	if _, err := os.Stat(folderDirectory); err != nil {
		if mkErr := os.MkdirAll(folderDirectory, os.ModePerm); mkErr != nil {
			return mkErr
		}
	}

	fullPath := folderDirectory + string(pokeNames[lastIndex]) + ".json"

	file, err := json.MarshalIndent(mpscms, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullPath, file, 0644)
}