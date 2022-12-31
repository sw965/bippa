package bippa

import (
	"fmt"
	"encoding/json"
	//"os"
	"io/ioutil"
	"github.com/sw965/omw"
	"math/rand"
)

type PokemonBuildCommonKnowledge struct {
	MoveNames MoveNames
	Items     Items
	Natures   Natures

	HPIndividuals Individuals
	AtkIndividuals Individuals
	DefIndividuals Individuals
	SpAtkIndividuals Individuals
	SpDefIndividuals Individuals
	SpeedIndividuals Individuals

	HPEfforts Efforts
	AtkEfforts Efforts
	DefEfforts Efforts
	SpAtkEfforts Efforts
	SpDefEfforts Efforts
	SpeedEfforts Efforts
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

	individualErrMsg := fmt.Sprintf("PokemonBuildCommonKnowledge の 個体値 に %v と %v を 含めてはならない", MIN_INDIVIDUAL, MAX_INDIVIDUAL)

	if result.HPIndividuals.In(MIN_INDIVIDUAL) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(individualErrMsg)
	}

	if result.HPIndividuals.In(MAX_INDIVIDUAL) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(individualErrMsg)
	}

	if result.AtkIndividuals.In(MIN_INDIVIDUAL) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(individualErrMsg)
	}

	if result.AtkIndividuals.In(MAX_INDIVIDUAL) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(individualErrMsg)
	}

	if result.DefIndividuals.In(MIN_INDIVIDUAL) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(individualErrMsg)
	}

	if result.DefIndividuals.In(MAX_INDIVIDUAL) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(individualErrMsg)
	}

	if result.SpAtkIndividuals.In(MIN_INDIVIDUAL) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(individualErrMsg)
	}

	if result.SpAtkIndividuals.In(MAX_INDIVIDUAL) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(individualErrMsg)
	}

	if result.SpDefIndividuals.In(MIN_INDIVIDUAL) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(individualErrMsg)
	}

	if result.SpDefIndividuals.In(MAX_INDIVIDUAL) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(individualErrMsg)
	}

	if result.SpeedIndividuals.In(MIN_INDIVIDUAL) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(individualErrMsg)
	}

	if result.SpeedIndividuals.In(MAX_INDIVIDUAL) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(individualErrMsg)
	}

	result.HPIndividuals = append(result.HPIndividuals, MIN_INDIVIDUAL)
	result.HPIndividuals = append(result.HPIndividuals, MAX_INDIVIDUAL)
	result.AtkIndividuals = append(result.AtkIndividuals, MIN_INDIVIDUAL)
	result.AtkIndividuals = append(result.AtkIndividuals, MAX_INDIVIDUAL)
	result.DefIndividuals = append(result.DefIndividuals, MIN_INDIVIDUAL)
	result.DefIndividuals = append(result.DefIndividuals, MAX_INDIVIDUAL)
	result.SpAtkIndividuals = append(result.SpAtkIndividuals, MIN_INDIVIDUAL)
	result.SpAtkIndividuals = append(result.SpAtkIndividuals, MAX_INDIVIDUAL)
	result.SpDefIndividuals = append(result.SpDefIndividuals, MIN_INDIVIDUAL)
	result.SpDefIndividuals = append(result.SpDefIndividuals, MAX_INDIVIDUAL)
	result.SpeedIndividuals = append(result.SpeedIndividuals, MIN_INDIVIDUAL)
	result.SpeedIndividuals = append(result.SpeedIndividuals, MAX_INDIVIDUAL)

	effortErrMsg := fmt.Sprintf("PokemonBuildCommonKnowledge の 努力値 に %v と %v を 含めてはならない", MIN_EFFORT, MAX_EFFORT)

	if result.HPEfforts.In(MIN_EFFORT) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(effortErrMsg)
	}

	if result.HPEfforts.In(MAX_EFFORT) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(effortErrMsg)
	}

	if result.AtkEfforts.In(MIN_EFFORT) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(effortErrMsg)
	}

	if result.AtkEfforts.In(MAX_EFFORT) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(effortErrMsg)
	}

	if result.DefEfforts.In(MIN_EFFORT) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(effortErrMsg)
	}

	if result.DefEfforts.In(MAX_EFFORT) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(effortErrMsg)
	}

	if result.SpAtkEfforts.In(MIN_EFFORT) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(effortErrMsg)
	}

	if result.SpAtkEfforts.In(MAX_EFFORT) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(effortErrMsg)
	}

	if result.SpDefEfforts.In(MIN_EFFORT) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(effortErrMsg)
	}

	if result.SpDefEfforts.In(MAX_EFFORT) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(effortErrMsg)
	}

	if result.SpeedEfforts.In(MIN_EFFORT) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(effortErrMsg)
	}

	if result.SpeedEfforts.In(MAX_EFFORT) {
		return PokemonBuildCommonKnowledge{}, fmt.Errorf(effortErrMsg)
	}

	result.HPEfforts = append(result.HPEfforts, MIN_EFFORT)
	result.HPEfforts = append(result.HPEfforts, MAX_EFFORT)
	result.AtkEfforts = append(result.AtkEfforts, MIN_EFFORT)
	result.AtkEfforts = append(result.AtkEfforts, MAX_EFFORT)
	result.DefEfforts = append(result.DefEfforts, MIN_EFFORT)
	result.DefEfforts = append(result.DefEfforts, MAX_EFFORT)
	result.SpAtkEfforts = append(result.SpAtkEfforts, MIN_EFFORT)
	result.SpAtkEfforts = append(result.SpAtkEfforts, MAX_EFFORT)
	result.SpDefEfforts = append(result.SpDefEfforts, MIN_EFFORT)
	result.SpDefEfforts = append(result.SpDefEfforts, MAX_EFFORT)
	result.SpeedEfforts = append(result.SpeedEfforts, MIN_EFFORT)
	result.SpeedEfforts = append(result.SpeedEfforts, MAX_EFFORT)

	return result, nil
}

func GetPokemonBuildCommonKnowledgeHPIndividuals(pbk *PokemonBuildCommonKnowledge) Individuals {
	return pbk.HPIndividuals
}

func GetPokemonBuildCommonKnowledgeAtkIndividuals(pbk *PokemonBuildCommonKnowledge) Individuals {
	return pbk.AtkIndividuals
}

func GetPokemonBuildCommonKnowledgeDefIndividuals(pbk *PokemonBuildCommonKnowledge) Individuals {
	return pbk.DefIndividuals
}

func GetPokemonBuildCommonKnowledgeSpAtkIndividuals(pbk *PokemonBuildCommonKnowledge) Individuals {
	return pbk.SpAtkIndividuals
}

func GetPokemonBuildCommonKnowledgeSpDefIndividuals(pbk *PokemonBuildCommonKnowledge) Individuals {
	return pbk.SpDefIndividuals
}

var GET_POKEMON_BUILD_COMMON_KNOWLEDGE_INDIVIDUALS = map[string]func(*PokemonBuildCommonKnowledge) Individuals {
	"HP":GetPokemonBuildCommonKnowledgeHPIndividuals, "Atk":GetPokemonBuildCommonKnowledgeAtkIndividuals,
	"Def":GetPokemonBuildCommonKnowledgeDefIndividuals, "SpAtk":GetPokemonBuildCommonKnowledgeSpAtkIndividuals,
	"SpDef":GetPokemonBuildCommonKnowledgeSpDefIndividuals, "Speed":GetPokemonBuildCommonKnowledgeSpeedIndividuals,
}

func GetPokemonBuildCommonKnowledgeSpeedIndividuals(pbk *PokemonBuildCommonKnowledge) Individuals {
	return pbk.SpeedIndividuals
}

func GetPokemonBuildCommonKnowledgeHPEfforts(pbk *PokemonBuildCommonKnowledge) Efforts {
	return pbk.HPEfforts
}

func GetPokemonBuildCommonKnowledgeAtkEfforts(pbk *PokemonBuildCommonKnowledge) Efforts {
	return pbk.AtkEfforts
}

func GetPokemonBuildCommonKnowledgeDefEfforts(pbk *PokemonBuildCommonKnowledge) Efforts {
	return pbk.DefEfforts
}

func GetPokemonBuildCommonKnowledgeSpAtkEfforts(pbk *PokemonBuildCommonKnowledge) Efforts {
	return pbk.SpAtkEfforts
}

func GetPokemonBuildCommonKnowledgeSpDefEfforts(pbk *PokemonBuildCommonKnowledge) Efforts {
	return pbk.SpDefEfforts
}

func GetPokemonBuildCommonKnowledgeSpeedEfforts(pbk *PokemonBuildCommonKnowledge) Efforts {
	return pbk.SpeedEfforts
}

var GET_POKEMON_BUILD_COMMON_KNOWLEDGE_EFFORTS = map[string]func(*PokemonBuildCommonKnowledge) Efforts {
	"HP":GetPokemonBuildCommonKnowledgeHPEfforts, "Atk":GetPokemonBuildCommonKnowledgeAtkEfforts,
	"Def":GetPokemonBuildCommonKnowledgeDefEfforts, "SpAtk":GetPokemonBuildCommonKnowledgeSpAtkEfforts,
	"SpDef":GetPokemonBuildCommonKnowledgeSpDefEfforts, "Speed":GetPokemonBuildCommonKnowledgeSpeedEfforts,
}

type PokemonStateCombination struct {
	MoveNames MoveNames
	Gender Gender
	Ability Ability
	Item Item
	Nature Nature

	HPIndividual Individual
	AtkIndividual Individual
	DefIndividual Individual
	SpAtkIndividual Individual
	SpDefIndividual Individual
	SpeedIndividual Individual

	HPEffort Effort
	AtkEffort Effort
	DefEffort Effort
	SpAtkEffort Effort
	SpDefEffort Effort
	SpeedEffort Effort
}

func NewInitPokemonStateCombination() PokemonStateCombination {
	result := PokemonStateCombination{}

	result.HPIndividual = EMPTY_INDIVIDUAL
	result.AtkIndividual = EMPTY_INDIVIDUAL
	result.DefIndividual = EMPTY_INDIVIDUAL
	result.SpAtkIndividual = EMPTY_INDIVIDUAL
	result.SpDefIndividual = EMPTY_INDIVIDUAL
	result.SpeedIndividual = EMPTY_INDIVIDUAL

	result.HPEffort = EMPTY_EFFORT
	result.AtkEffort = EMPTY_EFFORT
	result.DefEffort = EMPTY_EFFORT
	result.SpAtkEffort = EMPTY_EFFORT
	result.SpDefEffort = EMPTY_EFFORT
	result.SpeedEffort = EMPTY_EFFORT

	return result
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

	if psc.HPIndividual != EMPTY_INDIVIDUAL {
		if pokemon.IndividualState.HP == EMPTY_INDIVIDUAL {
			return false
		}
	}

	if psc.AtkIndividual != EMPTY_INDIVIDUAL {
		if pokemon.IndividualState.Atk == EMPTY_INDIVIDUAL {
			return false
		}
	}

	if psc.DefIndividual != EMPTY_INDIVIDUAL {
		if pokemon.IndividualState.Def == EMPTY_INDIVIDUAL {
			return false
		}
	}

	if psc.SpAtkIndividual != EMPTY_INDIVIDUAL {
		if pokemon.IndividualState.SpAtk == EMPTY_INDIVIDUAL {
			return false
		}
	}

	if psc.SpDefIndividual != EMPTY_INDIVIDUAL {
		if pokemon.IndividualState.SpDef == EMPTY_INDIVIDUAL {
			return false
		}
	}

	if psc.SpeedIndividual != EMPTY_INDIVIDUAL {
		if pokemon.IndividualState.Speed == EMPTY_INDIVIDUAL {
			return false
		}
	}

	if psc.HPEffort != EMPTY_EFFORT {
		if pokemon.EffortState.HP == EMPTY_EFFORT {
			return false
		}
	}

	if psc.AtkEffort != EMPTY_EFFORT {
		if pokemon.EffortState.Atk == EMPTY_EFFORT {
			return false
		}
	}

	if psc.DefEffort != EMPTY_EFFORT {
		if pokemon.EffortState.Def == EMPTY_EFFORT {
			return false
		}
	}

	if psc.SpAtkEffort != EMPTY_EFFORT {
		if pokemon.EffortState.SpAtk == EMPTY_EFFORT {
			return false
		}
	}

	if psc.SpDefEffort != EMPTY_EFFORT {
		if pokemon.EffortState.SpDef == EMPTY_EFFORT {
			return false
		}
	}

	if psc.SpeedEffort != EMPTY_EFFORT {
		if pokemon.EffortState.Speed == EMPTY_EFFORT {
			return false
		}
	}
	return true
}

type PokemonStateCombinations []PokemonStateCombination

func NewMoveNameCombinations(pokeName PokeName) PokemonStateCombinations {
	pokeData := POKEDEX[pokeName]
	learnset := pokeData.Learnset
	result := make(PokemonStateCombinations, len(learnset))
	for i, moveName := range learnset {
		psc := NewInitPokemonStateCombination()
		psc.MoveNames = MoveNames{moveName}
		result[i] = psc
	}
	return result
}

func NewGenderCombinations(pokeName PokeName) PokemonStateCombinations {
	validGenders := NewVaildGenders(pokeName)
	result := make(PokemonStateCombinations, len(validGenders))
	for i, gender := range validGenders {
		psc := NewInitPokemonStateCombination()
		psc.Gender = gender
		result[i] = psc
	}
	return result
}

func NewAbilityCombinations(pokeName PokeName) PokemonStateCombinations {
	pokeData := POKEDEX[pokeName]
	allAbilities := pokeData.AllAbilities
	result := make(PokemonStateCombinations, len(allAbilities))
	for i, ability := range allAbilities {
		psc := NewInitPokemonStateCombination()
		psc.Ability = ability
		result[i] = psc
	}
	return result
}

func NewItemCombinations() PokemonStateCombinations {
	result := make(PokemonStateCombinations, len(ALL_ITEMS))
	for i, item := range ALL_ITEMS {
		psc := NewInitPokemonStateCombination()
		psc.Item = item
		result[i] = psc
	}
	return result
}

func NewNatureCombinations() PokemonStateCombinations {
	result := make(PokemonStateCombinations, len(ALL_NATURES))
	for i, nature := range ALL_NATURES {
		psc := NewInitPokemonStateCombination()
		psc.Nature = nature
		result[i] = psc
	}
	return result
}

func NewIndividualCombinations(pbk *PokemonBuildCommonKnowledge, key string) PokemonStateCombinations {
	getter := GET_POKEMON_BUILD_COMMON_KNOWLEDGE_INDIVIDUALS[key]
	setter := SET_POKEMON_STATE_COMBINATION_INDIVIDUAL[key]
	pbkIndividuals := getter(pbk)
	result := make(PokemonStateCombinations, len(pbkIndividuals))

	for i, individual := range pbkIndividuals {
		psc := NewInitPokemonStateCombination()
		setter(&psc, individual)
		result[i] = psc
	}
	return result

}

func NewEffortCombinations(pbk *PokemonBuildCommonKnowledge, key string) PokemonStateCombinations {
	getter := GET_POKEMON_BUILD_COMMON_KNOWLEDGE_EFFORTS[key]
	setter := SET_POKEMON_STATE_COMBINATION_EFFORT[key]
	pbkEfforts := getter(pbk)
	result := make(PokemonStateCombinations, len(pbkEfforts))

	for i, effort := range pbkEfforts {
		psc := NewInitPokemonStateCombination()
		setter(&psc, effort)
		result[i] = psc
	}
	return result

}

func NewMoveNameAndAbilityCombinations(pokeName PokeName, pbk *PokemonBuildCommonKnowledge) PokemonStateCombinations {
	allAbilities :=  POKEDEX[pokeName].AllAbilities
	length := len(pbk.MoveNames) * len(allAbilities)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for _, ability := range allAbilities {
			psc := NewInitPokemonStateCombination()
			psc.MoveNames = MoveNames{moveName}
			psc.Ability = ability
			result = append(result, psc)
		}
	}
	return result
}

func NewMoveNameAndItemCombinations(pbk *PokemonBuildCommonKnowledge) PokemonStateCombinations {
	length := len(pbk.MoveNames) * len(pbk.Items)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for _, item := range pbk.Items {
			psc := NewInitPokemonStateCombination()
			psc.MoveNames = MoveNames{moveName}
			psc.Item = item
			result = append(result, psc)
		}
	}
	return result
}

func NewMoveNameAndNatureCombinations(pbk *PokemonBuildCommonKnowledge) PokemonStateCombinations {
	length := len(pbk.MoveNames) * len(pbk.Natures)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for _, nature := range pbk.Natures {
			psc := NewInitPokemonStateCombination()
			psc.MoveNames = MoveNames{moveName}
			psc.Nature = nature
			result = append(result, psc)
		}
	}
	return result
}

func NewMoveNameAndIndividualCombinations(pbk *PokemonBuildCommonKnowledge, key string) PokemonStateCombinations {
	getter := GET_POKEMON_BUILD_COMMON_KNOWLEDGE_INDIVIDUALS[key]
	setter := SET_POKEMON_STATE_COMBINATION_INDIVIDUAL[key]
	pbkIndividuals := getter(pbk)
	length := len(pbk.MoveNames) * len(pbkIndividuals)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for _, individual := range pbkIndividuals {
			psc := NewInitPokemonStateCombination()
			psc.MoveNames = MoveNames{moveName}
			setter(&psc, individual)
			result = append(result, psc)
		}
	}
	return result
}

func NewMoveNameAndEffortCombinations(pbk *PokemonBuildCommonKnowledge, key string) PokemonStateCombinations {
	getter := GET_POKEMON_BUILD_COMMON_KNOWLEDGE_EFFORTS[key]
	setter := SET_POKEMON_STATE_COMBINATION_EFFORT[key]
	pbkEfforts := getter(pbk)
	length := len(pbk.MoveNames) * len(pbkEfforts)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for _, effort := range pbkEfforts {
			psc := NewInitPokemonStateCombination()
			psc.MoveNames = MoveNames{moveName}
			setter(&psc, effort)
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
		psc := NewInitPokemonStateCombination()
		psc.MoveNames = moveNames
		result[i] = psc
	}
	return result, nil
}

func NewMoveNames2AndAbilityCombinations(pokeName PokeName, pbk *PokemonBuildCommonKnowledge) (PokemonStateCombinations, error) {
	combination2MoveNames, err := pbk.MoveNames.Combination(2)
	if err != nil {
		return PokemonStateCombinations{}, err
	}
	allAbilities := POKEDEX[pokeName].AllAbilities
	length := len(combination2MoveNames) * len(allAbilities)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveNames := range combination2MoveNames {
		for _, ability := range allAbilities {
			psc := NewInitPokemonStateCombination()
			psc.MoveNames = moveNames
			psc.Ability = ability
			result = append(result, psc)
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
			psc := NewInitPokemonStateCombination()
			psc.MoveNames = moveNames
			psc.Item = item
			result = append(result, psc)
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
			psc := NewInitPokemonStateCombination()
			psc.MoveNames = moveNames
			psc.Nature = nature
			result = append(result, psc)
		}
	}
	return result, nil
}

func NewMoveNames2AndIndividualCombinations(pbk *PokemonBuildCommonKnowledge, key string) (PokemonStateCombinations, error) {
	getter := GET_POKEMON_BUILD_COMMON_KNOWLEDGE_INDIVIDUALS[key]
	setter := SET_POKEMON_STATE_COMBINATION_INDIVIDUAL[key]
	pbkIndividuals := getter(pbk)
	combination2MoveNames, err := pbk.MoveNames.Combination(2)
	if err != nil {
		return PokemonStateCombinations{}, err
	}
	length := len(combination2MoveNames) * len(pbkIndividuals)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveNames := range combination2MoveNames {
		for _, individual := range pbkIndividuals {
			psc := NewInitPokemonStateCombination()
			psc.MoveNames = moveNames
			setter(&psc, individual)
			result = append(result, psc)
		}
	}
	return result, nil
}

func NewMoveNames2AndEffortCombinations(pbk *PokemonBuildCommonKnowledge, key string) (PokemonStateCombinations, error) {
	getter := GET_POKEMON_BUILD_COMMON_KNOWLEDGE_EFFORTS[key]
	setter := SET_POKEMON_STATE_COMBINATION_EFFORT[key]
	pbkEfforts := getter(pbk)
	combination2MoveNames, err := pbk.MoveNames.Combination(2)
	if err != nil {
		return PokemonStateCombinations{}, err
	}
	length := len(combination2MoveNames) * len(pbkEfforts)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveNames := range combination2MoveNames {
		for _, effort := range pbkEfforts {
			psc := NewInitPokemonStateCombination()
			psc.MoveNames = moveNames
			setter(&psc, effort)
			result = append(result, psc)
		}
	}
	return result, nil
}

func NewMoveNameAndNatureAndEffortCombinations(pbk *PokemonBuildCommonKnowledge, key string) PokemonStateCombinations {
	getter := GET_POKEMON_BUILD_COMMON_KNOWLEDGE_EFFORTS[key]
	setter := SET_POKEMON_STATE_COMBINATION_EFFORT[key]
	pbkEfforts := getter(pbk)
	length := len(pbk.MoveNames) * len(pbk.Natures) * len(pbkEfforts)
	result := make(PokemonStateCombinations, 0, length)

	for _, moveName := range pbk.MoveNames {
		for _, nature := range pbk.Natures {
			for _, effort := range pbkEfforts {
				psc := NewInitPokemonStateCombination()
				psc.MoveNames = MoveNames{moveName}
				psc.Nature = nature
				setter(&psc, effort)
				result = append(result, psc)
			}
		}
	}
	return result
}

func SetPokemonStateCombinationHPIndividual(psc *PokemonStateCombination, individual Individual) {
	psc.HPIndividual = individual
}

func SetPokemonStateCombinationAtkIndividual(psc *PokemonStateCombination, individual Individual) {
	psc.AtkIndividual = individual
}

func SetPokemonStateCombinationDefIndividual(psc *PokemonStateCombination, individual Individual) {
	psc.DefIndividual = individual
}

func SetPokemonStateCombinationSpAtkIndividual(psc *PokemonStateCombination, individual Individual) {
	psc.SpAtkIndividual = individual
}

func SetPokemonStateCombinationSpDefIndividual(psc *PokemonStateCombination, individual Individual) {
	psc.SpDefIndividual = individual
}

func SetPokemonStateCombinationSpeedIndividual(psc *PokemonStateCombination, individual Individual) {
	psc.SpeedIndividual = individual
}

var SET_POKEMON_STATE_COMBINATION_INDIVIDUAL = map[string]func(*PokemonStateCombination, Individual){
	"HP":SetPokemonStateCombinationHPIndividual, "Atk":SetPokemonStateCombinationAtkIndividual,
	"Def":SetPokemonStateCombinationDefIndividual, "SpAtk":SetPokemonStateCombinationSpAtkIndividual,
	"SpDef":SetPokemonStateCombinationSpDefIndividual, "Speed":SetPokemonStateCombinationSpeedIndividual,
}

func SetPokemonStateCombinationHPEffort(psc *PokemonStateCombination, effort Effort) {
	psc.HPEffort = effort
}

func SetPokemonStateCombinationAtkEffort(psc *PokemonStateCombination, effort Effort) {
	psc.AtkEffort = effort
}

func SetPokemonStateCombinationDefEffort(psc *PokemonStateCombination, effort Effort) {
	psc.DefEffort = effort
}

func SetPokemonStateCombinationSpAtkEffort(psc *PokemonStateCombination, effort Effort) {
	psc.SpAtkEffort = effort
}

func SetPokemonStateCombinationSpDefEffort(psc *PokemonStateCombination, effort Effort) {
	psc.SpDefEffort = effort
}

func SetPokemonStateCombinationSpeedEffort(psc *PokemonStateCombination, effort Effort) {
	psc.SpeedEffort = effort
}

var SET_POKEMON_STATE_COMBINATION_EFFORT = map[string]func(*PokemonStateCombination, Effort) {
	"HP":SetPokemonStateCombinationHPEffort, "Atk":SetPokemonStateCombinationAtkEffort,
	"Def":SetPokemonStateCombinationDefEffort, "SpAtk":SetPokemonStateCombinationSpAtkEffort,
	"SpDef":SetPokemonStateCombinationSpDefEffort, "Speed":SetPokemonStateCombinationSpeedEffort,
}

type PokemonStateCombinationModel struct {
	X PokemonStateCombination
	Policy float64
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
		policy, err := omw.RandomFloat64(0.01, 16.0, random)
		if err != nil {
			panic(err)
		}

		value, err := omw.RandomFloat64(0.01, 16.0, random)
		if err != nil {
			panic(err)
		}

		pscm := PokemonStateCombinationModel{X:psc, Policy:policy, Value:value}
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

// func (pscs PokemonStateCombinations) GetOK(pokemon *Pokemon) PokemonStateCombinations {
// 	result := make(PokemonStateCombinations, 0, len(pscs))
// 	for _, psc := range pscs {
// 		if psc.OK(pokemon) {
// 			result = append(result, psc)
// 		}
// 	}
// 	return result
// }

// func (pscs PokemonStateCombinations) GetNotOK(pokemon *Pokemon) PokemonStateCombinations {
// 	result := make(PokemonStateCombinations, 0, len(pscs))
// 	for _, psc := range pscs {
// 		if !psc.OK(pokemon) {
// 			result = append(result, psc)
// 		}
// 	}
// 	return result
// }

// func (pscs PokemonStateCombinations) DiffCalc(pokemon, nextPokemon *Pokemon) PokemonStateCombinations {
// 	//一つ前の状態(pokemon)で、満たしていない組み合わせを取り出す
// 	pscs = pscs.GetNotOK(pokemon)
// 	//次の状態(nextPokemon)で、満たしている組み合わせを取り出す
// 	pscs = pscs.GetOK(nextPokemon)
// 	return pscs
// }

// func (pscs *PokemonStateCombinations) MoveNameDiffCalc(moveNames MoveNames, pokemon Pokemon) map[MoveName]PokemonStateCombinations {
// 	moveNameDiffCalc := map[MoveName]PokemonStateCombinations{}
// 	nextPokemon := pokemon

// 	for _, moveName := range moveNames {
// 		_, ok := pokemon.Moveset[moveName]

// 		if ok {
// 			continue
// 		}

// 		moveset := pokemon.Moveset.Copy()
// 		moveset[moveName] = &PowerPoint{}
// 		nextPokemon.Moveset = moveset

// 		diffCalc := pscs.DiffCalc(&pokemon, &nextPokemon)
// 		moveNameDiffCalc[moveName] = diffCalc
// 	}
// 	return moveNameDiffCalc
// }

// func (pscs *PokemonStateCombinations) AbilityDiffCalc(abilities Abilities, pokemon Pokemon) map[Ability]PokemonStateCombinations {
// 	abilityDiffCalc := map[Ability]PokemonStateCombinations{}
// 	nextPokemon := pokemon

// 	for _, ability := range abilities {
// 		nextPokemon.Ability = ability
// 		diffCalc:= pscs.DiffCalc(&pokemon, &nextPokemon)
// 		abilityDiffCalc[ability] = diffCalc
// 	}
// 	return abilityDiffCalc
// }

// func (pscs *PokemonStateCombinations) ItemDiffCalc(items Items, pokemon Pokemon, team Team) map[Item]PokemonStateCombinations {
// 	itemDiffCalc := map[Item]PokemonStateCombinations{}
// 	nextPokemon := pokemon

// 	for _, item := range items {
// 		if team.Items().In(item) {
// 			continue
// 		}

// 		nextPokemon.Item = item
// 		diffCalc := pscs.DiffCalc(&pokemon, &nextPokemon)
// 		itemDiffCalc[item] = diffCalc
// 	}
// 	return itemDiffCalc
// }

// func (pscs *PokemonStateCombinations) NatureDiffCalc(natures Natures, pokemon Pokemon) map[Nature]PokemonStateCombinations {
// 	natureDiffCalc := map[Nature]PokemonStateCombinations{}
// 	nextPokemon := pokemon
// 	for _, nature := range natures {
// 		nextPokemon.Nature = nature
// 		diffCalc := pscs.DiffCalc(&pokemon, &nextPokemon)
// 		natureDiffCalc[nature] = diffCalc
// 	}
// 	return natureDiffCalc
// }

// type PokemonStateCombinationModel struct {
// 	X *PokemonStateCombination
// 	Policy float64
// 	Value float64
// }

// type PokemonStateCombinationModels []*PokemonStateCombinationModel

// func NewPokemonStateCombinationModels(pscs PokemonStateCombinations, random *rand.Rand) PokemonStateCombinationModels {
// 	length := len(pscs)
// 	result := make(PokemonStateCombinationModels, length)

// 	for i := 0; i < length; i++ {
// 		policy, err := omw.RandomFloat64(0.001, 1.0, random)
// 		if err != nil {
// 			panic(err)
// 		}

// 		value, err := omw.RandomFloat64(1.0, 16.0, random)
// 		if err != nil {
// 			panic(err)
// 		}

// 		result[i] = &PokemonStateCombinationModel{X:pscs[i], Policy:policy, Value:value}
// 	}
// 	return result
// }

// func (pscms PokemonStateCombinationModels) WriteJson(pokeName PokeName) error {
// 	filePath := PSCMS_PATH + string(pokeName) + ".json"

// 	file, err := json.MarshalIndent(pscms, "", " ")
// 	if err != nil {
// 		return err
// 	}
// 	return ioutil.WriteFile(filePath, file, 0644)
// }

// func (pscms PokemonStateCombinationModels) GetOK(pokemon *Pokemon) PokemonStateCombinationModels {
// 	result := make(PokemonStateCombinationModels, 0, len(pscms))
// 	for _, pscm := range pscms {
// 		if pscm.X.OK(pokemon) {
// 			result = append(result, pscm)
// 		}
// 	}
// 	return result
// }

// type MultiPokemonStateCombination map[PokeName]*PokemonStateCombination

// func (tc TeamCombination) Keys() PokeNames {
// 	result := make(PokeNames, 0, len(tc))
// 	for pokeName, _ := range tc {
// 		result = append(result, pokeName)
// 	}
// 	return result
// }

// func (tc TeamCombination) OK(team Team) bool {
// 	for pokeName, psc := range tc {
// 		pokemon, err := team.Find(pokeName)
// 		if err != nil {
// 			return false
// 		}
// 		if !psc.OK(&pokemon) {
// 			return false
// 		}
// 	}
// 	return true
// }

// type TeamCombinations []TeamCombination

// func NewTeamCombinations(pokeName1, pokeName2 PokeName) (TeamCombinations, error) {
// 	pokemon1BuildCommonKnowledge, err := LoadJsonPokemonBuildCommonKnowledge(pokeName1)
// 	if err != nil {
// 		return TeamCombinations{}, err
// 	}

// 	pokemon2BuildCommonKnowledge, err := LoadJsonPokemonBuildCommonKnowledge(pokeName2)
// 	if err != nil {
// 		return TeamCombinations{}, err
// 	}

// 	pbks := map[PokeName]*PokemonBuildCommonKnowledge{
// 		pokeName1: &pokemon1BuildCommonKnowledge,
// 		pokeName2: &pokemon2BuildCommonKnowledge,
// 	}

// 	result := make(TeamCombinations, 0, 51200)
// 	result = append(result, TeamCombination{pokeName1:&PokemonStateCombination{}, pokeName2:&PokemonStateCombination{}})
// 	allAbilities := map[PokeName]Abilities{pokeName1:POKEDEX[pokeName1].AllAbilities, pokeName2:POKEDEX[pokeName2].AllAbilities}


// 	get := func(pokeName1, pokeName2 PokeName) TeamCombinations {
// 		result := make(TeamCombinations, 0, 25600)

// 		combination2MoveNames, err := pbks[pokeName1].MoveNames.Combination(2)
// 		if err != nil {
// 			panic(err)
// 		}

// 		for _, moveName1 := range pbks[pokeName1].MoveNames {
// 			for _, moveName2 := range pbks[pokeName2].MoveNames {
// 				psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName1}}
// 				psc2 := PokemonStateCombination{MoveNames:MoveNames{moveName2}}
// 				result = append(result, TeamCombination{pokeName1:&psc1, pokeName2:&psc2})
// 			}
// 		}

// 		for _, moveName := range pbks[pokeName1].MoveNames {
// 			for _, item := range pbks[pokeName2].Items {
// 				psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName}}
// 				psc2 := PokemonStateCombination{Item:item}
// 				result = append(result, TeamCombination{pokeName1:&psc1, pokeName2:&psc2})
// 			}
// 		}
	
// 		for _, moveName := range pbks[pokeName1].MoveNames {
// 			for _, ability := range allAbilities[pokeName2] {
// 				psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName}}
// 				psc2 := PokemonStateCombination{Ability:ability}
// 				result = append(result, TeamCombination{pokeName1:&psc1, pokeName2:&psc2})
// 			}
// 		}
	
// 		for _, moveName := range pbks[pokeName1].MoveNames {
// 			for _, nature := range pbks[pokeName2].Natures {
// 				psc1 := PokemonStateCombination{MoveNames:MoveNames{moveName}}
// 				psc2 := PokemonStateCombination{Nature:nature}
// 				result = append(result, TeamCombination{pokeName1:&psc1, pokeName2:&psc2})
// 			}
// 		}

// 		for _, moveNames1 := range combination2MoveNames {
// 			for _, moveName2 := range pbks[pokeName2].MoveNames {
// 				psc1 := PokemonStateCombination{MoveNames:moveNames1}
// 				psc2 := PokemonStateCombination{MoveNames:MoveNames{moveName2}}
// 				result = append(result, TeamCombination{pokeName1:&psc1, pokeName2:&psc2})
// 			}
// 		}

// 		for _, moveNames := range combination2MoveNames {
// 			for _, item := range pbks[pokeName2].Items {
// 				psc1 := PokemonStateCombination{MoveNames:moveNames}
// 				psc2 := PokemonStateCombination{Item:item}
// 				result = append(result, TeamCombination{pokeName1:&psc1, pokeName2:&psc2})
// 			}
// 		}
	
// 		for _, moveNames := range combination2MoveNames {
// 			for _, ability := range allAbilities[pokeName2] {
// 				psc1 := PokemonStateCombination{MoveNames:moveNames}
// 				psc2 := PokemonStateCombination{Ability:ability}
// 				result = append(result, TeamCombination{pokeName1:&psc1, pokeName2:&psc2})
// 			}
// 		}
	
// 		for _, moveNames := range combination2MoveNames {
// 			for _, nature := range pbks[pokeName2].Natures {
// 				psc1 := PokemonStateCombination{MoveNames:moveNames}
// 				psc2 := PokemonStateCombination{Nature:nature}
// 				result = append(result, TeamCombination{pokeName1:&psc1, pokeName2:&psc2})
// 			}
// 		}

// 		for _, ability1 := range allAbilities[pokeName1] {
// 			for _, ability2 := range allAbilities[pokeName2] {
// 				psc1 := PokemonStateCombination{Ability:ability1}
// 				psc2 := PokemonStateCombination{Ability:ability2}
// 				result = append(result, TeamCombination{pokeName1:&psc1, pokeName2:&psc2})
// 			}
// 		}

// 		for _, item1 := range pbks[pokeName1].Items {
// 			for _, item2 := range pbks[pokeName2].Items {
// 				psc1 := PokemonStateCombination{Item:item1}
// 				psc2 := PokemonStateCombination{Item:item2}
// 				result = append(result, TeamCombination{pokeName1:&psc1, pokeName2:&psc2})
// 			}
// 		}

// 		for _, nature1 := range pbks[pokeName1].Natures {
// 			for _, nature2 := range pbks[pokeName2].Natures {
// 				psc1 := PokemonStateCombination{Nature:nature1}
// 				psc2 := PokemonStateCombination{Nature:nature2}
// 				result = append(result, TeamCombination{pokeName1:&psc1, pokeName2:&psc2})
// 			}
// 		}

// 		return result
// 	}

// 	result = append(result, get(pokeName1, pokeName2)...)
// 	result = append(result, get(pokeName2, pokeName1)...)
// 	return result, nil
// }

// func (tcs TeamCombinations) OKIndices(team Team) TeamCombinations {
// 	result := make(TeamCombinations, 0, len(tcs))
// 	for _, tc := range tcs {
// 		if tc.OK(team) {
// 			result = append(result, tc)
// 		}
// 	}
// 	return result
// }

// type TeamCombinationModel struct {
// 	X TeamCombination
// 	Value float64
// }

// func NewTeamCombinationModel(pokeName1, pokeName2 PokeName, random *rand.Rand) (TeamCombinationsModel, error) {
// 	tcs, err := NewTeamCombinations(pokeName1, pokeName2)
// 	if err != nil {
// 		return TeamCombinationsModel{}, err
// 	}

// 	length := len(tcs)
// 	values := make([]float64, length)

// 	for i := 0; i < length; i++ {
// 		v, err := omw.RandomFloat64(0.01, 1, random)
// 		if err != nil {
// 			panic(err)
// 		}
// 		values[i] = v
// 	}
// 	return TeamCombinationsModel{X:tcs, Values:values}, nil
// }

// func (tce TeamCombinationsModel) WriteJson() error {
// 	pokeNames := tce.X[0].Keys().Sort()	
// 	folderPath := TCE_PATH + string(pokeNames[0]) + "/"
// 	if _, err := os.Stat(folderPath); err != nil {
// 		return err
// 	}
// 	filePath := folderPath + string(pokeNames[1]) + ".json"

// 	file, err := json.MarshalIndent(tce, "", " ")
// 	if err != nil {
// 		return err
// 	}
// 	return ioutil.WriteFile(filePath, file, 0644)
// }