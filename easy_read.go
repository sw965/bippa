package bippa

type EasyReadMoveset map[string]*PowerPoint

type EasyReadPokemon struct {
	Name string

	Level Level
	Nature string
	Moveset EasyReadMoveset
	UnassignedLearnMoveCount int

	IVStat IVStat
	EVStat EVStat

	MaxHP int
	CurrentHP int
	Atk int
	Def int
	SpAtk int
	SpDef int
	Speed int
}

type EasyReadPokemons []EasyReadPokemon