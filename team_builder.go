package bippa

import (
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
	result := make(PokemonStateCombinations, len(ALL_INDIVIDUALS))
	for i, lowerLimit := range ALL_INDIVIDUALS {
		upperLimit := ALL_UPPER_LIMIT_INDIVIDUALS[i]
		psc := PokemonStateCombination{}
		setter(&psc, Individuals{lowerLimit, upperLimit})
		result[i] = psc
	}
	return result

}

func NewEffortCombinations(key string) PokemonStateCombinations {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[key]
	result := make(PokemonStateCombinations, len(ALL_VALID_EFFORTS))

	for i, lowerLimit := range ALL_VALID_EFFORTS {
		upperLimit := ALL_UPPER_LIMIT_EFFORTS[i]
		psc := PokemonStateCombination{}
		setter(&psc, Efforts{lowerLimit, upperLimit})
		result[i] = psc
	}
	return result

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
	length := len(pbk.MoveNames) * len(LOWER_LIMIT_INDIVIDUALS)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for i, lowerLimit := range LOWER_LIMIT_INDIVIDUALS {
			upperLimit := UPPER_LIMIT_INDIVIDUALS[i]
			psc := PokemonStateCombination{}
			psc.MoveNames = MoveNames{moveName}
			setter(&psc, Individuals{lowerLimit, upperLimit})
			result = append(result, psc)
		}
	}
	return result
}

func NewMoveNameAndEffortCombinations(pbk *PokemonBuildCommonKnowledge, key string) PokemonStateCombinations {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[key]
	length := len(pbk.MoveNames) * len(LOWER_LIMIT_EFFORTS)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for i, lowerLimit := range LOWER_LIMIT_EFFORTS {
			upperLimit := UPPER_LIMIT_EFFORTS[i]
			psc := PokemonStateCombination{}
			psc.MoveNames = MoveNames{moveName}
			setter(&psc, Efforts{lowerLimit, upperLimit})
			result = append(result, psc)
		}
	}
	return result
}

func NewMoveNames3Combinations(pbk *PokemonBuildCommonKnowledge) (PokemonStateCombinations, error) {
	combination3MoveNames, err := pbk.MoveNames.Combination(3)
	if err != nil {
		return PokemonStateCombinations{}, err
	}
	length := len(combination3MoveNames)
	result := make(PokemonStateCombinations, length)

	for i, moveNames := range combination3MoveNames {
		result[i] = PokemonStateCombination{MoveNames:moveNames}
	}
	return result, nil
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
	length := len(combination2MoveNames) * len(LOWER_LIMIT_INDIVIDUALS)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveNames := range combination2MoveNames {
		for i, lowerLimit := range LOWER_LIMIT_INDIVIDUALS {
			upperLimit := UPPER_LIMIT_INDIVIDUALS[i]
			psc := PokemonStateCombination{}
			psc.MoveNames = moveNames
			setter(&psc, Individuals{lowerLimit, upperLimit})
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
	length := len(combination2MoveNames) * len(LOWER_LIMIT_EFFORTS)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveNames := range combination2MoveNames {
		for i, lowerLimit := range LOWER_LIMIT_EFFORTS {
			upperLimit := UPPER_LIMIT_EFFORTS[i]
			psc := PokemonStateCombination{}
			psc.MoveNames = moveNames
			setter(&psc, Efforts{lowerLimit, upperLimit})
			result = append(result, psc)
		}
	}
	return result, nil
}

func NewMoveNameAndNatureAndEffortCombinations(pbk *PokemonBuildCommonKnowledge, key string) PokemonStateCombinations {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[key]
	length := len(pbk.MoveNames) * len(pbk.Natures) * len(LOWER_LIMIT_EFFORTS)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for _, nature := range pbk.Natures {
			for i, lowerLimit := range LOWER_LIMIT_EFFORTS {
				upperLimit := UPPER_LIMIT_EFFORTS[i]
				psc := PokemonStateCombination{}
				psc.MoveNames = MoveNames{moveName}
				psc.Nature = nature
				setter(&psc, Efforts{lowerLimit, upperLimit})
				result = append(result, psc)
			}
		}
	}
	return result
}

func NewNatureAndIndividualAndEffortCombinations(pbk *PokemonBuildCommonKnowledge, individualKey, effortKey string) PokemonStateCombinations {
	individualSetter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_INDIVIDUALS[individualKey]
	effortSetter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[effortKey]
	length := len(pbk.Natures) * len(LOWER_LIMIT_INDIVIDUALS) * len(LOWER_LIMIT_EFFORTS)
	result := make(PokemonStateCombinations, 0, length)

	for _, nature := range pbk.Natures {
		for i, lowerLimitIndividual := range LOWER_LIMIT_INDIVIDUALS {
			upperLimitIndividual := UPPER_LIMIT_INDIVIDUALS[i]
			for j, lowerLimitEffort := range LOWER_LIMIT_EFFORTS {
				upperLimitEffort := UPPER_LIMIT_EFFORTS[j]
				psc := PokemonStateCombination{}
				psc.Nature = nature
				individualSetter(&psc, Individuals{lowerLimitIndividual, upperLimitIndividual})
				effortSetter(&psc, Efforts{lowerLimitEffort, upperLimitEffort})
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

func (pscms PokemonStateCombinationModels) InitNumber() {
	length := len(pscms)
	for i := 0; i < length; i++ {
		pscms[i].Number = i
	}
}

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

func (pscms PokemonStateCombinationModels) WriteJson(pokeName PokeName, fileName string) error {
	filePath := PSCMS_PATH + string(pokeName) + "/" + fileName + ".json"
	file, err := json.MarshalIndent(pscms, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, file, 0644)
}

type MultiplePokemonStateCombination map[PokeName]*PokemonStateCombination

func (mpsc MultiplePokemonStateCombination) OK(pokemons ...*Pokemon) bool {
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

func NewPokemon1MoveNameAndPokemon2Ability(pbk1, pbk2 *PokemonBuildCommonKnowledge) MultiplePokemonStateCombinations {
	allAbilities := POKEDEX[pbk2.PokeName].AllAbilities
	length := len(pbk1.MoveNames) * len(allAbilities)
	result := make(MultiplePokemonStateCombinations, 0, length)

	for _, moveName := range pbk1.MoveNames {
		for _, ability := range allAbilities {
			pbc1 := PokemonStateCombination{MoveNames:MoveNames{moveName}}
			pbc2 := PokemonStateCombination{Ability:ability}
			mpsc := MultiplePokemonStateCombination{pbk1.PokeName:&pbc1, pbk2.PokeName:&pbc2}
			result = append(result, mpsc)
		}
	}
	return result
}

func NewPokemon1MoveNameAndPokemon2Item(pbk1, pbk2 *PokemonBuildCommonKnowledge) MultiplePokemonStateCombinations {
	length := len(pbk1.MoveNames) * len(pbk2.Items)
	result := make(MultiplePokemonStateCombinations, 0, length)
	for _, moveName := range pbk1.MoveNames {
		for _, item := range pbk2.Items {
			pbc1 := PokemonStateCombination{MoveNames:MoveNames{moveName}}
			pbc2 := PokemonStateCombination{Item:item}
			mpsc := MultiplePokemonStateCombination{pbk1.PokeName:&pbc1, pbk2.PokeName:&pbc2}
			result = append(result, mpsc)
		}
	}
	return result
}

func NewPokemon1MoveNameAndPokemon2Nature(pbk1, pbk2 *PokemonBuildCommonKnowledge) MultiplePokemonStateCombinations {
	length := len(pbk1.MoveNames) * len(pbk2.Natures)
	result := make(MultiplePokemonStateCombinations, 0, length)
	for _, moveName := range pbk1.MoveNames {
		for _, nature := range pbk2.Natures {
			pbc1 := PokemonStateCombination{MoveNames:MoveNames{moveName}}
			pbc2 := PokemonStateCombination{Nature:nature}
			mpsc := MultiplePokemonStateCombination{pbk1.PokeName:&pbc1, pbk2.PokeName:&pbc2}
			result = append(result, mpsc)
		}
	}
	return result
}

func NewPokemon1MoveNameAndPokemon2Effort(pbk1, pbk2 *PokemonBuildCommonKnowledge, key string) MultiplePokemonStateCombinations {
	setter := SET_POKEMON_STATE_COMBINATIONL_LOWER_AND_UPPER_LIMIT_EFFORTS[key]
	length := len(pbk1.MoveNames) * len(LOWER_LIMIT_EFFORTS)
	result := make(MultiplePokemonStateCombinations, 0, length)

	for _, moveName := range pbk1.MoveNames {
		for i, lowerLimit := range LOWER_LIMIT_EFFORTS {
			upperLimit := UPPER_LIMIT_EFFORTS[i]
			pbc1 := PokemonStateCombination{MoveNames:MoveNames{moveName}}
			pbc2 := PokemonStateCombination{}
			setter(&pbc2, Efforts{lowerLimit, upperLimit})
			mpsc := MultiplePokemonStateCombination{pbk1.PokeName:&pbc1, pbk2.PokeName:&pbc2}
			result = append(result, mpsc)
		}
	}
	return result
}

type MultiplePokemonStateCombinationModel struct {
	X MultiplePokemonStateCombination
	Value float64
}