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

	ret, err := StringsToPokeNames(buff)
	if err != nil {
		panic(err)
	}
	return ret
}()

type PokeData struct {
	Types Types
	BaseHP int
	BaseAtk int
	BaseDef int
	BaseSpAtk int
	BaseSpDef int
	BaseSpeed int
	Abilities Abilities
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
		Abilities:p.Abilities.ToStrings(),
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
	return buff.From()
}

type Pokedex map[PokeName]*PokeData

var POKEDEX = func() Pokedex {
	ret := Pokedex{}
	for _, name := range ALL_POKE_NAMES {
		data, err := LoadPokeData(name)
		if err != nil {
			panic(err)
		}
		ret[name] = &data
	}
	return ret
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

	ret, err := StringsToMoveNames(buff)
	if err != nil {
		panic(err)
	}
	return ret
}()

type MoveData struct {
	Type Type
	Category MoveCategory
	Power int
	Accuracy int
	BasePP int
	IsContact bool
    PriorityRank int
    CriticalRank CriticalRank
    Target TargetRange
}

func LoadMoveData(moveName MoveName) (MoveData, error) {
	path := MOVE_DATA_PATH + MOVE_NAME_TO_STRING[moveName] + omwjson.EXTENSION
	buff, err := omwjson.Load[EasyReadMoveData](path)
	if err != nil {
		return MoveData{}, err
	}
	return buff.From()
}

func (m *MoveData) ToEasyRead() EasyReadMoveData {
	return EasyReadMoveData{
		Type:m.Type.ToString(),
		Category:m.Category.ToString(),
		Power:m.Power,
		Accuracy:m.Accuracy,
		BasePP:m.BasePP,
		IsContact:m.IsContact,
		PriorityRank:m.PriorityRank,
		CriticalRank:m.CriticalRank,
		Target:m.Target.ToString(),
	}
}

type Movedex map[MoveName]*MoveData

var MOVEDEX = func() Movedex {
	ret := Movedex{}
	for _, name := range ALL_MOVE_NAMES {
		data, err := LoadMoveData(name)
		if err != nil {
			panic(err)
		}
		ret[name] = &data
	}
	return ret
}()

func (m Movedex) ToEasyRead() EasyReadMovedex {
	ret := EasyReadMovedex{}
	for name, data := range m {
		ret[name.ToString()] = data.ToEasyRead()
	}
	return ret
}

type DefTypeData map[Type]float64

func (t DefTypeData) ToEasyRead() EasyReadDefTypeData {
	ret := EasyReadDefTypeData{}
	for t, data := range t {
		ret[t.ToString()] = data
	}
	return ret
}

type Typedex map[Type]DefTypeData

var TYPEDEX = func() Typedex {
	buff, err := omwjson.Load[EasyReadTypedex](TYPEDEX_PATH)
	if err != nil {
		panic(err)
	}

	ret := Typedex{}
	for atkTypeStr, defTypeData := range buff {
		atkType, err := StringToType(atkTypeStr)
		if err != nil {
			panic(err)
		}
		ret[atkType] = DefTypeData{}
		for defTypeStr, effect := range defTypeData {
			defType, err := StringToType(defTypeStr)
			if err != nil {
				panic(err)
			}
			ret[atkType][defType] = effect
		}
	}
	return ret
}()

func (t Typedex) ToEasyRead() EasyReadTypedex {
	ret := EasyReadTypedex{}
	for k, v := range t {
		ret[k.ToString()] = v.ToEasyRead()
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

type Naturedex map[Nature]*NatureData

var NATUREDEX = func() Naturedex {
	buff, err := omwjson.Load[EasyReadNaturedex](NATUREDEX_PATH)
	if err != nil {
		panic(err)
	}
	
	ret := Naturedex{}
	for k := range buff {
		nature, err := StringToNature(k)
		if err != nil {
			panic(err)
		}
		v := buff[k]
		ret[nature] = &v
	}
	return ret
}()

func (n Naturedex) ToEasyRead() EasyReadNaturedex {
	ret := EasyReadNaturedex{}
	for k, v := range n {
		ret[k.ToString()] = *v
	}
	return ret
}

var ALL_NATURES = func() Natures {
	buff, err := omwjson.Load[[]string](ALL_NATURES_PATH)
	if err != nil {
		panic(err)
	}

	ret, err := StringsToNatures(buff)
	if err != nil {
		panic(err)
	}
	return ret
}()

var ALL_ITEMS = func() Items {
	buff, err := omwjson.Load[[]string](ALL_ITEMS_PATH)
	if err != nil {
		panic(err)
	}

	ret, err := StringsToItems(buff)
	if err != nil {
		panic(err)
	}
	return ret
}

type CriticalRank int