package bippa

import (
	"math/rand"
)

type DamageR float64

func NewDamageR(random *rand.Rand) DamageR {
	index := random.Intn(DAMAGE_RS_LENGTH)
	return DAMAGE_RS[index]
}

type DamageRs []DamageR

var DAMAGE_RS = []DamageR{
	0.85, 0.86, 0.87, 0.88, 0.89, 0.9, 0.91, 0.92, 0.93, 0.94, 0.95, 0.96, 0.97, 0.98, 0.99, 1.0,
}

var DAMAGE_RS_LENGTH = len(DAMAGE_RS)

type DamageBonus int

const (
  INIT_DAMAGE_BONUS = DamageBonus(4096)
)

func NewDamageBonus(spovb *SelfPointOfViewBattle) DamageBonus {
	result := INIT_DAMAGE_BONUS
	if spovb.SelfFighters[0].Item == "いのちのたま" {
		result = result.MulLifeOrb()
	}
	return result
}

func (damageBonus DamageBonus) MulLifeOrb() DamageBonus {
  result := RoundingZeroPointFiveOrMore(float64(damageBonus) * 5324.0 / 4096.0)
  return DamageBonus(result)
}

type FinalDamage int

func NewFinalDamage(spovb *SelfPointOfViewBattle, moveName MoveName, isCritical bool, damageR DamageR) (FinalDamage, error) {
	moveData := MOVEDEX[moveName]

	finalPower, err := NewFinalPower(spovb, moveName)
	if err != nil {
		return 0, err
	}

	finalAttack, err := NewFinalAttack(spovb, moveName, isCritical)
	if err != nil {
		return 0, err
	}

	opovb := spovb.SwitchPointOfView()
	finalDefense, err := NewFinalDefense(&opovb, moveName, isCritical)
	if err != nil {
		return 0, err
	}

	criticalBonus := BOOL_TO_CRITICAL_BONUS[isCritical]

	sameTypeAttackBonus := spovb.SelfFighters[0].NewSameTypeAttackBonus(moveName)
	effectivenessBonus := spovb.OpponentFighters[0].NewEffectivenessBonus(moveName)

	isBurn := spovb.SelfFighters[0].StatusAilment.Type == BURN
	isPhysics := moveData.Category == PHYSICS
	isBurnValid := isBurn && isPhysics
	burnBonus := BOOL_TO_BURN_BONUS[isBurnValid]

	damageBonus := NewDamageBonus(spovb)

	result := int(MAX_LEVEL)*2/5 + 2
	result = int(float64(result) * float64(finalPower) * float64(finalAttack) / float64(finalDefense))
	result = result/50 + 2
	result = RoundingZeroPointFiveOver(float64(result) * float64(criticalBonus))
	result = int(float64(result) * float64(damageR))
	result = RoundingZeroPointFiveOver(float64(result) * float64(sameTypeAttackBonus))
	result = int(float64(result) * float64(effectivenessBonus))
	result = RoundingZeroPointFiveOver(float64(result) * float64(burnBonus))
	result = RoundingZeroPointFiveOver(float64(result) * float64(damageBonus) / 4096.0)
	return FinalDamage(result), nil
}

type RealDamage int

func NewRealDamage(spovb *SelfPointOfViewBattle, finalDamage FinalDamage) RealDamage {
	if spovb.OpponentFighters[0].IsFocusSashOk(int(finalDamage)) {
		result := int(spovb.OpponentFighters[0].State.MaxHP) - 1
		return RealDamage(result)
	} else {
		return RealDamage(finalDamage)
	}
}
