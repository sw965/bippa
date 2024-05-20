package bippa

import (
	ojson "github.com/sw965/omw/json"
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
	PSYCHIC
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
	"エスパー":PSYCHIC,
	"むし":BUG,
	"いわ":ROCK,
	"ゴースト":GHOST,
	"ドラゴン":DRAGON,
	"あく":DARK,
	"はがね":STEEL,
	"フェアリー":FAIRY,
}

type Types []Type

var ALL_TYPES = func() Types {
	buff, err := ojson.Load[[]string](ALL_TYPES_PATH)
	if err != nil {
		panic(err)
	}
	ret := make(Types, len(buff))
	for i, s := range buff {
		ret[i] = STRING_TO_TYPE[s]
	}
	return ret
}()