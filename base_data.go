package bippa

import (
	"encoding/json"
	"github.com/sw965/omw"
	"io/ioutil"
	"os"
	"strings"
)

var (
	DATA_PATH      = os.Getenv("GOPATH") + "/raticate/data/"
	POKEDEX_PATH   = DATA_PATH + "pokedex/"
	MOVEDEX_PATH   = DATA_PATH + "movedex/"
	NATUREDEX_PATH = DATA_PATH + "naturedex.json"
	TYPEDEX_PATH   = DATA_PATH + "typedex.json"

	ALL_POKE_NAMES_PATH = DATA_PATH + "all_poke_names.txt"
	ALL_NATURES_PATH    = DATA_PATH + "all_natures.txt"
	ALL_MOVE_NAMES_PATH = DATA_PATH + "all_move_names.txt"
	ALL_ITEMS_PATH      = DATA_PATH + "all_items.txt"
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

func LoadPokeData(filePath string) PokeData {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	result := PokeData{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		panic(err)
	}
	return result
}

type Pokedex map[PokeName]*PokeData

var POKEDEX = func() Pokedex {
	listDir, err := omw.ListDir(POKEDEX_PATH)
	if err != nil {
		panic(err)
	}

	result := Pokedex{}
	for _, fileName := range listDir {
		fullPath := POKEDEX_PATH + fileName
		pokeName := strings.TrimRight(fileName, ".json")
		pokeData := LoadPokeData(fullPath)
		result[PokeName(pokeName)] = &pokeData
	}
	return result
}()

var ALL_POKE_NAMES = func() PokeNames {
	allPokeNames, err := omw.ReadTextLines(ALL_POKE_NAMES_PATH)
	if err != nil {
		panic(err)
	}

	result := make(PokeNames, len(allPokeNames))
	for i, pokeName := range allPokeNames {
		result[i] = PokeName(pokeName)
	}
	return result
}()

var ALL_ABILITIES = func() Abilities {
	result := make(Abilities, 0)
	for _, pokeData := range POKEDEX {
		for _, ability := range pokeData.AllAbilities {
			if !result.In(ability) {
				result = append(result, ability)
			}
		}
	}
	return result
}

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

func LoadMoveData(fullPath string) MoveData {
	bytes, err := ioutil.ReadFile(fullPath)
	if err != nil {
		panic(err)
	}

	result := MoveData{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		panic(err)
	}
	return result
}

type Movedex map[MoveName]*MoveData

var MOVEDEX = func() Movedex {
	result := Movedex{}
	listDir, err := omw.ListDir(MOVEDEX_PATH)
	if err != nil {
		panic(err)
	}

	for _, fileName := range listDir {
		moveName := strings.TrimRight(fileName, ".json")
		moveData := LoadMoveData(MOVEDEX_PATH + fileName)
		result[MoveName(moveName)] = &moveData
	}
	return result
}()

var ALL_MOVE_NAMES = func() MoveNames {
	allMoveNames, err := omw.ReadTextLines(ALL_MOVE_NAMES_PATH)
	if err != nil {
		panic(err)
	}

	result := make(MoveNames, len(allMoveNames))
	for i, moveName := range allMoveNames {
		result[i] = MoveName(moveName)
	}
	return result
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
	bytes, err := ioutil.ReadFile(NATUREDEX_PATH)

	if err != nil {
		panic(err)
	}

	result := Naturedex{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		panic(err)
	}
	return result
}()

var ALL_NATURES = func() Natures {
	allNatures, err := omw.ReadTextLines(ALL_NATURES_PATH)
	if err != nil {
		panic(err)
	}

	result := make(Natures, len(allNatures))
	for i, nature := range allNatures {
		result[i] = Nature(nature)
	}
	return result
}()

type TypeData map[Type]float64
type Typedex map[Type]TypeData

var TYPEDEX = func() Typedex {
	bytes, err := ioutil.ReadFile(TYPEDEX_PATH)

	if err != nil {
		panic(err)
	}

	result := Typedex{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		panic(err)
	}
	return result
}()

const (
	NO_ITEM = Item("??????")
)

var ALL_ITEMS = func() Items {
	allItems, err := omw.ReadTextLines(ALL_ITEMS_PATH)
	if err != nil {
		panic(err)
	}

	result := make(Items, 0, len(allItems))
	for _, item := range allItems {
		result = append(result, Item(item))
	}
	return result
}()
