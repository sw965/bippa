package bippa

import (
	omwjson "github.com/sw965/omw/json"
	omwmaps "github.com/sw965/omw/maps"
	"golang.org/x/exp/slices"
	"github.com/sw965/omw/fn"
	omwslices "github.com/sw965/omw/slices"
)

type TypeData map[Type]float64

func (t TypeData) ToEasyRead() EasyReadTypeData {
	e := EasyReadTypeData{}
	for k, v := range t {
		e[k.ToString()] = v
	}
	return e
}

type EasyReadTypeData map[string]float64

func(e EasyReadTypeData) From() (TypeData, error) {
	d := TypeData{}
	for k, v := range e {
		t, err := StringToType(k)
		if err != nil {
			return TypeData{}, err
		}
		d[t] = v
	}
	return d, nil
}

type Typedex map[Type]TypeData

var TYPEDEX = func() Typedex {
	e, err := omwjson.Load[EasyReadTypedex](TYPEDEX_PATH)
	if err != nil {
		panic(err)
	}
	d, err := e.From()
	if err != nil {
		panic(err)
	}
	return d
}()

func (t Typedex) Effectiveness(atk Type, def Types) float64 {
	v := 1.0
	for _, e := range def {
		v *= t[atk][e]
	}
	return v
}

func (t Typedex) Effective(atk Type, def Types) TypeEffective {
	v := t.Effectiveness(atk, def)
	if v == 1.0 {
		return NEUTRAL_EFFECTIVE
	} else if v > 1.0 {
		return SUPER_EFFECTIVE
	} else if v == 0.0 {
		return NO_EFFECTIVE
	} else {
		return NOT_VERY_EFFECTIVE
	}
}

func (t Typedex) ToEasyRead() EasyReadTypedex {
	e := EasyReadTypedex{}
	for k, v := range t {
		e[k.ToString()] = v.ToEasyRead()
	}
	return e
}

type EasyReadTypedex map[string]EasyReadTypeData

func (e EasyReadTypedex) From() (Typedex, error) {
	d := Typedex{}
	for k, v := range e {
		t, err := StringToType(k)
		if err != nil {
			return Typedex{}, err
		}
		d[t], err = v.From()
		if err != nil {
			return Typedex{}, err
		}
	}
	return d, nil
}

type Type int

const (
	NORMAL Type  = iota
	FIRE
	WATER
	GRASS
	ELECTRIC
	ICE
	FIGHTING
	POISON
	GROUND
	FLYING
	PSYCHIC_TYPE
	BUG
	ROCK
	GHOST
	DRAGON
	DARK
	STEEL
	FAIRY
)

var STRING_TO_TYPE = map[string]Type{
	"ノーマル":NORMAL,
	"ほのお":FIRE,
	"みず":WATER,
	"くさ":GRASS,
	"でんき":ELECTRIC,
	"こおり":ICE,
	"かくとう":FIGHTING,
	"どく":POISON,
	"じめん":GROUND,
	"ひこう":FLYING,
	"エスパー":PSYCHIC_TYPE,
	"むし":BUG,
	"いわ":ROCK,
	"ゴースト":GHOST,
	"ドラゴン":DRAGON,
	"あく":DARK,
	"はがね":STEEL,
	"フェアリー":FAIRY,
}

var TYPE_TO_STRING = omwmaps.Invert[map[Type]string](STRING_TO_TYPE)

func (t Type) ToString() string {
	return TYPE_TO_STRING[t]
}

type Types []Type

var ALL_TYPES = func() Types {
	buff, err := omwjson.Load[[]string](ALL_TYPES_PATH)
	if err != nil {
		panic(err)
	}
	ret := make(Types, len(buff))
	for i, s := range buff {
		ret[i] = STRING_TO_TYPE[s]
	}
	return ret
}()

func (ts Types) ToStrings() []string {
	ret := make([]string, len(ts))
	for i, t := range ts {
		ret[i] = t.ToString()
	}
	return ret
}

func (ts Types) Sort() Types {
	ret := slices.Clone(ts)
	slices.SortFunc(ret, func(t1, t2 Type) bool { return slices.Index(ALL_TYPES, t1) < slices.Index(ALL_TYPES, t2) } )
	return ret
}

type TypesSlice []Types

var ALL_TWO_TYPES_SLICE = func() TypesSlice {
	return omwslices.Concat(
		fn.Map[TypesSlice](ALL_TYPES, func(t Type) Types { return Types{t} }),
		omwslices.Combination[TypesSlice, Types](ALL_TYPES, 2),
	)
}()

type TypeEffective int

const (
	NEUTRAL_EFFECTIVE TypeEffective = iota
	SUPER_EFFECTIVE
	NOT_VERY_EFFECTIVE
	NO_EFFECTIVE
)