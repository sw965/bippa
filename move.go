package bippa

import (
  "math/rand"
)

type Move struct {
  Effect func(SelfPointOfViewBattle, rand.Rand) SelfPointOfViewBattle
  IsSelfEffect bool
  EffectAfterMissAttack func(SelfPointOfViewBattle) SelfPointOfViewBattle
}
