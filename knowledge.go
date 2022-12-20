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

type PokemonBuildCombinationKnowledge struct {
	MoveNames  MoveNames
	Ability    Ability
	Item       Item
	Nature     Nature

	CombinationNum int
	Policy         float64
	Value float64
}

func (pbCombK *PokemonBuildCombinationKnowledge) Init() {
	combinationNum := 0

	if pbCombK.Ability != "" {
		combinationNum += 1
	}

	if pbCombK.Item != "" {
		combinationNum += 1
	}

	combinationNum += len(pbCombK.MoveNames)

	if pbCombK.Nature != "" {
		combinationNum += 1
	}

	pbCombK.CombinationNum = combinationNum
	pbCombK.Policy = 0.5
	pbCombK.Value = 1.0
}

func (pbCombK *PokemonBuildCombinationKnowledge) All(pokemon *Pokemon, team Team) bool {
	if pbCombK.Ability != "" {
		if pbCombK.Ability != pokemon.Ability {
			return false
		}
	}

	if pbCombK.Item != "" {
		if pbCombK.Item != pokemon.Item {
			return false
		}
	}

	if len(pbCombK.MoveNames) != 0 {
		for _, moveName := range pbCombK.MoveNames {
			_, ok := pokemon.Moveset[moveName]
			if !ok {
				return false
			}
		}
	}

	if pbCombK.Nature != "" {
		if pbCombK.Nature != pokemon.Nature {
			return false
		}
	}
	return true
}

func (pbCombK1 *PokemonBuildCombinationKnowledge) NearlyEqual(pbCombK2 *PokemonBuildCombinationKnowledge) bool {
	if pbCombK1.Ability != pbCombK2.Ability {
		return false
	}

	if pbCombK1.Item != pbCombK2.Item {
		return false
	}

	sortedMoveNames1 := pbCombK1.MoveNames.Sort()
	sortedMoveNames2 := pbCombK2.MoveNames.Sort()

	if !sortedMoveNames1.Equal(sortedMoveNames2) {
		return false
	}

	if pbCombK1.Nature != pbCombK2.Nature {
		return false
	}
	return true
}

func (pbCombK1 *PokemonBuildCombinationKnowledge) IsInclusion(pbCombK2 *PokemonBuildCombinationKnowledge) bool {
	if pbCombK1.CombinationNum < pbCombK2.CombinationNum {
		return false
	}

	pbCombK2Ability := pbCombK2.Ability
	pbCombK2Item := pbCombK2.Item
	pbCombK2Nature := pbCombK2.Nature

	if pbCombK2Ability != "" {
		if pbCombK1.Ability != pbCombK2Ability {
			return false
		}
	}

	if pbCombK2Item != "" {
		if pbCombK1.Item != pbCombK2Item {
			return false
		}
	}

	if pbCombK2Nature != "" {
		if pbCombK1.Nature != pbCombK2Nature {
			return false
		}
	}

	if !pbCombK1.MoveNames.InAll(pbCombK2.MoveNames...) {
		return false
	}

	return true
}

type PokemonBuildCombinationKnowledgeList []*PokemonBuildCombinationKnowledge

func NewInitPokemonBuildCombinationKnowledgeList(pokeName PokeName, pbCommonK *PokemonBuildCommonKnowledge) PokemonBuildCombinationKnowledgeList {
	pokeData := POKEDEX[pokeName]

	allAbilities := pokeData.AllAbilities
	learnset := pokeData.Learnset

	pbCommonKMoveNames := pbCommonK.MoveNames
	pbCommonKItems := pbCommonK.Items
	pbCommonKNatures := pbCommonK.Natures

	result := make(PokemonBuildCombinationKnowledgeList, 0, 6400)

	for _, moveName := range learnset {
		result = append(result, &PokemonBuildCombinationKnowledge{MoveNames: MoveNames{moveName}})
	}

	for _, ability := range allAbilities {
		result = append(result, &PokemonBuildCombinationKnowledge{Ability: ability})
	}

	for _, item := range ALL_ITEMS {
		result = append(result, &PokemonBuildCombinationKnowledge{Item: item})
	}

	for _, nature := range ALL_NATURES {
		result = append(result, &PokemonBuildCombinationKnowledge{Nature: nature})
	}

	combination2MoveNames, err := pbCommonKMoveNames .Combination(2)
	if err != nil {
		panic(err)
	}

	for _, moveNames := range combination2MoveNames {
		result = append(result, &PokemonBuildCombinationKnowledge{MoveNames:moveNames})
	}

	for _, moveName := range pbCommonKMoveNames  {
		for _, ability := range allAbilities {
			result = append(result, &PokemonBuildCombinationKnowledge{MoveNames:MoveNames{moveName}, Ability:ability})
		}
	}

	for _, moveName := range pbCommonKMoveNames  {
		for _, item := range pbCommonKItems  {
			result = append(result, &PokemonBuildCombinationKnowledge{MoveNames:MoveNames{moveName}, Item:item})
		}
	}

	for _, moveName := range pbCommonKMoveNames  {
		for _, nature := range pbCommonKNatures {
			result = append(result, &PokemonBuildCombinationKnowledge{MoveNames:MoveNames{moveName}, Nature:nature})
		}
	}

	combination3MoveNames, err := pbCommonKMoveNames .Combination(3)
	if err != nil {
		panic(err)
	}

	for _, moveNames := range combination3MoveNames {
		result = append(result, &PokemonBuildCombinationKnowledge{MoveNames: moveNames})
	}

	for _, moveNames := range combination2MoveNames {
		for _, ability := range allAbilities {
			result = append(result, &PokemonBuildCombinationKnowledge{MoveNames:moveNames, Ability:ability})
		}
	}

	for _, moveNames := range combination2MoveNames {
		for _, item := range pbCommonKItems  {
			result = append(result, &PokemonBuildCombinationKnowledge{MoveNames:moveNames, Item:item})
		}
	}

	for _, moveNames := range combination2MoveNames {
		for _, nature := range pbCommonKNatures {
			result = append(result, &PokemonBuildCombinationKnowledge{MoveNames:moveNames, Nature:nature})
		}
	}

	return result
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) Init() {
	for i := 0; i < len(pbCombKList); i++ {
		pbCombKList[i].Init()
	}
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) Index(pbCombK *PokemonBuildCombinationKnowledge) int {
	for i, iPBCombK := range pbCombKList {
		if iPBCombK.NearlyEqual(pbCombK) {
			return i
		}
	}
	return -1
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) RandomChoice(random *rand.Rand) *PokemonBuildCombinationKnowledge {
	index := random.Intn(len(pbCombKList))
	return pbCombKList[index]
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) Policies() []float64 {
	result := make([]float64, len(pbCombKList))
	for i, pbCombK := range pbCombKList {
		result[i] = pbCombK.Policy
	}
	return result
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) MeanPolicy() float64 {
	return omw.SumFloat64(pbCombKList.Policies()...) / float64(len(pbCombKList))
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) SumValue() float64 {
	result := 0.0
	for _, pbCombK := range pbCombKList {
		result += pbCombK.Value
	}
	return result
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) MaxCombinationNum() int {
	result := pbCombKList[0].CombinationNum
	for _, pbCombK := range pbCombKList[1:] {
		maxCombinationNum := pbCombK.CombinationNum
		if maxCombinationNum > result {
			result = maxCombinationNum
		}
	}
	return result
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) MaxCombinationNumKnowledgeList() PokemonBuildCombinationKnowledgeList {
	maxCombinationNum := pbCombKList.MaxCombinationNum()
	result := make(PokemonBuildCombinationKnowledgeList, 0, len(pbCombKList))
	for _, pbCombK := range pbCombKList {
		if pbCombK.CombinationNum == maxCombinationNum {
			result = append(result, pbCombK)
		}
	}
	return result
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) Usable(pokemon *Pokemon, team Team) PokemonBuildCombinationKnowledgeList {
	result := make(PokemonBuildCombinationKnowledgeList, 0, len(pbCombKList))
	for _, pbCombK := range pbCombKList {
		if pbCombK.All(pokemon, team) {
			result = append(result, pbCombK)
		}
	}
	return result
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) NotUsing(pokemon *Pokemon, team Team) PokemonBuildCombinationKnowledgeList {
	result := make(PokemonBuildCombinationKnowledgeList, 0, len(pbCombKList))
	for _, pbCombK := range pbCombKList {
		if !pbCombK.All(pokemon, team) {
			result = append(result, pbCombK)
		}
	}
	return result
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) PolicyKnowledgeList(pokemon, nextPokemon *Pokemon, team Team) PokemonBuildCombinationKnowledgeList {
	//一つ前の状態(pokemon)で、活用不可能な知識を取り出す(差分を見る為に)
	pbCombKList = pbCombKList.NotUsing(pokemon, team)

	//次の状態(nextPokemon)において、活用可能な知識を取り出す。
	pbCombKList = pbCombKList.Usable(nextPokemon, team)

	//組み合わせ数が最も多い知識を取り出す
	return pbCombKList.MaxCombinationNumKnowledgeList()
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) MoveNameWithPolicyData(moveNames MoveNames, pokemon Pokemon, team Team) map[MoveName]PBCombKPolicyData {
	moveNameWithPolicyData := map[MoveName]PBCombKPolicyData{}
	nextPokemon := pokemon

	for _, moveName := range moveNames {
		_, ok := pokemon.Moveset[moveName]

		if ok {
			continue
		}

		moveset := pokemon.Moveset.Copy()
		moveset[moveName] = &PowerPoint{}
		nextPokemon.Moveset = moveset

		policyKnowledgeList := pbCombKList.PolicyKnowledgeList(&pokemon, &nextPokemon, team)
		meanPolicy := policyKnowledgeList.MeanPolicy()
		policyData := PBCombKPolicyData{Mean:meanPolicy, KnowledgeList:policyKnowledgeList}

		moveNameWithPolicyData[moveName] = policyData
	}
	return moveNameWithPolicyData
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) AbilityWithPolicyData(abilities Abilities, pokemon Pokemon, team Team) (map[Ability]PBCombKPolicyData) {
	abilityWithPolicyData:= map[Ability]PBCombKPolicyData{}
	nextPokemon := pokemon

	for _, ability := range abilities {
		nextPokemon.Ability = ability

		policyKnowledgeList := pbCombKList.PolicyKnowledgeList(&pokemon, &nextPokemon, team)
		meanPolicy := policyKnowledgeList.MeanPolicy()
		policyData := PBCombKPolicyData{Mean:meanPolicy, KnowledgeList:policyKnowledgeList}

		abilityWithPolicyData[ability] = policyData
	}
	return abilityWithPolicyData
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) ItemWithPolicyData(items Items, pokemon Pokemon, team Team) map[Item]PBCombKPolicyData {
	itemWithPolicyData := map[Item]PBCombKPolicyData{}
	nextPokemon := pokemon

	for _, item := range items {
		if team.Items().In(item) {
			continue
		}

		nextPokemon.Item = item

		policyKnowledgeList := pbCombKList.PolicyKnowledgeList(&pokemon, &nextPokemon, team)
		meanPolicy := policyKnowledgeList.MeanPolicy()
		policyData := PBCombKPolicyData{Mean:meanPolicy, KnowledgeList:policyKnowledgeList}

		itemWithPolicyData[item] = policyData
	}
	return itemWithPolicyData
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) NatureWithPolicyData(natures Natures, pokemon Pokemon, team Team) map[Nature]PBCombKPolicyData {
	natureWithPolicyData := map[Nature]PBCombKPolicyData{}
	nextPokemon := pokemon

	for _, nature := range natures {
		nextPokemon.Nature = nature

		policyKnowledgeList := pbCombKList.PolicyKnowledgeList(&pokemon, &nextPokemon, team)
		meanPolicy := policyKnowledgeList.MeanPolicy()
		policyData := PBCombKPolicyData{Mean:meanPolicy, KnowledgeList:policyKnowledgeList}

		natureWithPolicyData[nature] = policyData
	}
	return natureWithPolicyData
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) BuildMoveset(pokeName PokeName, moveNames MoveNames, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Moveset, MoveName, map[MoveName]PBCombKPolicyData, error) {
	moveset := pokemon.Moveset.Copy()
	movesetLength := len(pokemon.Moveset)
	learnsetLength := len(POKEDEX[pokeName].Learnset)

	if movesetLength == MAX_MOVESET_LENGTH || movesetLength == learnsetLength {
		return pokemon.Moveset, "", map[MoveName]PBCombKPolicyData{}, nil
	}

	moveNameWithPolicyData := pbCombKList.MoveNameWithPolicyData(moveNames, pokemon, team)

	if len(moveNameWithPolicyData) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v pokemon.Moveset.Keys() = %v の状態で、次の技の組み合わせが見つからなかった",
			pokemon.Name, pokemon.Moveset.Keys(),
		)
		return Moveset{}, "", map[MoveName]PBCombKPolicyData{}, fmt.Errorf(errMsg)
	}

	length := len(moveNameWithPolicyData)
	selectableMoveNames := make(MoveNames, 0, length)
	meanPolicies := make([]float64, 0, length)

	for k, v := range moveNameWithPolicyData {
		selectableMoveNames = append(selectableMoveNames, k)
		meanPolicies = append(meanPolicies, v.Mean)
	}

	finalPoliciesY := finalPolicies(meanPolicies)
	selectIndex := omw.RandomIntWithWeight(finalPoliciesY, random)
	selectMoveName := selectableMoveNames[selectIndex]

	powerPoint := NewPowerPoint(MOVEDEX[selectMoveName].BasePP, MAX_POINT_UP)
	moveset[selectMoveName] = &powerPoint

	return moveset, selectMoveName, moveNameWithPolicyData, nil
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) BuildAbility(abilities Abilities, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Ability, map[Ability]PBCombKPolicyData, error) {
	if pokemon.Ability != "" {
		return pokemon.Ability, map[Ability]PBCombKPolicyData{}, nil
	}

	abilityWithPolicyData := pbCombKList.AbilityWithPolicyData(abilities, pokemon, team)

	if len(abilityWithPolicyData) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次の特性の組み合わせが見つからなかった", pokemon.Name)
		return "", map[Ability]PBCombKPolicyData{}, fmt.Errorf(errMsg)
	}

	length := len(abilityWithPolicyData)
	selectableAbilities := make(Abilities, 0, length)
	meanPolicies := make([]float64, 0, length)

	for k, v := range abilityWithPolicyData {
		selectableAbilities = append(selectableAbilities, k)
		meanPolicies = append(meanPolicies, v.Mean)
	}

	finalPoliciesY := finalPolicies(meanPolicies)
	selectIndex := omw.RandomIntWithWeight(finalPoliciesY, random)
	selectAbility := selectableAbilities[selectIndex]

	return selectAbility, abilityWithPolicyData, nil
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) BuildItem(items Items, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Item, map[Item]PBCombKPolicyData, error) {
	if pokemon.Item != "" {
		return pokemon.Item, map[Item]PBCombKPolicyData{}, nil
	}

	itemWithPolicyData := pbCombKList.ItemWithPolicyData(items, pokemon, team)

	if len(itemWithPolicyData) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次のアイテムの組み合わせが見つからなかった", pokemon.Name)
		return "", map[Item]PBCombKPolicyData{}, fmt.Errorf(errMsg)
	}

	length := len(itemWithPolicyData)
	selectableItems := make(Items, 0, length)
	meanPolicies := make([]float64, 0, length)

	for k, v := range itemWithPolicyData {
		selectableItems = append(selectableItems, k)
		meanPolicies = append(meanPolicies, v.Mean)
	}

	finalPoliciesY := finalPolicies(meanPolicies)
	selectIndex := omw.RandomIntWithWeight(finalPoliciesY, random)
	selectItem := selectableItems[selectIndex]

	return selectItem, itemWithPolicyData, nil
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) BuildNature(natures Natures, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Nature, map[Nature]PBCombKPolicyData, error) {
	if pokemon.Nature != "" {
		return pokemon.Nature, map[Nature]PBCombKPolicyData{}, nil
	}

	natureWithPolicyData := pbCombKList.NatureWithPolicyData(natures, pokemon, team)

	if len(natureWithPolicyData) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次の性格の組み合わせが見つからなかった", pokemon.Name)
		return "", map[Nature]PBCombKPolicyData{}, fmt.Errorf(errMsg)
	}

	length := len(natureWithPolicyData)
	selectableNatures := make(Natures, 0, length)
	meanPolicies := make([]float64, 0, length)

	for k, v := range natureWithPolicyData {
		selectableNatures = append(selectableNatures, k)
		meanPolicies = append(meanPolicies, v.Mean)
	}

	finalPoliciesY := finalPolicies(meanPolicies)
	selectIndex := omw.RandomIntWithWeight(finalPoliciesY, random)
	selectNature := selectableNatures[selectIndex]

	return selectNature, natureWithPolicyData, nil	
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) Run(pokemon Pokemon, team Team, pbCommonK *PokemonBuildCommonKnowledge, finalPolicies func([]float64) []float64, random *rand.Rand) (Pokemon, PBCombKRunHistory, error) {
	if pokemon.Name == "" {
		return Pokemon{}, PBCombKRunHistory{}, fmt.Errorf("pokemon.Name が 空の状態で、PokemonBuildCombinationKnowledgeList.Run は 実行出来ない")
	}

	pbCombKRunHistory := PBCombKRunHistory{
		SelectMoveNames:make(MoveNames, 0, MAX_MOVESET_LENGTH),
		MoveNameWithPolicyDataList:make([]map[MoveName]PBCombKPolicyData, 0, MAX_MOVESET_LENGTH),
	}

	pokeData := POKEDEX[pokemon.Name]
	pbCommonKMoveNames := pbCommonK.MoveNames
	allAbilities := pokeData.AllAbilities
	pbCommonKItems := pbCommonK.Items
	pbCommonKNatures := pbCommonK.Natures

	ability, abilityWithPolicyData, err := pbCombKList.BuildAbility(allAbilities, pokemon, team, finalPolicies, random)
	if err != nil {
		return Pokemon{}, PBCombKRunHistory{}, err
	}

	pokemon.Ability = ability
	pbCombKRunHistory.SelectAbility = ability
	pbCombKRunHistory.AbilityWithPolicyData = abilityWithPolicyData

	item, itemWithPolicyData, err := pbCombKList.BuildItem(pbCommonKItems, pokemon, team, finalPolicies, random)
	if err != nil {
		return Pokemon{}, PBCombKRunHistory{}, err
	}

	pokemon.Item = item
	pbCombKRunHistory.SelectItem = item
	pbCombKRunHistory.ItemWithPolicyData = itemWithPolicyData

	for i := 0; i < MAX_MOVESET_LENGTH; i++ {
		moveset, selectMoveName, moveNameWithPolicyData, err := pbCombKList.BuildMoveset(pokemon.Name, pbCommonKMoveNames, pokemon, team, finalPolicies, random)
		if err != nil {
			return Pokemon{}, PBCombKRunHistory{}, err
		}

		if selectMoveName == "" {
			continue
		}

		pokemon.Moveset = moveset
		pbCombKRunHistory.SelectMoveNames = append(pbCombKRunHistory.SelectMoveNames, selectMoveName)
		pbCombKRunHistory.MoveNameWithPolicyDataList = append(pbCombKRunHistory.MoveNameWithPolicyDataList, moveNameWithPolicyData)
	}

	nature, natureWithPolicyData, err := pbCombKList.BuildNature(pbCommonKNatures, pokemon, team, finalPolicies, random)
	if err != nil {
		return Pokemon{}, PBCombKRunHistory{}, err
	}

	pokemon.Nature = nature
	pbCombKRunHistory.SelectNature = nature
	pbCombKRunHistory.NatureWithPolicyData = natureWithPolicyData

	return pokemon, pbCombKRunHistory, nil
}


type PokemonBuildCombinationKnowledgeLists []PokemonBuildCombinationKnowledgeList

func (pbCombKLists PokemonBuildCombinationKnowledgeLists) IndicesAccess(indices []int) PokemonBuildCombinationKnowledgeLists {
	result := make(PokemonBuildCombinationKnowledgeLists, len(indices))
	for i, index := range indices {
		result[i] = pbCombKLists[index]
	}
	return result
}

func (pbCombKLists PokemonBuildCombinationKnowledgeLists) RandomChoices(size int, random *rand.Rand) (PokemonBuildCombinationKnowledgeLists, []int) {
	length := len(pbCombKLists)
	indices := make([]int, size)

	for i := 0; i < size; i++ {
		index := random.Intn(length)
		indices[i] = index
	}
	return pbCombKLists.IndicesAccess(indices), indices
}

type PBCombKPolicyData struct {
	Mean float64
	KnowledgeList PokemonBuildCombinationKnowledgeList
}

type PBCombKRunHistory struct {
	SelectMoveNames MoveNames
	SelectAbility Ability
	SelectItem Item
	SelectNature Nature

	MoveNameWithPolicyDataList []map[MoveName]PBCombKPolicyData
	AbilityWithPolicyData map[Ability]PBCombKPolicyData
	ItemWithPolicyData map[Item]PBCombKPolicyData
	NatureWithPolicyData map[Nature]PBCombKPolicyData
}