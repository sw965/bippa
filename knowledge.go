package bippa

import (
	"fmt"
	"github.com/sw965/omw"
	"math"
	"math/rand"
)

type Tier int

const (
	TIER1 = Tier(1)
	TIER2 = Tier(2)
	TIER3 = Tier(3)
	TIER4 = Tier(4)
	TIER5 = Tier(5)
)

var TIER_SELECT_PERCENT = map[Tier]float64{
	TIER1: 0.7,
	TIER2: 0.15,
	TIER3: 0.10,
	TIER4: 0.04,
	TIER5: 0.01,
}

type Tiers []Tier

var ALL_TIERS = Tiers{TIER1, TIER2, TIER3, TIER4, TIER5}

func (tiers Tiers) RandomChoice(random *rand.Rand) Tier {
	index := random.Intn(len(tiers))
	return tiers[index]
}

func (tiers Tiers) SelectPercents() []float64 {
	result := make([]float64, 0, len(tiers))
	for _, tier := range tiers {
		result = append(result, TIER_SELECT_PERCENT[tier])
	}
	return result
}

type TierCounter map[Tier]int

func (tierWithCounter TierCounter) Values() []int {
	result := make([]int, 0, len(tierWithCounter))
	for _, v := range tierWithCounter {
		result = append(result, v)
	}
	return result
}

func (tierWithCounter TierCounter) Tiers() Tiers {
	length := omw.SumInt(tierWithCounter.Values()...)
	result := make(Tiers, 0, length)
	for tier, count := range tierWithCounter {
		for i := 0; i < count; i++ {
			result = append(result, tier)
		}
	}
	return result
}

type PokemonBuildCombinationKnowledge struct {
	SelfAbility    Ability
	SelfItem       Item
	SelfMoveNames  MoveNames
	SelfNature     Nature
	CombinationNum int
	SelectPercent float64
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
	pbck.SelectPercent = 0.5
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
			result = append(result, PokemonBuildCombinationKnowledge{SelfMoveNames:MoveNames{moveName1, moveName2}})
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

func (pb PokemonBuilder) SelectPercents() []float64 {
	result := make([]float64, len(pb))
	for i, pbck := range pb {
		result[i] = pbck.SelectPercent
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

func (pb PokemonBuilder) BuildMoveset(moveNames MoveNames, pokemon Pokemon, team Team, getSelectPercents func(PokemonBuilder) []float64, random *rand.Rand) (Moveset, PokemonBuilderActionHistory, error) {
	action := func(nextPokemon Pokemon) (MoveName, PokemonBuildCombinationKnowledge, PokemonBuilder, error) {
		moveNameWithSelectPercent := MoveNameWithFloat64{}
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
			selectPercent := getSelectPercents(diffCalcPB)
			selectIndices, err := omw.MakeSliceInt(0, len(diffCalcPB), 1)

			if err != nil {
				return "", PokemonBuildCombinationKnowledge{}, PokemonBuilder{}, err
			}

			selectIndex := omw.RandomChoiceInt(random, selectIndices...)

			moveNameWithSelectPercent[moveName] = selectPercent[selectIndex]
			moveNameWithSelectPBCK[moveName] = diffCalcPB[selectIndex]
			legalPB = append(legalPB, diffCalcPB...)
		}

		if len(moveNameWithSelectPercent) == 0 {
			errMsg := fmt.Sprintf("pokemon.Name = %v pokemon.Moveset.Keys() = %v の状態で、次の技の組み合わせが見つからなかった",
				pokemon.Name, pokemon.Moveset.Keys(),
			)
			return "", PokemonBuildCombinationKnowledge{}, PokemonBuilder{}, fmt.Errorf(errMsg)
		}
		
		moveNames, selectPercents := moveNameWithSelectPercent.KeysAndValues()
		index := omw.RandomIntWithWeight(selectPercents, random)
		selectMoveName := moveNames[index]

		return selectMoveName, moveNameWithSelectPBCK[selectMoveName], legalPB, nil
	}

	padNum := MAX_MOVESET_LENGTH - len(pokemon.Moveset)
	selectPokemonBuilder := make(PokemonBuilder, padNum)
	legalPBs := make([]PokemonBuilder, padNum)

	for i := 0; i < padNum; i++ {
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

func (pb PokemonBuilder) Optimizer(teacher *PokemonBuilderTeacher, target, learningRate float64) error {
	if !teacher.OK(teacher.Pokemon, teacher.Team) {
		return nil
	}

	actionHistory := teacher.ActionHistory
	actionNum := len(actionHistory.SelectPokemonBuilder)

	for i, selectPBCK := range actionHistory.SelectPokemonBuilder {
		selectIndex := pb.Index(&selectPBCK)
		if selectIndex == -1 {
			return fmt.Errorf("selectPBCKが見つからなかった")
		}

		legalPB := actionHistory.LegalPokemonBuilders[i]
		legalActionNum := len(legalPB)
		legalIndices := pb.Indices(legalPB)

		for _, legalIndex := range legalIndices {
			if legalIndex == -1 {
				return fmt.Errorf("legalPBが見つからなかった")
			}

			finalTarget := math.Pow(target, 1.0 / float64(actionNum))
			var moveingV float64

			if legalIndex == selectIndex {
				moveingV = pb[legalIndex].SelectPercent - finalTarget
				moveingV = moveingV * learningRate
			} else {
				moveingV = 1.0 - finalTarget
				moveingV =  pb[legalIndex].SelectPercent - (moveingV / float64(legalActionNum - 1))
				moveingV *= learningRate
			}
			pb[legalIndex].SelectPercent -= moveingV
		}
	}
	return nil
}

type PokemonBuilderActionHistory struct {
	SelectPokemonBuilder PokemonBuilder
	LegalPokemonBuilders []PokemonBuilder
}

type PokemonBuilderTeacher struct {
	ActionHistory PokemonBuilderActionHistory
	Pokemon Pokemon
	Team Team
	OK func(Pokemon, Team) bool
}
