package bippa

import (
	omwslices "github.com/sw965/omw/slices"
	omwjson "github.com/sw965/omw/json"
	omwos "github.com/sw965/omw/os"
	omwmath "github.com/sw965/omw/math"
	"golang.org/x/exp/slices"
	"strings"
	"fmt"
)

type BaseState int

type pokeData struct {
	BaseHP    BaseState
	BaseAtk   BaseState
	BaseDef   BaseState
	BaseSpAtk BaseState
	BaseSpDef BaseState
	BaseSpeed BaseState

	Weight    float64

	Types  []string
	Genders []string
	Abilities   Abilities

	Learnset MoveNames	
}

type PokeData struct {
	BaseHP    BaseState
	BaseAtk   BaseState
	BaseDef   BaseState
	BaseSpAtk BaseState
	BaseSpDef BaseState
	BaseSpeed BaseState

	Weight    float64

	Types  Types
	Genders Genders
	Abilities    Abilities

	Learnset MoveNames
}

func LoadPokeData(path string) PokeData {
	d, err := omwjson.Load[pokeData](path)
	if err != nil {
		panic(err)
	}

	types, err := NewTypes(d.Types)
	if err != nil {
		panic(err)
	}

	genders, err := NewGenders(d.Genders)
	if err != nil {
		panic(err)
	}

	y := PokeData{
		BaseHP:d.BaseHP,
		BaseAtk:d.BaseAtk,
		BaseDef:d.BaseDef,
		BaseSpAtk:d.BaseSpAtk,
		BaseSpeed:d.BaseSpeed,
		Weight:d.Weight,
		Types:types,
		Genders:genders,
		Abilities:d.Abilities,
		Learnset:d.Learnset,
	}
	return y
}

type Pokedex map[PokeName]*PokeData

var POKEDEX = func() Pokedex {
	entries, err := omwos.NewDirEntries(POKEDEX_PATH)
	if err != nil {
		panic(err)
	}

	y := Pokedex{}
	for _, name := range entries.Names() {
		if name == "テンプレート.json" {
			continue
		}
		full := POKEDEX_PATH + name
		pokeName := strings.TrimRight(name, ".json")
		pokeData := LoadPokeData(full)
		y[PokeName(pokeName)] = &pokeData
	}
	return y
}()

var ALL_POKE_NAMES = func() PokeNames {
	y, err := omwjson.Load[PokeNames](ALL_POKE_NAMES_PATH)
	if err != nil {
		panic(err)
	}
	return y
}()

var ALL_ABILITIES = func() Abilities {
	y := make(Abilities, len(ALL_POKE_NAMES) * 3)
	for _, pokeName := range ALL_POKE_NAMES {
		pokeData := POKEDEX[pokeName]
		for _, ability := range pokeData.Abilities {
			if !slices.Contains(y, ability) {
				y = append(y, ability)
			}
		}
	}
	return y
}()

type moveData struct {
	Type     string
	Category string
	Power    int
	Accuracy int
	BasePP   int
	Target   string

	IsContact    bool
	PriorityRank int
	CriticalRank CriticalRank
	FlinchPercent int

	MinAttackNum int
	MaxAttackNum int
}

type MoveData struct {
	Type     Type
	Category MoveCategory
	Power    int
	Accuracy int
	BasePP   int
	Target   Target

	IsContact    bool
	PriorityRank int
	CriticalRank CriticalRank
	FlinchPercent int

	MinAttackNum int
	MaxAttackNum int	
}

func LoadMoveData(path string) MoveData {
	d, err := omwjson.Load[moveData](path)
	if err != nil {
		panic(err)
	}

	type_, err := NewType(d.Type)
	if err != nil {
		panic(err)
	}

	c, err := NewMoveCategory(d.Category)
	if err != nil {
		panic(err)
	}

	target, err := NewTarget(d.Target)
	if err != nil {
		panic(err)
	}

	minAttackNum := omwmath.Max(d.MinAttackNum, 1)
	maxAttackNum := omwmath.Max(d.MaxAttackNum, 1)

	y := MoveData{
		Type:type_,
		Category:c,
		Power:d.Power,
		Accuracy:d.Accuracy,
		BasePP:d.BasePP,
		Target:target,
		IsContact:d.IsContact,
		PriorityRank:d.PriorityRank,
		CriticalRank:d.CriticalRank,
		FlinchPercent:d.FlinchPercent,
		MinAttackNum:minAttackNum,
		MaxAttackNum:maxAttackNum,
	}
	return y
}

type Movedex map[MoveName]*MoveData

var MOVEDEX = func() Movedex {
	y := Movedex{}
	entries, err := omwos.NewDirEntries(MOVEDEX_PATH)
	if err != nil {
		panic(err)
	}
	for _, name := range entries.Names() {
		if name == "テンプレート.json" {
			continue
		}
		moveName := strings.TrimRight(name, ".json")
		fmt.Println(name + " 読み込み開始")
		moveData := LoadMoveData(MOVEDEX_PATH + name)
		if moveData.BasePP == 0 {
			fmt.Println(name, " の BasePPが0になっている")
		}
		fmt.Println(name + " 読み込み完了")
		y[MoveName(moveName)] = &moveData
	}
	return y
}()

var NEVER_MISS_HIT_MOVE_NAMES = func() MoveNames{
	y, err := omwjson.Load[MoveNames](NEVER_MISS_HIT_MOVE_NAMES_PATH)
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
	y, err := omwjson.Load[Naturedex](NATUREDEX_PATH)
	if err != nil {
		panic(err)
	}
	return y
}()

var ALL_NATURES = func() Natures {
	y, err := omwjson.Load[Natures](ALL_NATURES_PATH)
	if err != nil {
		panic(err)
	}
	return y
}()

type typeData map[string]float64
type typedex map[string]typeData

type TypeData map[Type]float64
type Typedex map[Type]TypeData

var TYPEDEX = func() Typedex {
	d, err := omwjson.Load[typedex](TYPEDEX_PATH)
	if err != nil {
		panic(err)
	}

	y := Typedex{}
	for t1, data := range d {
		type1, err := NewType(t1)
		if err != nil {
			panic(err)
		}
		if _, ok := y[type1]; !ok {
			y[type1] = TypeData{}
		}

		for t2, v := range data {
			type2, err := NewType(t2)
			if err != nil {
				panic(err)
			}
			if _, ok := y[type1][type2]; !ok {
				y[type1][type2] = v
			}
		}
	}
	return y
}()

func init() {
	for pokeName, pokeData := range POKEDEX {
		if !omwslices.IsSubset(ALL_GENDERS, pokeData.Genders) {
			msg := fmt.Sprintf("%v の 性別に %v 以外 の 要素 が 含まれている", pokeName, ALL_GENDERS)
			fmt.Println(msg)
		}

		for _, moveName := range pokeData.Learnset {
			if omwslices.Count(pokeData.Learnset, moveName) != 1 {
				fmt.Println(pokeName, " の ", moveName, " が 重複している")
			}
		}
	}
}