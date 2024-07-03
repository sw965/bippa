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

func (p *PokeData) ToEasyRead() EasyReadPokeData {
	return EasyReadPokeData{
		Types:p.Types.ToStrings(),
		BaseHP:p.BaseHP,
		BaseAtk:p.BaseAtk,
		BaseDef:p.BaseDef,
		BaseSpAtk:p.BaseSpAtk,
		BaseSpDef:p.BaseSpDef,
		BaseSpeed:p.BaseSpeed,
		Learnset:p.Learnset.ToStrings(),
	}
}

func LoadPokeData(pokeName PokeName) (PokeData, error) {
	if _, ok := POKE_NAME_TO_STRING[pokeName]; !ok {
		msg := fmt.Sprintf("%s が POKE_NAME_TO_STRING の中に存在しない", pokeName.ToString())
		return PokeData{}, fmt.Errorf(msg)
	}

	path := POKE_DATA_PATH + POKE_NAME_TO_STRING[pokeName] + omwjson.EXTENSION
	buff, err := omwjson.Load[EasyReadPokeData](path)
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

func (p Pokedex) ToEasyRead() EasyReadPokedex {
	ret := EasyReadPokedex{}
	for k, v := range p {
		ret[k.ToString()] = v.ToEasyRead()
	}
	return ret
}

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

type MoveData struct {
	Type Type
	Category MoveCategory
	Power int
	Accuracy int
	BasePP int
}

func LoadMoveData(moveName MoveName) (MoveData, error) {
	if _, ok := MOVE_NAME_TO_STRING[moveName]; !ok {
		msg := fmt.Sprintf("%s が MOVE_NAME_TO_STRING の中に存在しない", moveName.ToString())
		return MoveData{}, fmt.Errorf(msg)
	}

	path := MOVE_DATA_PATH + MOVE_NAME_TO_STRING[moveName] + omwjson.EXTENSION
	buff, err := omwjson.Load[EasyReadMoveData](path)
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

func (m *MoveData) ToEasyRead() EasyReadMoveData {
	return EasyReadMoveData{
		Type:m.Type.ToString(),
		Category:m.Category.ToString(),
		Power:m.Power,
		Accuracy:m.Accuracy,
		BasePP:m.BasePP,
	}
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

func (m Movedex) ToEasyRead() EasyReadMovedex {
	ret := EasyReadMovedex{}
	for moveName, moveData := range m {
		ret[moveName.ToString()] = moveData.ToEasyRead()
	}
	return ret
}

type DefTypeData map[Type]float64

func (t DefTypeData) ToEasyRead() EasyReadDefTypeData {
	ret := EasyReadDefTypeData{}
	for k, v := range t {
		ret[k.ToString()] = v
	}
	return ret
}

type Typedex map[Type]DefTypeData

var TYPEDEX = func() Typedex {
	buff, err := omwjson.Load[EasyReadTypedex](TYPEDEX_PATH)
	if err != nil {
		panic(err)
	}
	typedex := Typedex{}
	for atkTypeStr, defTypeData := range buff {
		atkType := STRING_TO_TYPE[atkTypeStr]
		typedex[atkType] = DefTypeData{}
		for defTypeStr, effect := range defTypeData {
			defType := STRING_TO_TYPE[defTypeStr]
			typedex[atkType][defType] = effect
		}
	}
	return typedex
}()

func (t Typedex) ToEasyRead() EasyReadTypedex {
	ret := EasyReadTypedex{}
	for k, typeData := range t {
		ret[k.ToString()] = typeData.ToEasyRead()
	}
	return ret
}

type NatureData struct {
	AtkBonus NatureBonus
	DefBonus NatureBonus
	SpAtkBonus NatureBonus
	SpDefBonus NatureBonus
	SpeedBonus NatureBonus
}

type Naturedex map[Nature]NatureData

var NATUREDEX = func() Naturedex {
	buff, err := omwjson.Load[EasyReadNaturedex](NATUREDEX_PATH)
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

func (n Naturedex) ToEasyRead() EasyReadNaturedex {
	ret := EasyReadNaturedex{}
	for k, v := range n {
		ret[k.ToString()] = v
	}
	return ret
}

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