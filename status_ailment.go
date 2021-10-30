package bippa

type StatusAilment string

const (
	BAD_POISON = StatusAilment("もうどく")
	SLEEP      = StatusAilment("ねむり")
	BURN       = StatusAilment("やけど")
)

type StatusAilmentDetail struct {
}
