package battle

import (
	"fmt"
	omwmaps "github.com/sw965/omw/maps"
)

type Weather int

const (
	NORMAL_WEATHER Weather = iota
	RAIN
)

var STRING_TO_WEATHER = map[string]Weather{
	"":NORMAL_WEATHER,
	"雨":RAIN,
}

func StringToWeather(s string) (Weather, error) {
	if weather, ok := STRING_TO_WEATHER[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_WEATHER に含まれていない為、Weatherに変換出来ません。", s)
		return weather, fmt.Errorf(msg)
	} else {
		return weather, nil
	}
}

var WEATHER_TO_STRING = omwmaps.Invert[map[Weather]string](STRING_TO_WEATHER)

func (w Weather) ToString() string {
	return WEATHER_TO_STRING[w]
}