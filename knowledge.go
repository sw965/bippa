package bippa

import (
	"fmt"
	"github.com/sw965/omw"
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

type PokemonBuildKnowledge struct {
	SelfAbility    Ability
	SelfItem       Item
	SelfMoveNames  MoveNames
	SelfNature     Nature
	CombinationNum int
	Tier           Tier
}

func (pbk PokemonBuildKnowledge) Init() PokemonBuildKnowledge {
	combinationNum := 0

	if pbk.SelfAbility != "" {
		combinationNum += 1
	}

	if pbk.SelfItem != "" {
		combinationNum += 1
	}

	combinationNum += len(pbk.SelfMoveNames)

	if pbk.SelfNature != "" {
		combinationNum += 1
	}
	pbk.CombinationNum = combinationNum
	return pbk
}

func (pbk *PokemonBuildKnowledge) All(pokemon *Pokemon, team Team) bool {
	results := make([]bool, 0, pbk.CombinationNum)

	if pbk.SelfAbility != "" {
		results = append(results, pokemon.Ability == pbk.SelfAbility)
	}

	if pbk.SelfItem != "" {
		results = append(results, pokemon.Item == pbk.SelfItem)
	}

	if len(pbk.SelfMoveNames) != 0 {
		for _, moveName := range pbk.SelfMoveNames {
			_, ok := pokemon.Moveset[moveName]
			results = append(results, ok)
		}
	}

	if pbk.SelfNature != "" {
		results = append(results, pokemon.Nature == pbk.SelfNature)
	}

	return omw.All(results...)
}

func (pbk1 *PokemonBuildKnowledge) NearlyEqual(pbk2 *PokemonBuildKnowledge) bool {
	if pbk1.SelfAbility != pbk2.SelfAbility {
		return false
	}

	if pbk1.SelfItem != pbk2.SelfItem {
		return false
	}

	sortedMoveNames1 := pbk1.SelfMoveNames.Sort()
	sortedMoveNames2 := pbk2.SelfMoveNames.Sort()

	if !sortedMoveNames1.Equal(sortedMoveNames2) {
		return false
	}

	if pbk1.SelfNature != pbk2.SelfNature {
		return false
	}
	return true
}

type PokemonBuilder []PokemonBuildKnowledge

func NewInitPokemonBuilder(selfAbilities Abilities, selfItems Items, selfMoveNames MoveNames, selfNatures Natures, random *rand.Rand) PokemonBuilder {
	result := PokemonBuilder{}

	for _, ability := range selfAbilities {
		result = append(result, PokemonBuildKnowledge{SelfAbility: ability, Tier: ALL_TIERS.RandomChoice(random)}.Init())
	}

	for _, item := range selfItems {
		result = append(result, PokemonBuildKnowledge{SelfItem: item, Tier: ALL_TIERS.RandomChoice(random)}.Init())
	}

	for _, moveName := range selfMoveNames {
		result = append(result, PokemonBuildKnowledge{SelfMoveNames: MoveNames{moveName}, Tier: ALL_TIERS.RandomChoice(random)}.Init())
	}

	for _, nature := range selfNatures {
		result = append(result, PokemonBuildKnowledge{SelfNature: nature, Tier: ALL_TIERS.RandomChoice(random)}.Init())
	}
	return result
}

func (pb PokemonBuilder) Init() PokemonBuilder {
	result := PokemonBuilder{}
	for tier, pbks := range pb {
		result[tier] = pbks.Init()
	}
	return result
}

func (pb1 PokemonBuilder) Cross(pb2 PokemonBuilder, random *rand.Rand) PokemonBuilder {
	length := len(pb1)
	result := make(PokemonBuilder, length)
	for i := 0; i < length; i++ {
		if omw.RandomBool(random) {
			result[i] = pb1[i]
		} else {
			result[i] = pb2[i]
		}
	}
	return result
}

func (pb PokemonBuilder) Mutation(num int, random *rand.Rand) (PokemonBuilder, error) {
	length := len(pb)
	mutationIndices, err := omw.MakeRandomSliceInt(num, 0, length, random)

	if err != nil {
		return PokemonBuilder{}, err
	}

	result := make(PokemonBuilder, length)
	for i := 0; i < length; i++ {
		result[i] = pb[i]
	}

	for _, index := range mutationIndices {
		result[index].Tier = ALL_TIERS.RandomChoice(random)
	}

	return result, nil
}

func (pb PokemonBuilder) Filter(f func(*PokemonBuildKnowledge) bool) PokemonBuilder {
	result := make(PokemonBuilder, 0, len(pb))
	for _, pbk := range pb {
		if f(&pbk) {
			result = append(result, pbk)
		}
	}
	return result
}

func (pb PokemonBuilder) MaxCombinationNum() int {
	result := pb[0].CombinationNum
	for _, pbk := range pb[1:] {
		maxCombinationNum := pbk.CombinationNum
		if maxCombinationNum > result {
			result = maxCombinationNum
		}
	}
	return result
}

func (pb PokemonBuilder) DiffCalcTier(pokemon, nextPokemon *Pokemon, team Team, random *rand.Rand) Tier {
	//一つ前の状態(pokemon)が満たしている組み合わせを排除する(差分を見る為に)
	pb = pb.Filter(func(pbk *PokemonBuildKnowledge) bool { return !pbk.All(pokemon, team) })

	//次の状態(nextPokemon)が満たしている組み合わせを取り出す
	pb = pb.Filter(func(pbk *PokemonBuildKnowledge) bool { return pbk.All(nextPokemon, team) })

	//組み合わせ数が最も多い知識を取り出す
	maxCombinationNum := pb.MaxCombinationNum()
	pb = pb.Filter(func(pbk *PokemonBuildKnowledge) bool { return pbk.CombinationNum == maxCombinationNum })

	tierCounter := TierCounter{}
	for _, pbk := range pb {
		tierCounter[pbk.Tier] += 1
	}

	results := tierCounter.Tiers()
	return results.RandomChoice(random)
}

func (pb PokemonBuilder) BuildMoveset(moveNames MoveNames, pokemon Pokemon, team Team, random *rand.Rand) (Moveset, error) {
	getMoveName := func(nextPokemon Pokemon) (MoveName, error) {
		moveNameWithTier := MoveNameWithTier{}
		for _, moveName := range moveNames {
			moveset := pokemon.Moveset.Copy()

			if _, ok := moveset[moveName]; ok {
				continue
			}

			moveset[moveName] = &PowerPoint{}
			nextPokemon.Moveset = moveset

			tier := pb.DiffCalcTier(&pokemon, &nextPokemon, team, random)
			moveNameWithTier[moveName] = tier
		}

		if len(moveNameWithTier) == 0 {
			errMsg := fmt.Sprintf("pokemon.Name = %v pokemon.Moveset.Keys() = %v の状態で、次の技の組み合わせが見つからなかった",
				pokemon.Name, pokemon.Moveset.Keys(),
			)
			return "", fmt.Errorf(errMsg)
		}

		moveNames, tiers := moveNameWithTier.KeysAndValues()
		index := omw.RandomIntWithWeight(tiers.SelectPercents(), random)
		return moveNames[index], nil
	}

	padNum := MAX_MOVESET_LENGTH - len(pokemon.Moveset)
	for i := 0; i < padNum; i++ {
		moveName, err := getMoveName(pokemon)
		if err != nil {
			return Moveset{}, err
		}
		moveset := pokemon.Moveset.Copy()
		powerPoint := NewPowerPoint(MOVEDEX[moveName].BasePP, MAX_POINT_UP)
		moveset[moveName] = &powerPoint
		pokemon.Moveset = moveset
	}
	return pokemon.Moveset, nil
}

type PokemonBuilders []PokemonBuilder

func (pbs PokemonBuilders) AccuracyYs(accuracy func(PokemonBuilder) float64) []float64 {
	result := make([]float64, len(pbs))
	for i, pb := range pbs {
		result[i] = accuracy(pb)
	}
	return result
}

func (pbs PokemonBuilders) RandomChoice(random *rand.Rand) PokemonBuilder {
	index := random.Intn(len(pbs))
	return pbs[index]
}

func (pbs PokemonBuilders) Elite(accuracyYs []float64) PokemonBuilders {
	maxAccuracy := omw.MaxFloat64(accuracyYs...)
	result := make(PokemonBuilders, 0, len(pbs))
	for i, pb := range pbs {
		accuracyY := accuracyYs[i]
		if accuracyY == maxAccuracy {
			result = append(result, pb)
		}
	}
	return result
}

func (pbs PokemonBuilders) RouletteSelect(accuracyYs []float64, random *rand.Rand) PokemonBuilder {
	index := omw.RandomIntWithWeight(accuracyYs, random)
	return pbs[index]
}

func (pbs PokemonBuilders) NextGeneration(accuracy func(PokemonBuilder) float64, tournamentSize int, crossPercent, mutationPercent float64, mutationNum int, random *rand.Rand) (PokemonBuilders, error) {
	selectPercent := 1.0 - (crossPercent + mutationPercent)
	if selectPercent < 0.0 {
		return PokemonBuilders{}, fmt.Errorf("交叉確率 + 突然変異確率 <= 1.0 でなければならない")
	}

	length := len(pbs)
	accuracyYs := pbs.AccuracyYs(accuracy)
	weight := []float64{selectPercent, crossPercent, mutationPercent}

	result := make(PokemonBuilders, 0, length)
	result = append(result, pbs.Elite(accuracyYs).RandomChoice(random))

	for i := 0; i < length-1; i++ {
		index := omw.RandomIntWithWeight(weight, random)
		switch index {
		case 0:
			result = append(result, pbs.RouletteSelect(accuracyYs, random))
		case 1:
			pb1 := pbs.RouletteSelect(accuracyYs, random)
			pb2 := pbs.RouletteSelect(accuracyYs, random)
			result = append(result, pb1.Cross(pb2, random))
		default:
			pb := pbs.RandomChoice(random)
			pb, err := pb.Mutation(mutationNum, random)
			if err != nil {
				return PokemonBuilders{}, err
			}
			result = append(result, pb)
		}
	}
	return result, nil
}
