package bippa

import (
  "fmt"
)

//https://latest.pokewiki.net/%E3%83%80%E3%83%A1%E3%83%BC%E3%82%B8%E8%A8%88%E7%AE%97%E5%BC%8F
type PowerBonus int

const (
  INIT_POWER_BONUS = PowerBonus(4096)
)

type AlphaMovePower int

func NewAlphaMovePower(spovb *SelfPointOfViewBattle, moveName MoveName) AlphaMovePower {
  movePower := MOVEDEX[moveName].Power

  alphaMovePower := AlphaMovePower(movePower)

	if moveName == "アシストパワー" {
		alphaMovePower = alphaMovePower.AddStoredPower(&spovb.SelfFighters[0].Rank)
	}
  return alphaMovePower
}

//アシストパワー
func (alphaMovePower AlphaMovePower) AddStoredPower(rank *Rank) AlphaMovePower {
  return alphaMovePower + AlphaMovePower(20 * rank.TotalRise())
}

type FinalPower int

func NewFinalPower(spovb *SelfPointOfViewBattle, moveName MoveName) (FinalPower, error) {
	if MOVEDEX[moveName].Category == STATUS {
		return 0, fmt.Errorf("変化技以外でなければならない")
	}

	powerBonus := INIT_POWER_BONUS
	alphaMovePower := NewAlphaMovePower(spovb, moveName)

	result := RoundingZeroPointFiveOver(float64(alphaMovePower) * float64(powerBonus) / 4096.0)
	if result < 1 {
		return 1, nil
	}
	return FinalPower(result), nil
}
