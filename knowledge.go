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

type PokemonBuilder struct {
	PSCFs PokemonStateCombinationFeatures
	Policies []float64
	Values []float64
}

func(pokemonBuilder *PokemonBuilder) Eval() float64 {
	return omw.SumFloat64(pokemonBuilder.Values...)
}

func (pokemonBuilder *PokemonBuilder) MoveNameWithPolicyData(moveNames MoveNames, pokemon Pokemon, team Team) map[MoveName]PokemonBuilderPolicyData {
	moveNameWithPolicyData := map[MoveName]PokemonBuilderPolicyData{}
	nextPokemon := pokemon

	for _, moveName := range moveNames {
		_, ok := pokemon.Moveset[moveName]

		if ok {
			continue
		}

		moveset := pokemon.Moveset.Copy()
		moveset[moveName] = &PowerPoint{}
		nextPokemon.Moveset = moveset

		policyPSCFs := pokemonBuilder.PSCFs.Policy(&pokemon, &nextPokemon)
		initIndices := policyPSCFs.InitIndices()
		sumPolicy := omw.SumFloat64(omw.SliceFloat64IndicesAccess(pokemonBuilder.Policies, initIndices)...)
		policyData := PokemonBuilderPolicyData{Sum:sumPolicy, PSCFs:policyPSCFs}

		moveNameWithPolicyData[moveName] = policyData
	}
	return moveNameWithPolicyData
}

func (pokemonBuilder *PokemonBuilder) AbilityWithPolicyData(abilities Abilities, pokemon Pokemon, team Team) (map[Ability]PokemonBuilderPolicyData) {
	abilityWithPolicyData := map[Ability]PokemonBuilderPolicyData{}
	nextPokemon := pokemon

	for _, ability := range abilities {
		nextPokemon.Ability = ability

		policyPSCFs := pokemonBuilder.PSCFs.Policy(&pokemon, &nextPokemon)
		initIndices := policyPSCFs.InitIndices()
		sumPolicy := omw.SumFloat64(omw.SliceFloat64IndicesAccess(pokemonBuilder.Policies, initIndices)...)
		policyData := PokemonBuilderPolicyData{Sum:sumPolicy, PSCFs:policyPSCFs}

		abilityWithPolicyData[ability] = policyData
	}
	return abilityWithPolicyData
}

func (pokemonBuilder *PokemonBuilder) ItemWithPolicyData(items Items, pokemon Pokemon, team Team) map[Item]PokemonBuilderPolicyData {
	itemWithPolicyData := map[Item]PokemonBuilderPolicyData{}
	nextPokemon := pokemon

	for _, item := range items {
		if team.Items().In(item) {
			continue
		}

		nextPokemon.Item = item

		policyPSCFs := pokemonBuilder.PSCFs.Policy(&pokemon, &nextPokemon)
		initIndices := policyPSCFs.InitIndices()
		sumPolicy := omw.SumFloat64(omw.SliceFloat64IndicesAccess(pokemonBuilder.Policies, initIndices)...)
		policyData := PokemonBuilderPolicyData{Sum:sumPolicy, PSCFs:policyPSCFs}

		itemWithPolicyData[item] = policyData
	}
	return itemWithPolicyData
}

func (pokemonBuilder *PokemonBuilder) NatureWithPolicyData(natures Natures, pokemon Pokemon, team Team) map[Nature]PokemonBuilderPolicyData {
	natureWithPolicyData := map[Nature]PokemonBuilderPolicyData{}
	nextPokemon := pokemon

	for _, nature := range natures {
		nextPokemon.Nature = nature
		
		policyPSCFs := pokemonBuilder.PSCFs.Policy(&pokemon, &nextPokemon)
		initIndices := policyPSCFs.InitIndices()		
		sumPolicy := omw.SumFloat64(omw.SliceFloat64IndicesAccess(pokemonBuilder.Policies, initIndices)...)
		policyData := PokemonBuilderPolicyData{Sum:sumPolicy, PSCFs:policyPSCFs}

		natureWithPolicyData[nature] = policyData
	}
	return natureWithPolicyData
}

func (pokemonBuilder *PokemonBuilder) BuildMoveset(pokeName PokeName, moveNames MoveNames, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Moveset, MoveName, map[MoveName]PokemonBuilderPolicyData, error) {
	moveset := pokemon.Moveset.Copy()
	movesetLength := len(pokemon.Moveset)
	learnsetLength := len(POKEDEX[pokeName].Learnset)

	if movesetLength == MAX_MOVESET_LENGTH || movesetLength == learnsetLength {
		return pokemon.Moveset, "", map[MoveName]PokemonBuilderPolicyData{}, nil
	}

	moveNameWithPolicyData := pokemonBuilder.MoveNameWithPolicyData(moveNames, pokemon, team)

	if len(moveNameWithPolicyData) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v pokemon.Moveset.Keys() = %v の状態で、次の技の組み合わせが見つからなかった",
			pokemon.Name, pokemon.Moveset.Keys(),
		)
		return Moveset{}, "", map[MoveName]PokemonBuilderPolicyData{}, fmt.Errorf(errMsg)
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

func (pokemonBuilder *PokemonBuilder) BuildAbility(abilities Abilities, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Ability, map[Ability]PokemonBuilderPolicyData, error) {
	if pokemon.Ability != "" {
		return pokemon.Ability, map[Ability]PokemonBuilderPolicyData{}, nil
	}

	abilityWithPolicyData := pokemonBuilder.AbilityWithPolicyData(abilities, pokemon, team)

	if len(abilityWithPolicyData) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次の特性の組み合わせが見つからなかった", pokemon.Name)
		return "", map[Ability]PokemonBuilderPolicyData{}, fmt.Errorf(errMsg)
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

func (pokemonBuilder *PokemonBuilder) BuildItem(items Items, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Item, map[Item]PokemonBuilderPolicyData, error) {
	if pokemon.Item != "" {
		return pokemon.Item, map[Item]PokemonBuilderPolicyData{}, nil
	}

	itemWithPolicyData := pokemonBuilder.ItemWithPolicyData(items, pokemon, team)

	if len(itemWithPolicyData) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次のアイテムの組み合わせが見つからなかった", pokemon.Name)
		return "", map[Item]PokemonBuilderPolicyData{}, fmt.Errorf(errMsg)
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

func (pokemonBuilder *PokemonBuilder) BuildNature(natures Natures, pokemon Pokemon, team Team, finalPolicies func([]float64) []float64, random *rand.Rand) (Nature, map[Nature]PokemonBuilderPolicyData, error) {
	if pokemon.Nature != "" {
		return pokemon.Nature, map[Nature]PokemonBuilderPolicyData{}, nil
	}

	natureWithPolicyData := pokemonBuilder.NatureWithPolicyData(natures, pokemon, team)

	if len(natureWithPolicyData) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次の性格の組み合わせが見つからなかった", pokemon.Name)
		return "", map[Nature]PokemonBuilderPolicyData{}, fmt.Errorf(errMsg)
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

func (pokemonBuilder *PokemonBuilder) Run(pokemon Pokemon, team Team, pbCommonK *PokemonBuildCommonKnowledge, finalPolicies func([]float64) []float64, random *rand.Rand) (Pokemon, PokemonBuilderRunHistory, error) {
	if pokemon.Name == "" {
		return Pokemon{}, PokemonBuilderRunHistory{}, fmt.Errorf("pokemon.Name が 空の状態で、PokemonBuilder.Run は 実行出来ない")
	}

	pokemonBuilderRunHistory := PokemonBuilderRunHistory{
		SelectMoveNames:make(MoveNames, 0, MAX_MOVESET_LENGTH),
		MoveNameWithPolicyDataList:make([]map[MoveName]PokemonBuilderPolicyData, 0, MAX_MOVESET_LENGTH),
	}

	pokeData := POKEDEX[pokemon.Name]
	pbCommonKMoveNames := pbCommonK.MoveNames
	allAbilities := pokeData.AllAbilities
	pbCommonKItems := pbCommonK.Items
	pbCommonKNatures := pbCommonK.Natures

	ability, abilityWithPolicyData, err := pokemonBuilder.BuildAbility(allAbilities, pokemon, team, finalPolicies, random)
	if err != nil {
		return Pokemon{}, PokemonBuilderRunHistory{}, err
	}

	pokemon.Ability = ability
	pokemonBuilderRunHistory.SelectAbility = ability
	pokemonBuilderRunHistory.AbilityWithPolicyData = abilityWithPolicyData

	item, itemWithPolicyData, err := pokemonBuilder.BuildItem(pbCommonKItems, pokemon, team, finalPolicies, random)
	if err != nil {
		return Pokemon{}, PokemonBuilderRunHistory{}, err
	}

	pokemon.Item = item
	pokemonBuilderRunHistory.SelectItem = item
	pokemonBuilderRunHistory.ItemWithPolicyData = itemWithPolicyData

	for i := 0; i < MAX_MOVESET_LENGTH; i++ {
		moveset, selectMoveName, moveNameWithPolicyData, err := pokemonBuilder.BuildMoveset(pokemon.Name, pbCommonKMoveNames, pokemon, team, finalPolicies, random)
		if err != nil {
			return Pokemon{}, PokemonBuilderRunHistory{}, err
		}

		if selectMoveName == "" {
			continue
		}

		pokemon.Moveset = moveset
		pokemonBuilderRunHistory.SelectMoveNames = append(pokemonBuilderRunHistory.SelectMoveNames, selectMoveName)
		pokemonBuilderRunHistory.MoveNameWithPolicyDataList = append(pokemonBuilderRunHistory.MoveNameWithPolicyDataList, moveNameWithPolicyData)
	}

	nature, natureWithPolicyData, err := pokemonBuilder.BuildNature(pbCommonKNatures, pokemon, team, finalPolicies, random)
	if err != nil {
		return Pokemon{}, PokemonBuilderRunHistory{}, err
	}

	pokemon.Nature = nature
	pokemonBuilderRunHistory.SelectNature = nature
	pokemonBuilderRunHistory.NatureWithPolicyData = natureWithPolicyData

	return pokemon, pokemonBuilderRunHistory, nil
}

type PokemonBuilderPolicyData struct {
	Sum float64
	PSCFs PokemonStateCombinationFeatures
}

type PokemonBuilderRunHistory struct {
	SelectMoveNames MoveNames
	SelectAbility Ability
	SelectItem Item
	SelectNature Nature

	MoveNameWithPolicyDataList []map[MoveName]PokemonBuilderPolicyData
	AbilityWithPolicyData map[Ability]PokemonBuilderPolicyData
	ItemWithPolicyData map[Item]PokemonBuilderPolicyData
	NatureWithPolicyData map[Nature]PokemonBuilderPolicyData
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

func (tcfs TeamCombinationFeatures) GetOks(team Team) TeamCombinationFeatures {
	result := make(TeamCombinationFeatures, 0, len(tcfs))
	for _, tcf := range tcfs {
		if tcf.OK(team) {
			result = append(result, tcf)
		}
	}
	return result
}

type TeamBuilder struct {
	TCFs TeamCombinationFeatures
	Values []float64	
}

func NewTeamBuilder(pokeName1, pokeName2 PokeName, pbCommonKs map[PokeName]*PokemonBuildCommonKnowledge, random *rand.Rand) TeamBuilder {
	tcfs := NewTeamCombinationFeatures(pokeName1, pokeName2, pbCommonKs)
	values, err := omw.MakeRandomSliceFloat64(len(tcfs), 1, 10, random)
	if err != nil {
		panic(err)
	}
	return TeamBuilder{TCFs:tcfs, Values:values}
}

func (teamBuilder TeamBuilder) Eval() float64 {
	return omw.SumFloat64(teamBuilder.Values...)
}