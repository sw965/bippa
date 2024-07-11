package bippa

import (
	//"os"
	//"path/filepath"
)

var (
	DATA_PATH = "C:/Go/project/bippa/main/data/fourth-generation/"
	// DATA_PATH = func() string {
	// 	exePath, err := os.Executable()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	path := filepath.Dir(exePath) + "/bippa/"
	// 	return path
	// }()

	ALL_POKE_NAMES_PATH = DATA_PATH + "all-poke-names.json"
	POKE_DATA_PATH = DATA_PATH + "poke-data/"

	MOVE_DATA_PATH = DATA_PATH + "move-data/"
	ALL_MOVE_NAMES_PATH = DATA_PATH + "all-move-names.json"

	ALL_TYPES_PATH = DATA_PATH + "all-types.json"
	TYPEDEX_PATH = DATA_PATH + "typedex.json"

	ALL_NATURE_PATH = DATA_PATH + "all-natures.json"
	NATUREDEX_PATH = DATA_PATH + "naturedex.json"
	ALL_NATURES_PATH = DATA_PATH + "all-natures.json"
)