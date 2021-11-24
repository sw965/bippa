package bippa

//https://wiki.xn--rckteqa2e.com/wiki/%E3%81%99%E3%81%B0%E3%82%84%E3%81%95#.E8.A9.B3.E7.B4.B0.E3.81.AA.E4.BB.95.E6.A7.98
type SpeedBonus int

const (
  INIT_SPEED_BONUS = SpeedBonus(4096)
)

func NewSpeedBonus(spovb *SelfPointOfViewBattle) SpeedBonus {
	result := INIT_SPEED_BONUS

  ability := spovb.SelfFighters[0].Ability
  weatherType := spovb.ShareField.Weather.Type

  switch ability {
    case "すいすい":
      if weatherType == RAIN {
        result = result.MulSwiftSwim()
      }
    case "ゆきかき":
      if weatherType == HAIL {
        result = result.MulSlushRush()
      }
    case "すなかき":
      if weatherType == SANDSTORM {
        result = result.MulSandRush()
      }
    case "ようりょくそ":
      if weatherType == SUNNY_DAY {
        result = result.MulChlorophyll()
      }
  }

	if spovb.SelfFighters[0].Item == "こだわりスカーフ" {
		result = result.MulChoiceScarf()
	}
	return result
}

//こだわりスカーフ
func (speedBonus SpeedBonus) MulChoiceScarf() SpeedBonus {
  result := RoundingZeroPointFiveOrMore(float64(speedBonus) * 6144.0 / 4096.0)
  return SpeedBonus(result)
}

//すいすい
func (speedBonus SpeedBonus) MulSwiftSwim() SpeedBonus {
  result := RoundingZeroPointFiveOrMore(float64(speedBonus) * 8192.0 / 4096.0)
  return SpeedBonus(result)
}

//ゆきかき
func (speedBonus SpeedBonus) MulSlushRush() SpeedBonus {
  result := RoundingZeroPointFiveOrMore(float64(speedBonus) * 8192.0 / 4096.0)
  return SpeedBonus(result)
}

//すなかき
func (speedBonus SpeedBonus) MulSandRush() SpeedBonus {
  result := RoundingZeroPointFiveOrMore(float64(speedBonus) * 8192.0 / 4096.0)
  return SpeedBonus(result)
}

//ようりょくそ
func (speedBonus SpeedBonus) MulChlorophyll() SpeedBonus {
  result := RoundingZeroPointFiveOrMore(float64(speedBonus) * 8192.0 / 4096.0)
  return SpeedBonus(result)
}

type FinalSpeed float64

func NewFinalSpeed(spovb *SelfPointOfViewBattle) FinalSpeed {
	speed := spovb.SelfFighters[0].State.Speed
	rankBonus := RANK__TO_RANK_BONUS[spovb.SelfFighters[0].Rank.Speed]
	speedBonus := NewSpeedBonus(spovb)

	result := RoundingZeroPointFiveOrMore(float64(speed) * float64(rankBonus))
	result = RoundingZeroPointFiveOver(float64(result) * float64(speedBonus) / 4096.0)
  result = int(float64(result) * float64(PARALYSIS_BONUS))
  return FinalSpeed(result)
}
