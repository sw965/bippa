package battle

import (
	bp "github.com/sw965/bippa"
	omwmath "github.com/sw965/omw/math"
	"golang.org/x/exp/slices"
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

var DAMAGE_RAND_BONUSES = []float64{
	0.85, 0.86, 0.87, 0.88, 0.89, 0.90,
	0.91, 0.92, 0.93, 0.94, 0.95,
	0.96, 0.97, 0.98, 0.99, 1.0,
}

const (
	SAME_TYPE_ATTACK_BONUS = 6144.0
	NO_SAME_TYPE_ATTACK_BONUS = 4096.0
)

type AttackerInfo struct {
	PokeName bp.PokeName
	Level bp.Level
	Types bp.Types
	CurrentHP int
	Atk int
	AtkRank bp.Rank
	SpAtk int
	SpAtkRank bp.Rank
	Speed int
	Ability bp.Ability
}

func NewAttackerInfo(p *bp.Pokemon) AttackerInfo {
	return AttackerInfo{
		PokeName:p.Name,
		Level:p.Level,
		Types:p.Types,
		CurrentHP:p.Stat.CurrentHP,
		Atk:p.Stat.Atk,
		AtkRank:p.Rank.Atk,
		SpAtk:p.Stat.SpAtk,
		SpAtkRank:p.Rank.SpAtk,
		Speed:p.Stat.Speed,
		Ability:p.Ability,
	}
}

type DefenderInfo struct {
	PokeName bp.PokeName
	Level bp.Level
	Types bp.Types
	CurrentHP int
	Def int
	DefRank bp.Rank
	SpDef int
	SpDefRank bp.Rank
	Speed int
	Ability bp.Ability
}

func NewDefenderInfo(p *bp.Pokemon) DefenderInfo {
	return DefenderInfo{
		PokeName:p.Name,
		Level:p.Level,
		Types:p.Types,
		CurrentHP:p.Stat.CurrentHP,
		Def:p.Stat.Def,
		DefRank:p.Rank.Def,
		SpDef:p.Stat.SpDef,
		SpDefRank:p.Rank.SpDef,
		Speed:p.Stat.Speed,
		Ability:p.Ability,
	}
}

type DamageCalculator struct {
	Attacker AttackerInfo
	Defender DefenderInfo
	IsSingleDamage bool
	IsCritical bool
	RandBonus float64
	IsDamageCappedByCurrentHP bool
}

// https://latest.pokewiki.net/%E3%83%80%E3%83%A1%E3%83%BC%E3%82%B8%E8%A8%88%E7%AE%97%E5%BC%8F
func (c *DamageCalculator) Calculation(moveName bp.MoveName) DamageDetailResult {
	attacker := c.Attacker
	defender := c.Defender
	moveData := bp.MOVEDEX[moveName]
	effective := bp.TYPEDEX.Effective(moveData.Type, defender.Types)

	switch moveName {
		//がむしゃら
		case bp.ENDEAVOR:
			var dmg int
			if effective == bp.NO_EFFECTIVE {
				dmg = 0
			} else {
				dmg = defender.CurrentHP - attacker.CurrentHP
			}

			//ダメージがマイナスならば、0にする。
			dmg = omwmath.Max(dmg, 0)
			return DamageDetailResult{
				Damage:dmg,	
				IsEndeavorFailure:dmg == 0 || effective == bp.NO_EFFECTIVE,
			}
	}

	if effective == bp.NO_EFFECTIVE {
		return DamageDetailResult{TypeEffective:bp.NO_EFFECTIVE}
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
	var power int
	if moveName == bp.GYRO_BALL {
		// https://wiki.xn--rckteqa2e.com/wiki/%E3%82%B8%E3%83%A3%E3%82%A4%E3%83%AD%E3%83%9C%E3%83%BC%E3%83%AB
		power = int(25.0 * float64(defender.Speed) / float64(attacker.Speed))+1
		power = omwmath.Min(power, 150)
	} else {
		power = moveData.Power
	}

	lv := int(attacker.Level)
	dmg := (lv*2/5) + 2
	dmg = dmg * power * atkVal / defVal
	dmg = (dmg/50) + 2

	//ダブルバトルの複数攻撃
	var rangeAtkBonus float64
	if c.IsSingleDamage {
		rangeAtkBonus = NO_RANGE_ATTACK_BONUS
	} else {
		rangeAtkBonus = RANGE_ATTACK_BONUS
	}
	dmg = RoundOverHalf(float64(dmg) * rangeAtkBonus / NO_RANGE_ATTACK_BONUS)

	//急所
	var critBonus float64
	if c.IsCritical {
		critBonus = CRITICAL_BONUS
	} else {
		critBonus = NO_CRITICAL_BONUS
	}
	dmg = RoundOverHalf(float64(dmg) * critBonus / NO_CRITICAL_BONUS)

	//乱数
	dmg = int(float64(dmg) * float64(c.RandBonus))

	//タイプ一致
	var stab float64
	if slices.Contains(attacker.Types, moveData.Type) {
		stab = SAME_TYPE_ATTACK_BONUS
	} else {
		stab = NO_SAME_TYPE_ATTACK_BONUS
	}
	dmg = RoundOverHalf(float64(dmg) * stab / NO_SAME_TYPE_ATTACK_BONUS)

	//タイプ相性
	effectiveness := bp.TYPEDEX.Effectiveness(moveData.Type, defender.Types)
	dmg = int(float64(dmg) * effectiveness)

	//ダメージが0ならば、1にする。ただしタイプ相性が無効であれば、0のまま。
	if effective != bp.NO_EFFECTIVE {
		dmg = omwmath.Max(dmg, 1)
	}

	if c.IsDamageCappedByCurrentHP {
		dmg = omwmath.Min(dmg, defender.CurrentHP)
	}

	return DamageDetailResult{
		Damage:omwmath.Max(dmg, 1),
		TypeEffective:effective,
	}
}

type DamageDetailResult struct {
	Damage int
	TypeEffective bp.TypeEffective
	IsEndeavorFailure bool
}