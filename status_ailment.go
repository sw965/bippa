package bippa

type StatusAilment int

const (
	EMPTY_STATUS_AILMENT StatusAilment = iota
	BURN //やけど
	FREEZE //こおり
	PARALYSIS //まひ
	NORMAL_POISON //どく
	BAD_POISON //もうどく
	SLEEP //ねむり
)

func (sa StatusAilment) ToString() string {
	switch sa {
		case BURN:
			return "やけど"
		case FREEZE:
			return "こおり"
		case PARALYSIS:
			return "まひ"
		case NORMAL_POISON:
			return "どく"
		case BAD_POISON:
			return "もうどく"
		case SLEEP:
			return "ねむり"
	}
	return ""
}