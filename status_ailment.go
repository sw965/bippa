package bippa

import (
	"math/rand"
)

type StatusAilment_ string

const (
	NORMAL_POISON = StatusAilment_("どく")
	BAD_POISON = StatusAilment_("もうどく")
	SLEEP      = StatusAilment_("ねむり")
	BURN       = StatusAilment_("やけど")
	PARALYSIS = StatusAilment_("まひ")
	FREEZE = StatusAilment_("こおり")
)

func NewFreeze(percent int, random *rand.Rand) StatusAilment_ {
	if IsHit(percent, random) {
		return FREEZE
	} else {
		return StatusAilment_("")
	}
}

type ParalysisBonus float64

const (
	PARALYSIS_BONUS = ParalysisBonus(2048.0 / 4096.0)
)

type StatusAilment struct {
	Type StatusAilment_
	SleepRemainingTurn int
	BadPoisonElapsedTurn int
}
