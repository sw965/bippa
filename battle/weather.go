package battle

type Weather int

const (
	NORMAL_WEATHER Weather = iota
	RAIN
)

func (w Weather) ToString() string {
	switch w {
		case RAIN:
			return "雨"
		default:
			return "無天候"
	}
}