package bippa

import (
	"math/rand"
)

type Ability string

func (ability Ability) IsValid(pokeName PokeName) bool {
	for _, iAbility := range POKEDEX[pokeName].AllAbilities {
		if iAbility == ability {
			return true
		}
	}
	return false
}

type Abilities []Ability

func (abilities Abilities) RandomChoice(random *rand.Rand) Ability {
	index := random.Intn(len(abilities))
	return abilities[index]
}
