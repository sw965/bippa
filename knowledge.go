package bippa

import (
	"fmt"
	"github.com/sw965/omw"
	"math/rand"
)

type PokemonBuildCommonKnowledge struct {
	Items     Items
	MoveNames MoveNames
	Natures   Natures
}

type PokemonBuildCombinationKnowledge struct {
	SelfAbility    Ability
	SelfItem       Item
	SelfMoveNames  MoveNames
	SelfNature     Nature
	CombinationNum int
	Policy         float64
	Value float64
}

func (pbCombK *PokemonBuildCombinationKnowledge) Init() {
	combinationNum := 0

	if pbCombK.SelfAbility != "" {
		combinationNum += 1
	}

	if pbCombK.SelfItem != "" {
		combinationNum += 1
	}

	combinationNum += len(pbCombK.SelfMoveNames)

	if pbCombK.SelfNature != "" {
		combinationNum += 1
	}

	pbCombK.CombinationNum = combinationNum
	pbCombK.Policy = 0.5
	pbCombK.Value = 1.0
}

func (pbCombK *PokemonBuildCombinationKnowledge) All(pokemon *Pokemon, team Team) bool {
	if pbCombK.SelfAbility != "" {
		if pbCombK.SelfAbility != pokemon.Ability {
			return false
		}
	}

	if pbCombK.SelfItem != "" {
		if pbCombK.SelfItem != pokemon.Item {
			return false
		}
	}

	if len(pbCombK.SelfMoveNames) != 0 {
		for _, moveName := range pbCombK.SelfMoveNames {
			_, ok := pokemon.Moveset[moveName]
			if !ok {
				return false
			}
		}
	}

	if pbCombK.SelfNature != "" {
		if pbCombK.SelfNature != pokemon.Nature {
			return false
		}
	}

	return true
}

func (pbCombK1 *PokemonBuildCombinationKnowledge) NearlyEqual(pbCombK2 *PokemonBuildCombinationKnowledge) bool {
	if pbCombK1.SelfAbility != pbCombK2.SelfAbility {
		return false
	}

	if pbCombK1.SelfItem != pbCombK2.SelfItem {
		return false
	}

	sortedMoveNames1 := pbCombK1.SelfMoveNames.Sort()
	sortedMoveNames2 := pbCombK2.SelfMoveNames.Sort()

	if !sortedMoveNames1.Equal(sortedMoveNames2) {
		return false
	}

	if pbCombK1.SelfNature != pbCombK2.SelfNature {
		return false
	}
	return true
}

func (pbCombK1 *PokemonBuildCombinationKnowledge) IsInclusion(pbCombK2 *PokemonBuildCombinationKnowledge) bool {
	if pbCombK1.CombinationNum < pbCombK2.CombinationNum {
		return false
	}

	pbCombK2SelfAbility := pbCombK2.SelfAbility
	pbCombK2SelfItem := pbCombK2.SelfItem
	pbCombK2SelfNature := pbCombK2.SelfNature

	if pbCombK2SelfAbility != "" {
		if pbCombK1.SelfAbility != pbCombK2SelfAbility {
			return false
		}
	}

	if pbCombK2SelfItem != "" {
		if pbCombK1.SelfItem != pbCombK2SelfItem {
			return false
		}
	}

	if pbCombK2SelfNature != "" {
		if pbCombK1.SelfNature != pbCombK2SelfNature {
			return false
		}
	}

	return pbCombK1.SelfMoveNames.InAll(pbCombK2.SelfMoveNames...)
}

type PokemonBuildCombinationKnowledgeList []*PokemonBuildCombinationKnowledge

func NewInitPokemonBuildCombinationKnowledgeList(pokeName PokeName, pbCommonK *PokemonBuildCommonKnowledge) PokemonBuildCombinationKnowledgeList {
	pokeData := POKEDEX[pokeName]

	abilities := pokeData.AllAbilities
	learnset := pokeData.Learnset

	items := pbCommonK.Items
	moveNames := pbCommonK.MoveNames
	natures := pbCommonK.Natures

	result := make(PokemonBuildCombinationKnowledgeList, 0, 6400)

	for _, ability := range abilities {
		result = append(result, &PokemonBuildCombinationKnowledge{SelfAbility: ability})
	}

	for _, item := range ALL_ITEMS {
		result = append(result, &PokemonBuildCombinationKnowledge{SelfItem: item})
	}

	for _, moveName := range learnset {
		result = append(result, &PokemonBuildCombinationKnowledge{SelfMoveNames: MoveNames{moveName}})
	}

	for _, nature := range ALL_NATURES {
		result = append(result, &PokemonBuildCombinationKnowledge{SelfNature: nature})
	}

	for _, ability := range abilities {
		for _, item := range items {
			result = append(result, &PokemonBuildCombinationKnowledge{SelfAbility: ability, SelfItem: item})
		}
	}

	for _, ability := range abilities {
		for _, moveName := range moveNames {
			result = append(result, &PokemonBuildCombinationKnowledge{SelfAbility: ability, SelfMoveNames: MoveNames{moveName}})
		}
	}

	for _, ability := range abilities {
		for _, nature := range natures {
			result = append(result, &PokemonBuildCombinationKnowledge{SelfAbility: ability, SelfNature: nature})
		}
	}

	for _, item := range items {
		for _, moveName := range moveNames {
			result = append(result, &PokemonBuildCombinationKnowledge{SelfItem: item, SelfMoveNames: MoveNames{moveName}})
		}
	}

	for _, item := range items {
		for _, nature := range natures {
			result = append(result, &PokemonBuildCombinationKnowledge{SelfItem: item, SelfNature: nature})
		}
	}

	combination2MoveNames, err := moveNames.Combination(2)
	if err != nil {
		panic(err)
	}

	for _, moveNames := range combination2MoveNames {
		result = append(result, &PokemonBuildCombinationKnowledge{SelfMoveNames: moveNames})
	}

	for _, moveName := range moveNames {
		for _, nature := range natures {
			result = append(result, &PokemonBuildCombinationKnowledge{SelfMoveNames: MoveNames{moveName}, SelfNature: nature})
		}
	}

	combination3MoveNames, err := moveNames.Combination(3)
	if err != nil {
		panic(err)
	}

	for _, moveNames := range combination3MoveNames {
		result = append(result, &PokemonBuildCombinationKnowledge{SelfMoveNames: moveNames})
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

func (pbCombKList PokemonBuildCombinationKnowledgeList) Filter(f func(*PokemonBuildCombinationKnowledge) bool) PokemonBuildCombinationKnowledgeList {
	result := make(PokemonBuildCombinationKnowledgeList, 0, len(pbCombKList))
	for _, pbCombK := range pbCombKList {
		if f(pbCombK) {
			result = append(result, pbCombK)
		}
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

func (pbCombKList PokemonBuildCombinationKnowledgeList) Usable(pokemon, nextPokemon *Pokemon, team Team) PokemonBuildCombinationKnowledgeList {
	//一つ前の状態(pokemon)が満たしている組み合わせを排除する(差分を見る為に)
	pbCombKList = pbCombKList.Filter(func(pbCombK *PokemonBuildCombinationKnowledge) bool { return !pbCombK.All(pokemon, team) })

	//次の状態(nextPokemon)が満たしている組み合わせを取り出す
	pbCombKList = pbCombKList.Filter(func(pbCombK *PokemonBuildCombinationKnowledge) bool { return pbCombK.All(nextPokemon, team) })

	//組み合わせ数が最も多い知識を取り出す
	maxCombinationNum := pbCombKList.MaxCombinationNum()
	return pbCombKList.Filter(func(pbCombK *PokemonBuildCombinationKnowledge) bool {
		return pbCombK.CombinationNum == maxCombinationNum
	})
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) AbilityMeanPolicyAndUsableKnowledgeList(abilities Abilities, pokemon Pokemon, team Team) (AbilityWithFloat64, map[Ability]PokemonBuildCombinationKnowledgeList) {
	abilityWithPolicy := AbilityWithFloat64{}
	abilityWithUsableKnowledgeList := map[Ability]PokemonBuildCombinationKnowledgeList{}
	nextPokemon := pokemon

	for _, ability := range abilities {
		nextPokemon.Ability = ability

		usableKnowledgeList:= pbCombKList.Usable(&pokemon, &nextPokemon, team)
		meanPolicy := usableKnowledgeList.MeanPolicy()

		abilityWithPolicy[ability] = meanPolicy
		abilityWithUsableKnowledgeList[ability] = usableKnowledgeList
	}
	return abilityWithPolicy, abilityWithUsableKnowledgeList
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) ItemMeanPolicyAndUsableKnowledgeList(items Items, pokemon Pokemon, team Team) (ItemWithFloat64, map[Item]PokemonBuildCombinationKnowledgeList) {
	itemWithPolicy := ItemWithFloat64{}
	itemWithUsableKnowledgeList := map[Item]PokemonBuildCombinationKnowledgeList{}
	nextPokemon := pokemon

	for _, item := range items {
		if team.Items().In(item) {
			continue
		}

		nextPokemon.Item = item

		usableKnowledgeList:= pbCombKList.Usable(&pokemon, &nextPokemon, team)
		meanPolicy := usableKnowledgeList.MeanPolicy()

		itemWithPolicy[item] = meanPolicy
		itemWithUsableKnowledgeList[item] = usableKnowledgeList
	}
	return itemWithPolicy, itemWithUsableKnowledgeList
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) MoveNameMeanPolicyAndUsableKnowledgeList(moveNames MoveNames, pokemon Pokemon, team Team) (MoveNameWithFloat64, map[MoveName]PokemonBuildCombinationKnowledgeList) {
	moveNameWithMeanPolicy := MoveNameWithFloat64{}
	moveNameWithUsableKnowledgeList := map[MoveName]PokemonBuildCombinationKnowledgeList{}
	nextPokemon := pokemon

	for _, moveName := range moveNames {
		_, ok := pokemon.Moveset[moveName]

		if ok {
			continue
		}

		moveset := pokemon.Moveset.Copy()
		moveset[moveName] = &PowerPoint{}
		nextPokemon.Moveset = moveset

		usableKnowledgeList := pbCombKList.Usable(&pokemon, &nextPokemon, team)
		meanPolicy := usableKnowledgeList.MeanPolicy()

		moveNameWithMeanPolicy[moveName] = meanPolicy
		moveNameWithUsableKnowledgeList[moveName] = usableKnowledgeList
	}
	return moveNameWithMeanPolicy, moveNameWithUsableKnowledgeList
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) NatureMeanPolicyAndUsableKnowledgeList(natures Natures, pokemon Pokemon, team Team) (NatureWithFloat64, map[Nature]PokemonBuildCombinationKnowledgeList) {
	natureWithMeanPolicy := NatureWithFloat64{}
	natureWithUsableKnowledgeList := map[Nature]PokemonBuildCombinationKnowledgeList{}
	nextPokemon := pokemon

	for _, nature := range natures {
		nextPokemon.Nature = nature

		usableKnowledgeList := pbCombKList.Usable(&pokemon, &nextPokemon, team)
		meanPolicy := usableKnowledgeList.MeanPolicy()

		natureWithMeanPolicy[nature] = meanPolicy
		natureWithUsableKnowledgeList[nature] = usableKnowledgeList
	}
	return natureWithMeanPolicy, natureWithUsableKnowledgeList
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) BuildAbility(abilities Abilities, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Ability, PokemonBuildCombinationKnowledgeList, PokemonBuildCombinationKnowledgeList, error) {
	if pokemon.Ability != "" {
		return pokemon.Ability, PokemonBuildCombinationKnowledgeList{}, PokemonBuildCombinationKnowledgeList{}, nil
	}

	abilityWithMeanPolicy, abilityWithUsableKnowledgeList := pbCombKList.AbilityMeanPolicyAndUsableKnowledgeList(abilities, pokemon, team)
	if len(abilityWithMeanPolicy) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次の特性の組み合わせが見つからなかった", pokemon.Name)
		return "", PokemonBuildCombinationKnowledgeList{}, PokemonBuildCombinationKnowledgeList{}, fmt.Errorf(errMsg)
	}

	abilities, meanPolicies := abilityWithMeanPolicy.KeysAndValues()
	finalPoliciesY := finalPolicies(meanPolicies)
	selectIndex := omw.RandomIntWithWeight(finalPoliciesY, random)
	selectAbility := abilities[selectIndex]
	selectKnowledgeList := abilityWithUsableKnowledgeList[selectAbility]

	usableKnowledgeList := make(PokemonBuildCombinationKnowledgeList, 0, len(pbCombKList))
	for _, iUsableKnowledgeList := range abilityWithUsableKnowledgeList {
		usableKnowledgeList = append(usableKnowledgeList, iUsableKnowledgeList...)
	}

	return selectAbility, selectKnowledgeList, usableKnowledgeList, nil
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) BuildItem(items Items, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Item, PokemonBuildCombinationKnowledgeList, PokemonBuildCombinationKnowledgeList, error) {
	if pokemon.Item != "" {
		return pokemon.Item, PokemonBuildCombinationKnowledgeList{}, PokemonBuildCombinationKnowledgeList{}, nil
	}

	itemWithMeanPolicy, itemWithUsableKnowledgeList := pbCombKList.ItemMeanPolicyAndUsableKnowledgeList(items, pokemon, team)

	if len(itemWithMeanPolicy) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次のアイテムの組み合わせが見つからなかった", pokemon.Name)
		return "", PokemonBuildCombinationKnowledgeList{}, PokemonBuildCombinationKnowledgeList{}, fmt.Errorf(errMsg)
	}

	items, meanPolicies := itemWithMeanPolicy.KeysAndValues()
	finalPoliciesY := finalPolicies(meanPolicies)
	selectIndex := omw.RandomIntWithWeight(finalPoliciesY, random)
	selectItem := items[selectIndex]
	selectKnowledgeList := itemWithUsableKnowledgeList[selectItem]

	usableKnowledgeList := make(PokemonBuildCombinationKnowledgeList, 0, len(pbCombKList))
	for _, iUsableKnowledgeList := range itemWithUsableKnowledgeList {
		usableKnowledgeList = append(usableKnowledgeList, iUsableKnowledgeList...)
	}

	return selectItem, selectKnowledgeList, usableKnowledgeList, nil
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) BuildMoveset(pokeName PokeName, moveNames MoveNames, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Moveset, PokemonBuildCombinationKnowledgeList, PokemonBuildCombinationKnowledgeList, error) {
	moveset := pokemon.Moveset.Copy()
	movesetLength := len(pokemon.Moveset)
	learnsetLength := len(POKEDEX[pokeName].Learnset)

	if movesetLength == MAX_MOVESET_LENGTH || movesetLength == learnsetLength {
		return pokemon.Moveset, PokemonBuildCombinationKnowledgeList{}, PokemonBuildCombinationKnowledgeList{}, nil
	}

	moveNameWithMeanPolicy, moveNameWithUsableKnowledgeList := pbCombKList.MoveNameMeanPolicyAndUsableKnowledgeList(moveNames, pokemon, team)

	if len(moveNameWithMeanPolicy) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v pokemon.Moveset.Keys() = %v の状態で、次の技の組み合わせが見つからなかった",
			pokemon.Name, pokemon.Moveset.Keys(),
		)
		return Moveset{}, PokemonBuildCombinationKnowledgeList{}, PokemonBuildCombinationKnowledgeList{}, fmt.Errorf(errMsg)
	}

	moveNames, meanPolicies := moveNameWithMeanPolicy.KeysAndValues()
	finalPoliciesY := finalPolicies(meanPolicies)
	selectIndex := omw.RandomIntWithWeight(finalPoliciesY, random)
	selectMoveName := moveNames[selectIndex]
	selectKnowledgeList := moveNameWithUsableKnowledgeList[selectMoveName]
	powerPoint := NewPowerPoint(MOVEDEX[selectMoveName].BasePP, MAX_POINT_UP)
	moveset[selectMoveName] = &powerPoint

	usableKnowledgeList := make(PokemonBuildCombinationKnowledgeList, 0, len(pbCombKList))
	for _, iUsableKnowledgeList := range moveNameWithUsableKnowledgeList {
		usableKnowledgeList = append(usableKnowledgeList, iUsableKnowledgeList...)
	}

	return moveset, selectKnowledgeList, usableKnowledgeList, nil
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) BuildNature(natures Natures, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Nature, PokemonBuildCombinationKnowledgeList, PokemonBuildCombinationKnowledgeList, error) {
	if pokemon.Nature != "" {
		return pokemon.Nature, PokemonBuildCombinationKnowledgeList{}, PokemonBuildCombinationKnowledgeList{}, nil
	}

	natureWithMeanPolicy, natureWithUsableKnowledgeList := pbCombKList.NatureMeanPolicyAndUsableKnowledgeList(natures, pokemon, team)

	if len(natureWithMeanPolicy) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次の性格の組み合わせが見つからなかった", pokemon.Name)
		return "", PokemonBuildCombinationKnowledgeList{}, PokemonBuildCombinationKnowledgeList{}, fmt.Errorf(errMsg)
	}

	natures, meanPolicies := natureWithMeanPolicy.KeysAndValues()
	finalPoliciesY := finalPolicies(meanPolicies)
	selectIndex := omw.RandomIntWithWeight(finalPoliciesY, random)
	selectNature := natures[selectIndex]
	selectKnowledgeList := natureWithUsableKnowledgeList[selectNature]

	usableKnowledgeList := make(PokemonBuildCombinationKnowledgeList, 0, len(pbCombKList))
	for _, iUsableKnowledgeList := range natureWithUsableKnowledgeList {
		usableKnowledgeList = append(usableKnowledgeList, iUsableKnowledgeList...)
	}

	return selectNature, selectKnowledgeList, usableKnowledgeList, nil	
}

func (pbCombKList PokemonBuildCombinationKnowledgeList) Run(pokemon Pokemon, team Team, pbCommonK *PokemonBuildCommonKnowledge, finalPolicies func([]float64) []float64, random *rand.Rand) (Pokemon, PokemonBuildCombinationKnowledgeLists, PokemonBuildCombinationKnowledgeLists, error) {
	if pokemon.Name == "" {
		return Pokemon{}, PokemonBuildCombinationKnowledgeLists{}, PokemonBuildCombinationKnowledgeLists{}, fmt.Errorf("pokemon.Name が 空の状態で、PokemonBuildCombinationKnowledgeList.Run は 実行出来ない")
	}

	pokeData := POKEDEX[pokemon.Name]
	abilities := pokeData.AllAbilities
	items := pbCommonK.Items
	moveNames := pbCommonK.MoveNames
	natures := pbCommonK.Natures

	selectKnowledgeLists := make(PokemonBuildCombinationKnowledgeLists, 0, 8)
	usableKnowledgeLists := make(PokemonBuildCombinationKnowledgeLists, 0, 8)

	ability, abilitySelectKnowledgeList, abilityUsableKnowledgeList, err := pbCombKList.BuildAbility(abilities, pokemon, team, finalPolicies, random)
	if err != nil {
		return Pokemon{}, PokemonBuildCombinationKnowledgeLists{}, PokemonBuildCombinationKnowledgeLists{}, err
	}

	selectKnowledgeLists = append(selectKnowledgeLists, abilitySelectKnowledgeList)
	usableKnowledgeLists = append(usableKnowledgeLists, abilityUsableKnowledgeList)
	pokemon.Ability = ability

	item, itemSelectKnowledgeList, itemUsableKnowledgeList, err := pbCombKList.BuildItem(items, pokemon, team, finalPolicies, random)
	if err != nil {
		return Pokemon{}, PokemonBuildCombinationKnowledgeLists{}, PokemonBuildCombinationKnowledgeLists{}, err
	}

	selectKnowledgeLists = append(selectKnowledgeLists, itemSelectKnowledgeList)
	usableKnowledgeLists = append(usableKnowledgeLists, itemUsableKnowledgeList)
	pokemon.Item = item

	for i := 0; i < MAX_MOVESET_LENGTH; i++ {
		moveset, moveNameSelectKnowledgeList, moveNameUsableKnowledgeList, err := pbCombKList.BuildMoveset(pokemon.Name, moveNames, pokemon, team, finalPolicies, random)
		if err != nil {
			return Pokemon{}, PokemonBuildCombinationKnowledgeLists{}, PokemonBuildCombinationKnowledgeLists{}, err
		}

		selectKnowledgeLists = append(selectKnowledgeLists, moveNameSelectKnowledgeList)
		usableKnowledgeLists = append(usableKnowledgeLists, moveNameUsableKnowledgeList)
		pokemon.Moveset = moveset
	}

	nature, natureSelectKnowledgeList, natureUsableKnowledgeList, err := pbCombKList.BuildNature(natures, pokemon, team, finalPolicies, random)
	if err != nil {
		return Pokemon{}, PokemonBuildCombinationKnowledgeLists{}, PokemonBuildCombinationKnowledgeLists{}, err
	}

	selectKnowledgeLists = append(selectKnowledgeLists, natureSelectKnowledgeList)
	usableKnowledgeLists = append(usableKnowledgeLists, natureUsableKnowledgeList)
	pokemon.Nature = nature

	return pokemon, selectKnowledgeLists, usableKnowledgeLists, nil
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
	result := make(PokemonBuildCombinationKnowledgeLists, size)
	indices := make([]int, size)

	for i := 0; i < size; i++ {
		index := random.Intn(length)
		result[i] = pbCombKLists[i]
		indices[i] = index
	}
	return result, indices
}

func (selectKnowledgeLists PokemonBuildCombinationKnowledgeLists) PolicyOptimizer(usableKnowledgeLists PokemonBuildCombinationKnowledgeLists, target, learningRate, lowerLimit float64) error {
	for i, selectKnowledgeList := range selectKnowledgeLists {
		usableKnowledgeList := usableKnowledgeLists[i]
		for _, selectKnowledge := range selectKnowledgeList {
			selectIndex := usableKnowledgeList.Index(selectKnowledge)
			
			if selectIndex == -1 {
				return fmt.Errorf("PolicyOptimizer の 引数 の 整合性が取れていない")
			}

			usableKnowledgeListLength := len(usableKnowledgeList)

			for j := 0; j < usableKnowledgeListLength; j++ {
				y := selectKnowledgeList.MeanPolicy()
				var t float64
				if j == selectIndex {
					t = target
				}  else {
					t = 1.0 - target
				}

				//二乗和誤差の微分 (0.5 * (y - t) ^ 2)
				grad := (y - t)
				//算術平均の微分
				grad *= 1.0 / float64(len(selectKnowledgeList))

				usableKnowledgeList[j].Policy -= (grad * learningRate)
				if usableKnowledgeList[j].Policy < lowerLimit {
					usableKnowledgeList[j].Policy = lowerLimit
				}
			}
		}
	}
	return nil
}