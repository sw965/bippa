package bippa

import (
	"fmt"
	"math/rand"
	"github.com/sw965/omw"
)

type PokemonBuildCombinationKnowledge struct {
	SelfAbility    Ability
	SelfItem       Item
	SelfMoveNames  MoveNames
	SelfNature     Nature
	CombinationNum int
	W float64
}

func (pbck PokemonBuildCombinationKnowledge) Init(random *rand.Rand) PokemonBuildCombinationKnowledge {
	combinationNum := 0

	if pbck.SelfAbility != "" {
		combinationNum += 1
	}

	if pbck.SelfItem != "" {
		combinationNum += 1
	}

	combinationNum += len(pbck.SelfMoveNames)

	if pbck.SelfNature != "" {
		combinationNum += 1
	}

	pbck.CombinationNum = combinationNum
	pbck.W = 0.5
	return pbck
}

func (pbck *PokemonBuildCombinationKnowledge) All(pokemon *Pokemon, team Team) bool {
	results := make([]bool, 0, pbck.CombinationNum)

	if pbck.SelfAbility != "" {
		results = append(results, pokemon.Ability == pbck.SelfAbility)
	}

	if pbck.SelfItem != "" {
		results = append(results, pokemon.Item == pbck.SelfItem)
	}

	if len(pbck.SelfMoveNames) != 0 {
		for _, moveName := range pbck.SelfMoveNames {
			_, ok := pokemon.Moveset[moveName]
			results = append(results, ok)
		}
	}

	if pbck.SelfNature != "" {
		results = append(results, pokemon.Nature == pbck.SelfNature)
	}

	return omw.All(results...)
}

func (pbck1 *PokemonBuildCombinationKnowledge) NearlyEqual(pbck2 *PokemonBuildCombinationKnowledge) bool {
	if pbck1.SelfAbility != pbck2.SelfAbility {
		return false
	}

	if pbck1.SelfItem != pbck2.SelfItem {
		return false
	}

	sortedMoveNames1 := pbck1.SelfMoveNames.Sort()
	sortedMoveNames2 := pbck2.SelfMoveNames.Sort()

	if !sortedMoveNames1.Equal(sortedMoveNames2) {
		return false
	}

	if pbck1.SelfNature != pbck2.SelfNature {
		return false
	}
	return true
}

type PokemonBuilder []PokemonBuildCombinationKnowledge

func NewInitPokemonBuilder(selfAbilities Abilities, selfItems Items, selfMoveNames MoveNames, selfNatures Natures, random *rand.Rand) PokemonBuilder {
	result := PokemonBuilder{}

	for _, ability := range selfAbilities {
		result = append(result, PokemonBuildCombinationKnowledge{SelfAbility: ability})
	}

	for _, item := range selfItems {
		result = append(result, PokemonBuildCombinationKnowledge{SelfItem: item})
	}

	for _, moveName := range selfMoveNames {
		result = append(result, PokemonBuildCombinationKnowledge{SelfMoveNames: MoveNames{moveName}})
	}

	for _, nature := range selfNatures {
		result = append(result, PokemonBuildCombinationKnowledge{SelfNature: nature})
	}

	for _, ability := range selfAbilities {
		for _, item := range selfItems {
			result = append(result, PokemonBuildCombinationKnowledge{SelfAbility:ability, SelfItem:item})
		}
	}

	for _, ability := range selfAbilities {
		for _, moveName := range selfMoveNames {
			result = append(result, PokemonBuildCombinationKnowledge{SelfAbility:ability, SelfMoveNames:MoveNames{moveName}})
		}
	}

	for _, ability := range selfAbilities {
		for _, nature := range selfNatures {
			result = append(result, PokemonBuildCombinationKnowledge{SelfAbility:ability, SelfNature:nature})
		}
	}

	for _, item := range selfItems {
		for _, moveName := range selfMoveNames {
			result = append(result, PokemonBuildCombinationKnowledge{SelfItem:item, SelfMoveNames:MoveNames{moveName}})
		}
	}

	for _, item := range selfItems {
		for _, nature := range selfNatures {
			result = append(result, PokemonBuildCombinationKnowledge{SelfItem:item, SelfNature:nature})
		}
	}

	for _, moveName1 := range selfMoveNames {
		for _, moveName2 := range selfMoveNames {
			if moveName1 == moveName2 {
				continue
			}
			pb := PokemonBuildCombinationKnowledge{SelfMoveNames:MoveNames{moveName1, moveName2}}
			if result.In(&pb) {
				continue
			}
			result = append(result, pb)
		}
	}

	for _, moveName := range selfMoveNames {
		for _, nature := range selfNatures {
			result = append(result, PokemonBuildCombinationKnowledge{SelfMoveNames:MoveNames{moveName}, SelfNature:nature})
		}
	}

	return result
}

func (pb PokemonBuilder) Init(random *rand.Rand) PokemonBuilder {
	result := make(PokemonBuilder, len(pb))
	for i, pbck := range pb {
		result[i] = pbck.Init(random)
	}
	return result
}

func (pb PokemonBuilder) Copy() PokemonBuilder {
	result := make(PokemonBuilder, len(pb))
	for i, pbck := range pb {
		result[i] = pbck
	}
	return result
}

func (pb PokemonBuilder) In(pbck *PokemonBuildCombinationKnowledge) bool {
	for _, iPBCK := range pb {
		if iPBCK.NearlyEqual(pbck) {
			return true
		}
	}
	return false
}

func (pb PokemonBuilder) Index(pbck *PokemonBuildCombinationKnowledge) int {
	for i, iPBCK := range pb {
		if iPBCK.NearlyEqual(pbck) {
			return i
		}
	}
	return -1
}

func (pb1 PokemonBuilder) Indices(pb2 PokemonBuilder) []int {
	result := make([]int, 0, len(pb1))
	for _, pbck := range pb2 {
		result = append(result, pb1.Index(&pbck))
	}
	return result
}

func (pb PokemonBuilder) Ws() []float64 {
	result := make([]float64, len(pb))
	for i, pbck := range pb {
		result[i] = pbck.W
	}
	return result
}

func (pb PokemonBuilder) Filter(f func(*PokemonBuildCombinationKnowledge) bool) PokemonBuilder {
	result := make(PokemonBuilder, 0, len(pb))
	for _, pbck := range pb {
		if f(&pbck) {
			result = append(result, pbck)
		}
	}
	return result
}

func (pb PokemonBuilder) MaxCombinationNum() int {
	result := pb[0].CombinationNum
	for _, pbck := range pb[1:] {
		maxCombinationNum := pbck.CombinationNum
		if maxCombinationNum > result {
			result = maxCombinationNum
		}
	}
	return result
}

func (pb PokemonBuilder) DiffCalc(pokemon, nextPokemon *Pokemon, team Team, random *rand.Rand) PokemonBuilder {
	//一つ前の状態(pokemon)が満たしている組み合わせを排除する(差分を見る為に)
	pb = pb.Filter(func(pbck *PokemonBuildCombinationKnowledge) bool { return !pbck.All(pokemon, team) })

	//次の状態(nextPokemon)が満たしている組み合わせを取り出す
	pb = pb.Filter(func(pbck *PokemonBuildCombinationKnowledge) bool { return pbck.All(nextPokemon, team) })

	//組み合わせ数が最も多い知識を取り出す
	maxCombinationNum := pb.MaxCombinationNum()
	return pb.Filter(func(pbck *PokemonBuildCombinationKnowledge) bool { return pbck.CombinationNum == maxCombinationNum })
}

func (pb PokemonBuilder) BuildMoveset(moveNames MoveNames, pokemon Pokemon, team Team, getWs func(PokemonBuilder) []float64, actionNum int, random *rand.Rand) (Moveset, PokemonBuilderActionHistory, error) {
	if actionNum > MAX_MOVESET_LENGTH {
		return Moveset{}, PokemonBuilderActionHistory{}, fmt.Errorf("BuildMovesetのactionNumは、4以下でなければならない")
	}

	actionNum = omw.MinInt(actionNum, MAX_MOVESET_LENGTH - len(pokemon.Moveset))

	action := func(nextPokemon Pokemon) (MoveName, PokemonBuildCombinationKnowledge, PokemonBuilder, error) {
		moveNameWithW := MoveNameWithFloat64{}
		moveNameWithSelectPBCK := map[MoveName]PokemonBuildCombinationKnowledge{}
		legalPB := make(PokemonBuilder, 0, len(pb))

		for _, moveName := range moveNames {
			moveset := pokemon.Moveset.Copy()

			if _, ok := moveset[moveName]; ok {
				continue
			}

			moveset[moveName] = &PowerPoint{}
			nextPokemon.Moveset = moveset
			
			diffCalcPB := pb.DiffCalc(&pokemon, &nextPokemon, team, random)
			//ある一つの技について、複数の知識を活用できる場合は、ランダムに選択する
			//例：「ギガドレイン, ヘドロばくだん」を既に覚えているとし、まもるについての知識として、「ギガドレイン, まもる」 [ヘドロばくだん, まもる] の 組み合わせ知識を活用できる場合、
			//ギガドレインとの組み合わせを考慮するか、ヘドロばくだんとの組み合わせを考慮するかはランダムで選ぶ
			ws := getWs(diffCalcPB)
			selectIndex := random.Intn(len(diffCalcPB))

			moveNameWithW[moveName] = ws[selectIndex]
			moveNameWithSelectPBCK[moveName] = diffCalcPB[selectIndex]
			legalPB = append(legalPB, diffCalcPB...)
		}

		if len(moveNameWithW) == 0 {
			errMsg := fmt.Sprintf("pokemon.Name = %v pokemon.Moveset.Keys() = %v の状態で、次の技の組み合わせが見つからなかった",
				pokemon.Name, pokemon.Moveset.Keys(),
			)
			return "", PokemonBuildCombinationKnowledge{}, PokemonBuilder{}, fmt.Errorf(errMsg)
		}
		
		//それそれの技で、活用した知識の重みに応じて、技を選ぶ
		moveNames, ws := moveNameWithW.KeysAndValues()
		index := omw.RandomIntWithWeight(ws, random)
		selectMoveName := moveNames[index]

		return selectMoveName, moveNameWithSelectPBCK[selectMoveName], legalPB, nil
	}

	selectPokemonBuilder := make(PokemonBuilder, actionNum)
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

		selectPokemonBuilder[i] = selectPBCK
		legalPBs[i] = legalPB
	}

	actionHistory := PokemonBuilderActionHistory{SelectPokemonBuilder:selectPokemonBuilder, LegalPokemonBuilders:legalPBs}
	return pokemon.Moveset, actionHistory, nil
}

func (pb PokemonBuilder) Optimizer(pbah *PokemonBuilderActionHistory, learningRate float64) error {
	for i, selectPBCK := range pbah.SelectPokemonBuilder {
		selectIndex := pb.Index(&selectPBCK)
		if selectIndex == -1 {
			return fmt.Errorf("selectPBCKが見つからなかった")
		}

		legalPB := pbah.LegalPokemonBuilders[i]
		legalIndices := pb.Indices(legalPB)

		for _, legalIndex := range legalIndices {
			if legalIndex == -1 {
				return fmt.Errorf("legalPBが見つからなかった")
			}

			var moveingV float64

			if legalIndex == selectIndex {
				//1,0に近似
				moveingV = pb[legalIndex].W - 1.0
				moveingV = moveingV * learningRate
			} else {
				//0.0に近似
				moveingV = pb[legalIndex].W * learningRate
			}
			pb[legalIndex].W -= moveingV
		}
	}
	return nil
}

type PokemonBuilderActionHistory struct {
	SelectPokemonBuilder PokemonBuilder
	LegalPokemonBuilders []PokemonBuilder
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