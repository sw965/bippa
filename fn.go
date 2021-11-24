package bippa

import (
	"math/rand"
)

func CanLearn(pokeName PokeName, moveName MoveName) bool {
	pokeData, _ := POKEDEX[pokeName]
	for _, iMoveName := range pokeData.Learnset {
		if iMoveName == moveName {
			return true
		}
	}
	return false
}

func AfterTheDecimalPoint(x float64) float64 {
	return float64(x) - float64(int(x))
}

func RoundingZeroPointFiveOrMore(x float64) int {
	afterTheDecimalPoint := AfterTheDecimalPoint(x)
	if afterTheDecimalPoint >= 0.5 {
		return int(x + 1)
	}
	return int(x)
}

func RoundingZeroPointFiveOver(x float64) int {
	afterTheDecimalPoint := AfterTheDecimalPoint(x)
	if afterTheDecimalPoint > 0.5 {
		return int(x + 1)
	}
	return int(x)
}

func IsHit(percent int, random *rand.Rand) bool {
	return random.Intn(100) < percent
}

func IsCritical(random *rand.Rand) bool {
	return random.Intn(24) == 0
}

func GetSleepRemainingTurn(random *rand.Rand) int {
	return random.Intn(3) + 2
}
