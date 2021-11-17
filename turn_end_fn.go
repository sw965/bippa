package bippa

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

	if spovb.SelfFighters[0].Types.In(POISON) {
		if spovb.SelfFighters[0].IsFullHP() {
			return spovb
		}

		heal := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 16.0)
		spovb = spovb.Heal(heal)
	} else {
		damage := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 8.0)
		spovb, sitrusBerryHeal := spovb.ToDamage(damage)
		spovb = sitrusBerryHeal(spovb)
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

	if !spovb.OpponentFighters[0].IsLeechSeed {
		return spovb
	}

	damageAndHealValue := int(float64(spovb.OpponentFighters[0].State.MaxHP) * 1.0 / 8.0)
	opovb := spovb.SwitchPointOfView()
	opovb, sitrusBerryHeal := opovb.ToDamage(damageAndHealValue)
	opovb = sitrusBerryHeal(opovb)
	spovb = opovb.SwitchPointOfView()
	spovb = spovb.Heal(damageAndHealValue)
	return spovb
}

func TurnEndNormalPoison(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	if spovb.SelfFighters[0].StatusAilment.Type != NORMAL_POISON {
		return spovb
	}

	return spovb
}

func TurnEndBadPoison(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	if spovb.SelfFighters[0].StatusAilment.Type != BAD_POISON {
		return spovb
	}

	if spovb.SelfFighters[0].StatusAilment.BadPoisonElapsedTurn < 15 {
		spovb.SelfFighters[0].StatusAilment.BadPoisonElapsedTurn += 1
	}

	damage := spovb.SelfFighters[0].BadPoisonDamage()
	spovb, sitrusBerryHeal := spovb.ToDamage(damage)
	return sitrusBerryHeal(spovb)
}

func TurnEndBurn(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	if spovb.SelfFighters[0].StatusAilment.Type != BURN {
		return spovb
	}

	damage := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 16.0)
	spovb, sitrusBerryHeal := spovb.ToDamage(damage)
	return sitrusBerryHeal(spovb)
}
