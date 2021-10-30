import os
import json
import boa
import base_data
import poketetu_movedex as pt_movedex

class AttackNum:
    def __init__(self, min_, max_):
        self.min = min_
        self.max = max_

MIN_AND_MAX_ATTACK_NUM = {
    "ギアソーサー":AttackNum(2, 2),
    "スイープビンタ":AttackNum(2, 5),
    "すいりゅうれんだ":AttackNum(3, 3),
    "スケイルショット":AttackNum(2, 5),
    "タネマシンガン":AttackNum(2, 5),
    "ダブルアタック":AttackNum(2, 2),
    "ダブルウイング":AttackNum(2, 2),
    "ダブルチョップ":AttackNum(2, 2),
    "ダブルパンツァー":AttackNum(2, 2),
    "つっぱり":AttackNum(2, 5),
    "つららばり":AttackNum(2, 5),
    "トリプルアクセル":AttackNum(3, 3),
    "トリプルキック":AttackNum(3, 3),
    "ドラゴンアロー":AttackNum(2, 2),
    "にどげり":AttackNum(2, 2),
    "ホネブーメラン":AttackNum(2, 2),
    "ボーンラッシュ":AttackNum(2, 5),
    "ミサイルばり":AttackNum(2, 5),
    "みずしゅりけん":AttackNum(2, 5),
    "みだれづき":AttackNum(2, 5),
    "みだれひっかき":AttackNum(2, 5),
    "ロックブラスト":AttackNum(2, 5),
}

def main():
    for japanese_syllabary_folder_name in os.listdir(base_data.POKETETU_MOVEDEX_PATH):
        folder_path = base_data.POKETETU_MOVEDEX_PATH + japanese_syllabary_folder_name + "/"
        for file_name in os.listdir(folder_path):
            move_name = file_name[:-4]
            print(move_name)
            full_path = folder_path + file_name
            move_data = pt_movedex.parse_move_data(full_path)

            if move_name in MIN_AND_MAX_ATTACK_NUM:
                attack_num = MIN_AND_MAX_ATTACK_NUM[move_name]
                move_data["MinAttackNum"] = attack_num.min
                move_data["MaxAttackNum"] = attack_num.max
            else:
                move_data["MinAttackNum"] = 1
                move_data["MaxAttackNum"] = 1

            json_move_data = json.dumps(move_data, ensure_ascii=False, indent=4)
            boa.write_txt(base_data.MOVEDEX_PATH + move_name + ".json", json_move_data)

if __name__ == "__main__":
    main()
