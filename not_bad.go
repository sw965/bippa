package bippa

import (
  "fmt"
)

type NotBad struct {
  Policy Policy
  Statistical []float64
}

func NewNotBad(battle *Battle) error {
  reverseBattle := battle.Reverse()

  for moveName, _ := range reverseBattle.P1Fighters[0].Moveset {
    if MOVEDEX[moveName].Category == STATUS {
      continue
    }

    fmt.Println(moveName)
    adpd, err := reverseBattle.AttackDamageProbabilityDistribution(moveName)
    if err != nil {
      panic(err)
    }

    fmt.Println(battle.AttackDamageMaxHPRatioExpected(adpd))
    fmt.Println(battle.AttackDamageCurrentHPRatioExpected(adpd))

  }
  return nil
}
