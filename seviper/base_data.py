import os
import random
import boa as boa

BASE_DIRECTORY = os.path.dirname(os.path.abspath(__file__))
DATA_FOLDER_PATH = BASE_DIRECTORY + "/data/"

ALL_POKE_NAMES_PATH = DATA_FOLDER_PATH + "all_poke_names.txt"
ALL_MOVE_NAMES_PATH = DATA_FOLDER_PATH + "all_move_names.txt"

POKEDEX_PATH = DATA_FOLDER_PATH + "pokedex/"
MOVEDEX_PATH = DATA_FOLDER_PATH + "movedex/"
NATUREDEX_PATH = DATA_FOLDER_PATH + "naturedex.json"
TYPEDEX_PATH = DATA_FOLDER_PATH + "typedex.json"

POKETETU_DATA_FOLDER_PATH = BASE_DIRECTORY + "/poketetu_data/"
POKETETU_POKEDEX_PATH = POKETETU_DATA_FOLDER_PATH + "pokedex/"
POKETETU_MOVEDEX_PATH = POKETETU_DATA_FOLDER_PATH + "movedex/"
POKETETU_LEARNSET_PATH = POKETETU_DATA_FOLDER_PATH + "learnset/"

IMAGE_FOLDER_PATH = BASE_DIRECTORY + "/image/"
POKE_GIF_PATH = IMAGE_FOLDER_PATH + "gif/"
POKE_GIF_NORMAL_PATH = POKE_GIF_PATH + "normal/"
POKE_GIF_MIRROR_PATH = POKE_GIF_PATH + "mirror/"

GOOD_EFFECTIVE_BATTLE_MSG = "効果は抜群だ！"
BAD_EFFECTIVE_BATTLE_MSG = "効果はいまひとつのようだ..."
NO_EFFECTIVE_BATTLE_MSG = "効果はないようだ..."

EMPTY = "なし"

def load_naturedex():
    return boa.load_json(NATUREDEX_PATH)

def load_typedex():
    return boa.load_json(TYPEDEX_PATH)

ALL_POKE_NAMES = boa.readlines_txt(ALL_POKE_NAMES_PATH, True)
ALL_MOVE_NAMES = [file_name[:-4] for folder_name in os.listdir(POKETETU_MOVEDEX_PATH)\
                      for file_name in os.listdir(POKETETU_MOVEDEX_PATH + folder_name)]
NATUREDEX = load_naturedex()
TYPEDEX = load_typedex()

IVS = [i for i in range(32)]
EVS = [i for i in range(253)]
VALID_EVS = [ev for ev in EVS if ev%4 == 0]

MAX_IV = max(IVS)
MAX_EV = max(EVS)

STATS_NUM = 6
MAX_MOVESET_NUM = 4
MIN_TEAM_NUM = 3
MAX_TEAM_NUM = 6
FIGHTERS_NUM = 3
MAX_SUM_EV = 510
