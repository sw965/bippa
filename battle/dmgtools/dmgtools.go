package dmgtools

import (
	"golang.org/x/exp/slices"
	bp "github.com/sw965/bippa"
	omwmath "github.com/sw965/omw/math"
	"github.com/sw965/omw/fn"
)

// 小数点以下が0.5より大きいならば、切り上げ
func RoundOverHalf(x float64) int {
	decimalPart := float64(x) - float64(int(x))
	if decimalPart > 0.5 {
		return int(x + 1)
	}
	return int(x)
}

func Effectiveness(atkType bp.Type, defTypes bp.Types) float64 {
	ret := 1.0
	for _, defType := range defTypes {
		ret *= bp.TYPEDEX[atkType][defType]
	}
	return ret
}

type RandBonus float64
type RandBonuses []RandBonus

var RAND_BONUSES = RandBonuses{
	0.85, 0.86, 0.87, 0.88, 0.89, 0.90,
	0.91, 0.92, 0.93, 0.94, 0.95,
	0.96, 0.97, 0.98, 0.99, 1.0,
}

var MEAN_RAND_BONUS = omwmath.Mean(RAND_BONUSES...)

type Attacker struct {
	PokeName bp.PokeName
	Level bp.Level
	Atk int
	SpAtk int
}

type Defender struct {
	PokeName bp.PokeName
	Level bp.Level
	Def int
	SpDef int
}

type Calculator struct {
	Attacker Attacker
	Defender Defender
}

// https://latest.pokewiki.net/%E3%83%80%E3%83%A1%E3%83%BC%E3%82%B8%E8%A8%88%E7%AE%97%E5%BC%8F
func (c *Calculator) Calculation(moveName bp.MoveName, randBonus RandBonus) int {
	attacker := c.Attacker
	defender := c.Defender
	attackerPokeData := bp.POKEDEX[attacker.PokeName]
	defenderPokeData := bp.POKEDEX[defender.PokeName]
	moveData := bp.MOVEDEX[moveName]

	var atkVal int
	var defVal int
	if moveData.Category == bp.PHYSICS {
		atkVal = attacker.Atk
		defVal = defender.Def
	} else {
		atkVal = attacker.SpAtk
		defVal = defender.SpDef
	}
	power := moveData.Power

	lv := int(attacker.Level)
	dmg := (lv*2/5) + 2
	dmg = dmg * power * atkVal / defVal
	dmg = (dmg/50) + 2

	var sameTypeAttackBonus float64
	if moveName == bp.STRUGGLE {
		sameTypeAttackBonus = 1.0
	} else if slices.Contains(attackerPokeData.Types, moveData.Type) {
		sameTypeAttackBonus = 6144.0/4096.0
	} else {
		sameTypeAttackBonus = 1.0
	}

	var effect float64
	if moveName == bp.STRUGGLE {
		effect = 1.0
	} else {
		effect = Effectiveness(moveData.Type, defenderPokeData.Types)
		if effect == 0.0 {
			return 0
		}
	}

	dmg = RoundOverHalf(float64(dmg) * sameTypeAttackBonus)
	dmg = int(float64(dmg) * effect)
	dmg = int(float64(dmg) * float64(randBonus))
	return omwmath.Max(dmg, 1)
}

func (c *Calculator) Calculations(moveName bp.MoveName) []int {
	dmgs := make([]int, len(RAND_BONUSES))
	for i, r := range RAND_BONUSES {
		dmgs[i] = c.Calculation(moveName, r)
	}
	return dmgs
}

func (c *Calculator) Expected(moveName bp.MoveName) float64 {
	dmgs := fn.Map[[]float64](c.Calculations(moveName), fn.IntToFloat64[int, float64])
	return omwmath.Mean(dmgs...)
}