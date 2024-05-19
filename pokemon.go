package bippa

import (
	omaps "github.com/sw965/omw/maps"
)

type PokeName int

const (
	EMPTY_POKE_NAME PokeName = iota
	BULBASAUR
	CHARMANDER
	SQUIRTLE
	GARCHOMP
)

var STRING_TO_POKE_NAME = map[string]PokeName{
	"":EMPTY_POKE_NAME,
	"フシギダネ":BULBASAUR,
	"ヒトカゲ":CHARMANDER,
	"ゼニガメ":SQUIRTLE,
	"ガブリアス":GARCHOMP,
}

var POKE_NAME_TO_STRING = omaps.Invert[map[PokeName]string](STRING_TO_POKE_NAME)

type PokeNames []PokeName
type PokeNamess []PokeNames

type Level int

const (
	DEFAULT_LEVEL Level = 50
)

type NatureBonus float64

const (
	GOOD_NATURE_BONUS = 1.1
	NEUTRAL_NATURE_BONUS = 1.0
	BAD_NATURE_BONUS = 0.9
)

type IV int

const (
	MIN_IV IV = 0
	MAX_IV IV = 31
)

type IVStats struct {
	HP IV
	Atk IV
	Def IV
	SpAtk IV
	SpDef IV
	Speed IV
}

type EV int

const (
	MIN_EV EV = 0
)

type EVStats struct {
	HP EV
	Atk EV
	Def EV
	SpAtk EV
	SpDef EV
	Speed EV
}

func CalcHP(base int, lv Level, iv IV, ev EV) int {
	ivi := int(iv)
	evi := int(ev)
	lvi := int(lv)
	return ((base*2 + ivi + evi/4) * lvi / 100) + lvi + 10
}

func CalcOtherStat(base int, lv Level, iv IV, ev EV, bonus NatureBonus) int {
	ivi := int(iv)
	evi := int(ev)
	lvi := int(lv)
	baseStat := float64((base*2 + ivi + evi/4) * lvi / 100 + 5)
	return int(baseStat * float64(bonus))
}

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

func NewPokemon(name PokeName, movesetNames MoveNames) (Pokemon, error) {
	pokeData := POKEDEX[name]
	hp := CalcHP(pokeData.BaseHP, DEFAULT_LEVEL, MAX_IV, MIN_EV)
	atk := CalcOtherStat(pokeData.BaseAtk, DEFAULT_LEVEL, MAX_IV, MIN_EV, NEUTRAL_NATURE_BONUS)
	def := CalcOtherStat(pokeData.BaseDef, DEFAULT_LEVEL, MAX_IV, MIN_EV, NEUTRAL_NATURE_BONUS)
	spAtk := CalcOtherStat(pokeData.BaseSpAtk, DEFAULT_LEVEL, MAX_IV, MIN_EV, NEUTRAL_NATURE_BONUS)
	spDef := CalcOtherStat(pokeData.BaseSpDef, DEFAULT_LEVEL, MAX_IV, MIN_EV, NEUTRAL_NATURE_BONUS)
	speed := CalcOtherStat(pokeData.BaseSpeed, DEFAULT_LEVEL, MAX_IV, MIN_EV, NEUTRAL_NATURE_BONUS)

	moveset, err := NewMoveset(name, movesetNames)
	if err != nil {
		return Pokemon{}, err
	}

	return Pokemon{
		Name:name,
		Level:DEFAULT_LEVEL,
		MaxHP:hp,
		CurrentHP:hp,
		Atk:atk,
		Def:def,
		SpAtk:spAtk,
		SpDef:spDef,
		Speed:speed,
		Moveset:moveset,
	}, err
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

func (p Pokemon) Clone() Pokemon {
	p.Moveset = p.Moveset.Clone()
	return p
}

func (p *Pokemon) HPPercentage() float64 {
	return float64(p.CurrentHP) / float64(p.MaxHP)
}

func (p *Pokemon) IsFaint() bool {
	return p.CurrentHP <= 0
}

func NewTemplateBulbasaur() Pokemon {
	pokemon, err := NewPokemon(BULBASAUR, MoveNames{TACKLE, VINE_WHIP})
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewTemplateCharmander() Pokemon {
	pokemon, err := NewPokemon(CHARMANDER, MoveNames{EMBER})
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewTemplateSquirtle() Pokemon {
	pokemon, err := NewPokemon(SQUIRTLE, MoveNames{WATER_GUN})
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewTemplateGarchomp() Pokemon {
	pokemon, err := NewPokemon(GARCHOMP, MoveNames{STONE_EDGE})
	if err != nil {
		panic(err)
	}
	return pokemon
}