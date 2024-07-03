package bippa

import (
	"fmt"
	"golang.org/x/exp/slices"
	osliecs "github.com/sw965/omw/slices"
	omaps "github.com/sw965/omw/maps"
)

type MoveName int

const (
	EMPTY_MOVE_NAME MoveName = iota
	THUNDERBOLT // 10まんボルト
	STONE_EDGE  // ストーンエッジ
	EMBER       // ひのこ
	TACKLE      // たいあたり
	VINE_WHIP   // みずでっぽう
	SURF        // なみのり
	WATER_GUN   // みずでっぽう
	ICE_BEAM    // れいとうビーム
	STRUGGLE    // わるあがき
)

var STRING_TO_MOVE_NAME = map[string]MoveName{
	"":           EMPTY_MOVE_NAME,
	"10まんボルト":THUNDERBOLT,
	"ストーンエッジ":STONE_EDGE,
	"ひのこ":     EMBER,
	"たいあたり": TACKLE,
	"つるのムチ":VINE_WHIP,
	"なみのり":SURF,
	"みずでっぽう": WATER_GUN,
	"れいとうビーム":ICE_BEAM,
	"わるあがき":STRUGGLE,
}

var MOVE_NAME_TO_STRING = omaps.Invert[map[MoveName]string](STRING_TO_MOVE_NAME)

func (name MoveName) ToString() string {
	return MOVE_NAME_TO_STRING[name]
}

type MoveNames []MoveName

func (names MoveNames) ToStrings() []string {
	ret := make([]string, len(names))
	for i, name := range names {
		ret[i] = name.ToString()
	}
	return ret
}

func (names MoveNames) Sort() MoveNames {
	ret := make(MoveNames, len(names))
	for i := 0; i < osliecs.Count(names, EMPTY_MOVE_NAME); i++ {
		ret = append(ret, EMPTY_MOVE_NAME)
	}

	for _, name := range ALL_MOVE_NAMES {
		if slices.Contains(names, name) {
			ret = append(ret, name)
		}
	}
	return ret
}

type MoveNamess []MoveNames

type MoveCategory int

const (
	PHYSICS MoveCategory = iota
	SPECIAL
	STATUS
)

var STRING_TO_MOVE_CATEGORY = map[string]MoveCategory{
	"物理":PHYSICS,
	"特殊":SPECIAL,
	"変化":STATUS,
}

var MOVE_CATEGORY_TO_STRING = omaps.Invert[map[MoveCategory]string](STRING_TO_MOVE_CATEGORY)

func (m MoveCategory) ToString() string {
	return MOVE_CATEGORY_TO_STRING[m]
}

type PointUp int

const (
	MAX_POINT_UP = 3
)

type PointUps []PointUp

type PowerPoint struct {
	Max int
	Current int
}

const (
	MIN_MOVESET = 1
	MAX_MOVESET = 4
)

func NewPowerPoint(base int, up PointUp) PowerPoint {
    increment := int(float64(base) / 5.0)
    max := base + (increment * int(up))
	return PowerPoint{Max:max, Current:max}
}

type Moveset map[MoveName]*PowerPoint

func NewMoveset(pokeName PokeName, moveNames MoveNames) (Moveset, error) {
	if len(moveNames) == 0 {
		msg := fmt.Sprintf("覚えさせる技が指定されていません。ポケモンには、少なくとも%dつ以上の技を覚えさせる必要があります。", MIN_MOVESET)
		return Moveset{}, fmt.Errorf(msg)
	}

	if len(moveNames) > MAX_MOVESET {
		msg := fmt.Sprintf("覚えさせる技の数が、限度を超えています。技は最大で%dつまで覚えさせることが出来ます。", MAX_MOVESET)
		return Moveset{}, fmt.Errorf(msg)
	}

	learnset := POKEDEX[pokeName].Learnset
	moveset := Moveset{}
	for i := range moveNames {
		moveName := moveNames[i]
		if !slices.Contains(learnset, moveNames[i]) {
			pokeNameStr := POKE_NAME_TO_STRING[pokeName]
			moveNameStr := MOVE_NAME_TO_STRING[moveName]
			msg := fmt.Sprintf("「%s」 は 「%s」 を覚えることができません。覚えられる技を選択してください。", pokeNameStr, moveNameStr)
			return Moveset{}, fmt.Errorf(msg)
		}
		basePP := MOVEDEX[moveName].BasePP
		moveset[moveName] = &PowerPoint{Max:basePP, Current:basePP}
	}
	return moveset, nil
}

func (m Moveset) Equal(other Moveset) bool {
	for moveName, pp := range m {
		otherPP, ok := other[moveName]
		if !ok {
			return false
		}
		if *pp != *otherPP {
			return false
		}
	}
	return true
}

func (m Moveset) Clone() Moveset {
	clone := Moveset{}
	for moveName, pp := range m {
		clone[moveName] = &PowerPoint{Max:pp.Max, Current:pp.Current}
	}
	return clone
}

func (m Moveset) ToEasyRead() EasyReadMoveset {
	ret := EasyReadMoveset{}
	for moveName, pp := range m {
		ret[moveName.ToString()] = *pp
	}
	return ret
}