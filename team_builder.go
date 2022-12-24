package bippa

import (
	"fmt"
	"github.com/sw965/omw"
	"math/rand"
	"encoding/json"
	"io/ioutil"
)

type PokemonBuildCommonKnowledge struct {
	MoveNames MoveNames
	Items     Items
	Natures   Natures
}

func ReadJsonPokemonBuildCommonKnowledge(filePath string) PokemonBuildCommonKnowledge {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	result := PokemonBuildCommonKnowledge{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		panic(err)
	}
	return result
}

type PokemonStateCombinationFeature struct {
	MoveNames  MoveNames
	Ability    Ability
	Item       Item
	Nature     Nature
	InitIndex int
}

func (pscf *PokemonStateCombinationFeature) OK(pokemon *Pokemon) bool {
	if pscf.Ability != "" {
		if pscf.Ability != pokemon.Ability {
			return false
		}
	}

	if pscf.Item != "" {
		if pscf.Item != pokemon.Item {
			return false
		}
	}

	if len(pscf.MoveNames) != 0 {
		for _, moveName := range pscf.MoveNames {
			_, ok := pokemon.Moveset[moveName]
			if !ok {
				return false
			}
		}
	}

	if pscf.Nature != "" {
		if pscf.Nature != pokemon.Nature {
			return false
		}
	}
	return true
}

func (pscf1 *PokemonStateCombinationFeature) Equal(pscf2 *PokemonStateCombinationFeature) bool {
	if pscf1.Ability != pscf2.Ability {
		return false
	}

	if pscf1.Item != pscf2.Item {
		return false
	}

	sortedMoveNames1 := pscf1.MoveNames.Sort()
	sortedMoveNames2 := pscf2.MoveNames.Sort()

	if !sortedMoveNames1.Equal(sortedMoveNames2) {
		return false
	}

	if pscf1.Nature != pscf2.Nature {
		return false
	}
	return true
}

type PokemonStateCombinationFeatures []*PokemonStateCombinationFeature

func NewPokemonStateCombinationFeatures(pokeName PokeName, pbCommonK *PokemonBuildCommonKnowledge) PokemonStateCombinationFeatures {
	pokeData := POKEDEX[pokeName]

	allAbilities := pokeData.AllAbilities
	learnset := pokeData.Learnset

	pbCommonKMoveNames := pbCommonK.MoveNames
	pbCommonKItems := pbCommonK.Items
	pbCommonKNatures := pbCommonK.Natures

	result := make(PokemonStateCombinationFeatures, 0, 6400)

	for _, moveName := range learnset {
		result = append(result, &PokemonStateCombinationFeature{MoveNames: MoveNames{moveName}})
	}

	for _, ability := range allAbilities {
		result = append(result, &PokemonStateCombinationFeature{Ability: ability})
	}

	for _, item := range ALL_ITEMS {
		result = append(result, &PokemonStateCombinationFeature{Item: item})
	}

	for _, nature := range ALL_NATURES {
		result = append(result, &PokemonStateCombinationFeature{Nature: nature})
	}

	combination2MoveNames, err := pbCommonKMoveNames .Combination(2)
	if err != nil {
		panic(err)
	}

	for _, moveNames := range combination2MoveNames {
		result = append(result, &PokemonStateCombinationFeature{MoveNames:moveNames})
	}

	for _, moveName := range pbCommonKMoveNames  {
		for _, ability := range allAbilities {
			result = append(result, &PokemonStateCombinationFeature{MoveNames:MoveNames{moveName}, Ability:ability})
		}
	}

	for _, moveName := range pbCommonKMoveNames  {
		for _, item := range pbCommonKItems  {
			result = append(result, &PokemonStateCombinationFeature{MoveNames:MoveNames{moveName}, Item:item})
		}
	}

	for _, moveName := range pbCommonKMoveNames  {
		for _, nature := range pbCommonKNatures {
			result = append(result, &PokemonStateCombinationFeature{MoveNames:MoveNames{moveName}, Nature:nature})
		}
	}

	combination3MoveNames, err := pbCommonKMoveNames .Combination(3)
	if err != nil {
		panic(err)
	}

	for _, moveNames := range combination3MoveNames {
		result = append(result, &PokemonStateCombinationFeature{MoveNames: moveNames})
	}

	for _, moveNames := range combination2MoveNames {
		for _, ability := range allAbilities {
			result = append(result, &PokemonStateCombinationFeature{MoveNames:moveNames, Ability:ability})
		}
	}

	for _, moveNames := range combination2MoveNames {
		for _, item := range pbCommonKItems  {
			result = append(result, &PokemonStateCombinationFeature{MoveNames:moveNames, Item:item})
		}
	}

	for _, moveNames := range combination2MoveNames {
		for _, nature := range pbCommonKNatures {
			result = append(result, &PokemonStateCombinationFeature{MoveNames:moveNames, Nature:nature})
		}
	}

	return result
}

func (pscfs PokemonStateCombinationFeatures) Init(pokeName PokeName) {
	for i := 0; i < len(pscfs); i++ {
		pscfs[i].InitIndex = i
	}
}

func (pscfs PokemonStateCombinationFeatures) InitIndices() []int {
	result := make([]int, 0, len(pscfs))
	for i, pscf := range pscfs {
		result[i] = pscf.InitIndex
	}
	return result
}

func (pscfs PokemonStateCombinationFeatures) GetOKs(pokemon *Pokemon) PokemonStateCombinationFeatures {
	result := make(PokemonStateCombinationFeatures, 0, len(pscfs))
	for _, pscf := range pscfs {
		if pscf.OK(pokemon) {
			result = append(result, pscf)
		}
	}
	return result
}

func (pscfs PokemonStateCombinationFeatures) GetNotOKs(pokemon *Pokemon) PokemonStateCombinationFeatures {
	result := make(PokemonStateCombinationFeatures, 0, len(pscfs))
	for _, pscf := range pscfs {
		if !pscf.OK(pokemon) {
			result = append(result, pscf)
		}
	}
	return result
}

func (pscfs PokemonStateCombinationFeatures) Policy(pokemon, nextPokemon *Pokemon) PokemonStateCombinationFeatures {
	//一つ前の状態(pokemon)で、該当する特徴を満たしていない特徴量を取り出す
	pscfs = pscfs.GetNotOKs(pokemon)

	//次の状態(nextPokemon)で、該当する特徴を満たしている特徴量を取り出す
	return pscfs.GetOKs(nextPokemon)
}
                             
type PokemonStateCombinationEvaluator struct {
	X PokemonStateCombinationFeatures
	Policies []float64
	Values []float64
}

func(psce *PokemonStateCombinationEvaluator) Eval() float64 {
	result := 0.0
	for _, v := range psce.Values {
		result += v
	}
	return result
}

func (psce *PokemonStateCombinationEvaluator) MoveNameWithPolicyData(moveNames MoveNames, pokemon Pokemon, team Team) map[MoveName]PSCEPolicyData {
	moveNameWithPolicyData := map[MoveName]PSCEPolicyData{}
	nextPokemon := pokemon

	for _, moveName := range moveNames {
		_, ok := pokemon.Moveset[moveName]

		if ok {
			continue
		}

		moveset := pokemon.Moveset.Copy()
		moveset[moveName] = &PowerPoint{}
		nextPokemon.Moveset = moveset

		policyPSCFs := psce.X.Policy(&pokemon, &nextPokemon)
		initIndices := policyPSCFs.InitIndices()

		policies := make([]float64, len(initIndices))
		for i, index := range initIndices {
			policies[i] = psce.Policies[index]
		}

		sumPolicy := 0.0
		for _, v := range policies {
			sumPolicy += v
		}

		policyData := PSCEPolicyData{X:policyPSCFs, Sum:sumPolicy}
		moveNameWithPolicyData[moveName] = policyData
	}
	return moveNameWithPolicyData
}

func (psce *PokemonStateCombinationEvaluator) AbilityWithPolicyData(abilities Abilities, pokemon Pokemon, team Team) (map[Ability]PSCEPolicyData) {
	abilityWithPolicyData := map[Ability]PSCEPolicyData{}
	nextPokemon := pokemon

	for _, ability := range abilities {
		nextPokemon.Ability = ability

		policyPSCFs := psce.X.Policy(&pokemon, &nextPokemon)
		initIndices := policyPSCFs.InitIndices()

		policies := make([]float64, len(initIndices))
		for i, index := range initIndices {
			policies[i] = psce.Policies[index]
		}

		sumPolicy := 0.0
		for _, v := range policies {
			sumPolicy += v
		}

		policyData := PSCEPolicyData{X:policyPSCFs, Sum:sumPolicy}
		abilityWithPolicyData[ability] = policyData
	}
	return abilityWithPolicyData
}

func (psce *PokemonStateCombinationEvaluator) ItemWithPolicyData(items Items, pokemon Pokemon, team Team) map[Item]PSCEPolicyData {
	itemWithPolicyData := map[Item]PSCEPolicyData{}
	nextPokemon := pokemon

	for _, item := range items {
		if team.Items().In(item) {
			continue
		}

		nextPokemon.Item = item

		policyPSCFs := psce.X.Policy(&pokemon, &nextPokemon)
		initIndices := policyPSCFs.InitIndices()

		policies := make([]float64, len(initIndices))
		for i, index := range initIndices {
			policies[i] = psce.Policies[index]
		}

		sumPolicy := 0.0
		for _, v := range policies {
			sumPolicy += v
		}

		policyData := PSCEPolicyData{X:policyPSCFs, Sum:sumPolicy}
		itemWithPolicyData[item] = policyData
	}
	return itemWithPolicyData
}

func (psce *PokemonStateCombinationEvaluator) NatureWithPolicyData(natures Natures, pokemon Pokemon, team Team) map[Nature]PSCEPolicyData {
	natureWithPolicyData := map[Nature]PSCEPolicyData{}
	nextPokemon := pokemon

	for _, nature := range natures {
		nextPokemon.Nature = nature
		
		policyPSCFs := psce.X.Policy(&pokemon, &nextPokemon)
		initIndices := policyPSCFs.InitIndices()		

		policies := make([]float64, len(initIndices))
		for i, index := range initIndices {
			policies[i] = psce.Policies[index]
		}

		sumPolicy := 0.0
		for _, v := range policies {
			sumPolicy += v
		}

		policyData := PSCEPolicyData{X:policyPSCFs, Sum:sumPolicy}
		natureWithPolicyData[nature] = policyData
	}
	return natureWithPolicyData
}

func (psce *PokemonStateCombinationEvaluator) BuildMoveset(pokeName PokeName, moveNames MoveNames, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Moveset, MoveName, map[MoveName]PSCEPolicyData, error) {
	moveset := pokemon.Moveset.Copy()
	movesetLength := len(pokemon.Moveset)
	learnsetLength := len(POKEDEX[pokeName].Learnset)

	if movesetLength == MAX_MOVESET_LENGTH || movesetLength == learnsetLength {
		return pokemon.Moveset, "", map[MoveName]PSCEPolicyData{}, nil
	}

	moveNameWithPolicyData := psce.MoveNameWithPolicyData(moveNames, pokemon, team)

	if len(moveNameWithPolicyData) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v pokemon.Moveset.Keys() = %v の状態で、次の技の組み合わせが見つからなかった",
			pokemon.Name, pokemon.Moveset.Keys(),
		)
		return Moveset{}, "", map[MoveName]PSCEPolicyData{}, fmt.Errorf(errMsg)
	}

	length := len(moveNameWithPolicyData)
	selectableMoveNames := make(MoveNames, 0, length)
	sumPolicies := make([]float64, 0, length)

	for k, v := range moveNameWithPolicyData {
		selectableMoveNames = append(selectableMoveNames, k)
		sumPolicies = append(sumPolicies, v.Sum)
	}

	finalPoliciesY := finalPolicies(sumPolicies)
	selectIndex := omw.RandomIntWithWeight(finalPoliciesY, random)
	selectMoveName := selectableMoveNames[selectIndex]

	powerPoint := NewPowerPoint(MOVEDEX[selectMoveName].BasePP, MAX_POINT_UP)
	moveset[selectMoveName] = &powerPoint

	return moveset, selectMoveName, moveNameWithPolicyData, nil
}

func (psce *PokemonStateCombinationEvaluator) BuildAbility(abilities Abilities, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Ability, map[Ability]PSCEPolicyData, error) {
	if pokemon.Ability != "" {
		return pokemon.Ability, map[Ability]PSCEPolicyData{}, nil
	}

	abilityWithPolicyData := psce.AbilityWithPolicyData(abilities, pokemon, team)

	if len(abilityWithPolicyData) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次の特性の組み合わせが見つからなかった", pokemon.Name)
		return "", map[Ability]PSCEPolicyData{}, fmt.Errorf(errMsg)
	}

	length := len(abilityWithPolicyData)
	selectableAbilities := make(Abilities, 0, length)
	sumPolicies := make([]float64, 0, length)

	for k, v := range abilityWithPolicyData {
		selectableAbilities = append(selectableAbilities, k)
		sumPolicies = append(sumPolicies, v.Sum)
	}

	finalPoliciesY := finalPolicies(sumPolicies)
	selectIndex := omw.RandomIntWithWeight(finalPoliciesY, random)
	selectAbility := selectableAbilities[selectIndex]

	return selectAbility, abilityWithPolicyData, nil
}

func (psce *PokemonStateCombinationEvaluator) BuildItem(items Items, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Item, map[Item]PSCEPolicyData, error) {
	if pokemon.Item != "" {
		return pokemon.Item, map[Item]PSCEPolicyData{}, nil
	}

	itemWithPolicyData := psce.ItemWithPolicyData(items, pokemon, team)

	if len(itemWithPolicyData) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次のアイテムの組み合わせが見つからなかった", pokemon.Name)
		return "", map[Item]PSCEPolicyData{}, fmt.Errorf(errMsg)
	}

	length := len(itemWithPolicyData)
	selectableItems := make(Items, 0, length)
	sumPolicies := make([]float64, 0, length)

	for k, v := range itemWithPolicyData {
		selectableItems = append(selectableItems, k)
		sumPolicies = append(sumPolicies, v.Sum)
	}

	finalPoliciesY := finalPolicies(sumPolicies)
	selectIndex := omw.RandomIntWithWeight(finalPoliciesY, random)
	selectItem := selectableItems[selectIndex]

	return selectItem, itemWithPolicyData, nil
}

func (psce *PokemonStateCombinationEvaluator) BuildNature(natures Natures, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Nature, map[Nature]PSCEPolicyData, error) {
	if pokemon.Nature != "" {
		return pokemon.Nature, map[Nature]PSCEPolicyData{}, nil
	}

	natureWithPolicyData := psce.NatureWithPolicyData(natures, pokemon, team)

	if len(natureWithPolicyData) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次の性格の組み合わせが見つからなかった", pokemon.Name)
		return "", map[Nature]PSCEPolicyData{}, fmt.Errorf(errMsg)
	}

	length := len(natureWithPolicyData)
	selectableNatures := make(Natures, 0, length)
	sumPolicies := make([]float64, 0, length)

	for k, v := range natureWithPolicyData {
		selectableNatures = append(selectableNatures, k)
		sumPolicies = append(sumPolicies, v.Sum)
	}

	finalPoliciesY := finalPolicies(sumPolicies)
	selectIndex := omw.RandomIntWithWeight(finalPoliciesY, random)
	selectNature := selectableNatures[selectIndex]

	return selectNature, natureWithPolicyData, nil	
}

func (psce *PokemonStateCombinationEvaluator) Run(pokemon Pokemon, team Team, pbCommonK *PokemonBuildCommonKnowledge, finalPolicies func([]float64) []float64, random *rand.Rand) (Pokemon, PSCERunHistory, error) {
	if pokemon.Name == "" {
		return Pokemon{}, PSCERunHistory{}, fmt.Errorf("pokemon.Name が 空の状態で、PokemonStateCombinationEvaluator.Run は 実行出来ない")
	}

	psceRunHistory := PSCERunHistory{
		SelectMoveNames:make(MoveNames, 0, MAX_MOVESET_LENGTH),
		MoveNameWithPolicyDataList:make([]map[MoveName]PSCEPolicyData, 0, MAX_MOVESET_LENGTH),
	}

	pokeData := POKEDEX[pokemon.Name]
	pbCommonKMoveNames := pbCommonK.MoveNames
	allAbilities := pokeData.AllAbilities
	pbCommonKItems := pbCommonK.Items
	pbCommonKNatures := pbCommonK.Natures

	ability, abilityWithPolicyData, err := psce.BuildAbility(allAbilities, pokemon, team, finalPolicies, random)
	if err != nil {
		return Pokemon{}, PSCERunHistory{}, err
	}

	pokemon.Ability = ability
	psceRunHistory.SelectAbility = ability
	psceRunHistory.AbilityWithPolicyData = abilityWithPolicyData

	item, itemWithPolicyData, err := psce.BuildItem(pbCommonKItems, pokemon, team, finalPolicies, random)
	if err != nil {
		return Pokemon{}, PSCERunHistory{}, err
	}

	pokemon.Item = item
	psceRunHistory.SelectItem = item
	psceRunHistory.ItemWithPolicyData = itemWithPolicyData

	for i := 0; i < MAX_MOVESET_LENGTH; i++ {
		moveset, selectMoveName, moveNameWithPolicyData, err := psce.BuildMoveset(pokemon.Name, pbCommonKMoveNames, pokemon, team, finalPolicies, random)
		if err != nil {
			return Pokemon{}, PSCERunHistory{}, err
		}

		if selectMoveName == "" {
			continue
		}

		pokemon.Moveset = moveset
		psceRunHistory.SelectMoveNames = append(psceRunHistory.SelectMoveNames, selectMoveName)
		psceRunHistory.MoveNameWithPolicyDataList = append(psceRunHistory.MoveNameWithPolicyDataList, moveNameWithPolicyData)
	}

	nature, natureWithPolicyData, err := psce.BuildNature(pbCommonKNatures, pokemon, team, finalPolicies, random)
	if err != nil {
		return Pokemon{}, PSCERunHistory{}, err
	}

	pokemon.Nature = nature
	psceRunHistory.SelectNature = nature
	psceRunHistory.NatureWithPolicyData = natureWithPolicyData

	return pokemon, psceRunHistory, nil
}

type PSCEPolicyData struct {
	X PokemonStateCombinationFeatures
	Sum float64
}

type PSCERunHistory struct {
	SelectMoveNames MoveNames
	SelectAbility Ability
	SelectItem Item
	SelectNature Nature

	MoveNameWithPolicyDataList []map[MoveName]PSCEPolicyData
	AbilityWithPolicyData map[Ability]PSCEPolicyData
	ItemWithPolicyData map[Item]PSCEPolicyData
	NatureWithPolicyData map[Nature]PSCEPolicyData
}

type TeamCombinationFeature map[PokeName]*PokemonStateCombinationFeature

func (tcf TeamCombinationFeature) OK(team Team) bool {
	for k, v := range tcf {
		pokemon, err := team.Find(k)
		if err != nil {
			return false
		}
		if !v.OK(&pokemon) {
			return false
		}
	}
	return true
}

type TeamCombinationFeatures []TeamCombinationFeature

func NewTeamCombinationFeatures(pokeName1, pokeName2 PokeName, pbCommonKs map[PokeName]*PokemonBuildCommonKnowledge) TeamCombinationFeatures {
	result := make(TeamCombinationFeatures, 0, 51200)
	allAbilities := map[PokeName]Abilities{pokeName1:POKEDEX[pokeName1].AllAbilities, pokeName2:POKEDEX[pokeName2].AllAbilities}

	get := func(pokeName1, pokeName2 PokeName) TeamCombinationFeatures {
		result := make(TeamCombinationFeatures, 0, 25600)

		combination2MoveNames, err := pbCommonKs[pokeName1].MoveNames.Combination(2)
		if err != nil {
			panic(err)
		}

		for _, moveName1 := range pbCommonKs[pokeName1].MoveNames {
			for _, moveName2 := range pbCommonKs[pokeName2].MoveNames {
				pscf1 := PokemonStateCombinationFeature{MoveNames:MoveNames{moveName1}}
				pscf2 := PokemonStateCombinationFeature{MoveNames:MoveNames{moveName2}}
				result = append(result, TeamCombinationFeature{pokeName1:&pscf1, pokeName2:&pscf2})
			}
		}

		for _, moveName := range pbCommonKs[pokeName1].MoveNames {
			for _, item := range pbCommonKs[pokeName2].Items {
				pscf1 := PokemonStateCombinationFeature{MoveNames:MoveNames{moveName}}
				pscf2 := PokemonStateCombinationFeature{Item:item}
				result = append(result, TeamCombinationFeature{pokeName1:&pscf1, pokeName2:&pscf2})
			}
		}
	
		for _, moveName := range pbCommonKs[pokeName1].MoveNames {
			for _, ability := range allAbilities[pokeName2] {
				pscf1 := PokemonStateCombinationFeature{MoveNames:MoveNames{moveName}}
				pscf2 := PokemonStateCombinationFeature{Ability:ability}
				result = append(result, TeamCombinationFeature{pokeName1:&pscf1, pokeName2:&pscf2})
			}
		}
	
		for _, moveName := range pbCommonKs[pokeName1].MoveNames {
			for _, nature := range pbCommonKs[pokeName2].Natures {
				pscf1 := PokemonStateCombinationFeature{MoveNames:MoveNames{moveName}}
				pscf2 := PokemonStateCombinationFeature{Nature:nature}
				result = append(result, TeamCombinationFeature{pokeName1:&pscf1, pokeName2:&pscf2})
			}
		}

		for _, moveNames1 := range combination2MoveNames {
			for _, moveName2 := range pbCommonKs[pokeName2].MoveNames {
				pscf1 := PokemonStateCombinationFeature{MoveNames:moveNames1}
				pscf2 := PokemonStateCombinationFeature{MoveNames:MoveNames{moveName2}}
				result = append(result, TeamCombinationFeature{pokeName1:&pscf1, pokeName2:&pscf2})
			}
		}

		for _, moveNames := range combination2MoveNames {
			for _, item := range pbCommonKs[pokeName2].Items {
				pscf1 := PokemonStateCombinationFeature{MoveNames:moveNames}
				pscf2 := PokemonStateCombinationFeature{Item:item}
				result = append(result, TeamCombinationFeature{pokeName1:&pscf1, pokeName2:&pscf2})
			}
		}
	
		for _, moveNames := range combination2MoveNames {
			for _, ability := range allAbilities[pokeName2] {
				pscf1 := PokemonStateCombinationFeature{MoveNames:moveNames}
				pscf2 := PokemonStateCombinationFeature{Ability:ability}
				result = append(result, TeamCombinationFeature{pokeName1:&pscf1, pokeName2:&pscf2})
			}
		}
	
		for _, moveNames := range combination2MoveNames {
			for _, nature := range pbCommonKs[pokeName2].Natures {
				pscf1 := PokemonStateCombinationFeature{MoveNames:moveNames}
				pscf2 := PokemonStateCombinationFeature{Nature:nature}
				result = append(result, TeamCombinationFeature{pokeName1:&pscf1, pokeName2:&pscf2})
			}
		}

		for _, ability1 := range allAbilities[pokeName1] {
			for _, ability2 := range allAbilities[pokeName2] {
				pscf1 := PokemonStateCombinationFeature{Ability:ability1}
				pscf2 := PokemonStateCombinationFeature{Ability:ability2}
				result = append(result, TeamCombinationFeature{pokeName1:&pscf1, pokeName2:&pscf2})
			}
		}

		for _, item1 := range pbCommonKs[pokeName1].Items {
			for _, item2 := range pbCommonKs[pokeName2].Items {
				pscf1 := PokemonStateCombinationFeature{Item:item1}
				pscf2 := PokemonStateCombinationFeature{Item:item2}
				result = append(result, TeamCombinationFeature{pokeName1:&pscf1, pokeName2:&pscf2})
			}
		}

		for _, nature1 := range pbCommonKs[pokeName1].Natures {
			for _, nature2 := range pbCommonKs[pokeName2].Natures {
				pscf1 := PokemonStateCombinationFeature{Nature:nature1}
				pscf2 := PokemonStateCombinationFeature{Nature:nature2}
				result = append(result, TeamCombinationFeature{pokeName1:&pscf1, pokeName2:&pscf2})
			}
		}

		return result
	}

	result = append(result, get(pokeName1, pokeName2)...)
	result = append(result, get(pokeName2, pokeName1)...)
	return result
}

func (tcfs TeamCombinationFeatures) Init() {
	for i, tcf := range tcfs {
		for k, _ := range tcf {
			tcfs[i][k].InitIndex = i
		}
	}
}

func (tcfs TeamCombinationFeatures) UtilInitIndices() []int {
	result := make([]int, 0, len(tcfs))
	for _, tcf := range tcfs {
		for _, v := range tcf {
			result = append(result, v.InitIndex)
			break
		}
	}
	return result
}

func (tcfs TeamCombinationFeatures) GetOKs(team Team) TeamCombinationFeatures {
	result := make(TeamCombinationFeatures, 0, len(tcfs))
	for _, tcf := range tcfs {
		if tcf.OK(team) {
			result = append(result, tcf)
		}
	}
	return result
}

type TeamCombinationEvaluator struct {
	X TeamCombinationFeatures
	Values []float64	
}

func NewTeamCombinationEvaluator(pokeName1, pokeName2 PokeName, pbCommonKs map[PokeName]*PokemonBuildCommonKnowledge, random *rand.Rand) TeamCombinationEvaluator {
	tcfs := NewTeamCombinationFeatures(pokeName1, pokeName2, pbCommonKs)
	length := len(tcfs)
	values := make([]float64, length)
	for i := 0; i < length; i++ {
		v, err := omw.RandomFloat64(0.01, 1, random)
		if err != nil {
			panic(err)
		}
		values[i] = v
	}
	return TeamCombinationEvaluator{X:tcfs, Values:values}
}

func (tce TeamCombinationEvaluator) Eval(team Team) float64 {
	okXs := tce.X.GetOKs(team)
	initIndices := okXs.UtilInitIndices()
	result := 0.0

	for _, index := range initIndices {
		result += tce.Values[index]
	}
	return result
}