package bippa

type StatusAilment int

const (
	EMPTY_STATUS_AILMENT StatusAilment = iota
	SLEEP //ねむり
	PARALYSIS //まひ
	NORMAL_POISON //どく
	BAD_POISON //もうどく
	BURN //やけど
	FREEZE //こおり
)

func (a StatusAilment) ToString() string {
	switch a {
		case SLEEP:
			return "ねむり"
		case PARALYSIS:
			return "まひ"
		case NORMAL_POISON:
			return "どく"
		case BAD_POISON:
			return "もうどく"
		case BURN:
			return "やけど"
		case FREEZE:
			return "こおり"
	}
	return ""
}