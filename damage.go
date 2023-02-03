package bippa

import (
	"fmt"
	"github.com/sw965/omw"
	"math/rand"
)

// 小数点以下がが0.5以上ならば、繰り上げ
func FiveOrMoreRounding(x float64) int {
	afterTheDecimalPoint := float64(x) - float64(int(x))
	if afterTheDecimalPoint >= 0.5 {
		return int(x + 1)
	}
	return int(x)
}

// 小数点以下が0.5より大きいならば、繰り上げ
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

func NewPhysicsAttackBonus(pokemon *Pokemon) PhysicsAttackBonus {
	result := int(INIT_PHYSICS_ATTACK_BONUS)
	if pokemon.Item == "こだわりハチマキ" {
		result = FiveOrMoreRounding(float64(result) * 6144.0 / 4096.0)
	}
	return PhysicsAttackBonus(result)
}

type SpecialAttackBonus int

const (
	INIT_SPECIAL_ATTACK_BONUS = SpecialAttackBonus(4096)
)

func NewSpecialAttackBonus(pokemon *Pokemon) SpecialAttackBonus {
	result := int(INIT_SPECIAL_ATTACK_BONUS)
	if pokemon.Item == "こだわりメガネ" {
		result = FiveOrMoreRounding(float64(result) * 6144.0 / 4096.0)
	}
	return SpecialAttackBonus(result)
}

type AttackBonus int

func NewAttackBonus(pokemon *Pokemon, moveName MoveName) (AttackBonus, error) {
	moveData := MOVEDEX[moveName]
	switch moveData.Category {
	case PHYSICS:
		physicsAttackBonus := NewPhysicsAttackBonus(pokemon)
		return AttackBonus(physicsAttackBonus), nil
	case SPECIAL:
		specialAttackBonus := NewSpecialAttackBonus(pokemon)
		return AttackBonus(specialAttackBonus), nil
	default:
		return 0, fmt.Errorf("物理/特殊技でなければならない")
	}
}

type FinalAttack int

func NewFinalAttack(pokemon *Pokemon, moveName MoveName, isCritical bool) (FinalAttack, error) {
	moveData := MOVEDEX[moveName]

	var attack int
	var attackRank Rank

	switch moveData.Category {
	case PHYSICS:
		attack = pokemon.Atk
		attackRank = pokemon.RankState.Atk
	case SPECIAL:
		attack = pokemon.SpAtk
		attackRank = pokemon.RankState.SpAtk
	}

	//変化技の場合、ここでエラーが起きるので、上のswitch文ではチェック不要
	attackBonus, err := NewAttackBonus(pokemon, moveName)

	if err != nil {
		return 0, err
	}

	if attackRank < 0 && isCritical {
		attackRank = 0
	}

	rankBonus := attackRank.ToBonus()

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

func NewDefenseBonus(pokemon *Pokemon) DefenseBonus {
	result := INIT_DEFENSE_BONUS
	if pokemon.Item == "とつげきチョッキ" {
		tmp := FiveOrMoreRounding(float64(result) * (6144.0 / 4096.0))
		result = DefenseBonus(tmp)
	}
	return result
}

type FinalDefense int

func NewFinalDefense(pokemon *Pokemon, moveName MoveName, isCritical bool) (FinalDefense, error) {
	moveData := MOVEDEX[moveName]

	var defense int
	var defenseRank Rank

	switch moveData.Category {
	case PHYSICS:
		defense = pokemon.Def
		defenseRank = pokemon.RankState.Def
	case SPECIAL:
		defense = pokemon.SpDef
		defenseRank = pokemon.RankState.SpDef
	default:
		return 0, fmt.Errorf("物理/特殊技でなければならない")
	}

	if defenseRank > 0 && isCritical {
		defenseRank = 0
	}

	rankBonus := defenseRank.ToBonus()
	result := int(float64(defense) * float64(rankBonus))

	if result < 1 {
		return 1, nil
	}
	return FinalDefense(result), nil
}

// https://latest.pokewiki.net/%E3%83%80%E3%83%A1%E3%83%BC%E3%82%B8%E8%A8%88%E7%AE%97%E5%BC%8F
type PowerBonus int

const (
	INIT_POWER_BONUS = PowerBonus(4096)
)

type FinalPower int

func NewFinalPower(moveName MoveName) (FinalPower, error) {
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

func NewCriticalBonus(x bool) CriticalBonus {
	if x {
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

func NewSameTypeAttackBonus(x bool) SameTypeAttackBonus {
	if x {
		return SAME_TYPE_ATTACK_BONUS
	} else {
		return NO_SAME_TYPE_ATTACK_BONUS
	}
}

type EffectivenessBonus float64

// https://latest.pokewiki.net/%E3%83%80%E3%83%A1%E3%83%BC%E3%82%B8%E8%A8%88%E7%AE%97%E5%BC%8F
type RandomDamageBonus float64

func NewRandomDamageBonus(random *rand.Rand) RandomDamageBonus {
	index := random.Intn(RANDOM_DAMAGE_BONUSES_LENGTH)
	return RANDOM_DAMAGE_BONUSES[index]
}

type RandomDamageBonuses []RandomDamageBonus

var RANDOM_DAMAGE_BONUSES = RandomDamageBonuses{
	0.85, 0.86, 0.87, 0.88, 0.89, 0.9, 0.91, 0.92, 0.93, 0.94, 0.95, 0.96, 0.97, 0.98, 0.99, 1.0,
}

var RANDOM_DAMAGE_BONUSES_LENGTH = len(RANDOM_DAMAGE_BONUSES)
var MAX_RANDOM_DAMAGE_BONUS = omw.Max(RANDOM_DAMAGE_BONUSES...)
var MEAN_RANDOM_DAMAGE_BONUS = omw.Mean(RANDOM_DAMAGE_BONUSES...)

type DamageBonus int

const (
	INIT_DAMAGE_BONUS = DamageBonus(4096)
)

func NewDamageBonus(pokemon *Pokemon) DamageBonus {
	result := INIT_DAMAGE_BONUS
	if pokemon.Item == "いのちのたま" {
		tmp := FiveOrMoreRounding(float64(result) * 5324.0 / 4096.0)
		result = DamageBonus(tmp)
	}
	return result
}

type FinalDamage int

func NewFinalDamage(attackPokemon, defensePokemon *Pokemon, moveName MoveName, isCritical bool, randomDamageBonus RandomDamageBonus) (FinalDamage, error) {
	finalPower, err := NewFinalPower(moveName)
	if err != nil {
		return 0, err
	}

	finalAttack, err := NewFinalAttack(attackPokemon, moveName, isCritical)
	if err != nil {
		return 0, err
	}

	finalDefense, err := NewFinalDefense(defensePokemon, moveName, isCritical)
	if err != nil {
		return 0, err
	}

	criticalBonus := NewCriticalBonus(isCritical)
	sameTypeAttackBonus := attackPokemon.SameTypeAttackBonus(moveName)
	effectivenessBonus := defensePokemon.EffectivenessBonus(moveName)

	damageBonus := NewDamageBonus(attackPokemon)

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

type DamageProbabilityDistribution map[int]float64

func NewAttackDamageProbabilityDistribution(attackPokemon, defensePokemon *Pokemon, moveName MoveName, accuracy, criticalN int) (DamageProbabilityDistribution, error) {
	accuracyPercent := float64(accuracy) / 100.0
	randomDamageBonusPercent := 1.0 / float64(RANDOM_DAMAGE_BONUSES_LENGTH)
	criticalPercent := 1.0 / float64(criticalN)
	boolToCriticalPercent := map[bool]float64{true: criticalPercent, false: 1.0 - criticalPercent}
	result := map[int]float64{0: 1.0 - accuracyPercent}

	for _, randomDamageBonus := range RANDOM_DAMAGE_BONUSES {
		for _, isCritical := range []bool{true, false} {
			finalDamage, err := NewFinalDamage(attackPokemon, defensePokemon, moveName, isCritical, randomDamageBonus)

			if err != nil {
				return result, err
			}

			p := accuracyPercent * randomDamageBonusPercent * boolToCriticalPercent[isCritical]

			if _, ok := result[int(finalDamage)]; ok {
				//確率の加法定理
				result[int(finalDamage)] += p
			} else {
				result[int(finalDamage)] = p
			}
		}
	}
	return result, nil
}

func (dpd DamageProbabilityDistribution) RatioExpected(v float64) float64 {
	result := 0.0
	for damage, percent := range dpd {
		attackDamageRatio := float64(damage) / float64(v)
		if attackDamageRatio > 1.0 {
			attackDamageRatio = 1.0
		}
		result += attackDamageRatio * percent
	}
	return result
}

func (dpd DamageProbabilityDistribution) StandardFeatureValue(v float64, featureSize int) []float64 {
	width := 1.0 / float64(featureSize)
	under := 0.0
	upper := width

	adpdRatioExpected := dpd.RatioExpected(v)
	result := make([]float64, featureSize)

	for i := 0; i < featureSize; i++ {
		if i == (featureSize - 1) {
			result[i] = 1.0
			break
		}

		if adpdRatioExpected >= under && adpdRatioExpected < upper {
			result[i] = 1.0
			break
		}

		under += width
		upper += width
	}
	return result
}
