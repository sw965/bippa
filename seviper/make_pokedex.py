import copy
import json
import boa
import base_data
import poketetu_pokedex as pt_pokedex
import poketetu_learnset as pt_learnset

def main():
    for poke_name in base_data.ALL_POKE_NAMES:
        print(poke_name)
        poke_data = pt_pokedex.parse_poke_data(base_data.POKETETU_POKEDEX_PATH + poke_name + ".txt")
        poketetu_learnset = boa.read_txt(base_data.POKETETU_LEARNSET_PATH + poke_name + ".txt")
        learnset = pt_learnset.parse_learnset(poketetu_learnset)
        poke_data["Learnset"] = learnset
        json_data = json.dumps(poke_data, ensure_ascii=False, indent=4)
        boa.write_txt(base_data.POKEDEX_PATH + poke_name + ".json", json_data)

if __name__ == "__main__":
    main()
