package bippa

type SameTypeAttackBonus float64

const (
  SAME_TYPE_ATTACK_BONUS = SameTypeAttackBonus(6144.0 / 4096.0)
  NO_SAME_TYPE_ATTACK_BONUS = SameTypeAttackBonus(4096.0 / 4096.0)
)

var BOOL_TO_SAME_TYPE_ATTACK_BONUS = map[bool]SameTypeAttackBonus{
  true:SAME_TYPE_ATTACK_BONUS, false:NO_SAME_TYPE_ATTACK_BONUS,
}
