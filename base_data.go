package bippa

import (
	"fmt"
	ojson "github.com/sw965/omw/json"
	opath "github.com/sw965/omw/path"
)

var (
	DATA_PATH = opath.SW965 + "bippa/"
	ALL_POKE_NAMES_PATH = DATA_PATH + "all_poke_names.json"
	POKE_DATA_PATH = DATA_PATH + "poke_data/"

	MOVE_DATA_PATH = DATA_PATH + "move_data/"
	ALL_MOVE_NAMES_PATH = DATA_PATH + "all_move_names.json"

	TYPEDEX_PATH = DATA_PATH + "typedex.json"
)

var ALL_POKE_NAMES = func() PokeNames {
	buff, err := ojson.Load[[]string](ALL_POKE_NAMES_PATH)
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
	path := POKE_DATA_PATH + POKE_NAME_TO_STRING[pokeName] + ojson.EXTENSION
	buff, err := ojson.Load[pokeDataJSONBuffer](path)
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
	buff, err := ojson.Load[[]string](ALL_MOVE_NAMES_PATH)
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
	path := MOVE_DATA_PATH + MOVE_NAME_TO_STRING[moveName] + ojson.EXTENSION
	buff, err := ojson.Load[moveDataJSONBuffer](path)
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
	buff, err := ojson.Load[typedexJSONBuffer](TYPEDEX_PATH)
	if err != nil {
		panic(err)
	}
	typedex := Typedex{}
	for strAtkT, defData := range buff {
		atkType := STRING_TO_TYPE[strAtkT]
		typedex[atkType] = DefTypeData{}
		for strDefT, effect := range defData {
			defType := STRING_TO_TYPE[strDefT]
			typedex[atkType][defType] = effect
		}
	}
	return typedex
}()