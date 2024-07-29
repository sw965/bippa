package single

import (
	bp "github.com/sw965/bippa"
)

type EasyReadBattle struct {
	SelfLeadPokemons bp.EasyReadPokemons
	SelfBenchPokemons bp.EasyReadPokemons

	OpponentLeadPokemons bp.EasyReadPokemons
	OpponentBenchPokemons bp.EasyReadPokemons

	Turn int
	IsPlayer1 bool

	Weather string
	RemainingTurnWeather int
}