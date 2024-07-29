package dmgtools

import (
	bp "github.com/sw965/bippa"	
	"golang.org/x/exp/slices"
)

const (
	RANGE_ATTACK_BONUS = 3072.0
	NO_RANGE_ATTACK_BONUS = 4096.0
)

func RangeAttackBonus(isSingleDmg bool) float64 {
	if isSingleDmg {
		return NO_RANGE_ATTACK_BONUS
	} else {
		return RANGE_ATTACK_BONUS
	}
}

const (
	CRITICAL_BONUS = 6144.0
	NO_CRITICAL_BONUS = 4096.0
)

func CriticalBonus(isCrit bool) float64 {
	if isCrit {
		return CRITICAL_BONUS
	} else {
		return NO_CRITICAL_BONUS
	}
}

var RAND_BONUSES = []float64{
	0.85, 0.86, 0.87, 0.88, 0.89, 0.90,
	0.91, 0.92, 0.93, 0.94, 0.95,
	0.96, 0.97, 0.98, 0.99, 1.0,
}

const (
	SAME_TYPE_ATTACK_BONUS = 6144.0
	NO_SAME_TYPE_ATTACK_BONUS = 4096.0
)

func SameTypeAttackBonus(moveName bp.MoveName, pokeTypes bp.Types, moveType bp.Type) float64 {
	if moveName == bp.STRUGGLE {
		return NO_SAME_TYPE_ATTACK_BONUS
	} else if slices.Contains(pokeTypes, moveType) {
		return SAME_TYPE_ATTACK_BONUS
	} else {
		return NO_SAME_TYPE_ATTACK_BONUS
	}
}

func EffectivenessBonus(moveName bp.MoveName, atkMoveType bp.Type, defTypes bp.Types) float64 {
	if moveName == bp.STRUGGLE {
		return 1.0
	} else {
		return bp.TYPEDEX.EffectivenessValue(atkMoveType, defTypes)
	}
}