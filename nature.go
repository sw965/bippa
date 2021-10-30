package bippa

import (
	"math/rand"
)

type Nature string

func (nature Nature) IsValid() bool {
	_, ok := NATUREDEX[nature]
	return ok
}

type Natures []Nature

func (natures Natures) RandomChoice(random *rand.Rand) Nature {
	index := random.Intn(len(natures))
	return natures[index]
}
