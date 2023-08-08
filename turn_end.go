package bippa

import (
	"github.com/sw965/omw"
	"golang.org/x/exp/slices"
)

// https://wiki.xn--rckteqa2e.com/wiki/%E3%82%BF%E3%83%BC%E3%83%B3#5..E3.82.BF.E3.83.BC.E3.83.B3.E7.B5.82.E4.BA.86.E6.99.82.E3.81.AE.E5.87.A6.E7.90.86
func TurnEndLeftovers(bt *Battle) {
	if bt.P1Fighters[0].Item != "たべのこし" {
		return
	}

	if bt.P1Fighters[0].IsFullHP() {
		return
	}

	heal := int(float64(bt.P1Fighters[0].MaxHP) * 1.0 / 16.0)
	bt.Heal(heal)
}

func TurnEndBlackSludge(bt *Battle) {
	if bt.P1Fighters[0].Item != "くろいヘドロ" {
		return
	}

	if slices.Contains(bt.P1Fighters[0].Types, POISON) {
		heal := int(float64(bt.P1Fighters[0].MaxHP) * 1.0 / 16.0)
		heal = omw.Max(heal, 1)
		bt.Heal(heal)
	} else {
		dmg := int(float64(bt.P1Fighters[0].MaxHP) * 1.0 / 8.0)
		dmg = omw.Max(dmg, 1)
		bt.Damage(dmg)
	}
}

func TurnEndLeechSeed(bt *Battle) {
	if bt.P1Fighters[0].IsFaint() {
		return
	}

	if bt.P2Fighters[0].IsFaint() {
		return
	}

	if !bt.P2Fighters[0].IsLeechSeed {
		return
	}

	dmg := int(float64(bt.P2Fighters[0].MaxHP) * 1.0 / 8.0)
	heal := dmg

	bt.Reverse()
	bt.Damage(dmg)
	bt.Reverse()
	bt.Heal(heal)
}

func TurnEndBadPoison(bt *Battle) {
	if bt.P1Fighters[0].StatusAilment != BAD_POISON {
		return
	}

	if bt.P1Fighters[0].BadPoisonElapsedTurn < 16 {
		bt.P1Fighters[0].BadPoisonElapsedTurn += 1
	}

	dmg := bt.P1Fighters[0].BadPoisonDamage()
	bt.Damage(dmg)
}