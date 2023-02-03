package bippa

import (
	"github.com/sw965/omw"
	"os"
	"strings"
)

var (
	SW965_PATH = os.Getenv("GOPATH") + "sw965/"

	DATA_PATH      = SW965_PATH + "arbok/data/"
	POKEDEX_PATH   = DATA_PATH + "pokedex/"
	MOVEDEX_PATH   = DATA_PATH + "movedex/"
	NATUREDEX_PATH = DATA_PATH + "naturedex.json"
	TYPEDEX_PATH   = DATA_PATH + "typedex.json"

	ALL_POKE_NAMES_PATH = DATA_PATH + "all_poke_names.json"
	ALL_NATURES_PATH    = DATA_PATH + "all_natures.json"
	ALL_MOVE_NAMES_PATH = DATA_PATH + "all_move_names.json"
	ALL_ITEMS_PATH      = DATA_PATH + "all_items.json"

	RATTA_PATH               = SW965_PATH + "ratta/"
	BIPARAM_INDIVIDUALS_PATH = RATTA_PATH + "individuals.json"
	BIPARAM_EFFORTS_PATH     = RATTA_PATH + "efforts.json"
)

type PokeData struct {
	NormalAbilities Abilities
	HiddenAbility   Ability
	AllAbilities    Abilities

	Gender string
	Types  Types

	BaseHP    int
	BaseAtk   int
	BaseDef   int
	BaseSpAtk int
	BaseSpDef int
	BaseSpeed int

	Height    float64
	Weight    float64
	EggGroups []string
	Category  string

	Learnset MoveNames
}

func LoadPokeData(path string) PokeData {
	y := PokeData{}
	if err := omw.LoadJson(&y, path); err != nil {
		panic(err)
	}
	return y
}

type Pokedex map[PokeName]*PokeData

var POKEDEX = func() Pokedex {
	listDir, err := omw.ListDir(POKEDEX_PATH)
	if err != nil {
		panic(err)
	}
	y := Pokedex{}
	for _, fileName := range listDir {
		fullPath := POKEDEX_PATH + fileName
		pokeName := strings.TrimRight(fileName, ".json")
		pokeData := LoadPokeData(fullPath)
		y[PokeName(pokeName)] = &pokeData
	}
	return y
}()

var ALL_POKE_NAMES = func() PokeNames {
	y := make(PokeNames, 0)
	if err := omw.LoadJson(&y, ALL_POKE_NAMES_PATH); err != nil {
		panic(err)
	}
	return y
}()

var ALL_ABILITIES = func() Abilities {
	y := make(Abilities, 0)
	for _, pokeData := range POKEDEX {
		for _, ability := range pokeData.AllAbilities {
			if !omw.Contains(y, ability) {
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
	y := MoveData{}
	omw.LoadJson(&y, path)
	return y
}

type Movedex map[MoveName]*MoveData

var MOVEDEX = func() Movedex {
	y := Movedex{}
	listDir, err := omw.ListDir(MOVEDEX_PATH)
	if err != nil {
		panic(err)
	}
	for _, fileName := range listDir {
		moveName := strings.TrimRight(fileName, ".json")
		moveData := LoadMoveData(MOVEDEX_PATH + fileName)
		y[MoveName(moveName)] = &moveData
	}
	return y
}()

var ALL_MOVE_NAMES = func() MoveNames {
	y := make(MoveNames, 0)
	if err := omw.LoadJson(&y, ALL_MOVE_NAMES_PATH); err != nil {
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
	y := Naturedex{}
	if err := omw.LoadJson(&y, NATUREDEX_PATH); err != nil {
		panic(err)
	}
	return y
}()

var ALL_NATURES = func() Natures {
	y := make(Natures, 0)
	if err := omw.LoadJson(&y, ALL_NATURES_PATH); err != nil {
		panic(err)
	}
	return y
}()

type TypeData map[Type]float64
type Typedex map[Type]TypeData

var TYPEDEX = func() Typedex {
	y := Typedex{}
	if err := omw.LoadJson(&y, TYPEDEX_PATH); err != nil {
		panic(err)
	}
	return y
}()

var ALL_ITEMS = func() Items {
	y := make(Items, 0)
	if err := omw.LoadJson(&y, ALL_ITEMS_PATH); err != nil {
		panic(err)
	}
	return y
}()

var ALL_POKE_NAMES_LENGTH = len(ALL_POKE_NAMES)
var ALL_ABILITIES_LENGTH = len(ALL_ABILITIES)
var ALL_MOVE_NAMES_LENGTH = len(ALL_MOVE_NAMES)
var ALL_NATURES_LENGTH = len(ALL_NATURES)
var ALL_ITEMS_LENGTH = len(ALL_ITEMS)
