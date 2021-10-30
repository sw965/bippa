package bippa

import (
	"math/rand"
)

type MoveName string

const (
	STRUGGLE = MoveName("わるあがき")
)

func (moveName MoveName) IsValid() bool {
	_, ok := MOVEDEX[MoveName(moveName)]
	return ok
}

type MoveNames []MoveName

func (moveNames MoveNames) In(moveName MoveName) bool {
	for _, iMoveName := range moveNames {
		if iMoveName == moveName {
			return true
		}
	}
	return false
}

func (moveNames MoveNames) Copy() MoveNames {
	result := make(MoveNames, len(moveNames))
	for i, moveName := range moveNames {
		result[i] = moveName
	}
	return result
}

func (moveNames MoveNames) Shuffle(random *rand.Rand) MoveNames {
	result := moveNames.Copy()
	for i := len(result) - 1; i >= 0; i-- {
		j := random.Intn(i + 1)
		result[i], result[j] = result[j], result[i]
	}
	return result
}

func (moveNames MoveNames) RandomChoice(random *rand.Rand) MoveName {
	index := random.Intn(len(moveNames))
	return moveNames[index]
}
