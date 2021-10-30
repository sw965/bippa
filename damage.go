package bippa

type RandomDamage float64

var RANDOM_DAMAGES = []RandomDamage{
	0.85, 0.86, 0.87, 0.88, 0.89, 0.9, 0.91, 0.92, 0.93, 0.94, 0.95, 0.96, 0.97, 0.98, 0.99, 1.0,
}

type DamageBonus int

const (
  INIT_DAMAGE_BONUS = 4096
)

func NewDamageBonus(spovb *SelfPointOfViewBattle) DamageBonus {
	damageBonus := INIT_DAMAGE_BONUS
	if spovb.SelfFighters[0].Item == "いのちのたま" {
		damageBonus = damageBonus.MulLifeOrb()
	}
	return damageBonus
}

func (damageBonus DamageBonus) MulLifeOrb() DamageBonus {
  result := RoundingZeroPointFiveOrMore(float64(damageBonus) * 5324.0 / 4096.0)
  return DamageBonus(result)
}

type FinalDamage int

func FinalDamageCalc(spovb *SelfPointOfViewBattle, moveName MoveName, isCritical bool, randomDamage RandomDamage) (FinalDamage, error) {
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
	finalDefense, err := NewFinalDefense(spovb, moveName, isCritical)
	if err != nil {
		return 0, err
	}

	var criticalBonus float64

	if isCritical {
		criticalBonus = 6144.0 / 4096.0
	} else {
		criticalBonus = 4096.0 / 4096.0
	}

	sameTypeAttackBonus := spovb.SelfFighters[0].SameTypeAttackBonus(moveName)
	effectivenessBonus := spovb.OpponentFighters[0].EffectivenessBonus(moveName)

	var burnBonus float64
	if spovb.SelfFighters[0].StatusAilmentParam.StatusAilment == BURN && moveData.Category == PHYSICS {
		burnBonus = 2048.0 / 4096.0
	} else {
		burnBonus = 4096.0 / 4096.0
	}

	damageBonus := spovb.DamageBonus(moveName)

	result := LEVEL*2/5 + 2
	result = result * finalPower * finalAttack / finalDefense
	result = result/50 + 2
	result = RoundingZeroPointFiveOver(float64(result) * criticalBonus)
	result = int(float64(result) * randomNum)
	result = RoundingZeroPointFiveOver(float64(result) * float64(sameTypeAttackBonus))
	result = int(float64(result) * effectivenessBonus)
	result = RoundingZeroPointFiveOver(float64(result) * burnBonus)
	result = RoundingZeroPointFiveOver(float64(result) * float64(damageBonus) / 4096.0)
	return result, nil
}
