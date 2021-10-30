import os
import boa
import base_data

KEYWORDS = ["進化時", "基本", "遺伝", "特別タマゴ技", "両方", "鎧の孤島", "教え技", "特別"] \
         + ["Lv." + str(i + 1) for i in range(100)] \
         + [("マシン" + str(i).zfill(2)) for i in range(100)] \
         + [("レコ." + str(i).zfill(2)) for i in range(100)]

def parse_move_name(move_name):
    if "[遺伝経路]New" in move_name:
        return move_name[:-len("[遺伝経路]New")]
    if "New" in move_name:
        return move_name[:-len("New")]
    elif "[遺伝経路]" in move_name:
        return move_name[:-len("[遺伝経路]")]
    else:
        return move_name

def parse_learnset(poketetu_learnset):
    split_data = poketetu_learnset.split()
    result = []
    for i, split_data_ele in enumerate(split_data):
        if split_data_ele in KEYWORDS:
            assert split_data[i + 2] in base_data.TYPEDEX
            move_name = parse_move_name(split_data[i + 1])

            if move_name == "ドレインパンチv1.1.1まで":
                assert poketetu_learnset.split()[1] == "マーシャドーがレベルアップで覚える技"
                continue

            assert move_name in base_data.ALL_MOVE_NAMES, move_name
            if move_name not in result:
                result.append(move_name)
    return result
