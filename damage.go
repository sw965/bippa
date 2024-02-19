package bippa

import (
	omwmath "github.com/sw965/omw/math"
)
type EffectivenessBonus float64

// https://latest.pokewiki.net/%E3%83%80%E3%83%A1%E3%83%BC%E3%82%B8%E8%A8%88%E7%AE%97%E5%BC%8F
type RandomDamageBonus float64
type RandomDamageBonuses []RandomDamageBonus

var RANDOM_DAMAGE_BONUSES = RandomDamageBonuses{
	0.85, 0.86, 0.87, 0.88, 0.89, 0.9, 0.91, 0.92, 0.93, 0.94, 0.95, 0.96, 0.97, 0.98, 0.99, 1.0,
}
var MAX_RANDOM_DAMAGE_BONUS = omwmath.Max(RANDOM_DAMAGE_BONUSES...)
var MEAN_RANDOM_DAMAGE_BONUS = omwmath.Mean(RANDOM_DAMAGE_BONUSES...)

type WeatherBonus float64

const (
	NO_WEATHER_BONUS = 4096.0 / 4096.0
	GOOD_WEATHER_BONUS = 6144.0 / 4096.0
	BAD_WEATHER_BONUS = 2048.0 / 4096.0
)

// https://latest.pokewiki.net/%E3%83%80%E3%83%A1%E3%83%BC%E3%82%B8%E8%A8%88%E7%AE%97%E5%BC%8F
func NewFinalPower(bt *Battle, moveName MoveName) float64 {
	attacker := bt.P1Fighters[0]
	power := NewPower(bt, moveName)
	bonus := 4096
	y := FiveOverRounding(float64(power) * float64(bonus) / 4096.0)
	return FinalPower(y)
}

type AtkValForCalc int

func NewAtkValForCalc(p *Pokemon, c Category) AtkValForCalc {
	if c == PHYSICS {
		return AtkValForCalc(p.Atk)
	} else {
		return AtkValForCalc(p.SpAtk)
	}
}

type DefValForCalc int

func NewDefValForCalc(p *Pokemon, c Category) DefValForCalc {
	if c =- PHYSICS {
		return DefValForCalc(p.Def)
	} else {
		return DefValForCalc(p.SpDef)
	}
}

type AtkRankValForCalc int

func NewAtkRankValForCalc(p *Pokemon, c Category) AtkRankValForCalc {
	if c == PHYSICS {
		return AtkRankValForCalc(p.Rank.Atk)
	} else {
		return AtkRankValForCalc(p.Rank.SpAtk)
	}
}

func NewDefRankValForCalc(p *Pokemon, c Category) DefValRankForCalc {
	if c == PHYSICS {
		return DefValRankForCalc(p.Rank.Def)
	} else {
		return DefValRankForCalc(p.Rank.SpDef)
	}
}

type AtkBonus int

const (
	INIT_ATK_BONUS = 4096
)

//物理技のボーナス
func NewChoiceBandBonus(p *Pokemon) AtkBonus {
	if p.Item == CHOICE_BAND {
		return AtkBonus(6144.0 / 4096.0)
	} else {
		return bonus
	}
}

//特殊技のボーナス
func NewChoiceSpecsBonus(p *Pokemon) AtkBonus {
	if p.Item == CHOICE_SPECS {
		return AtkBonus(6144.0 / 4096.0)
	} else {
		return bonus
	}
}

type AtkBonuses []AtkBonus

func NewPhysicsAtkBonuses(p *Pokemon) AtkBonuses {
	return AtkBonuses{
		NewChoiceBandBonus(p),
	}
}

func NewSpecialAtkBonuses(p *Pokemon) AtkBonuses {
	return AtkBonuses{
		NewChoiceSpecsBonus(p),
	}
}

type FinalAttack int

func NewFinalAttack(attacker *Pokemon, moveName MoveName, isCrit bool) FinalAttack {
	moveData := MOVEDEX[moveName]
	atk := NewAtkValForCalc(attacker, moveData.Category)
	rank := map[bool]AtkRankValForCalc{
		true:0,
		false:NewAtkRankValForCalc(attacker, moveData.Category),
	}[isCrit]
	atkBonus := fn.Reduce[AtkBonuses](
		NewAtkBonuses(moveData.Category, attacker),
		fn.Mul,
		INIT_ATK_BONUS,
	)
	rankBonus := rank.ToBonus()

	y := int(float64(atk) * rankBonus)
	y = FiveOverRounding(y * float64(atkBonus) / 4096.0)
	return FinalAttack(omwmath.Max(y, 1))
}

func NewFinalDefense(defender *Pokemon, moveName MoveName, isCrit bool) float64 {
	moveData := MOVEDEX[moveName]

	var def int
	var rank Rank

	switch moveData.Category {
		case PHYSICS:
			def = defender.Def
			rank = defender.Rank.Def
		case SPECIAL:
			def = defender.SpDef
			rank = defender.Rank.SpDef
	}

	if rank > 0 && isCrit {
		rank = 0
	}

	defBonus := 4096
	switch moveData.Category {
		case PHYSICS:
			defBonus = 4096
		case SPECIAL:
			if defender.Item == ASSAULT_VEST {
				defBonus = FiveOrMoreRounding(defBonus * (6144.0 / 4096.0))
			}

	rankBonus := rank.ToBonus()
	y := int(float64(def) * rankBonus)
	y = FiveOverRounding(y * float64(defBonus) / 4096.0)
	return omwmath.Max(FinalDefense(y), 1)
}

var CRITICAL_N = map[CriticalRank]int{0: 24, 1: 8, 2: 2, 3: 1}

type FinalDamage int

func NewFinalDamage(bt *Battle, moveName MoveName, isCrit bool, randBonus float64) FinalDamage {
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

	dmgBonus := 4096
	if attacker.Item == LIFE_ORB {
		dmgBonus = FiveOrMoreRounding(float64(y) * 5324.0 / 4096.0)
	}

	var critBonus float64
	if isCrit {
		critBonus = 6144.0 / 4096.0
	} else {
		critBonus = 4096.0 / 4096.0
	}

	var stab float64
	if slices.Contains(attacker.Types, moveData.Type) {
		stab = 6144.0 / 4096.0
	} else {
		stab = 4096.0 / 4096.0
	}

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