package bippa

//https://wiki.xn--rckteqa2e.com/wiki/%E3%82%BF%E3%83%BC%E3%83%B3#5..E3.82.BF.E3.83.BC.E3.83.B3.E7.B5.82.E4.BA.86.E6.99.82.E3.81.AE.E5.87.A6.E7.90.86

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
	return spovb.ToDamage(damage)
}

func TurnEndBurn(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	if spovb.SelfFighters[0].StatusAilment.Type != BURN {
		return spovb
	}

	damage := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 16.0)
	return spovb.ToDamage(damage)
}

func TurnEndWeather(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.ShareField.Weather.RemainingTurn > 0 {
		spovb.ShareField.Weather.RemainingTurn -= 1
	}

	if spovb.ShareField.Weather.RemainingTurn == 0 {
		spovb.ShareField.Weather.Type = ""
	}
	return spovb
}

func TurnEndHailDamage(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	if spovb.ShareField.Weather.Type != HAIL {
		return spovb
	}

	if spovb.SelfFighters[0].Types.In(ICE) {
		return spovb
	}

	if spovb.SelfFighters[0].Ability == "アイスボディ" {
		return spovb
	}

	damage := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 16.0)
	if damage < 1 {
		damage = 1
	}
	return spovb.ToDamage(damage)
}

func TurnEndSandstoremDamage(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	if spovb.ShareField.Weather.Type != SANDSTORM {
		return spovb
	}

	pokeTypes := spovb.SelfFighters[0].Types
	if pokeTypes.In(STEEL) || pokeTypes.In(ROCK) || pokeTypes.In(GROUND) {
		return spovb
	}

	if spovb.SelfFighters[0].Ability == "すなかき" {
		return spovb
	}

	damage := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 16.0)
	if damage < 1 {
		damage = 1
	}
	return spovb.ToDamage(damage)
}

//かんそうはだ
func TurnEndDrySkin(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	if spovb.SelfFighters[0].Ability != "かんそうはだ" {
		return spovb
	}

	weatherType := spovb.ShareField.Weather.Type
	switch weatherType {
		case RAIN:
			heal := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 8.0)
			return spovb.Heal(heal)
		case SUNNY_DAY:
			damage := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 8.0)
			return spovb.ToDamage(damage)
	}
	return spovb
}

//サンパワー
func TurnEndSolarPower(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	if spovb.SelfFighters[0].Ability != "サンパワー" {
		return spovb
	}

	if spovb.ShareField.Weather.Type != SUNNY_DAY {
		return spovb
	}

	damage := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 /8.0)
	return spovb.ToDamage(damage)
}

//あめうけざら
func TurnEndRainDish(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	if spovb.SelfFighters[0].Ability != "あめうけざら" {
		return spovb
	}

	if spovb.ShareField.Weather.Type != RAIN {
		return spovb
	}

	heal := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 16.0)
	return spovb.Heal(heal)
}

func TurnEndIceBody(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	if spovb.SelfFighters[0].Ability != "アイスボディ" {
		return spovb
	}

	if spovb.ShareField.Weather.Type != HAIL {
		return spovb
	}

	heal := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 16.0)
	return spovb.Heal(heal)
}
