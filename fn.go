package bippa

import (
	"math/rand"
)

func CanLearn(pokeName PokeName, moveName MoveName) bool {
	pokeData, _ := POKEDEX[pokeName]
	for _, iMoveName := range pokeData.Learnset {
		if iMoveName == moveName {
			return true
		}
	}
	return false
}

func HpStateCalc(baseHP int, individual Individual, effort Effort) int {
	return ((baseHP*2)+int(individual)+(int(effort)/4))*LEVEL/100 + LEVEL + 10
}

func StateCalc(baseState int, individual Individual, effort Effort, natureBonus float64) int {
	result := ((baseState*2)+int(individual)+(int(effort)/4))*LEVEL/100 + 5
	return int(float64(result) * natureBonus)
}

func PowerPointCalc(basePP int, pointUp PointUp) int {
	v := ((5.0 + float64(pointUp)) / 5.0)
	return int(float64(basePP) * v)
}

func MakeMaxPointUps(length int) []PointUp {
	result := make([]PointUp, length)
	for i := 0; i < length; i++ {
		result[i] = MAX_POINT_UP
	}
	return result
}

func AfterTheDecimalPoint(x float64) float64 {
	return float64(x) - float64(int(x))
}

func RoundingZeroPointFiveOrMore(x float64) int {
	afterTheDecimalPoint := AfterTheDecimalPoint(x)
	if afterTheDecimalPoint >= 0.5 {
		return int(x + 1)
	}
	return int(x)
}

func RoundingZeroPointFiveOver(x float64) int {
	afterTheDecimalPoint := AfterTheDecimalPoint(x)
	if afterTheDecimalPoint > 0.5 {
		return int(x + 1)
	}
	return int(x)
}

func IsHit(percent int, random *rand.Rand) bool {
	return random.Intn(100) < percent
}

func IsCritical(random *rand.Rand) bool {
	return random.Intn(24) == 0
}

func DamageRandomNum(random *rand.Rand) float64 {
	index := random.Intn(len(DAMAGE_RANDOM_NUMS))
	return DAMAGE_RANDOM_NUMS[index]
}

func TurnEndLeftovers(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.SelfFighters[0].Item != "たべのこし" {
		return spovb
	}

	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	if spovb.SelfFighters[0].IsFullHP() {
		return spovb
	}

	heal := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 16.0)
	spovb = spovb.Heal(heal)
	return spovb
}

func TurnEndBlackSludge(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.SelfFighters[0].Item != "くろいヘドロ" {
		return spovb
	}

	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	if spovb.SelfFighters[0].InType(POISON) {
		if spovb.SelfFighters[0].IsFullHP() {
			return spovb
		}

		heal := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 16.0)
		spovb = spovb.Heal(heal)
	} else {
		damage := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 8.0)
		spovb = spovb.ToDamage(damage)
	}
	return spovb
}

func TurnEndLeechSeed(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	if spovb.OpponentFighters[0].IsFaint() {
		return spovb
	}

	if !spovb.OpponentFighters[0].IsLeechSeedState {
		return spovb
	}

	damageAndHealValue := int(float64(spovb.OpponentFighters[0].State.MaxHP) * 1.0 / 8.0)
	opovb := spovb.SwitchPointOfView()
	opovb = opovb.ToDamage(damageAndHealValue)
	spovb = opovb.SwitchPointOfView()
	spovb = spovb.Heal(damageAndHealValue)
	return spovb
}

func TurnEndBadPoison(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	if spovb.SelfFighters[0].StatusAilmentParam.StatusAilment != BAD_POISON {
		return spovb
	}

	if spovb.SelfFighters[0].StatusAilmentParam.BadPoisonElapsedTurn < 15 {
		spovb.SelfFighters[0].StatusAilmentParam.BadPoisonElapsedTurn += 1
	}

	damage := spovb.SelfFighters[0].BadPoisonDamage()
	spovb = spovb.ToDamage(damage)
	return spovb
}

func TurnEndBurn(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	if spovb.SelfFighters[0].StatusAilmentParam.StatusAilment != BURN {
		return spovb
	}

	damage := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 16.0)
	spovb = spovb.ToDamage(damage)
	return spovb
}
