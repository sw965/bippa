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
	TIER6 = Tier(6)
)

type Tiers []Tier

var ALL_TIERS = Tiers{TIER1, TIER2, TIER3, TIER4, TIER5, TIER6}

func (tiers Tiers) RandomChoice(random *rand.Rand) Tier {
	index := random.Intn(len(tiers))
	return tiers[index]
}

type TierWithInt map[Tier]int

func (tierWithInt TierWithInt) Values() []int {
	result := make([]int, 0, len(tierWithInt))
	for _, v := range result {
		result = append(result, v)
	}
	return result
}

func (tierWithInt TierWithInt) TiersOfMaxValue() Tiers {
	result := make(Tiers, 0, len(tierWithInt))
	maxValue := omw.MaxInt(tierWithInt.Values()...)

	for tier, v := range tierWithInt {
		if v == maxValue {
			result = append(result, tier)
		}
	}
	return result
}

type TierWithFloat64 map[Tier]float64

var TIER_WITH_SELECT_PERCENT = TierWithFloat64{
	TIER1: 0.7,
	TIER2: 0.2,
	TIER3: 0.065,
	TIER4: 0.025,
	TIER5: 0.01,
}

func (tierWithFloat64 TierWithFloat64) Keys() Tiers {
	result := make(Tiers, 0, len(tierWithFloat64))
	for k, _ := range tierWithFloat64 {
		result = append(result, k)
	}
	return result
}

func (tierWithFloat64 TierWithFloat64) Items() (Tiers, []float64) {
	keys := tierWithFloat64.Keys()
	values := make([]float64, len(tierWithFloat64))
	for i, key := range keys {
		values[i] = tierWithFloat64[key]
	}
	return keys, values
}

func (tierWithFloat64 TierWithFloat64) TierRandomChoiceWithWeight(random *rand.Rand) Tier {
	keys, values := tierWithFloat64.Items()
	index := omw.RandomIntWithWeight(values, random)
	return keys[index]
}

type PokemonBuildEvent struct {
	Combination    func(*Pokemon, Team) []bool
	All            bool
	CombinationNum int
}

func NewPokemonBuildSelfMoveNameEvent(moveName ...MoveName) PokemonBuildEvent {
	combination := func(pokemon *Pokemon, team Team) []bool {
		result := make([]bool, len(moveName))
		for i, iMoveName := range moveName {
			_, ok := pokemon.Moveset[iMoveName]
			result[i] = ok
		}
		return result
	}
	return PokemonBuildEvent{Combination: combination}
}

func (pbe *PokemonBuildEvent) Output(pokemon *Pokemon, team Team) PokemonBuildEvent {
	y := pbe.Combination(pokemon, team)
	return PokemonBuildEvent{Combination: pbe.Combination, All: omw.All(y...), CombinationNum: len(y)}
}

type PokemonBuildEvents []PokemonBuildEvent

func (pbes PokemonBuildEvents) Output(pokemon *Pokemon, team Team) PokemonBuildEvents {
	result := make(PokemonBuildEvents, len(pbes))
	for i, pbe := range pbes {
		result[i] = pbe.Output(pokemon, team)
	}
	return result
}

func (pbes PokemonBuildEvents) Filter(f func(*PokemonBuildEvent) bool) PokemonBuildEvents {
	result := make(PokemonBuildEvents, 0, len(pbes))
	for _, pbe := range pbes {
		if f(&pbe) {
			result = append(result, pbe)
		}
	}
	return result
}

func (pbes PokemonBuildEvents) AnyAll() bool {
	for _, pbe := range pbes {
		if pbe.All {
			return true
		}
	}
	return false
}

func (pbes PokemonBuildEvents) MaxCombinationNum() int {
	result := 0
	for _, pbe := range pbes {
		combinationNum := pbe.CombinationNum
		if combinationNum > result {
			result = combinationNum
		}
	}
	return result
}

type PBEsWithTier map[Tier]PokemonBuildEvents

func (pbesWithTier PBEsWithTier) Output(pokemon *Pokemon, team Team) PBEsWithTier {
	result := PBEsWithTier{}
	for tier, pbes := range pbesWithTier {
		result[tier] = pbes.Output(pokemon, team)
	}
	return result
}

func (pbesWithTier PBEsWithTier) Filter(f func(Tier, *PokemonBuildEvent) bool) PBEsWithTier {
	result := PBEsWithTier{}
	for tier, pbes := range pbesWithTier {
		newPBEs := make(PokemonBuildEvents, 0, len(pbes))
		for _, pbe := range pbes {
			if f(tier, &pbe) {
				newPBEs = append(newPBEs, pbe)
			}
		}
		result[tier] = newPBEs
	}
	return result
}

func (pbesWithTier PBEsWithTier) MaxCombinationNum() int {
	result := 0
	for _, pbes := range pbesWithTier {
		maxCombinationNum := pbes.MaxCombinationNum()
		if maxCombinationNum > result {
			result = maxCombinationNum
		}
	}
	return result
}

type PokemonBuildKnowledge struct {
	Natures        Natures
	MoveNames      MoveNames
	EventsWithTier PBEsWithTier
}

func NewVenusaurBuildKnowledge() PokemonBuildKnowledge {
	moveNames := MoveNames{"ギガドレイン", "ヘドロばくだん", "だいちのちから", "やどりぎのタネ", "どくどく", "まもる", "こうごうせい"}
	natures := Natures{"しんちょう", "ずぶとい", "ひかえめ"}

	tier1 := PokemonBuildEvents{
		NewPokemonBuildSelfMoveNameEvent("ギガドレイン"),
		NewPokemonBuildSelfMoveNameEvent("ギガドレイン", "ヘドロばくだん"),
		NewPokemonBuildSelfMoveNameEvent("やどりぎのタネ", "まもる"),
		NewPokemonBuildSelfMoveNameEvent("どくどく", "まもる"),
		NewPokemonBuildSelfMoveNameEvent("ギガドレイン", "ヘドロばくだん", "どくどく", "まもる"),
		NewPokemonBuildSelfMoveNameEvent("ギガドレイン", "ヘドロばくだん", "やどりぎのタネ", "まもる"),
	}

	tier2 := PokemonBuildEvents{
		NewPokemonBuildSelfMoveNameEvent("ギガドレイン", "ヘドロばくだん", "こうごうせい", "どくどく"),
		NewPokemonBuildSelfMoveNameEvent("ギガドレイン", "だいちのちから"),
		NewPokemonBuildSelfMoveNameEvent("ヘドロばくだん", "だいちのちから"),
		NewPokemonBuildSelfMoveNameEvent("ギガドレイン", "やどりぎのタネ"),
		NewPokemonBuildSelfMoveNameEvent("ギガドレイン", "どくどく"),
	}

	tier3 := PokemonBuildEvents{
		NewPokemonBuildSelfMoveNameEvent("ギガドレイン", "ヘドロばくだん", "だいちのちから", "こうごうせい"),
	}

	tier4 := PokemonBuildEvents{}

	tier5 := PokemonBuildEvents{
		NewPokemonBuildSelfMoveNameEvent("やどりぎのタネ", "どくどく", "まもる"),
		NewPokemonBuildSelfMoveNameEvent("やどりぎのタネ", "どくどく"),
	}

	tier6 := PokemonBuildEvents{
		NewPokemonBuildSelfMoveNameEvent("ギガドレイン", "だいちのちから", "ヘドロばくだん", "やどりぎのタネ"),
	}

	eventsWithTier := PBEsWithTier{TIER1: tier1, TIER2: tier2, TIER3: tier3, TIER4: tier4, TIER5: tier5, TIER6: tier6}
	return PokemonBuildKnowledge{Natures: natures, MoveNames: moveNames, EventsWithTier: eventsWithTier}
}

func (pbk PokemonBuildKnowledge) DiffCalcTier(pokemon *Pokemon, nextPokemon *Pokemon, team Team, random *rand.Rand) Tier {
	pbesWithTier := pbk.EventsWithTier.Output(pokemon, team)

	//既に条件を満たしている組み合わせを排除する
	pbesWithTier = pbesWithTier.Filter(func(tier Tier, pbe *PokemonBuildEvent) bool { return !pbe.All })
	pbesWithTier = pbesWithTier.Output(nextPokemon, team)

	if pbesWithTier[TIER6].AnyAll() {
		return TIER6
	}

	pbesWithTier = pbesWithTier.Filter(func(tier Tier, pbe *PokemonBuildEvent) bool { return pbe.All })
	maxCombinationNum := pbesWithTier.MaxCombinationNum()
	pbesWithTier = pbesWithTier.Filter(func(tier Tier, pbe *PokemonBuildEvent) bool { return pbe.CombinationNum == maxCombinationNum })

	tiersLength := 0
	for _, pbes := range pbesWithTier {
		tiersLength += len(pbes)
	}

	tiers := make(Tiers, 0, tiersLength)
	for tier, pbes := range pbesWithTier {
		for i := 0; i < len(pbes); i++ {
			tiers = append(tiers, tier)
		}
	}

	if len(tiers) == 0 {
		return TIER6
	}

	return tiers.RandomChoice(random)
}

func (pbk *PokemonBuildKnowledge) BuildMoveset(pokemon Pokemon, team Team, random *rand.Rand) (Pokemon, error) {
	getMoveName := func(nextPokemon Pokemon) (MoveName, error) {
		moveNameWithTier := MoveNameWithTier{}

		for _, moveName := range pbk.MoveNames {
			moveset := pokemon.Moveset.Copy()

			if _, ok := moveset[moveName]; ok {
				continue
			}

			moveset[moveName] = &PowerPoint{}
			nextPokemon.Moveset = moveset

			tier := pbk.DiffCalcTier(&pokemon, &nextPokemon, team, random)
			if tier == TIER6 {
				continue
			}

			moveNameWithTier[moveName] = tier
		}

		if len(moveNameWithTier) == 0 {
			errMsg := fmt.Sprintf("movesetKeys = %v の状態で、次の組み合わせが見つからなかった", pokemon.Moveset.Keys())
			return "", fmt.Errorf(errMsg)
		}
		return moveNameWithTier.MoveNameRandomChoiceWithTierWeight(random), nil
	}

	for i := 0; i < MAX_MOVESET_LENGTH; i++ {
		moveName, err := getMoveName(pokemon)
		if err != nil {
			return Pokemon{}, err
		}
		moveset := pokemon.Moveset.Copy()
		powerPoint := NewPowerPoint(MOVEDEX[moveName].BasePP, MAX_POINT_UP)
		moveset[moveName] = &powerPoint
		pokemon.Moveset = moveset
	}
	return pokemon, nil
}

// func (pbk PokemonBuildKnowledge) BuildPokemon(team Team) Pokemon {
// 	pokemon := Pokemon{}
// 	for i := 0; i < MAX_MOVESET_LENGTH; i++ {
// 		moveNameWithTier := MoveNameWithTier{}

// 		for _, moveName := range pbk.MoveNames {
// 			if _, ok := pokemon.Moveset[moveName]; ok {
// 				continue
// 			}

// 			powerPoint := NewPowerPoint(MOVEDEX[moveName].BasePP, MAX_POINT_UP)
// 			pokemon.Moveset[moveName] = &powerPoint
// 			tiers := []int{}

// 			if pbk.Tier1EventCombinations.AnyAll(&pokemon, team) {
// 				tiers = append(tiers, 1)
// 			}

// 			if pbk.Tier1EventCombinations.AnyAll(&pokemon, team) {
// 				tiers = append(tiers, 2)
// 			}

// 			if pbk.Tier1EventCombinations.AnyAll(&pokemon, team) {
// 				tiers = append(tiers, 3)
// 			}

// 			if pbk.Tier1EventCombinations.AnyAll(&pokemon, team) {
// 				tiers = append(tiers, 4)
// 			}

// 			if pbk.Tier1EventCombinations.AnyAll(&pokemon, team) {
// 				tiers = append(tiers, 5)
// 			}

// 			if len(tiers) != 0 {
// 				fmt.Println(moveNameWithTier)
// 			}
// 		}
// 	}
// 	return pokemon
// }
