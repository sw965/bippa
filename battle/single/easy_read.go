package single

import (
	bp "github.com/sw965/bippa"
)

type EasyReadBattle struct {
	SelfFighters bp.EasyReadPokemons
	OpponentFighters bp.EasyReadPokemons
	Turn int
	IsRealSelf bool
}