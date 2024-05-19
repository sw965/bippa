package dmgtools

import (
	"golang.org/x/exp/slices"
	bp "github.com/sw965/bippa"
	omath "github.com/sw965/omw/math"
)

// 小数点以下が0.5より大きいならば、切り上げ
func RoundOverHalf(x float64) int {
	decimalPart := float64(x) - float64(int(x))
	if decimalPart > 0.5 {
		return int(x + 1)
	}
	return int(x)
}

type RandBonus float64
type RandBonuses []RandBonus

var RAND_BONUSES = RandBonuses{
	0.85, 0.86, 0.87, 0.88, 0.89, 0.90,
	0.91, 0.92, 0.93, 0.94, 0.95,
	0.96, 0.97, 0.98, 0.99, 1.0,
}

type Attacker struct {
	PokeName bp.PokeName
	Level bp.Level
	Atk int
	SpAtk int
	MoveName bp.MoveName
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
func (c *Calculator) Calculation(randBonus RandBonus) int {
	attacker := c.Attacker
	defender := c.Defender
	attackerPokeData := bp.POKEDEX[attacker.PokeName]
	defenderPokeData := bp.POKEDEX[defender.PokeName]
	moveData := bp.MOVEDEX[attacker.MoveName]

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

	var typeMatchBonus float64 
	if slices.Contains(attackerPokeData.Types, moveData.Type) {
		typeMatchBonus = 6144.0/4096.0
	} else {
		typeMatchBonus = 1.0
	}

	dmg = RoundOverHalf(float64(dmg) * typeMatchBonus)
	for _, defType := range defenderPokeData.Types {
		dmg = int(float64(dmg) * bp.TYPEDEX[moveData.Type][defType])
	}
	dmg = int(float64(dmg) * float64(randBonus))
	return omath.Max(dmg, 1)
}