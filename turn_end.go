package bippa

import (
	"github.com/sw965/omw"
)

// https://wiki.xn--rckteqa2e.com/wiki/%E3%82%BF%E3%83%BC%E3%83%B3#5..E3.82.BF.E3.83.BC.E3.83.B3.E7.B5.82.E4.BA.86.E6.99.82.E3.81.AE.E5.87.A6.E7.90.86
func TurnEndLeftovers(battle Battle) Battle {
	if battle.P1Fighters[0].Item != "たべのこし" {
		return battle
	}

	if battle.P1Fighters[0].IsFullHP() {
		return battle
	}

	heal := int(float64(battle.P1Fighters[0].MaxHP) * 1.0 / 16.0)
	battle = battle.Heal(heal)
	return battle
}

func TurnEndBlackSludge(battle Battle) Battle {
	if battle.P1Fighters[0].Item != "くろいヘドロ" {
		return battle
	}

	if omw.Contains(battle.P1Fighters[0].Types, POISON) {
		heal := int(float64(battle.P1Fighters[0].MaxHP) * 1.0 / 16.0)
		heal = omw.Max(heal, 1)
		battle = battle.Heal(heal)
	} else {
		damage := int(float64(battle.P1Fighters[0].MaxHP) * 1.0 / 8.0)
		damage = omw.Max(damage, 1)
		battle = battle.Damage(damage)
	}
	return battle
}

func TurnEndLeechSeed(battle Battle) Battle {
	if battle.P1Fighters[0].IsFaint() {
		return battle
	}

	if battle.P2Fighters[0].IsFaint() {
		return battle
	}

	if !battle.P2Fighters[0].IsLeechSeed {
		return battle
	}

	damage := int(float64(battle.P2Fighters[0].MaxHP) * 1.0 / 8.0)
	heal := damage

	battle = battle.Reverse()
	battle = battle.Damage(damage)
	battle = battle.Reverse()
	battle = battle.Heal(heal)
	return battle
}

func TurnEndBadPoison(battle Battle) Battle {
	if battle.P1Fighters[0].StatusAilment != BAD_POISON {
		return battle
	}

	if battle.P1Fighters[0].BadPoisonElapsedTurn < 16 {
		battle.P1Fighters[0].BadPoisonElapsedTurn += 1
	}

	damage := battle.P1Fighters[0].BadPoisonDamage()
	return battle.Damage(damage)
}
