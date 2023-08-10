package bippa

import (
	"fmt"
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

func NewType(s string) (Type, error) {
	switch s {
		case "ノーマル":
			return NORMAL, nil
		case "ほのお":
			return FIRE, nil
		case "みず":
			return WATER, nil
		case "くさ":
			return GRASS, nil
		case "でんき":
			return ELECTRIC, nil
		case "こおり":
			return ICE, nil
		case "かくとう":
			return FIGHTING, nil
		case "どく":
			return POISON, nil
		case "じめん":
			return GROUND, nil
		case "ひこう":
			return FLYING, nil
		case "エスパー":
			return PSYCHIC, nil
		case "むし":
			return BUG, nil
		case "いわ":
			return ROCK, nil
		case "ゴースト":
			return GHOST, nil
		case "ドラゴン":
			return DRAGON, nil
		case "あく":
			return DARK, nil
		case "はがね":
			return STEEL, nil
		case "フェアリー":
			return FAIRY, nil
		default:
			return -1, fmt.Errorf("不適なtype")
	}
}

type Types []Type

func NewTypes(ss []string) (Types, error) {
	ys := make(Types, len(ss))
	for i, s := range ss {
		y, err := NewType(s)
		if err != nil {
			return ys, err
		}
		ys[i] = y
	}
	return ys, nil
}