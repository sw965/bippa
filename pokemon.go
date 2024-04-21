package bippa

import (
	"github.com/sw965/omw"
)

type PokeName int

const (
	EMPTY_POKE_NAME PokeName = iota
	BULBASAUR
)

var STRING_TO_POKE_NAME = map[string]PokeName{
	"":EMPTY_POKE_NAME,
	"フシギダネ":BULBASAUR,
}

var POKE_NAME_TO_STRING = omw.InvertMap[map[PokeName]string](STRING_TO_POKE_NAME)

type PokeNames []PokeName
type PokeNamess []PokeNames

type Level int

const (
	DEFAULT_LEVEL Level = 50
)

type Pokemon struct {
	Name PokeName
	Level Level
	MaxHP int
	CurrentHP int
	Atk int
	Def int
	SpAtk int
	SpDef int
	Speed int
	Moveset Moveset
}

func (p *Pokemon) Equal(other *Pokemon) bool {
	if p.Name != other.Name {
		return false
	}

	if p.Level != other.Level {
		return false
	}

	if p.MaxHP != other.MaxHP {
		return false
	}

	if p.CurrentHP != other.CurrentHP {
		return false
	}

	if p.Atk != other.Atk {
		return false
	}

	if p.Def != other.Def {
		return false
	}

	if p.SpAtk != other.SpAtk {
		return false
	}

	if p.SpDef != other.SpDef {
		return false
	}

	if p.Speed != other.Speed {
		return false
	}
	return p.Moveset.Equal(other.Moveset)
}