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

	if !spovb.OpponentFighters[0].IsLeechSeed {
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

	if spovb.SelfFighters[0].StatusAilmentDetail.StatusAilment != BAD_POISON {
		return spovb
	}

	if spovb.SelfFighters[0].StatusAilmentDetail.BadPoisonElapsedTurn < 15 {
		spovb.SelfFighters[0].StatusAilmentDetail.BadPoisonElapsedTurn += 1
	}

	damage := spovb.SelfFighters[0].BadPoisonDamage()
	spovb = spovb.ToDamage(damage)
	return spovb
}

func TurnEndBurn(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	if spovb.SelfFighters[0].StatusAilmentDetail.StatusAilment != BURN {
		return spovb
	}

	damage := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 16.0)
	spovb = spovb.ToDamage(damage)
	return spovb
}
