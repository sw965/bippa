package bippa

type Weather_ string

const (
	RAIN = Weather_("あめ")
	HAIL = Weather_("あられ")
	SANDSTORM = Weather_("すなあらし")
	SUNNY_DAY = Weather_("にほんばれ")
)

type Weather_s []Weather_

func (weather_s Weather_s) In(weather_ Weather_) bool {
	for _, iWeather_ := range weather_s {
		if iWeather_ == weather_ {
			return true
		}
	}
	return false
}

type Weather struct {
	Type Weather_
	RemainingTurn int
}

type WeatherBonus float64

const (
  BAD_WEATHER_BONUS = 2048.0 / 4096.0
  GOOD_WEATHER_BONUS = 6144.0 / 4096.0
)

func NewWeatherBonus(spovb *SelfPointOfViewBattle, moveType Type) WeatherBonus {
  weatherType := spovb.ShareField.Weather.Type
  if weatherType == "" {
    return 1.0
  }

  switch weatherType {
    case RAIN:
      if moveType == FIRE {
        return BAD_WEATHER_BONUS
      } else if moveType == WATER {
        return GOOD_WEATHER_BONUS
      }
    case SUNNY_DAY:
      if moveType == WATER {
        return BAD_WEATHER_BONUS
      } else if moveType == FIRE {
        return GOOD_WEATHER_BONUS
      }
  }
  return 1.0
}
