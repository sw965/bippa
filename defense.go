package bippa

import (
	omwmath "github.com/sw965/omw/math"
)

type DefenseBonus int

const (
	INIT_DEFENSE_BONUS = DefenseBonus(4096)
)

func NewDefenseBonus(poke *Pokemon) DefenseBonus {
	y := INIT_DEFENSE_BONUS
	if poke.Item == ASSAULT_VEST {
		v := FiveOrMoreRounding(float64(y) * (6144.0 / 4096.0))
		y = DefenseBonus(v)
	}
	return y
}

type FinalDefense int

func NewFinalDefense(poke *Pokemon, moveName MoveName, isCrit bool) FinalDefense {
	moveData := MOVEDEX[moveName]

	var def State
	var rank Rank

	switch moveData.Category {
	case PHYSICS:
		def = poke.Def
		rank = poke.RankState.Def
	case SPECIAL:
		def = poke.SpDef
		rank = poke.RankState.SpDef
	}

	if rank > 0 && isCrit {
		rank = 0
	}

	bonus := rank.ToBonus()
	y := int(float64(def) * float64(bonus))
	return omwmath.Max(FinalDefense(y), 1)
}