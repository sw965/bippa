package bippa

import (
	"os"
)

var (
	SW965_PATH = os.Getenv("GOPATH") + "sw965/"

	DATA_PATH      = SW965_PATH + "pokemon_sv/"
	POKEDEX_PATH   = DATA_PATH + "pokedex/"
	MOVEDEX_PATH   = DATA_PATH + "movedex/"
	NATUREDEX_PATH = DATA_PATH + "naturedex.json"
	TYPEDEX_PATH   = DATA_PATH + "typedex.json"

	ALL_POKE_NAMES_PATH = DATA_PATH + "all_poke_names.json"
	ALL_NATURES_PATH    = DATA_PATH + "all_natures.json"
	ALL_MOVE_NAMES_PATH = DATA_PATH + "all_move_names.json"

	ALL_ITEMS_PATH      = DATA_PATH + "all_items.json"
)