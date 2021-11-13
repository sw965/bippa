package bippa

import (
  "math/rand"
)

type Move func(SelfPointOfViewBattle) SelfPointOfViewBattle

//にほんばれ
func NewMorningSun(selfPointOfViewBattle SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  maxHP := float64(selfPointOfViewBattle.SelfFighters[0].State.MaxHP)
  weather_ := selfPointOfViewBattle.ShareField.Weather.Type
  badWeather_s := Weather_s{RAIN, HAIL, SANDSTORM}

  var heal int
  if badWeather_s.In(weather_) {
    heal = int(float64(maxHP) * (1.0 / 4.0))
  } else if weather_ == SUNNY_DAY {
    heal = int(float64(maxHP) * (2.0 / 3.0))
  } else {
    heal = int(float64(maxHP) * (1.0 / 2.0))
  }
  return selfPointOfViewBattle.Heal(heal)
}

//こうごうせい
func NewSynthesis(selfPointOfViewBattle SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  maxHP := float64(selfPointOfViewBattle.SelfFighters[0].State.MaxHP)
  weather_ := selfPointOfViewBattle.ShareField.Weather.Type
  badWeather_s := Weather_s{RAIN, HAIL, SANDSTORM}

  var heal int
  if badWeather_s.In(weather_) {
    heal = int(float64(maxHP) * (1.0 / 4.0))
  } else if weather_ == SUNNY_DAY {
    heal = int(float64(maxHP) * (2.0 / 3.0))
  } else {
    heal = int(float64(maxHP) * (1.0 / 2.0))
  }
  return selfPointOfViewBattle.Heal(heal)
}

//じこさいせい
func NewRecover(selfPointOfViewBattle SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  maxHP := selfPointOfViewBattle.SelfFighters[0].State.MaxHP
  heal := int(float64(maxHP) * (1.0 / 2.0))
  return selfPointOfViewBattle.Heal(heal)
}

//すなあつめ
func NewShoreUp(selfPointOfViewBattle SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  maxHP := selfPointOfViewBattle.SelfFighters[0].State.MaxHP
  var heal int
  if selfPointOfViewBattle.ShareField.Weather.Type == SANDSTORM {
    heal = int(float64(maxHP) * (2.0 / 3.0))
  } else {
    heal = int(float64(maxHP) * (1.0 / 2.0))
  }
  return selfPointOfViewBattle.Heal(heal)
}

//タマゴうみ
func NewSoftBoiled(selfPointOfViewBattle SelfPointOfViewBattle, _*rand.Rand) SelfPointOfViewBattle {
  maxHP := selfPointOfViewBattle.SelfFighters[0].State.MaxHP
  heal := int(float64(maxHP) * (1.0 / 2.0))
  return selfPointOfViewBattle.Heal(heal)
}

//つきのひかり
func NewMoonlight(selfPointOfViewBattle SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  maxHP := float64(selfPointOfViewBattle.SelfFighters[0].State.MaxHP)
  weather_ := selfPointOfViewBattle.ShareField.Weather.Type
  badWeather_s := Weather_s{RAIN, HAIL, SANDSTORM}

  var heal int
  if badWeather_s.In(weather_) {
    heal = int(float64(maxHP) * (1.0 / 4.0))
  } else if weather_ == SUNNY_DAY {
    heal = int(float64(maxHP) * (2.0 / 3.0))
  } else {
    heal = int(float64(maxHP) * (1.0 / 2.0))
  }
  return selfPointOfViewBattle.Heal(heal)
}

//なまける
func NewSlackOff(selfPointOfViewBattle SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  maxHP := selfPointOfViewBattle.SelfFighters[0].State.MaxHP
  heal := int(float64(maxHP) * (1.0 / 2.0))
  return selfPointOfViewBattle.Heal(heal)
}

//はねやすめ
func NewRoost(selfPointOfViewBattle SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  selfPointOfViewBattle.IsRoost = true
  maxHP := selfPointOfViewBattle.SelfFighters[0].State.MaxHP
  heal := int(float64(maxHP) * (1.0 / 2.0))
  return selfPointOfViewBattle.Heal(heal)
}

//ミルクのみ
func NewMilkDrink(selfPointOfViewBattle SelfPointOfViewBattle, _ *rand.Rand) SelfPointOfViewBattle {
  maxHP := selfPointOfViewBattle.SelfFighters[0].State.MaxHP
  heal := int(float64(maxHP) * (1.0 / 2.0))
  return selfPointOfViewBattle.Heal(heal)
}
