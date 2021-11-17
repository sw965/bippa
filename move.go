package bippa

type Move func(SelfPointOfViewBattle) SelfPointOfViewBattle

//あさのひざし
func NewMorningSun(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
  maxHP := float64(spovb.SelfFighters[0].State.MaxHP)
  weather_ := spovb.ShareField.Weather.Type
  badWeather_s := Weather_s{RAIN, HAIL, SANDSTORM}

  var heal int
  if badWeather_s.In(weather_) {
    heal = int(float64(maxHP) * (1.0 / 4.0))
  } else if weather_ == SUNNY_DAY {
    heal = int(float64(maxHP) * (2.0 / 3.0))
  } else {
    heal = int(float64(maxHP) * (1.0 / 2.0))
  }
  return spovb.Heal(heal)
}

//こうごうせい
func NewSynthesis(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
  maxHP := float64(spovb.SelfFighters[0].State.MaxHP)
  weather_ := spovb.ShareField.Weather.Type
  badWeather_s := Weather_s{RAIN, HAIL, SANDSTORM}

  var heal int
  if badWeather_s.In(weather_) {
    heal = int(float64(maxHP) * (1.0 / 4.0))
  } else if weather_ == SUNNY_DAY {
    heal = int(float64(maxHP) * (2.0 / 3.0))
  } else {
    heal = int(float64(maxHP) * (1.0 / 2.0))
  }
  return spovb.Heal(heal)
}

//じこさいせい
func NewRecover(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
  maxHP := spovb.SelfFighters[0].State.MaxHP
  heal := int(float64(maxHP) * (1.0 / 2.0))
  return spovb.Heal(heal)
}

//すなあつめ
func NewShoreUp(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
  maxHP := spovb.SelfFighters[0].State.MaxHP
  var heal int
  if spovb.ShareField.Weather.Type == SANDSTORM {
    heal = int(float64(maxHP) * (2.0 / 3.0))
  } else {
    heal = int(float64(maxHP) * (1.0 / 2.0))
  }
  return spovb.Heal(heal)
}

//タマゴうみ
func NewSoftBoiled(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
  maxHP := spovb.SelfFighters[0].State.MaxHP
  heal := int(float64(maxHP) * (1.0 / 2.0))
  return spovb.Heal(heal)
}

//つきのひかり
func NewMoonlight(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
  maxHP := float64(spovb.SelfFighters[0].State.MaxHP)
  weather_ := spovb.ShareField.Weather.Type
  badWeather_s := Weather_s{RAIN, HAIL, SANDSTORM}

  var heal int
  if badWeather_s.In(weather_) {
    heal = int(float64(maxHP) * (1.0 / 4.0))
  } else if weather_ == SUNNY_DAY {
    heal = int(float64(maxHP) * (2.0 / 3.0))
  } else {
    heal = int(float64(maxHP) * (1.0 / 2.0))
  }
  return spovb.Heal(heal)
}

//なまける
func NewSlackOff(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
  maxHP := spovb.SelfFighters[0].State.MaxHP
  heal := int(float64(maxHP) * (1.0 / 2.0))
  return spovb.Heal(heal)
}

//はねやすめ
func NewRoost(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
  spovb.SelfFighters[0].IsRoost = true
  maxHP := spovb.SelfFighters[0].State.MaxHP
  heal := int(float64(maxHP) * (1.0 / 2.0))
  return spovb.Heal(heal)
}

//ミルクのみ
func NewMilkDrink(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
  maxHP := spovb.SelfFighters[0].State.MaxHP
  heal := int(float64(maxHP) * (1.0 / 2.0))
  return spovb.Heal(heal)
}

//ねむる
func NewRest(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
  if spovb.SelfFighters[0].IsFullHP() {
    return spovb
  }

  maxHP := spovb.SelfFighters[0].State.MaxHP
  spovb.SelfFighters[0].State.CurrentHP = maxHP
  spovb.SelfFighters[0].StatusAilment = StatusAilment{}
  spovb.SelfFighters[0].StatusAilment.Type = SLEEP
  spovb.SelfFighters[0].StatusAilment.SleepRemainingTurn = 2
  return spovb
}

//どくどく
func NewToxic(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
  if spovb.OpponentFighters[0].StatusAilment.Type != "" {
    return spovb
  }

  isCorrosion := spovb.OpponentFighters[0].Ability == "ふしょく"
  if spovb.OpponentFighters[0].Types.In(POISON) && !isCorrosion {
    return spovb
  }

  if spovb.OpponentFighters[0].Types.In(STEEL) && !isCorrosion {
    return spovb
  }

  spovb.OpponentFighters[0].StatusAilment.Type = BAD_POISON
  return spovb
}

//やどりぎのタネ
func NewLeechSeed(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
  if spovb.OpponentFighters[0].Types.In(GRASS) {
    return spovb
  }

  spovb.OpponentFighters[0].IsLeechSeed = true
  return spovb
}

//まきびし
func NewSpikes(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
  if spovb.OpponentField.SpikesCount < MAX_SPIKES {
    spovb.OpponentField.SpikesCount += 1
  }
  return spovb
}

//どくびし
func NewToxicSpikes(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
  if spovb.OpponentField.ToxicSpikesCount < MAX_TOXIC_SPIKES {
    spovb.OpponentField.ToxicSpikesCount += 1
  }
  return spovb
}

//ステルスロック
func NewStealthRock(spovb SelfPointOfViewBattle) SelfPointOfViewBattle {
  spovb.OpponentField.IsStealthRock = true
  return spovb
}

var STATUS_MOVES = map[MoveName]Move{
  "あさのひざし":NewMoonlight,
  "こうごうせい":NewSynthesis,
  "じこさいせい":NewRecover,
  "すなあつめ":NewShoreUp,
  "タマゴうみ":NewSoftBoiled,
  "つきのひかり":NewMoonlight,
  "なまける":NewSlackOff,
  "はねやすめ":NewRoost,
  "ミルクのみ":NewMorningSun,
  "どくどく":NewToxic,
  "やどりぎのタネ":NewLeechSeed,
  "まきびし":NewSpikes,
  "どくびし":NewToxicSpikes,
  "ステルスロック":NewStealthRock,
}
