package bippa

import (
  "fmt"
)

//https://latest.pokewiki.net/%E3%83%80%E3%83%A1%E3%83%BC%E3%82%B8%E8%A8%88%E7%AE%97%E5%BC%8F
type PowerBonus int

const (
  INIT_POWER_BONUS = PowerBonus(4096)
)

type MovePowerAlpha int

func NewMovePowerAlpha(moveName MoveName) MovePowerAlpha {
  movePower := MOVEDEX[moveName].Power

  movePowerAlpha := MovePowerAlpha(movePower)

	if moveName == "アシストパワー" {
		movePowerAlpha = movePowerAlpha.AddStoredPower()
	}
  return movePowerAlpha
}

//アシストパワー
func (movePowerAlpha MovePowerAlpha) AddStoredPower() MovePowerAlpha {
  return movePowerAlpha + MovePowerAlpha(20 * spovb.SelfFighters[0].RankState.TotalRise())
}

type FinalPower int

func NewFinalPower(spovb *SelfPointOfViewBattle, moveName MoveName) (FinalMovePower, error) {
	if MOVEDEX[moveName].Category == STATUS {
		return 0, fmt.Errorf("変化技以外でなければならない")
	}

	powerBonus := INIT_POWER_BONUS
	movePowerAlpha := NewMovePowerAlpha(moveName)

	result := RoundingZeroPointFiveOver(float64(movePower) * float64(powerBonus) / 4096.0)
	if result < 1 {
		return 1, nil
	}
	return result, nil
}
