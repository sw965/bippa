package bippa

type BurnBonus float64

const (
  BURN_BONUS = BurnBonus(2048.0 / 4096.0)
  NO_BURN_BONUS = BurnBonus(4096.0 / 4096.0)
)

var BOOL_TO_BURN_BONUS = map[bool]BurnBonus{
  true:BURN_BONUS, false:NO_BURN_BONUS,
}
