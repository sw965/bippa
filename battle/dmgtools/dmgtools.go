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
	decimalPart := x - float64(int(x))
	if decimalPart > 0.5 {
		return int(x + 1)
	}
	return int(x)
}

const (
	RANGE_ATTACK_BONUS = 3072.0
	NO_RANGE_ATTACK_BONUS = 4096.0
)

const (
	CRITICAL_BONUS = 6144.0
	NO_CRITICAL_BONUS = 4096.0
)

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

const (
	SAME_TYPE_ATTACK_BONUS = 6144.0
	NO_SAME_TYPE_ATTACK_BONUS = 4096.0
)

type RandBonus float64
type RandBonuses []RandBonus

var RAND_BONUSES = RandBonuses{
	0.85, 0.86, 0.87, 0.88, 0.89, 0.90,
	0.91, 0.92, 0.93, 0.94, 0.95,
	0.96, 0.97, 0.98, 0.99, 1.0,
}

func EffectivenessBonus(atkType bp.Type, defTypes bp.Types) float64 {
	ret := 1.0
	for _, defType := range defTypes {
		ret *= bp.TYPEDEX[atkType][defType]
	}
	return ret
}

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
	moveData := bp.MOVEDEX[moveName]

	//無効タイプならば、ダメージは0
	for _, defType := range defenderPokeData.Types {
		if bp.TYPEDEX[moveData.Type][defType] == 0.0 {
			return 0
		}
	}

	var atkVal int
	var defVal int
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

	//基礎ダメージ
	power := moveData.Power
	lv := int(attacker.Level)
	dmg := (lv*2/5) + 2
	dmg = dmg * power * atkVal / defVal
	dmg = (dmg/50) + 2

	//ダブルバトルの範囲攻撃
	var rangeAttackBonus float64
	if c.IsSingleDamage {
		rangeAttackBonus = NO_RANGE_ATTACK_BONUS
	} else {
		rangeAttackBonus = RANGE_ATTACK_BONUS
	}

	dmg = RoundOverHalf(float64(dmg) * rangeAttackBonus / NO_RANGE_ATTACK_BONUS)

	//急所
	var critBonus float64
	if c.IsCritical {
		critBonus = CRITICAL_BONUS
	} else {
		critBonus = NO_CRITICAL_BONUS
	}
	dmg = RoundOverHalf(float64(dmg) * critBonus / NO_CRITICAL_BONUS)

	dmg = int(float64(dmg) * float64(randBonus))

	//タイプ一致
	var sameTypeAttackBonus float64
	if moveName == bp.STRUGGLE {
		sameTypeAttackBonus = NO_SAME_TYPE_ATTACK_BONUS
	} else if slices.Contains(attackerPokeData.Types, moveData.Type) {
		sameTypeAttackBonus = SAME_TYPE_ATTACK_BONUS
	} else {
		sameTypeAttackBonus = NO_SAME_TYPE_ATTACK_BONUS
	}
	dmg = RoundOverHalf(float64(dmg) * sameTypeAttackBonus / NO_SAME_TYPE_ATTACK_BONUS)

	//タイプ相性
	var effect float64
	if moveName == bp.STRUGGLE {
		effect = 1.0
	} else {
		effect = EffectivenessBonus(moveData.Type, defenderPokeData.Types)
	}
	dmg = int(float64(dmg) * effect)

	//ダメージが0ならば、1にする
	return omwmath.Max(dmg, 1)
}