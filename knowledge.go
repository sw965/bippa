package bippa

import (
	"fmt"
)

type Tier int

const (
	TIER1 = 1 
	TIER2 = 2
	TIER3 = 3
	TIER4 = 4
	TIER5 = 5
	TIER6 = 6
)

func (tier Tier) LinearFunctionSelectPercent(min, max float64) float64 {
	//一次関数の連立方程式を解く
	a := (min - max) / (float64(TIER6) - float64(TIER1))
	b := min - (float64(TIER6) * a)
	return (a * float64(tier)) + b
}

type PokemonBuildEvent func(*Pokemon, Team) bool

func NewPokemonBuildSelfMoveNameEvent(moveName MoveName) PokemonBuildEvent {
	result := func(pokemon *Pokemon, team Team) bool {
		_, ok := pokemon.Moveset[moveName]
		return ok
	}
	return result
}

type PokemonBuildEventCombination []PokemonBuildEvent

func NewPokemonBuildSelfMoveNameEventCombination(moveName ...MoveName) PokemonBuildEventCombination {
	result := make(PokemonBuildEventCombination, len(moveName))
	for i, iMoveName := range moveName {
		result[i] = NewPokemonBuildSelfMoveNameEvent(iMoveName)
	}
	return result
}

func (pbec PokemonBuildEventCombination) IsAll(pokemon *Pokemon, team Team) bool {
	for _, pbe := range pbec {
		if pbe(pokemon, team) {
			return true
		}
	}
	return false
}

type PokemonBuildEventCombinations []PokemonBuildEventCombination

func (pbecs PokemonBuildEventCombinations) InAll(pokemon *Pokemon, team Team) bool {
	for _, pbec := range pbecs {
		if pbec.IsAll(pokemon, team) {
			return true
		}
	}
	return false
}

type PokemonBuildKnowledge struct {
	Natures Natures
	MoveNames MoveNames
	Tier1EventCombinations PokemonBuildEventCombinations
	Tier2EventCombinations PokemonBuildEventCombinations
	Tier3EventCombinations PokemonBuildEventCombinations
	Tier4EventCombinations PokemonBuildEventCombinations
	Tier5EventCombinations PokemonBuildEventCombinations
	Tier6EventCombinations PokemonBuildEventCombinations
}

func NewVenusaurCombinationKnowledge() PokemonBuildKnowledge {
	moveNames := MoveNames{"ギガドレイン", "ヘドロばくだん", "やどりぎのタネ", "まもる", "どくどく", "だいちのちから"}
	natures := Natures{"しんちょう", "ずぶとい", "ひかえめ"}

	tier1 := PokemonBuildEventCombinations{
		NewPokemonBuildSelfMoveNameEventCombination("ギガドレイン"),
		NewPokemonBuildSelfMoveNameEventCombination("ギガドレイン", "ヘドロばくだん"),
		NewPokemonBuildSelfMoveNameEventCombination("やどりぎのタネ", "まもる"),
		NewPokemonBuildSelfMoveNameEventCombination("どくどく", "まもる"),
	}

	tier2 := PokemonBuildEventCombinations{
		NewPokemonBuildSelfMoveNameEventCombination("ギガドレイン", "だいちのちから"),
		NewPokemonBuildSelfMoveNameEventCombination("だいちのちから", "ヘドロばくだん"),
		NewPokemonBuildSelfMoveNameEventCombination("ギガドレイン", "やどりぎのタネ"),
		NewPokemonBuildSelfMoveNameEventCombination("ギガドレイン", "どくどく"),
	}

	tier4 := PokemonBuildEventCombinations{
		NewPokemonBuildSelfMoveNameEventCombination("ヘドロばくだん", "まもる"),
	}

	return PokemonBuildKnowledge{Natures:natures, MoveNames:moveNames,
		Tier1EventCombinations:tier1, Tier2EventCombinations:tier2, Tier4EventCombinations:tier4}
}

func (pbk PokemonBuildKnowledge) BuildPokemon(team Team) Pokemon {
	pokemon := Pokemon{}
	for i := 0; i < MAX_MOVESET_LENGTH; i++ {
		moveNameWithTier := MoveNameWithTier{}

		for _, moveName := range pbk.MoveNames {
			if _, ok := pokemon.Moveset[moveName]; ok {
				continue
			}

			powerPoint := NewPowerPoint(MOVEDEX[moveName].BasePP, MAX_POINT_UP)
			pokemon.Moveset[moveName] = &powerPoint
			tiers := []int{}

			if pbk.Tier1EventCombinations.InAll(&pokemon, team) {
				tiers = append(tiers, 1)
			}

			if pbk.Tier1EventCombinations.InAll(&pokemon, team) {
				tiers = append(tiers, 2)
			}

			if pbk.Tier1EventCombinations.InAll(&pokemon, team) {
				tiers = append(tiers, 3)
			}

			if pbk.Tier1EventCombinations.InAll(&pokemon, team) {
				tiers = append(tiers, 4)
			}

			if pbk.Tier1EventCombinations.InAll(&pokemon, team) {
				tiers = append(tiers, 5)
			}

			if len(tiers) != 0 {
				fmt.Println(moveNameWithTier)
			}
		}
	}
	return pokemon
}