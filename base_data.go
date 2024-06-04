package bippa

import (
	"fmt"
	omwjson "github.com/sw965/omw/json"
)

var ALL_POKE_NAMES = func() PokeNames {
	buff, err := omwjson.Load[[]string](ALL_POKE_NAMES_PATH)
	if err != nil {
		panic(err)
	}
	names := make(PokeNames, len(buff))
	for i := range names {
		 names[i] = STRING_TO_POKE_NAME[buff[i]]
	}
	return names
}()

type pokeDataJSONBuffer struct {
	Types []string
	BaseHP int
	BaseAtk int
	BaseDef int
	BaseSpAtk int
	BaseSpDef int
	BaseSpeed int
	Learnset []string
}

type PokeData struct {
	Types Types
	BaseHP int
	BaseAtk int
	BaseDef int
	BaseSpAtk int
	BaseSpDef int
	BaseSpeed int
	Learnset MoveNames
}

func LoadPokeData(pokeName PokeName) (PokeData, error) {
	path := POKE_DATA_PATH + POKE_NAME_TO_STRING[pokeName] + omwjson.EXTENSION
	buff, err := omwjson.Load[pokeDataJSONBuffer](path)
	if err != nil {
		return PokeData{}, err
	}

	types := make(Types, len(buff.Types))
	for i := range types {
		types[i] = STRING_TO_TYPE[buff.Types[i]]
	}

	learnset := make(MoveNames, len(buff.Learnset))
	for i := range learnset {
		learnset[i] = STRING_TO_MOVE_NAME[buff.Learnset[i]]
	}

	return PokeData{
		Types:types,
		BaseHP:buff.BaseHP,
		BaseAtk:buff.BaseAtk,
		BaseDef:buff.BaseDef,
		BaseSpAtk:buff.BaseSpAtk,
		BaseSpDef:buff.BaseSpDef,
		BaseSpeed:buff.BaseSpeed,
		Learnset:learnset,
	}, nil
}

type Pokedex map[PokeName]*PokeData

var POKEDEX = func() Pokedex {
	dex := Pokedex{}
	for i := range ALL_POKE_NAMES {
		name := ALL_POKE_NAMES[i]
		data, err := LoadPokeData(name)
		if err != nil {
			panic(err)
		}
		dex[name] = &data
	}
	return dex
}()

var ALL_MOVE_NAMES = func() MoveNames {
	buff, err := omwjson.Load[[]string](ALL_MOVE_NAMES_PATH)
	if err != nil {
		panic(err)
	}
	names := make(MoveNames, len(buff))
	for i := range names {
		names[i] = STRING_TO_MOVE_NAME[buff[i]]
	}
	return names
}()

type moveDataJSONBuffer struct {
    Type        string
    Category    string
    Power       int
    Accuracy    int
    BasePP      int
}

type MoveData struct {
	Type Type
	Category MoveCategory
	Power int
	Accuracy int
	BasePP int
}

func LoadMoveData(moveName MoveName) (MoveData, error) {
	path := MOVE_DATA_PATH + MOVE_NAME_TO_STRING[moveName] + omwjson.EXTENSION
	buff, err := omwjson.Load[moveDataJSONBuffer](path)
	if err != nil {
		return MoveData{}, err
	}
	if buff.BasePP == 0 {
		moveNameStr := MOVE_NAME_TO_STRING[moveName]
		msg := fmt.Sprintf("%s の MoveDataのjsonファイルのbasePPが0になっている。", moveNameStr)
		return MoveData{}, fmt.Errorf(msg)
	}

	return MoveData{
		Type:STRING_TO_TYPE[buff.Type],
		Category:STRING_TO_MOVE_CATEGORY[buff.Category],
		Power:buff.Power,
		Accuracy:buff.Accuracy,
		BasePP:buff.BasePP,
	}, nil
}

type Movedex map[MoveName]*MoveData

var MOVEDEX = func() Movedex {
	dex := Movedex{}
	for i := range ALL_MOVE_NAMES {
		name := ALL_MOVE_NAMES[i]
		data, err := LoadMoveData(name)
		if err != nil {
			panic(err)
		}
		dex[name] = &data
	}
	return dex
}()

type typedexJSONBuffer map[string]map[string]float64
type DefTypeData map[Type]float64
type Typedex map[Type]DefTypeData

var TYPEDEX = func() Typedex {
	buff, err := omwjson.Load[typedexJSONBuffer](TYPEDEX_PATH)
	if err != nil {
		panic(err)
	}
	typedex := Typedex{}
	for atkTStr, defData := range buff {
		atkType := STRING_TO_TYPE[atkTStr]
		typedex[atkType] = DefTypeData{}
		for defTStr, effect := range defData {
			defType := STRING_TO_TYPE[defTStr]
			typedex[atkType][defType] = effect
		}
	}
	return typedex
}()

type NatureData struct {
	AtkBonus NatureBonus
	DefBonus NatureBonus
	SpAtkBonus NatureBonus
	SpDefBonus NatureBonus
	SpeedBonus NatureBonus
}

type naturedexJSONBuffer map[string]NatureData
type Naturedex map[Nature]NatureData

var NATUREDEX = func() Naturedex {
	buff, err := omwjson.Load[naturedexJSONBuffer](NATUREDEX_PATH)
	if err != nil {
		panic(err)
	}
	ret := Naturedex{}
	for natureStr, data := range buff {
		nature := STRING_TO_NATURE[natureStr]
		ret[nature] = data
	}
	return ret
}()

var ALL_NATURES = func() Natures {
	buff, err := omwjson.Load[[]string](ALL_NATURES_PATH)
	if err != nil {
		panic(err)
	}
	ret := make(Natures, len(buff))
	for i, natureStr := range buff {
		ret[i] = STRING_TO_NATURE[natureStr]
	}
	return ret
}()