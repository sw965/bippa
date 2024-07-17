package bippa

import (
	omwmaps "github.com/sw965/omw/maps"
	omwjson "github.com/sw965/omw/json"
	"golang.org/x/exp/slices"
	omwslices "github.com/sw965/omw/slices"
	"github.com/sw965/omw/fn"
)

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

var ALL_TWO_TYPESS = func() TypesSlice {
	return omwslices.Concat(
		fn.Map[TypesSlice](ALL_TYPES, func(t Type) Types { return Types{t} }),
		omwslices.Combination[TypesSlice, Types](ALL_TYPES, 2),
	)
}()