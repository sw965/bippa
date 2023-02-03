package bippa

// import (
// 	"encoding/json"
// 	"github.com/sw965/omw"
// 	"io/ioutil"
// 	"math/rand"
// 	"os"
// )

// type PokemonBuildCommonKnowledge struct {
// 	PokeName  PokeName
// 	MoveNames MoveNames
// 	Items     Items
// 	Natures   Natures
// }

// func LoadJsonPokemonBuildCommonKnowledge(pokeName PokeName) (PokemonBuildCommonKnowledge, error) {
// 	filePath := POKEMON_BUILD_COMMON_KNOWLEDGE_PATH + string(pokeName) + ".json"

// 	bytes, err := ioutil.ReadFile(filePath)
// 	if err != nil {
// 		return PokemonBuildCommonKnowledge{}, err
// 	}

// 	result := PokemonBuildCommonKnowledge{}
// 	if err := json.Unmarshal(bytes, &result); err != nil {
// 		return PokemonBuildCommonKnowledge{}, err
// 	}

// 	return result, nil
// }

// type PokemonBuildCommonKnowledgeList []PokemonBuildCommonKnowledge

// func (pbkl PokemonBuildCommonKnowledgeList) Access(indices []int) PokemonBuildCommonKnowledgeList {
// 	result := make(PokemonBuildCommonKnowledgeList, len(indices))
// 	for i, index := range indices {
// 		result[i] = pbkl[index]
// 	}
// 	return result
// }

// func (pbkl PokemonBuildCommonKnowledgeList) Permutation(r int) ([]PokemonBuildCommonKnowledgeList, error) {
// 	n := len(pbkl)
// 	permutationTotalNum := omw.PermutationTotalNum(n, r)
// 	permutationNumbers, err := omw.PermutationNumbers(n, r)
// 	if err != nil {
// 		return []PokemonBuildCommonKnowledgeList{}, err
// 	}
// 	result := make([]PokemonBuildCommonKnowledgeList, permutationTotalNum)

// 	for i, indices := range permutationNumbers {
// 		result[i] = pbkl.Access(indices)
// 	}
// 	return result, nil
// }

// type PokemonStateCombination struct {
// 	MoveNames MoveNames
// 	Gender    Gender
// 	Ability   Ability
// 	Item      Item
// 	Nature    Nature

// 	HPIndividuals Individuals
// 	AtkIndividuals Individuals
// 	DefIndividuals Individuals
// 	SpAtkIndividuals Individuals
// 	SpDefIndividuals Individuals
// 	SpeedIndividuals Individuals

// 	HPEfforts Efforts
// 	AtkEfforts Efforts
// 	DefEfforts Efforts
// 	SpAtkEfforts Efforts
// 	SpDefEfforts Efforts
// 	SpeedEfforts Efforts
// }

// func (psc *PokemonStateCombination) OK(pokemon *Pokemon) bool {
// 	if psc.Ability != "" {
// 		if psc.Ability != pokemon.Ability {
// 			return false
// 		}
// 	}

// 	if psc.Item != "" {
// 		if psc.Item != pokemon.Item {
// 			return false
// 		}
// 	}

// 	if len(psc.MoveNames) != 0 {
// 		for _, moveName := range psc.MoveNames {
// 			_, ok := pokemon.Moveset[moveName]
// 			if !ok {
// 				return false
// 			}
// 		}
// 	}

// 	if psc.Nature != "" {
// 		if psc.Nature != pokemon.Nature {
// 			return false
// 		}
// 	}

// 	if !omw.Contains(psc.HPIndividuals, pokemon.IndividualsState.HP) {
// 		return false
// 	}

// 	if !omw.Contains(psc.AtkIndividuals, pokemon.IndividualState.Atk) {
// 		return false
// 	}

// 	if !omw.Contains(psc.DefIndividuals, pokemon.IndividualState.Def) {
// 		return false
// 	}

// 	if !omw.Contains(psc.SpAtkIndividuals, pokemon.IndividualState.SpAtk) {
// 		return false
// 	}

// 	if !omw.Contains(psc.SpDefIndividuals, pokemon.IndividualsState.SpDef) {
// 		return false
// 	}

// 	if !omw.Contains(psc.SpeedIndividuals, pokemon.IndividualsState.Speed) {
// 		return false
// 	}

// 	if !omw.Contains(psc.HPEfforts, pokemon.IndividualsState.HP) {
// 		return false
// 	}

// 	if !omw.Contains(psc.HP)

// 	isIndividualOK := func(individual Individual, lowerAndUpperLimitIndividuals Individuals) bool {
// 		if len(lowerAndUpperLimitIndividuals) == 0 {
// 			return true
// 		}
// 		return (individual >= lowerAndUpperLimitIndividuals[0]) && (individual < lowerAndUpperLimitIndividuals[1])
// 	}

// 	if !isIndividualOK(pokemon.IndividualState.HP, psc.LowerAndUpperLimitHPIndividuals) {
// 		return false
// 	}

// 	if !isIndividualOK(pokemon.IndividualState.Atk, psc.LowerAndUpperLimitAtkIndividuals) {
// 		return false
// 	}

// 	if !isIndividualOK(pokemon.IndividualState.Def, psc.LowerAndUpperLimitDefIndividuals) {
// 		return false
// 	}

// 	if !isIndividualOK(pokemon.IndividualState.SpAtk, psc.LowerAndUpperLimitSpAtkIndividuals) {
// 		return false
// 	}

// 	if !isIndividualOK(pokemon.IndividualState.SpDef, psc.LowerAndUpperLimitSpDefIndividuals) {
// 		return false
// 	}

// 	if !isIndividualOK(pokemon.IndividualState.Speed, psc.LowerAndUpperLimitSpeedIndividuals) {
// 		return false
// 	}

// 	isEffortOK := func(effort Effort, lowerAndUpperLimitEfforts Efforts) bool {
// 		if len(lowerAndUpperLimitEfforts) == 0 {
// 			return true
// 		}
// 		return (effort >= lowerAndUpperLimitEfforts[0]) && (effort < lowerAndUpperLimitEfforts[1])
// 	}

// 	if !isEffortOK(pokemon.EffortState.HP, psc.LowerAndUpperLimitHPEfforts) {
// 		return false
// 	}

// 	if !isEffortOK(pokemon.EffortState.Atk, psc.LowerAndUpperLimitAtkEfforts) {
// 		return false
// 	}

// 	if !isEffortOK(pokemon.EffortState.Def, psc.LowerAndUpperLimitDefEfforts) {
// 		return false
// 	}

// 	if !isEffortOK(pokemon.EffortState.SpAtk, psc.LowerAndUpperLimitSpAtkEfforts) {
// 		return false
// 	}

// 	if !isEffortOK(pokemon.EffortState.SpDef, psc.LowerAndUpperLimitSpDefEfforts) {
// 		return false
// 	}

// 	if !isEffortOK(pokemon.EffortState.Speed, psc.LowerAndUpperLimitSpeedEfforts) {
// 		return false
// 	}
// 	return true
// }

// type PokemonStateCombinations []PokemonStateCombination

// type PokemonStateCombinationModel struct {
// 	X      PokemonStateCombination
// 	Value  float64
// 	Number int
// }

// type PokemonStateCombinationModels []*PokemonStateCombinationModel

// func NewPokemonStateCombinationModels(pscs PokemonStateCombinations, random *rand.Rand) PokemonStateCombinationModels {
// 	result := make(PokemonStateCombinationModels, len(pscs))
// 	for i, psc := range pscs {
// 		value, err := omw.RandomFloat64(0, 16.0, random)
// 		if err != nil {
// 			panic(err)
// 		}
// 		pscm := PokemonStateCombinationModel{X: psc, Value: value}
// 		result[i] = &pscm
// 	}
// 	result.InitNumber()
// 	return result
// }

// func (pscms PokemonStateCombinationModels) InitNumber() {
// 	length := len(pscms)
// 	for i := 0; i < length; i++ {
// 		pscms[i].Number = i
// 	}
// }

// func (pscms PokemonStateCombinationModels) WriteJson(pokeName PokeName, modelsName string) error {
// 	folderDirectory := POKEMON_STATE_COMBINATION_MODELS_PATH + string(pokeName) + "/"

// 	if _, err := os.Stat(folderDirectory); err != nil {
// 		if mkErr := os.Mkdir(folderDirectory, os.ModePerm); mkErr != nil {
// 			return mkErr
// 		}
// 	}

// 	fullPath := folderDirectory + modelsName
// 	file, err := json.MarshalIndent(pscms, "", " ")
// 	if err != nil {
// 		return err
// 	}
// 	return ioutil.WriteFile(fullPath, file, 0644)
// }

// type MultiplePokemonStateCombination map[PokeName]*PokemonStateCombination

// func (mpsc MultiplePokemonStateCombination) OK(pokemons ...*Pokemon) bool {
// 	if len(mpsc) > len(pokemons) {
// 		return false
// 	}

// 	for _, pokemon := range pokemons {
// 		psc, ok := mpsc[pokemon.Name]
// 		if !ok {
// 			return false
// 		}
// 		if !psc.OK(pokemon) {
// 			return false
// 		}
// 	}
// 	return true
// }

// type MultiplePokemonStateCombinations []MultiplePokemonStateCombination

// type MultiplePokemonStateCombinationModel struct {
// 	X      MultiplePokemonStateCombination
// 	Value  float64
// 	Number int
// }

// type MultiplePokemonStateCombinationModels []MultiplePokemonStateCombinationModel

// func NewMultiplePokemonStateCombinationModels(mpscs MultiplePokemonStateCombinations, random *rand.Rand) MultiplePokemonStateCombinationModels {
// 	length := len(mpscs)
// 	result := make(MultiplePokemonStateCombinationModels, 0, length)
// 	for i, mpsc := range mpscs {
// 		value, err := omw.RandomFloat64(0.01, 16.0, random)
// 		if err != nil {
// 			panic(err)
// 		}
// 		result = append(result, MultiplePokemonStateCombinationModel{X: mpsc, Value: value, Number: i})
// 	}
// 	return result
// }

// func (mpscms MultiplePokemonStateCombinationModels) WriteJson(modelsName string, pokeNames ...PokeName) error {
// 	lastIndex := len(pokeNames) - 1
// 	folderDirectory := MULTIPLE_POKEMON_STATE_COMBINATION_MODELS_PATH + modelsName + "/"

// 	for _, pokeName := range pokeNames[:lastIndex] {
// 		folderDirectory += string(pokeName) + "/"
// 	}

// 	if _, err := os.Stat(folderDirectory); err != nil {
// 		if mkErr := os.MkdirAll(folderDirectory, os.ModePerm); mkErr != nil {
// 			return mkErr
// 		}
// 	}

// 	fullPath := folderDirectory + string(pokeNames[lastIndex]) + ".json"

// 	file, err := json.MarshalIndent(mpscms, "", " ")
// 	if err != nil {
// 		return err
// 	}
// 	return ioutil.WriteFile(fullPath, file, 0644)
// }
