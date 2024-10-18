package bippa

import (
	omwmaps "github.com/sw965/omw/maps"
)

var POKE_NAME_TO_STRING = omwmaps.Invert[map[PokeName]string](STRING_TO_POKE_NAME)
var GENDER_TO_STRING = omwmaps.Invert[map[Gender]string](STRING_TO_GENDER)
var NATURE_TO_STRING = omwmaps.Invert[map[Nature]string](STRING_TO_NATURE)
var ABILITY_TO_STRING = omwmaps.Invert[map[Ability]string](STRING_TO_ABILITY)
var ITEM_TO_STRING = omwmaps.Invert[map[Item]string](STRING_TO_ITEM)
var MOVE_NAME_TO_STRING = omwmaps.Invert[map[MoveName]string](STRING_TO_MOVE_NAME)
var MOVE_CATEGORY_TO_STRING = omwmaps.Invert[map[MoveCategory]string](STRING_TO_MOVE_CATEGORY)
