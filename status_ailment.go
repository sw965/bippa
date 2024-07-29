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