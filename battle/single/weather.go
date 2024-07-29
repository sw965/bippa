package single

type Weather int

const (
	NORMAL_WEATHER Weather = iota
	RAIN
)

func (w Weather) ToString() string {
	//Weatherは疑似Enumなので、wは正しい値である前提。
	switch w {
		case RAIN:
			return "雨"
		default:
			return ""
	}
}