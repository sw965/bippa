package bippa

import (
	omwmath "github.com/sw965/omw/math"
)

type CriticalBonus float64

var (
	CRITICAL_BONUS    = CriticalBonus(6144.0 / 4096.0)
	NO_CRITICAL_BONUS = CriticalBonus(4096.0 / 4096.0)
)

var CRITICAL_N = map[CriticalRank]int{0: 24, 1: 8, 2: 2, 3: 1}

func NewCriticalBonus(isCrit bool) CriticalBonus {
	if isCrit {
		return CRITICAL_BONUS
	} else {
		return NO_CRITICAL_BONUS
	}
}

type SameTypeAttackBonus float64

const (
	SAME_TYPE_ATTACK_BONUS    = SameTypeAttackBonus(6144.0 / 4096.0)
	NO_SAME_TYPE_ATTACK_BONUS = SameTypeAttackBonus(4096.0 / 4096.0)
)

func NewSameTypeAttackBonus(isSTAB bool) SameTypeAttackBonus {
	if isSTAB {
		return SAME_TYPE_ATTACK_BONUS
	} else {
		return NO_SAME_TYPE_ATTACK_BONUS
	}
}

type EffectivenessBonus float64

// https://latest.pokewiki.net/%E3%83%80%E3%83%A1%E3%83%BC%E3%82%B8%E8%A8%88%E7%AE%97%E5%BC%8F
type RandomDamageBonus float64
type RandomDamageBonuses []RandomDamageBonus

var RANDOM_DAMAGE_BONUSES = RandomDamageBonuses{
	0.85, 0.86, 0.87, 0.88, 0.89, 0.9, 0.91, 0.92, 0.93, 0.94, 0.95, 0.96, 0.97, 0.98, 0.99, 1.0,
}
var MAX_RANDOM_DAMAGE_BONUS = omwmath.Max(RANDOM_DAMAGE_BONUSES...)
var MEAN_RANDOM_DAMAGE_BONUS = omwmath.Mean(RANDOM_DAMAGE_BONUSES...)

type DamageBonus int

const (
	INIT_DAMAGE_BONUS = DamageBonus(4096)
)

func NewDamageBonus(poke *Pokemon) DamageBonus {
	y := INIT_DAMAGE_BONUS
	if poke.Item == LIFE_ORB {
		v := FiveOrMoreRounding(float64(y) * 5324.0 / 4096.0)
		y = DamageBonus(v)
	}
	return y
}

type WeatherBonus float64

const (
	NO_WEATHER_BONUS = 4096.0 / 4096.0
	GOOD_WEATHER_BONUS = 6144.0 / 4096.0
	BAD_WEATHER_BONUS = 2048.0 / 4096.0
)

type FinalDamage int

func NewFinalDamage(bt *Battle, moveName MoveName, isCrit bool, randBonus RandomDamageBonus) FinalDamage {
	attacker := bt.P1Fighters[0]
	defender := bt.P2Fighters[0]

	power := NewFinalPower(bt, moveName)
	atk := NewFinalAttack(&attacker, moveName, isCrit)
	def := NewFinalDefense(&defender, moveName, isCrit)

	moveData := MOVEDEX[moveName]
	var weatherBonus WeatherBonus
	if moveData.Type == FIRE && bt.Weather == SUNNY_DAY {
		weatherBonus = GOOD_WEATHER_BONUS
	} else if moveData.Type == WATER && bt.Weather == SUNNY_DAY {
		weatherBonus = BAD_WEATHER_BONUS
	} else {
		weatherBonus = NO_WEATHER_BONUS
	}

	critBonus := NewCriticalBonus(isCrit)
	stab := attacker.SameTypeAttackBonus(moveName)
	effeBonus := defender.EffectivenessBonus(moveName)
	dmgBonus := NewDamageBonus(&attacker)

	y := int(DEFAULT_LEVEL)*2/5 + 2
	y = int(float64(y) * float64(power) * float64(atk) / float64(def))
	y = y/50 + 2
	y = FiveOverRounding(float64(y) * float64(weatherBonus))
	y = FiveOverRounding(float64(y) * float64(critBonus))
	y = int(float64(y) * float64(randBonus))
	y = FiveOverRounding(float64(y) * float64(stab))
	y = int(float64(y) * float64(effeBonus))
	y = FiveOverRounding(float64(y) * float64(dmgBonus) / 4096.0)
	return FinalDamage(y)
}