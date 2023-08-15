package bippa

import (
	omwmath "github.com/sw965/omw/math"
)

type PhysicsAttackBonus int

const (
	INIT_PHYSICS_ATTACK_BONUS = PhysicsAttackBonus(4096)
)

func NewPhysicsAttackBonus(poke *Pokemon) PhysicsAttackBonus {
	y := int(INIT_PHYSICS_ATTACK_BONUS)
	if poke.Item == CHOICE_BAND {
		y = FiveOrMoreRounding(float64(y) * 6144.0 / 4096.0)
	}
	return PhysicsAttackBonus(y)
}

type SpecialAttackBonus int

const (
	INIT_SPECIAL_ATTACK_BONUS = SpecialAttackBonus(4096)
)

func NewSpecialAttackBonus(poke *Pokemon) SpecialAttackBonus {
	y := int(INIT_SPECIAL_ATTACK_BONUS)
	if poke.Item == CHOICE_SPECS {
		y = FiveOrMoreRounding(float64(y) * 6144.0 / 4096.0)
	}
	return SpecialAttackBonus(y)
}

type AttackBonus int

func NewAttackBonus(poke *Pokemon, moveName MoveName) AttackBonus {
	moveData := MOVEDEX[moveName]
	switch moveData.Category {
	case PHYSICS:
		bonus := NewPhysicsAttackBonus(poke)
		return AttackBonus(bonus)
	case SPECIAL:
		bonus := NewSpecialAttackBonus(poke)
		return AttackBonus(bonus)
	default:
		return -1
	}
}

type FinalAttack int

func NewFinalAttack(poke *Pokemon, moveName MoveName, isCrit bool) FinalAttack {
	moveData := MOVEDEX[moveName]

	var atk State
	var rank Rank

	switch moveData.Category {
	case PHYSICS:
		atk = poke.Atk
		rank = poke.RankState.Atk
	case SPECIAL:
		atk = poke.SpAtk
		rank = poke.RankState.SpAtk
	}

	atkBonus := NewAttackBonus(poke, moveName)

	if rank < 0 && isCrit {
		rank = 0
	}

	rankBonus := rank.ToBonus()

	y := int(float64(atk) * float64(rankBonus))
	y = FiveOverRounding(float64(y) * float64(atkBonus) / 4096.0)
	return omwmath.Max(FinalAttack(y), 1)
}