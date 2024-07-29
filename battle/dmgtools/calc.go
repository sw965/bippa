package dmgtools

import (
	bp "github.com/sw965/bippa"
	omwmath "github.com/sw965/omw/math"
)

// 小数点以下が0.5以下なら切り捨て、0.5よりも大きいなら切り上げ。
func RoundOverHalf(x float64) int {
	decimalPart := x - float64(int(x))
	if decimalPart > 0.5 {
		return int(x + 1)
	}
	return int(x)
}

type Attacker struct {
	PokeName bp.PokeName
	Level bp.Level
	Types bp.Types
	Atk int
	AtkRank bp.Rank
	SpAtk int
	SpAtkRank bp.Rank
	Ability bp.Ability
}

func NewAttacker(p *bp.Pokemon) Attacker {
	return Attacker{
		PokeName:p.Name,
		Level:p.Level,
		Types:p.Types,
		Atk:p.Atk,
		AtkRank:p.Rank.Atk,
		SpAtk:p.SpAtk,
		SpAtkRank:p.Rank.SpAtk,
		Ability:p.Ability,
	}
}

type Defender struct {
	PokeName bp.PokeName
	Level bp.Level
	Types bp.Types
	Def int
	DefRank bp.Rank
	SpDef int
	SpDefRank bp.Rank
	Ability bp.Ability
}

func NewDefender(p *bp.Pokemon) Defender {
	return Defender{
		PokeName:p.Name,
		Level:p.Level,
		Types:p.Types,
		Def:p.Def,
		DefRank:p.Rank.Def,
		SpDef:p.SpDef,
		SpDefRank:p.Rank.SpDef,
		Ability:p.Ability,
	}
}

type Calculator struct {
	Attacker Attacker
	Defender Defender
	IsSingleDamage bool
	IsCritical bool
	RandBonus float64
}

// https://latest.pokewiki.net/%E3%83%80%E3%83%A1%E3%83%BC%E3%82%B8%E8%A8%88%E7%AE%97%E5%BC%8F
func (c *Calculator) Calculation(moveName bp.MoveName) (int, bool) {
	attacker := c.Attacker
	defender := c.Defender
	moveData := bp.MOVEDEX[moveName]
	var isNoEffect bool

	fs := []func(int) int {
		//ダブルバトルの複数攻撃
		func(dmg int) int {
			bonus := RangeAttackBonus(c.IsSingleDamage)
			return RoundOverHalf(float64(dmg) * bonus / NO_RANGE_ATTACK_BONUS)
		},

		//急所
		func(dmg int) int {
			bonus := CriticalBonus(c.IsCritical)
			return RoundOverHalf(float64(dmg) * bonus / NO_CRITICAL_BONUS)
		},

		//乱数
		func(dmg int) int {
			return int(float64(dmg) * float64(c.RandBonus))
		},

		//タイプ一致
		func(dmg int) int {
			bonus := SameTypeAttackBonus(moveName, attacker.Types, moveData.Type)
			return RoundOverHalf(float64(dmg) * bonus / NO_SAME_TYPE_ATTACK_BONUS)	
		},

		//タイプ相性
		func(dmg int) int {
			bonus := EffectivenessBonus(moveName, moveData.Type, attacker.Types)
			isNoEffect = bonus == 0.0
			return int(float64(dmg) * bonus)
		},
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

	for _, f := range fs {
		dmg = f(dmg)
		if isNoEffect {
			return 0, isNoEffect
		}
	}

	//ダメージが0ならば、1にする
	return omwmath.Max(dmg, 1), isNoEffect
}