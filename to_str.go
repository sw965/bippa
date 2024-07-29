package bippa

import (
	omwmaps "github.com/sw965/omw/maps"
)

var POKE_NAME_TO_STRING = omwmaps.Invert[map[PokeName]string](STRING_TO_POKE_NAME)
var NATURE_TO_STRING = omwmaps.Invert[map[Nature]string](STRING_TO_NATURE)
var ITEM_TO_STRING = omwmaps.Invert[map[Item]string](STRING_TO_ITEM)