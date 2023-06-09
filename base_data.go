package bippa

import (
	"github.com/sw965/omw"
	"golang.org/x/exp/slices"
	"strings"
)

type BaseState int

type PokeData struct {
	NormalAbilities Abilities
	HiddenAbility   Ability
	AllAbilities    Abilities

	Gender string
	Types  Types

	BaseHP    BaseState
	BaseAtk   BaseState
	BaseDef   BaseState
	BaseSpAtk BaseState
	BaseSpDef BaseState
	BaseSpeed BaseState

	Height    float64
	Weight    float64
	EggGroups []string
	Category  string

	Learnset MoveNames
}

func LoadPokeData(path string) PokeData {
	y, err := omw.LoadJson[PokeData](path)
	if err != nil {
		panic(err)
	}
	return y
}

type Pokedex map[PokeName]*PokeData

var POKEDEX = func() Pokedex {
	names, err := omw.DirNames(POKEDEX_PATH)
	if err != nil {
		panic(err)
	}
	y := Pokedex{}
	for _, name := range names {
		full := POKEDEX_PATH + name
		pokeName := strings.TrimRight(name, ".json")
		pokeData := LoadPokeData(full)
		y[PokeName(pokeName)] = &pokeData
	}
	return y
}()

var ALL_POKE_NAMES = func() PokeNames {
	y, err := omw.LoadJson[PokeNames](ALL_POKE_NAMES_PATH)
	if err != nil {
		panic(err)
	}
	return y
}()

var ALL_ABILITIES = func() Abilities {
	y := make(Abilities, 0)
	for _, pokeData := range POKEDEX {
		for _, ability := range pokeData.AllAbilities {
			if !slices.Contains(y, ability) {
				y = append(y, ability)
			}
		}
	}
	return y
}()

type MoveData struct {
	Type     Type
	Category string
	Power    int
	Accuracy int
	BasePP   int
	Target   string

	Contact    string
	Protect    string
	MagicCoat  string
	Snatch     string
	MirrorMove string
	Substitute string

	GigantamaxMove  string
	GigantamaxPower int

	PriorityRank int
	CriticalRank CriticalRank

	MinAttackNum int
	MaxAttackNum int
}

func LoadMoveData(path string) MoveData {
	y, err := omw.LoadJson[MoveData](path)
	if err != nil {
		panic(err)
	}
	return y
}

type Movedex map[MoveName]*MoveData

var MOVEDEX = func() Movedex {
	y := Movedex{}
	names, err := omw.DirNames(MOVEDEX_PATH)
	if err != nil {
		panic(err)
	}
	for _, name := range names {
		moveName := strings.TrimRight(name, ".json")
		moveData := LoadMoveData(MOVEDEX_PATH + name)
		y[MoveName(moveName)] = &moveData
	}
	return y
}()

var ALL_MOVE_NAMES = func() MoveNames {
	y, err := omw.LoadJson[MoveNames](ALL_MOVE_NAMES_PATH)
	if err != nil {
		panic(err)
	}
	return y
}()

type NatureData struct {
	AtkBonus   NatureBonus
	DefBonus   NatureBonus
	SpAtkBonus NatureBonus
	SpDefBonus NatureBonus
	SpeedBonus NatureBonus
}

type Naturedex map[Nature]*NatureData

var NATUREDEX = func() Naturedex {
	y, err := omw.LoadJson[Naturedex](NATUREDEX_PATH)
	if err != nil {
		panic(err)
	}
	return y
}()

var ALL_NATURES = func() Natures {
	y, err := omw.LoadJson[Natures](ALL_NATURES_PATH)
	if err != nil {
		panic(err)
	}
	return y
}()

type TypeData map[Type]float64
type Typedex map[Type]TypeData

var TYPEDEX = func() Typedex {
	y, err := omw.LoadJson[Typedex](TYPEDEX_PATH)
	if err != nil {
		panic(err)
	}
	return y
}()

var ALL_ITEMS = func() Items {
	y, err := omw.LoadJson[Items](ALL_ITEMS_PATH)
	if err != nil {
		panic(err)
	}
	return y
}()
