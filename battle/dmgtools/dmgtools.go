package dmgtools

import (
	"fmt"
	"golang.org/x/exp/slices"
	bp "github.com/sw965/bippa"
	omwmath "github.com/sw965/omw/math"
	"math/rand"
)

// 小数点以下が0.5以下なら切り捨て、0.5よりも大きいなら切り上げ。
func RoundOverHalf(x float64) int {
	decimalPart := float64(x) - float64(int(x))
	if decimalPart > 0.5 {
		return int(x + 1)
	}
	return int(x)
}

func CriticalN(rank bp.CriticalRank) (int, error) {
	if rank < 0 {
		return 0, fmt.Errorf("急所ランクは0以上でなければならない")
	}

	var n int
	switch rank {
		case 0:
			n = 16
		case 1:
			n = 8
		case 2:
			n = 4
		case 3:
			n = 3
		default:
			n = 2
	}
	return n, nil
}

func IsCritical(rank bp.CriticalRank, r *rand.Rand) (bool, error) {
	n, err := CriticalN(rank)
	if err != nil {
		return false, err
	}
	return r.Intn(n) == 0, nil
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
	AtkRank bp.Rank
	SpAtk int
	SpAtkRank bp.Rank
	Ability bp.Ability
}

type Defender struct {
	PokeName bp.PokeName
	Level bp.Level
	Def int
	DefRank bp.Rank
	SpDef int
	SpDefRank bp.Rank
	Ability bp.Ability
}

type Calculator struct {
	Attacker Attacker
	Defender Defender
	IsCritical bool
	IsSingleDamage bool
}

// https://latest.pokewiki.net/%E3%83%80%E3%83%A1%E3%83%BC%E3%82%B8%E8%A8%88%E7%AE%97%E5%BC%8F
func (c *Calculator) Calculation(moveName bp.MoveName, randBonus RandBonus) int {
	attacker := c.Attacker
	defender := c.Defender
	attackerPokeData := bp.POKEDEX[attacker.PokeName]
	defenderPokeData := bp.POKEDEX[defender.PokeName]

	var atkVal int
	var defVal int

	moveData := bp.MOVEDEX[moveName]
	if moveData.Category == bp.PHYSICS {
		atkRank := c.Attacker.AtkRank
		if c.IsCritical {
			atkRank = omwmath.Max(atkRank, 0)
		}
		defRank := c.Defender.DefRank
		if c.IsCritical {
			defRank = omwmath.Min(defRank, 0)
		}
		atkVal = int(float64(attacker.Atk) * atkRank.Bonus())
		defVal = int(float64(defender.Def) * defRank.Bonus())
	} else {
		spAtkRank := c.Attacker.SpAtkRank
		if c.IsCritical {
			spAtkRank = omwmath.Max(spAtkRank, 0)
		}
		spDefRank := c.Defender.SpDefRank
		if c.IsCritical {
			spDefRank = omwmath.Min(spDefRank, 0)
		}
		atkVal = int(float64(attacker.SpAtk) * spAtkRank.Bonus())
		defVal = int(float64(defender.SpDef) * spDefRank.Bonus())
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
		sameTypeAttackBonus = 1.5
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

	if !c.IsSingleDamage {
		dmg = RoundOverHalf(float64(dmg) * 0.75)
	}

	dmg = int(float64(dmg) * float64(randBonus))

	if c.IsCritical {
		dmg = RoundOverHalf(float64(dmg) * 2.0)
	}

	dmg = RoundOverHalf(float64(dmg) * sameTypeAttackBonus)
	dmg = int(float64(dmg) * effect)
	return omwmath.Max(dmg, 1)
}

func (c *Calculator) Calculations(moveName bp.MoveName) []int {
	dmgs := make([]int, len(RAND_BONUSES))
	for i, r := range RAND_BONUSES {
		dmgs[i] = c.Calculation(moveName, r)
	}
	return dmgs
}