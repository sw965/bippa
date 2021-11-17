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

type SelfField struct {
	SpikesCount int
	ToxicSpikesCount int
	IsStealthRock bool
}

type ShareField struct {
	Weather Weather
}
