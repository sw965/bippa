package bippa

import (
	"golang.org/x/exp/slices"
)

type Power int

func NewWeightDependentPower(w Weight) Power {
	if w < 10 {
		return 20
	} else if w < 25 {
		return 40
	} else if w < 50 {
		return 60
	} else if w < 100 {
		return 80
	} else if w < 200 {
		return 100
	} else {
		return 120
	}
}

func NewPower(bt *Battle, moveName MoveName) Power {
	moveData := MOVEDEX[moveName]
	if moveName == A_KU_RO_BA_TTO {
		return moveData.Power * 2
	} else if slices.Contains(WEIGHT_DEPENDENT_ATTACK_MOVE_NAMES, moveName) {
		name := bt.P2Fighters[0].Name
		return NewWeightDependentPower(POKEDEX[name].Weight)
	} else {
		return moveData.Power
	}
}

// https://latest.pokewiki.net/%E3%83%80%E3%83%A1%E3%83%BC%E3%82%B8%E8%A8%88%E7%AE%97%E5%BC%8F
type PowerBonus int

const (
	INIT_POWER_BONUS = PowerBonus(4096)
)

func NewPowerBonus(attacker *Pokemon, moveName MoveName) PowerBonus {
	y := INIT_POWER_BONUS
	sas := StatusAilments{NORMAL_POISON, BAD_POISON, PARALYSIS, BURN}
	if moveName == KA_RA_GE_N_KI && slices.Contains(sas, attacker.StatusAilment) {
		bonus := 8192.0/4096.0
		y = PowerBonus(FiveOrMoreRounding(float64(y) * bonus))
	}
	return y
}

type FinalPower int

func NewFinalPower(bt *Battle, moveName MoveName) FinalPower {
	attacker := bt.P1Fighters[0]
	power := NewPower(bt, moveName)
	bonus := NewPowerBonus(&attacker, moveName)
	y := FiveOverRounding(float64(power) * float64(bonus) / 4096.0)
	return FinalPower(y)
}