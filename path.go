package bippa

import (
	//"os"
	//"path/filepath"
)

var (
	DATA_PATH = "C:/Go/project/bippa/main/bippa/"
	// DATA_PATH = func() string {
	// 	exePath, err := os.Executable()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	path := filepath.Dir(exePath) + "/bippa/"
	// 	return path
	// }()

	ALL_POKE_NAMES_PATH = DATA_PATH + "all_poke_names.json"
	POKE_DATA_PATH = DATA_PATH + "poke_data/"

	MOVE_DATA_PATH = DATA_PATH + "move_data/"
	ALL_MOVE_NAMES_PATH = DATA_PATH + "all_move_names.json"

	ALL_TYPES_PATH = DATA_PATH + "all_types.json"
	TYPEDEX_PATH = DATA_PATH + "typedex.json"

	ALL_NATURE_PATH = DATA_PATH + "all_natures.json"
	NATUREDEX_PATH = DATA_PATH + "naturedex.json"
	ALL_NATURES_PATH = DATA_PATH + "all_natures.json"
)