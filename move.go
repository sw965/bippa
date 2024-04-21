package bippa

import (
	"github.com/sw965/omw"
)

type MoveName int

const (
	EMPTY_MOVE_NAME MoveName = iota
	TACKLE
)

var STRING_TO_MOVE_NAME = map[string]MoveName{
	"":EMPTY_MOVE_NAME,
	"たいあたり":TACKLE,
}
var MOVE_NAME_TO_STRING = omw.InvertMap[map[MoveName]string](STRING_TO_MOVE_NAME)

type MoveNames []MoveName
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

var MOVE_CATEGORY_TO_STRING = omw.InvertMap[map[MoveCategory]string](STRING_TO_MOVE_CATEGORY)

type PowerPoint struct {
	Max int
	Current int
}

const (
	MAX_MOVESET_NUM = 4
)

type Moveset map[MoveName]*PowerPoint

func (moveset Moveset) Equal(other Moveset) bool {
	for moveName, pp := range moveset {
		if otherPP, ok := other[moveName]; !ok {
			if *pp != *otherPP {
				return false
			}
		}
	}
	return true
}