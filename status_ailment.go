package bippa

type StatusAilment_ string

const (
	NORMAL_POISON = StatusAilment_("どく")
	BAD_POISON = StatusAilment_("もうどく")
	SLEEP      = StatusAilment_("ねむり")
	BURN       = StatusAilment_("やけど")
)

type StatusAilment struct {
	Type StatusAilment_
	SleepRemainingTurn int
	BadPoisonElapsedTurn int
}
