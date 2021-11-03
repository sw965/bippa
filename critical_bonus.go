package bippa

type CriticalBonus float64

var (
  CRITICAL_BONUS = CriticalBonus(6144.0 / 4096.0)
  NO_CRITICAL_BONUS = CriticalBonus(4096.0 / 4096.0)
)

var BOOL_TO_CRITICAL_BONUS = map[bool]CriticalBonus{true:CRITICAL_BONUS, false:NO_CRITICAL_BONUS}
