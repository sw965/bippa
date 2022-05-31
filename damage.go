package bippa

import (
	"fmt"
	"math/rand"
)

//小数点以下がが0.5以上ならば、繰り上げ
func FiveOrMoreRounding(x float64) int {
	afterTheDecimalPoint := float64(x) - float64(int(x))
	if afterTheDecimalPoint >= 0.5 {
		return int(x + 1)
	}
	return int(x)
}

//小数点以下が0.5より大きいならば、繰り上げ
func FiveOverRounding(x float64) int {
	afterTheDecimalPoint := float64(x) - float64(int(x))
	if afterTheDecimalPoint > 0.5 {
		return int(x + 1)
	}
	return int(x)
}

type PhysicsAttackBonus int

const (
	INIT_PHYSICS_ATTACK_BONUS = PhysicsAttackBonus(4096)
)

func NewPhysicsAttackBonus(battle *Battle) PhysicsAttackBonus {
	result := int(INIT_PHYSICS_ATTACK_BONUS)
	if battle.P1Fighters[0].Item == "こだわりハチマキ" {
		result = FiveOrMoreRounding(float64(result) * 6144.0 / 4096.0)
	}
	return PhysicsAttackBonus(result)
}

type SpecialAttackBonus int

const (
	INIT_SPECIAL_ATTACK_BONUS = SpecialAttackBonus(4096)
)

func NewSpecialAttackBonus(battle *Battle) SpecialAttackBonus {
	result := int(INIT_SPECIAL_ATTACK_BONUS)
	if battle.P1Fighters[0].Item == "こだわりメガネ" {
		result = FiveOrMoreRounding(float64(result) * 6144.0 / 4096.0)
	}
	return SpecialAttackBonus(result)
}

type AttackBonus int

func NewAttackBonus(battle *Battle, moveName MoveName) (AttackBonus, error) {
	moveData := MOVEDEX[moveName]
	switch moveData.Category {
	case PHYSICS:
		physicsAttackBonus := NewPhysicsAttackBonus(battle)
		return AttackBonus(physicsAttackBonus), nil
	case SPECIAL:
		specialAttackBonus := NewSpecialAttackBonus(battle)
		return AttackBonus(specialAttackBonus), nil
	default:
		return 0, fmt.Errorf("物理/特殊技でなければならない")
	}
}

type FinalAttack int

func NewFinalAttack(battle *Battle, moveName MoveName, isCritical bool) (FinalAttack, error) {
	moveData := MOVEDEX[moveName]

	var attack int
	var attackRankVal RankVal

	switch moveData.Category {
	case PHYSICS:
		attack = battle.P1Fighters[0].Atk
		attackRankVal = battle.P1Fighters[0].Rank.Atk
	case SPECIAL:
		attack = battle.P1Fighters[0].SpAtk
		attackRankVal = battle.P1Fighters[0].Rank.SpAtk
	}

	//変化技の場合、ここでエラーが起きるので、上のswitch文ではチェック不要
	attackBonus, err := NewAttackBonus(battle, moveName)

	if err != nil {
		return 0, err
	}

	if attackRankVal < 0 && isCritical {
		attackRankVal = 0
	}

	rankBonus := attackRankVal.ToBonus()

	result := int(float64(attack) * float64(rankBonus))
	result = FiveOverRounding(float64(result) * float64(attackBonus) / 4096.0)
	if result < 1 {
		return 1, nil
	} else {
		return FinalAttack(result), nil
	}
}

type DefenseBonus int

const (
	INIT_DEFENSE_BONUS = DefenseBonus(4096)
)

func NewDefenseBonus(battle *Battle) DefenseBonus {
	result := INIT_DEFENSE_BONUS
	if battle.P1Fighters[0].Item == "とつげきチョッキ" {
		tmp := FiveOrMoreRounding(float64(result) * (6144.0 / 4096.0))
		result = DefenseBonus(tmp)
	}
	return result
}

type FinalDefense int

func NewFinalDefense(battle *Battle, moveName MoveName, isCritical bool) (FinalDefense, error) {
	moveData := MOVEDEX[moveName]

	var defense int
	var defenseRankVal RankVal

	switch moveData.Category {
	case PHYSICS:
		defense = battle.P1Fighters[0].Def
		defenseRankVal = battle.P1Fighters[0].Rank.Def
	case SPECIAL:
		defense = battle.P1Fighters[0].SpDef
		defenseRankVal = battle.P1Fighters[0].Rank.SpDef
	default:
		return 0, fmt.Errorf("物理/特殊技でなければならない")
	}

	if defenseRankVal > 0 && isCritical {
		defenseRankVal = 0
	}

	rankBonus := defenseRankVal.ToBonus()
	result := int(float64(defense) * float64(rankBonus))

	if result < 1 {
		return 1, nil
	}
	return FinalDefense(result), nil
}

//https://latest.pokewiki.net/%E3%83%80%E3%83%A1%E3%83%BC%E3%82%B8%E8%A8%88%E7%AE%97%E5%BC%8F
type PowerBonus int

const (
	INIT_POWER_BONUS = PowerBonus(4096)
)

type FinalPower int

func NewFinalPower(battle *Battle, moveName MoveName) (FinalPower, error) {
	moveData := MOVEDEX[moveName]

	if moveData.Category == STATUS {
		return 0, fmt.Errorf("物理/特殊技でなければならない")
	}

	power := moveData.Power
	powerBonus := INIT_POWER_BONUS

	result := FiveOverRounding(float64(power) * float64(powerBonus) / 4096.0)
	if result < 1 {
		return 1, nil
	}
	return FinalPower(result), nil
}

type CriticalBonus float64

var (
	CRITICAL_BONUS    = CriticalBonus(6144.0 / 4096.0)
	NO_CRITICAL_BONUS = CriticalBonus(4096.0 / 4096.0)
)

var CRITICAL_N = map[CriticalRank]int{0: 24, 1: 8, 2: 2, 3: 1}
var BOOL_TO_CRITICAL_BONUS = map[bool]CriticalBonus{true: CRITICAL_BONUS, false: NO_CRITICAL_BONUS}

type SameTypeAttackBonus float64

const (
	SAME_TYPE_ATTACK_BONUS    = SameTypeAttackBonus(6144.0 / 4096.0)
	NO_SAME_TYPE_ATTACK_BONUS = SameTypeAttackBonus(4096.0 / 4096.0)
)

var BOOL_TO_SAME_TYPE_ATTACK_BONUS = map[bool]SameTypeAttackBonus{
	true: SAME_TYPE_ATTACK_BONUS, false: NO_SAME_TYPE_ATTACK_BONUS,
}

type EffectivenessBonus float64

//https://latest.pokewiki.net/%E3%83%80%E3%83%A1%E3%83%BC%E3%82%B8%E8%A8%88%E7%AE%97%E5%BC%8F
type RandomDamageBonus float64

func NewRandomDamageBonus(random *rand.Rand) RandomDamageBonus {
	index := random.Intn(RANDOM_DAMAGE_BONUSES_LENGTH)
	return RANDOM_DAMAGE_BONUSES[index]
}

type RandomDamageBonuses []RandomDamageBonus

func (randomDamageBonuses RandomDamageBonuses) Average() RandomDamageBonus {
	sum := RandomDamageBonus(0.0)
	for _, randomDamageBonus := range randomDamageBonuses {
		sum += randomDamageBonus
	}
	return RandomDamageBonus(sum) / RandomDamageBonus(RANDOM_DAMAGE_BONUSES_LENGTH)
}

func (randomDamageBonuses RandomDamageBonuses) Max() RandomDamageBonus {
	result := randomDamageBonuses[0]
	for _, v := range randomDamageBonuses[1:] {
		if v > result {
			result = v
		}
	}
	return result
}

var RANDOM_DAMAGE_BONUSES = RandomDamageBonuses{
	0.85, 0.86, 0.87, 0.88, 0.89, 0.9, 0.91, 0.92, 0.93, 0.94, 0.95, 0.96, 0.97, 0.98, 0.99, 1.0,
}

var RANDOM_DAMAGE_BONUSES_LENGTH = len(RANDOM_DAMAGE_BONUSES)

var MAX_RANDOM_DAMAGE_BONUS = RANDOM_DAMAGE_BONUSES.Max()
var AVERAGE_RANDOM_DAMAGE_BONUS = RANDOM_DAMAGE_BONUSES.Average()

type DamageBonus int

const (
	INIT_DAMAGE_BONUS = DamageBonus(4096)
)

func NewDamageBonus(battle *Battle) DamageBonus {
	result := INIT_DAMAGE_BONUS
	if battle.P1Fighters[0].Item == "いのちのたま" {
		tmp := FiveOrMoreRounding(float64(result) * 5324.0 / 4096.0)
		result = DamageBonus(tmp)
	}
	return result
}

type FinalDamage int

func NewFinalDamage(battle *Battle, moveName MoveName, isCritical bool, randomDamageBonus RandomDamageBonus) (FinalDamage, error) {
	finalPower, err := NewFinalPower(battle, moveName)
	if err != nil {
		return 0, err
	}

	finalAttack, err := NewFinalAttack(battle, moveName, isCritical)
	if err != nil {
		return 0, err
	}

	reverseBattle := battle.Reverse()
	finalDefense, err := NewFinalDefense(&reverseBattle, moveName, isCritical)
	if err != nil {
		return 0, err
	}

	criticalBonus := BOOL_TO_CRITICAL_BONUS[isCritical]
	sameTypeAttackBonus := battle.P1Fighters[0].SameTypeAttackBonus(moveName)
	effectivenessBonus := battle.P2Fighters[0].EffectivenessBonus(moveName)

	damageBonus := NewDamageBonus(battle)

	result := int(DEFAULT_LEVEL)*2/5 + 2
	result = int(float64(result) * float64(finalPower) * float64(finalAttack) / float64(finalDefense))
	result = result/50 + 2
	result = FiveOverRounding(float64(result) * float64(criticalBonus))
	result = int(float64(result) * float64(randomDamageBonus))
	result = FiveOverRounding(float64(result) * float64(sameTypeAttackBonus))
	result = int(float64(result) * float64(effectivenessBonus))
	result = FiveOverRounding(float64(result) * float64(damageBonus) / 4096.0)
	return FinalDamage(result), nil
}
