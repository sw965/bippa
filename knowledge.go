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
	CombinationNum int
	InitIndex int
}

func (pscf *PokemonStateCombinationFeature) Init() {
	combinationNum := 0

	if pscf.Ability != "" {
		combinationNum += 1
	}

	if pscf.Item != "" {
		combinationNum += 1
	}

	combinationNum += len(pscf.MoveNames)

	if pscf.Nature != "" {
		combinationNum += 1
	}

	pscf.CombinationNum = combinationNum
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

func (pscf1 *PokemonStateCombinationFeature) IsInclusion(pscf2 *PokemonStateCombinationFeature) bool {
	if pscf1.CombinationNum < pscf2.CombinationNum {
		return false
	}

	pscf2Ability := pscf2.Ability
	pscf2Item := pscf2.Item
	pscf2Nature := pscf2.Nature

	if pscf2Ability != "" {
		if pscf1.Ability != pscf2Ability {
			return false
		}
	}

	if pscf2Item != "" {
		if pscf1.Item != pscf2Item {
			return false
		}
	}

	if pscf2Nature != "" {
		if pscf1.Nature != pscf2Nature {
			return false
		}
	}

	if !pscf1.MoveNames.InAll(pscf2.MoveNames...) {
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

func (pscfs PokemonStateCombinationFeatures) Init() {
	for i := 0; i < len(pscfs); i++ {
		pscfs[i].Init()
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

//ある要素に属している要素を取り除く
func (pscfs PokemonStateCombinationFeatures) Set() PokemonStateCombinationFeatures {
	result := make(PokemonStateCombinationFeatures, 0, len(pscfs))

	for _, pscfB := range pscfs {
		for _, pscfA := range pscfs {
			if pscfA.IsInclusion(pscfB) {
				continue
			}
		}
		result = append(result, pscfB)
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
	pscfs = pscfs.GetOKs(nextPokemon)

	//情報量が多い特徴量のみを取り出す
	return pscfs.Set()
}

type PokemonBuildModel struct {
	X PokemonStateCombinationFeatures
	Policies []float64
	Values []float64
}

func (pokemonBuildModel *PokemonBuildModel) MoveNameWithPolicyData(moveNames MoveNames, pokemon Pokemon, team Team) map[MoveName]PBMPolicyData {
	moveNameWithPolicyData := map[MoveName]PBMPolicyData{}
	nextPokemon := pokemon

	for _, moveName := range moveNames {
		_, ok := pokemon.Moveset[moveName]

		if ok {
			continue
		}

		moveset := pokemon.Moveset.Copy()
		moveset[moveName] = &PowerPoint{}
		nextPokemon.Moveset = moveset

		policyPSCFs := pokemonBuildModel.X.Policy(&pokemon, &nextPokemon)
		initIndices := policyPSCFs.InitIndices()
		meanPolicy := omw.SliceFloat64Mean(omw.SliceFloat64IndicesAccess(pokemonBuildModel.Policies, initIndices))
		policyData := PBMPolicyData{Mean:meanPolicy, PSCFs:policyPSCFs}

		moveNameWithPolicyData[moveName] = policyData
	}
	return moveNameWithPolicyData
}

func (pokemonBuildModel *PokemonBuildModel) AbilityWithPolicyData(abilities Abilities, pokemon Pokemon, team Team) (map[Ability]PBMPolicyData) {
	abilityWithPolicyData:= map[Ability]PBMPolicyData{}
	nextPokemon := pokemon

	for _, ability := range abilities {
		nextPokemon.Ability = ability

		policyPSCFs := pokemonBuildModel.X.Policy(&pokemon, &nextPokemon)
		initIndices := policyPSCFs.InitIndices()
		meanPolicy := omw.SliceFloat64Mean(omw.SliceFloat64IndicesAccess(pokemonBuildModel.Policies, initIndices))
		policyData := PBMPolicyData{Mean:meanPolicy, PSCFs:policyPSCFs}

		abilityWithPolicyData[ability] = policyData
	}
	return abilityWithPolicyData
}

func (pokemonBuildModel *PokemonBuildModel) ItemWithPolicyData(items Items, pokemon Pokemon, team Team) map[Item]PBMPolicyData {
	itemWithPolicyData := map[Item]PBMPolicyData{}
	nextPokemon := pokemon

	for _, item := range items {
		if team.Items().In(item) {
			continue
		}

		nextPokemon.Item = item

		policyPSCFs := pokemonBuildModel.X.Policy(&pokemon, &nextPokemon)
		initIndices := policyPSCFs.InitIndices()
		meanPolicy := omw.SliceFloat64Mean(omw.SliceFloat64IndicesAccess(pokemonBuildModel.Policies, initIndices))
		policyData := PBMPolicyData{Mean:meanPolicy, PSCFs:policyPSCFs}

		itemWithPolicyData[item] = policyData
	}
	return itemWithPolicyData
}

func (pokemonBuildModel *PokemonBuildModel) NatureWithPolicyData(natures Natures, pokemon Pokemon, team Team) map[Nature]PBMPolicyData {
	natureWithPolicyData := map[Nature]PBMPolicyData{}
	nextPokemon := pokemon

	for _, nature := range natures {
		nextPokemon.Nature = nature
		
		policyPSCFs := pokemonBuildModel.X.Policy(&pokemon, &nextPokemon)
		initIndices := policyPSCFs.InitIndices()		
		meanPolicy := omw.SliceFloat64Mean(omw.SliceFloat64IndicesAccess(pokemonBuildModel.Policies, initIndices))
		policyData := PBMPolicyData{Mean:meanPolicy, PSCFs:policyPSCFs}

		natureWithPolicyData[nature] = policyData
	}
	return natureWithPolicyData
}

func (pokemonBuildModel *PokemonBuildModel) BuildMoveset(pokeName PokeName, moveNames MoveNames, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Moveset, MoveName, map[MoveName]PBMPolicyData, error) {
	moveset := pokemon.Moveset.Copy()
	movesetLength := len(pokemon.Moveset)
	learnsetLength := len(POKEDEX[pokeName].Learnset)

	if movesetLength == MAX_MOVESET_LENGTH || movesetLength == learnsetLength {
		return pokemon.Moveset, "", map[MoveName]PBMPolicyData{}, nil
	}

	moveNameWithPolicyData := pokemonBuildModel.MoveNameWithPolicyData(moveNames, pokemon, team)

	if len(moveNameWithPolicyData) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v pokemon.Moveset.Keys() = %v の状態で、次の技の組み合わせが見つからなかった",
			pokemon.Name, pokemon.Moveset.Keys(),
		)
		return Moveset{}, "", map[MoveName]PBMPolicyData{}, fmt.Errorf(errMsg)
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

func (pokemonBuildModel *PokemonBuildModel) BuildAbility(abilities Abilities, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Ability, map[Ability]PBMPolicyData, error) {
	if pokemon.Ability != "" {
		return pokemon.Ability, map[Ability]PBMPolicyData{}, nil
	}

	abilityWithPolicyData := pokemonBuildModel.AbilityWithPolicyData(abilities, pokemon, team)

	if len(abilityWithPolicyData) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次の特性の組み合わせが見つからなかった", pokemon.Name)
		return "", map[Ability]PBMPolicyData{}, fmt.Errorf(errMsg)
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

func (pokemonBuildModel *PokemonBuildModel) BuildItem(items Items, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Item, map[Item]PBMPolicyData, error) {
	if pokemon.Item != "" {
		return pokemon.Item, map[Item]PBMPolicyData{}, nil
	}

	itemWithPolicyData := pokemonBuildModel.ItemWithPolicyData(items, pokemon, team)

	if len(itemWithPolicyData) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次のアイテムの組み合わせが見つからなかった", pokemon.Name)
		return "", map[Item]PBMPolicyData{}, fmt.Errorf(errMsg)
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

func (pokemonBuildModel *PokemonBuildModel) BuildNature(natures Natures, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Nature, map[Nature]PBMPolicyData, error) {
	if pokemon.Nature != "" {
		return pokemon.Nature, map[Nature]PBMPolicyData{}, nil
	}

	natureWithPolicyData := pokemonBuildModel.NatureWithPolicyData(natures, pokemon, team)

	if len(natureWithPolicyData) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次の性格の組み合わせが見つからなかった", pokemon.Name)
		return "", map[Nature]PBMPolicyData{}, fmt.Errorf(errMsg)
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

func (pokemonBuildModel *PokemonBuildModel) Run(pokemon Pokemon, team Team, pbCommonK *PokemonBuildCommonKnowledge, finalPolicies func([]float64) []float64, random *rand.Rand) (Pokemon, PBMRunHistory, error) {
	if pokemon.Name == "" {
		return Pokemon{}, PBMRunHistory{}, fmt.Errorf("pokemon.Name が 空の状態で、PokemonStateCombinationFeatures.Run は 実行出来ない")
	}

	pbmRunHistory := PBMRunHistory{
		SelectMoveNames:make(MoveNames, 0, MAX_MOVESET_LENGTH),
		MoveNameWithPolicyDataList:make([]map[MoveName]PBMPolicyData, 0, MAX_MOVESET_LENGTH),
	}

	pokeData := POKEDEX[pokemon.Name]
	pbCommonKMoveNames := pbCommonK.MoveNames
	allAbilities := pokeData.AllAbilities
	pbCommonKItems := pbCommonK.Items
	pbCommonKNatures := pbCommonK.Natures

	ability, abilityWithPolicyData, err := pokemonBuildModel.BuildAbility(allAbilities, pokemon, team, finalPolicies, random)
	if err != nil {
		return Pokemon{}, PBMRunHistory{}, err
	}

	pokemon.Ability = ability
	pbmRunHistory.SelectAbility = ability
	pbmRunHistory.AbilityWithPolicyData = abilityWithPolicyData

	item, itemWithPolicyData, err := pokemonBuildModel.BuildItem(pbCommonKItems, pokemon, team, finalPolicies, random)
	if err != nil {
		return Pokemon{}, PBMRunHistory{}, err
	}

	pokemon.Item = item
	pbmRunHistory.SelectItem = item
	pbmRunHistory.ItemWithPolicyData = itemWithPolicyData

	for i := 0; i < MAX_MOVESET_LENGTH; i++ {
		moveset, selectMoveName, moveNameWithPolicyData, err := pokemonBuildModel.BuildMoveset(pokemon.Name, pbCommonKMoveNames, pokemon, team, finalPolicies, random)
		if err != nil {
			return Pokemon{}, PBMRunHistory{}, err
		}

		if selectMoveName == "" {
			continue
		}

		pokemon.Moveset = moveset
		pbmRunHistory.SelectMoveNames = append(pbmRunHistory.SelectMoveNames, selectMoveName)
		pbmRunHistory.MoveNameWithPolicyDataList = append(pbmRunHistory.MoveNameWithPolicyDataList, moveNameWithPolicyData)
	}

	nature, natureWithPolicyData, err := pokemonBuildModel.BuildNature(pbCommonKNatures, pokemon, team, finalPolicies, random)
	if err != nil {
		return Pokemon{}, PBMRunHistory{}, err
	}

	pokemon.Nature = nature
	pbmRunHistory.SelectNature = nature
	pbmRunHistory.NatureWithPolicyData = natureWithPolicyData

	return pokemon, pbmRunHistory, nil
}

type PBMPolicyData struct {
	Mean float64
	PSCFs PokemonStateCombinationFeatures
}

type PBMRunHistory struct {
	SelectMoveNames MoveNames
	SelectAbility Ability
	SelectItem Item
	SelectNature Nature

	MoveNameWithPolicyDataList []map[MoveName]PBMPolicyData
	AbilityWithPolicyData map[Ability]PBMPolicyData
	ItemWithPolicyData map[Item]PBMPolicyData
	NatureWithPolicyData map[Nature]PBMPolicyData
}