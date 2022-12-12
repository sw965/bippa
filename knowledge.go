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
	W              float64
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
	pbCombK.W = 0.5
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

type PokemonBuilder []*PokemonBuildCombinationKnowledge

func NewInitPokemonBuilder(pokeName PokeName, pbCommonK *PokemonBuildCommonKnowledge) PokemonBuilder {
	pokeData := POKEDEX[pokeName]

	abilities := pokeData.AllAbilities
	learnset := pokeData.Learnset

	items := pbCommonK.Items
	moveNames := pbCommonK.MoveNames
	natures := pbCommonK.Natures

	result := make(PokemonBuilder, 0, 6400)

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

func (pb PokemonBuilder) Init() {
	for i := 0; i < len(pb); i++ {
		pb[i].Init()
	}
}

func (pb PokemonBuilder) Index(pbCombK *PokemonBuildCombinationKnowledge) int {
	for i, iPBCK := range pb {
		if iPBCK.NearlyEqual(pbCombK) {
			return i
		}
	}
	return -1
}

func (pb PokemonBuilder) Ws() []float64 {
	result := make([]float64, len(pb))
	for i, pbCombK := range pb {
		result[i] = pbCombK.W
	}
	return result
}

func (pb PokemonBuilder) AverageW() float64 {
	return omw.SumFloat64(pb.Ws()...) / float64(len(pb))
}

func (pb PokemonBuilder) Filter(f func(*PokemonBuildCombinationKnowledge) bool) PokemonBuilder {
	result := make(PokemonBuilder, 0, len(pb))
	for _, pbCombK := range pb {
		if f(pbCombK) {
			result = append(result, pbCombK)
		}
	}
	return result
}

func (pb PokemonBuilder) MaxCombinationNum() int {
	result := pb[0].CombinationNum
	for _, pbCombK := range pb[1:] {
		maxCombinationNum := pbCombK.CombinationNum
		if maxCombinationNum > result {
			result = maxCombinationNum
		}
	}
	return result
}

func (pb PokemonBuilder) DiffCalc(pokemon, nextPokemon *Pokemon, team Team, random *rand.Rand) PokemonBuilder {
	//一つ前の状態(pokemon)が満たしている組み合わせを排除する(差分を見る為に)
	pb = pb.Filter(func(pbCombK *PokemonBuildCombinationKnowledge) bool { return !pbCombK.All(pokemon, team) })

	//次の状態(nextPokemon)が満たしている組み合わせを取り出す
	pb = pb.Filter(func(pbCombK *PokemonBuildCombinationKnowledge) bool { return pbCombK.All(nextPokemon, team) })

	//組み合わせ数が最も多い知識を取り出す
	maxCombinationNum := pb.MaxCombinationNum()
	return pb.Filter(func(pbCombK *PokemonBuildCombinationKnowledge) bool {
		return pbCombK.CombinationNum == maxCombinationNum
	})
}

func (pb PokemonBuilder) BuildAbility(abilities Abilities, pokemon Pokemon, team Team, finalWsFunc func([]float64) []float64, random *rand.Rand) (Ability, PokemonBuilderActionHistory, error) {
	if pokemon.Ability != "" {
		return pokemon.Ability, PokemonBuilderActionHistory{}, nil
	}

	abilityWithAverageW := AbilityWithFloat64{}
	abilityWithSelectPBCK := map[Ability]*PokemonBuildCombinationKnowledge{}
	legalPB := make(PokemonBuilder, 0, len(pb))
	nextPokemon := pokemon

	for _, ability := range abilities {
		nextPokemon.Ability = ability

		diffCalcPB := pb.DiffCalc(&pokemon, &nextPokemon, team, random)
		averageW := diffCalcPB.AverageW()
		selectIndex := random.Intn(len(diffCalcPB))

		abilityWithAverageW[ability] = averageW
		abilityWithSelectPBCK[ability] = diffCalcPB[selectIndex]
		legalPB = append(legalPB, diffCalcPB...)
	}

	if len(abilityWithAverageW) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次の特性の組み合わせが見つからなかった", pokemon.Name)
		return "", PokemonBuilderActionHistory{}, fmt.Errorf(errMsg)
	}

	abilities, averageWs := abilityWithAverageW.KeysAndValues()
	finalWs := finalWsFunc(averageWs)
	selectIndex := omw.RandomIntWithWeight(finalWs, random)
	selectAbility := abilities[selectIndex]
	selectPBCK := abilityWithSelectPBCK[selectAbility]
	actionHistory := PokemonBuilderActionHistory{SelectPokemonBuilder:PokemonBuilder{selectPBCK}, LegalPokemonBuilders:[]PokemonBuilder{legalPB}}

	return selectAbility, actionHistory, nil
}

func (pb PokemonBuilder) BuildItem(items Items, pokemon Pokemon, team Team, finalWsFunc func([]float64) []float64, random *rand.Rand) (Item, PokemonBuilderActionHistory, error) {
	if pokemon.Item != "" {
		return pokemon.Item, PokemonBuilderActionHistory{}, nil
	}

	itemWithAverageW := ItemWithFloat64{}
	itemWithSelectPBCK := map[Item]*PokemonBuildCombinationKnowledge{}
	legalPB := make(PokemonBuilder, 0, len(pb))
	nextPokemon := pokemon

	for _, item := range items {
		if team.Items().In(item) {
			continue
		}

		nextPokemon.Item = item
		diffCalcPB := pb.DiffCalc(&pokemon, &nextPokemon, team, random)
		averageW := diffCalcPB.AverageW()
		selectIndex := random.Intn(len(diffCalcPB))

		itemWithAverageW[item] = averageW
		itemWithSelectPBCK[item] = diffCalcPB[selectIndex]
		legalPB = append(legalPB, diffCalcPB...)
	}

	if len(itemWithAverageW) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次のアイテムの組み合わせが見つからなかった", pokemon.Name)
		return "", PokemonBuilderActionHistory{}, fmt.Errorf(errMsg)
	}

	items, averageWs := itemWithAverageW.KeysAndValues()
	finalWs := finalWsFunc(averageWs)
	selectIndex := omw.RandomIntWithWeight(finalWs, random)
	selectItem := items[selectIndex]
	selectPBCK := itemWithSelectPBCK[selectItem]
	actionHistory := PokemonBuilderActionHistory{SelectPokemonBuilder:PokemonBuilder{selectPBCK}, LegalPokemonBuilders:[]PokemonBuilder{legalPB}}

	return selectItem, actionHistory, nil
}

func (pb PokemonBuilder) BuildMoveset(moveNames MoveNames, pokemon Pokemon, team Team, finalWsFunc func([]float64) []float64, random *rand.Rand) (Moveset, PokemonBuilderActionHistory, error) {
	initMovesetLength := len(pokemon.Moveset)
	learnsetLength := len(POKEDEX[pokemon.Name].Learnset)

	if initMovesetLength == MAX_MOVESET_LENGTH {
		return pokemon.Moveset, PokemonBuilderActionHistory{}, nil
	}

	if initMovesetLength == learnsetLength {
		return pokemon.Moveset, PokemonBuilderActionHistory{}, nil
	}

	actionNum := omw.MinInt([]int{MAX_MOVESET_LENGTH - initMovesetLength, learnsetLength}...)

	action := func(nextPokemon Pokemon) (MoveName, *PokemonBuildCombinationKnowledge, PokemonBuilder, error) {
		moveNameWithAverageW := MoveNameWithFloat64{}
		moveNameWithSelectPBCK := map[MoveName]*PokemonBuildCombinationKnowledge{}
		legalPB := make(PokemonBuilder, 0, len(pb))

		for _, moveName := range moveNames {
			moveset := pokemon.Moveset.Copy()

			if _, ok := moveset[moveName]; ok {
				continue
			}

			moveset[moveName] = &PowerPoint{}
			nextPokemon.Moveset = moveset

			diffCalcPB := pb.DiffCalc(&pokemon, &nextPokemon, team, random)
			//ある技について、複数の知識を活用出来る(考慮すべき重みが複数ある)場合は、重みを平均化する。
			averageW := diffCalcPB.AverageW()
			selectIndex := random.Intn(len(diffCalcPB))

			moveNameWithAverageW[moveName] = averageW
			moveNameWithSelectPBCK[moveName] = diffCalcPB[selectIndex]

			//他の選択(現在の選択も含む)によって得られる知識を記録する
			legalPB = append(legalPB, diffCalcPB...)
		}

		if len(moveNameWithAverageW) == 0 {
			errMsg := fmt.Sprintf("pokemon.Name = %v pokemon.Moveset.Keys() = %v の状態で、次の技の組み合わせが見つからなかった",
				pokemon.Name, pokemon.Moveset.Keys(),
			)
			return "", &PokemonBuildCombinationKnowledge{}, PokemonBuilder{}, fmt.Errorf(errMsg)
		}

		moveNames, averageWs := moveNameWithAverageW.KeysAndValues()
		finalWs := finalWsFunc(averageWs)
		selectIndex := omw.RandomIntWithWeight(finalWs, random)
		selectMoveName := moveNames[selectIndex]
		selectPBCK := moveNameWithSelectPBCK[selectMoveName]
		return selectMoveName, selectPBCK, legalPB, nil
	}

	selectPB := make(PokemonBuilder, actionNum)
	legalPBs := make([]PokemonBuilder, actionNum)

	for i := 0; i < actionNum; i++ {
		selectMoveName, selectPBCK, legalPB, err := action(pokemon)
		if err != nil {
			return Moveset{}, PokemonBuilderActionHistory{}, err
		}

		moveset := pokemon.Moveset.Copy()
		powerPoint := NewPowerPoint(MOVEDEX[selectMoveName].BasePP, MAX_POINT_UP)
		moveset[selectMoveName] = &powerPoint
		pokemon.Moveset = moveset

		selectPB[i] = selectPBCK
		legalPBs[i] = legalPB
	}

	actionHistory := PokemonBuilderActionHistory{SelectPokemonBuilder: selectPB, LegalPokemonBuilders: legalPBs}
	return pokemon.Moveset, actionHistory, nil
}

func (pb PokemonBuilder) BuildNature(natures Natures, pokemon Pokemon, team Team, finalWsFunc func([]float64) []float64, random *rand.Rand) (Nature, PokemonBuilderActionHistory, error) {
	if pokemon.Nature != "" {
		return pokemon.Nature, PokemonBuilderActionHistory{}, nil
	}

	natureWithAverageW := NatureWithFloat64{}
	natureWithSelectPBCK := map[Nature]*PokemonBuildCombinationKnowledge{}
	legalPB := make(PokemonBuilder, 0, len(pb))
	nextPokemon := pokemon

	for _, nature := range natures {
		nextPokemon.Nature = nature
		diffCalcPB := pb.DiffCalc(&pokemon, &nextPokemon, team, random)
		averageW := diffCalcPB.AverageW()
		selectIndex := random.Intn(len(diffCalcPB))

		natureWithAverageW[nature] = averageW
		natureWithSelectPBCK[nature] = diffCalcPB[selectIndex]
		legalPB = append(legalPB, diffCalcPB...)
	}

	if len(natureWithAverageW) == 0 {
		errMsg := fmt.Sprintf("pokemon.Name = %v の状態で、次の性格の組み合わせが見つからなかった", pokemon.Name)
		return "", PokemonBuilderActionHistory{}, fmt.Errorf(errMsg)
	}

	natures, averageWs := natureWithAverageW.KeysAndValues()
	finalWs := finalWsFunc(averageWs)
	selectIndex := omw.RandomIntWithWeight(finalWs, random)
	selectNature := natures[selectIndex]
	selectPBCK := natureWithSelectPBCK[selectNature]
	actionHistory := PokemonBuilderActionHistory{SelectPokemonBuilder:PokemonBuilder{selectPBCK}, LegalPokemonBuilders:[]PokemonBuilder{legalPB}}

	return selectNature, actionHistory, nil	
}

func (pb PokemonBuilder) Run(pokemon Pokemon, team Team, pbCommonK *PokemonBuildCommonKnowledge, finalWsFunc func([]float64) []float64, random *rand.Rand) (Pokemon, PokemonBuilderActionHistory, error) {
	if pokemon.Name == "" {
		return Pokemon{}, PokemonBuilderActionHistory{}, fmt.Errorf("pokemon.Name が 空の状態で、PokemonBuilder.Run は 実行出来ない")
	}

	pokeData := POKEDEX[pokemon.Name]
	abilities := pokeData.AllAbilities
	items := pbCommonK.Items
	moveNames := pbCommonK.MoveNames
	natures := pbCommonK.Natures

	selectPB := make(PokemonBuilder, 0, 8)
	legalPBs := make([]PokemonBuilder, 0, 8)

	ability, abilityPBAH, err := pb.BuildAbility(abilities, pokemon, team, finalWsFunc, random)
	if err != nil {
		return Pokemon{}, PokemonBuilderActionHistory{}, err
	}

	pokemon.Ability = ability

	item, itemPBAH, err := pb.BuildItem(items, pokemon, team, finalWsFunc, random)
	if err != nil {
		return Pokemon{}, PokemonBuilderActionHistory{}, err
	}

	pokemon.Item = item

	moveset, movesetPBAH, err := pb.BuildMoveset(moveNames, pokemon, team, finalWsFunc, random)
	if err != nil {
		return Pokemon{}, PokemonBuilderActionHistory{}, err
	}

	pokemon.Moveset = moveset

	nature, naturePBAH, err := pb.BuildNature(natures, pokemon, team, finalWsFunc, random)
	if err != nil {
		return Pokemon{}, PokemonBuilderActionHistory{}, err
	}

	pokemon.Nature = nature

	selectPB = append(selectPB, abilityPBAH.SelectPokemonBuilder...)
	legalPBs = append(legalPBs, abilityPBAH.LegalPokemonBuilders...)

	selectPB = append(selectPB, itemPBAH.SelectPokemonBuilder...)
	legalPBs = append(legalPBs, itemPBAH.LegalPokemonBuilders...)

	selectPB = append(selectPB, movesetPBAH.SelectPokemonBuilder...)
	legalPBs = append(legalPBs, movesetPBAH.LegalPokemonBuilders...)
	
	selectPB = append(selectPB, naturePBAH.SelectPokemonBuilder...)
	legalPBs = append(legalPBs, naturePBAH.LegalPokemonBuilders...)

	actionHistory := PokemonBuilderActionHistory{SelectPokemonBuilder:selectPB, LegalPokemonBuilders:legalPBs}
	return pokemon, actionHistory, nil
}

type PokemonBuilderActionHistory struct {
	SelectPokemonBuilder PokemonBuilder
	LegalPokemonBuilders []PokemonBuilder
}

func (pbah *PokemonBuilderActionHistory) Optimizer(learningRate float64) error {
	legalPBs := pbah.LegalPokemonBuilders

	for i, selectPBCK := range pbah.SelectPokemonBuilder {
		legalPB := legalPBs[i]
		selectIndex := legalPB.Index(selectPBCK)
		if selectIndex == -1 {
			return fmt.Errorf("PokemonBuilderActionHistory の 整合性 が とれていない")
		}
		legalPBLength := len(legalPB)

		for i := 0; i < legalPBLength; i++ {
			var moveingV float64
			if i == selectIndex {
				//1,0に近似
				moveingV = legalPB[i].W - 1.0
				moveingV = moveingV * learningRate
			} else {
				//0.0に近似
				moveingV = legalPB[i].W * learningRate
			}
			legalPB[i].W -= moveingV
		}
	}
	return nil
}

type PokemonBuilderActionHistories []PokemonBuilderActionHistory

func (pbahs PokemonBuilderActionHistories) RandomChoices(size int, random *rand.Rand) PokemonBuilderActionHistories {
	length := len(pbahs)
	if length == 0 {
		return PokemonBuilderActionHistories{}
	}

	result := make(PokemonBuilderActionHistories, size)
	for i := 0; i < size; i++ {
		index := random.Intn(length)
		result[i] = pbahs[index]
	}
	return result
}

func (actionHistories PokemonBuilderActionHistories) Optimizer(iterNum, miniBatchSize int, learningRate float64, random *rand.Rand) error {
	for i := 0; i < iterNum; i++ {
		for _, actionHistory := range actionHistories.RandomChoices(miniBatchSize, random) {
			err := actionHistory.Optimizer(learningRate)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
