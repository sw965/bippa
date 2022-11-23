package bippa

import (
	"github.com/sw965/omw"
	"math/rand"
)

type MoveName string

const (
	EMPTY_MOVE_NAME = MoveName("なし")
	STRUGGLE        = MoveName("わるあがき")
)

type MoveNames []MoveName

func (moveNames MoveNames) Count(moveName MoveName) int {
	result := 0
	for _, iMoveName := range moveNames {
		if iMoveName == moveName {
			result += 1
		}
	}
	return result
}

func (moveNames MoveNames) In(moveName MoveName) bool {
	for _, iMoveName := range moveNames {
		if iMoveName == moveName {
			return true
		}
	}
	return false
}

func (moveNames MoveNames) InAll(moveName ...MoveName) bool {
	for _, iMoveName := range moveName {
		if moveNames.In(iMoveName) {
			return true
		}
	}
	return false
}

func (moveNames MoveNames) Sort() MoveNames {
	result := make(MoveNames, 0, len(moveNames))
	for _, moveName := range ALL_MOVE_NAMES {
		if moveNames.In(moveName) {
			result = append(result, moveName)
		}
	}

	for i := 0; i < moveNames.Count(EMPTY_MOVE_NAME); i++ {
		result = append(result, EMPTY_MOVE_NAME)
	}
	return result
}

func (moveNames MoveNames) Index(moveName MoveName) int {
	for i, iMoveName := range moveNames {
		if iMoveName == moveName {
			return i
		}
	}
	return -1
}

type MoveNameWithTier map[MoveName]Tier

func (moveNameWithTier MoveNameWithTier) Keys() MoveNames {
	result := make(MoveNames, 0, len(moveNameWithTier))
	for k, _ := range moveNameWithTier {
		result = append(result, k)
	}
	return result
}

func (moveNameWithTier MoveNameWithTier) Items() (MoveNames, Tiers) {
	keys := moveNameWithTier.Keys()
	values := make(Tiers, len(moveNameWithTier))
	for i, key := range keys {
		values[i] = moveNameWithTier[key]
	}
	return keys, values
}

func (moveNameWithTier MoveNameWithTier) MoveNameRandomChoiceWithTierWeight(random *rand.Rand) MoveName {
	keys, values := moveNameWithTier.Items()
	weight := make([]float64, len(values))
	for i, tier := range values {
		weight[i] = TIER_WITH_SELECT_PERCENT[tier]
	}
	index := omw.RandomIntWithWeight(weight, random)
	return keys[index]
}
