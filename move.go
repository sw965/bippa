package bippa

import (
  "math/rand"
)

type StatusMove func(SelfPointOfViewBattle, *rand.Rand) SelfPointOfViewBattle

//あさのひざし
func NewMorningSun(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
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
func NewSynthesis(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
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
func NewRecover(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  maxHP := spovb.SelfFighters[0].State.MaxHP
  heal := int(float64(maxHP) * (1.0 / 2.0))
  return spovb.Heal(heal)
}

//すなあつめ
func NewShoreUp(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
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
func NewSoftBoiled(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  maxHP := spovb.SelfFighters[0].State.MaxHP
  heal := int(float64(maxHP) * (1.0 / 2.0))
  return spovb.Heal(heal)
}

//つきのひかり
func NewMoonlight(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
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
func NewSlackOff(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  maxHP := spovb.SelfFighters[0].State.MaxHP
  heal := int(float64(maxHP) * (1.0 / 2.0))
  return spovb.Heal(heal)
}

//はねやすめ
func NewRoost(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  spovb.SelfFighters[0].IsRoost = true
  maxHP := spovb.SelfFighters[0].State.MaxHP
  heal := int(float64(maxHP) * (1.0 / 2.0))
  return spovb.Heal(heal)
}

//ミルクのみ
func NewMilkDrink(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  maxHP := spovb.SelfFighters[0].State.MaxHP
  heal := int(float64(maxHP) * (1.0 / 2.0))
  return spovb.Heal(heal)
}

//ねむる
func NewRest(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
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

//キノコのほうし
func NewSpore(spovb SelfPointOfViewBattle, random *rand.Rand) SelfPointOfViewBattle {
  if spovb.OpponentFighters[0].StatusAilment.Type != "" {
    return spovb
  }

  if spovb.OpponentFighters[0].Types.In(GRASS) {
    return spovb
  }

  spovb.OpponentFighters[0].StatusAilment.Type = SLEEP
  spovb.OpponentFighters[0].StatusAilment.SleepRemainingTurn = GetSleepRemainingTurn(random)
  return spovb
}

//さいみんじゅつ
func NewHypnosis(spovb SelfPointOfViewBattle, random *rand.Rand) SelfPointOfViewBattle {
  if spovb.OpponentFighters[0].StatusAilment.Type != "" {
    return spovb
  }

  spovb.OpponentFighters[0].StatusAilment.Type = SLEEP
  spovb.OpponentFighters[0].StatusAilment.SleepRemainingTurn = GetSleepRemainingTurn(random)
  return spovb
}

//どくどく
func NewToxic(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
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

//でんじは
func NewThunderWave(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  if spovb.OpponentFighters[0].StatusAilment.Type != "" {
    return spovb
  }

  types := spovb.OpponentFighters[0].Types
  if types.In(ELECTRIC) || types.In(GROUND) {
    return spovb
  }

  spovb.OpponentFighters[0].StatusAilment.Type = PARALYSIS
  return spovb
}

//やどりぎのタネ
func NewLeechSeed(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  if spovb.OpponentFighters[0].Types.In(GRASS) {
    return spovb
  }

  spovb.OpponentFighters[0].IsLeechSeed = true
  return spovb
}

//まきびし
func NewSpikes(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  if spovb.OpponentField.SpikesCount < MAX_SPIKES {
    spovb.OpponentField.SpikesCount += 1
  }
  return spovb
}

//どくびし
func NewToxicSpikes(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  if spovb.OpponentField.ToxicSpikesCount < MAX_TOXIC_SPIKES {
    spovb.OpponentField.ToxicSpikesCount += 1
  }
  return spovb
}

//ステルスロック
func NewStealthRock(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  spovb.OpponentField.IsStealthRock = true
  return spovb
}

//あまごい
func NewRainDance(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  spovb.ShareField.Weather.Type = RAIN
  spovb.ShareField.Weather.RemainingTurn = spovb.SelfFighters[0].RainActiveTurn()
  return spovb
}

//あられ
func NewHail(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  spovb.ShareField.Weather.Type = HAIL
  spovb.ShareField.Weather.RemainingTurn = spovb.SelfFighters[0].HailActiveTurn()
  return spovb
}

//すなあらし
func NewSandstorm(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  spovb.ShareField.Weather.Type = SANDSTORM
  spovb.ShareField.Weather.RemainingTurn = spovb.SelfFighters[0].SandstormActiveTurn()
  return spovb
}

//にほんばれ
func NewSunnyDay(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  spovb.ShareField.Weather.Type = SUNNY_DAY
  spovb.ShareField.Weather.RemainingTurn = spovb.SelfFighters[0].SunnyDayActiveTurn()
  return spovb
}

//つるぎのまい
func NewSwordsDance(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  return spovb.RankFluctuation(&Rank{Atk:2})
}

//りゅうのまい
func NewDragonDance(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  return spovb.RankFluctuation(&Rank{Atk:1, Speed:1})
}

//からをやぶる
func NewShellSmash(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  return spovb.RankFluctuation(&Rank{Atk:2, Def:-1, SpAtk:2, SpDef:-1, Speed:2})
}

//てっぺき
func NewIronDefense(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  return spovb.RankFluctuation(&Rank{Def:2})
}

//めいそう
func NewCalmMind(spovb SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  return spovb.RankFluctuation(&Rank{SpAtk:1, SpDef:1})
}

var STATUS_MOVES = map[MoveName]StatusMove{
  "あさのひざし":NewMoonlight,
  "こうごうせい":NewSynthesis,
  "じこさいせい":NewRecover,
  "すなあつめ":NewShoreUp,
  "タマゴうみ":NewSoftBoiled,
  "つきのひかり":NewMoonlight,
  "なまける":NewSlackOff,
  "はねやすめ":NewRoost,
  "ミルクのみ":NewMorningSun,
  "キノコのほうし":NewSpore,
  "さいみんじゅつ":NewHypnosis,
  "どくどく":NewToxic,
  "でんじは":NewThunderWave,
  "やどりぎのタネ":NewLeechSeed,
  "まきびし":NewSpikes,
  "どくびし":NewToxicSpikes,
  "ステルスロック":NewStealthRock,
  "あまごい":NewRainDance,
  "あられ":NewHail,
  "すなあらし":NewSandstorm,
  "にほんばれ":NewSunnyDay,
  "つるぎのまい":NewSwordsDance,
  "りゅうのまい":NewDragonDance,
  "からをやぶる":NewShellSmash,
  "てっぺき":NewIronDefense,
}
