import re
import boa
import base_data

FOLDER_PATH = "C:/Python35/pyckage/seviper/poketetu_data/pokedex/"

def read_poketetu_poke_data_lines(path):
    return [re.sub("\t", " ", line) for line in boa.readlines_txt(path, True)]

def parse_base_statuses(poketetu_poke_data_lines):
    keys = ["BaseHP", "BaseAtk", "BaseDef", "BaseSpAtk", "BaseSpDef", "BaseSpeed"]
    expected_marks = ["HP", "こうげき", "ぼうぎょ", "とくこう", "とくぼう", "すばやさ"]

    result = {}
    for i, key in enumerate(keys):
        mark, status = poketetu_poke_data_lines[i + 1].split()
        assert mark == expected_marks[i], mark
        result[key] = int(status)
    return result

def parse_category(poketetu_poke_data_lines):
    for line in poketetu_poke_data_lines:
        data = line.split()
        if data[0] == "カテゴリー":
            return data[1]
    assert False

def parse_gender(poketetu_poke_data_lines):
    for line in poketetu_poke_data_lines:
        data = line.split()
        if data[0] == "性別":
            if "♂" in line and "♀" in line:
                return "♂♀両方"
            elif "♂のみ" in data[1]:
                return "♂のみ"
            elif "♀のみ" in data[1]:
                return "♀のみ"
            elif data[1] == "ふめい":
                return "不明"
            assert False
    assert False

def parse_egg_groups(poketetu_poke_data_lines):
    for line in poketetu_poke_data_lines:
        data = line.split()
        if data[0] == "タマゴグループ":
            if len(data) == 4:
                return [data[1], data[3]]
            elif len(data) == 2:
                return [data[1]]
            elif len(data) == 3:
                assert data[2] == "(性別不明)"
                return [data[1]]
            elif len(data) == 5:
                assert data[4] == "(性別不明)"
                return [data[1], data[3]]
            assert False
    assert False

def parse_abilities(poketetu_poke_data_lines):
    for i, line in enumerate(poketetu_poke_data_lines):
        data = line.split()
        if len(data) == 1:
            continue

        if data[0] == "ヘルプ◆" and "の特性(とくせい)" in data[1]:
            abilities = []
            abilities.append(poketetu_poke_data_lines[i + 1].split()[0])
            if poketetu_poke_data_lines[i + 2].split()[0] == "ヘルプ◆":
                return abilities
            else:
                abilities.append(poketetu_poke_data_lines[i + 2].split()[0])
                return abilities
    assert False

def parse_hidden_ability(poketetu_poke_data_lines):
    for i, line in enumerate(poketetu_poke_data_lines):
        data = line.split()
        if len(data) == 1:
            continue

        if data[0] == "ヘルプ◆" and "の隠れ特性(夢特性)" in data[1]:
            result = poketetu_poke_data_lines[i + 1].split()[0]
            if result == "なし":
                return result
            assert result[0] == "*"
            return result[1:]
    assert False

def parse_all_abilities(poketetu_poke_data_lines):
    result = parse_abilities(poketetu_poke_data_lines)
    hidden_ability = parse_hidden_ability(poketetu_poke_data_lines)
    if hidden_ability != "なし":
        result += [hidden_ability]
    return result

def parse_height(poketetu_poke_data_lines):
    for line in poketetu_poke_data_lines:
        data = line.split()
        if data[0] == "高さ":
            return float(data[1][:-1])
    assert False

def parse_weight(poketetu_poke_data_lines):
    for i, line in enumerate(poketetu_poke_data_lines):
        if line[:2] == "重さ":
            return float(poketetu_poke_data_lines[i + 1][:-2])
    assert False

def parse_types(poketetu_poke_data_lines):
    for i, line in enumerate(poketetu_poke_data_lines):
        if line[:3] == "タイプ":
            assert "けたぐり" in poketetu_poke_data_lines[i -1]
            types = [poketetu_poke_data_lines[i + 1]]
            if poketetu_poke_data_lines[i + 2] in base_data.TYPEDEX:
                types.append(poketetu_poke_data_lines[i + 2])
            return types
    assert False

def parse_poke_data(path):
    poketetu_poke_data_lines = read_poketetu_poke_data_lines(path)
    result = {}
    result.update(parse_base_statuses(poketetu_poke_data_lines))
    result["Category"] = parse_category(poketetu_poke_data_lines)
    result["Gender"] = parse_gender(poketetu_poke_data_lines)
    result["EggGroups"] = parse_egg_groups(poketetu_poke_data_lines)
    result["NormalAbilities"] = parse_abilities(poketetu_poke_data_lines)
    result["HiddenAbility"] = parse_hidden_ability(poketetu_poke_data_lines)
    result["AllAbilities"] = parse_all_abilities(poketetu_poke_data_lines)
    result["Height"] = parse_height(poketetu_poke_data_lines)
    result["Weight"] = parse_weight(poketetu_poke_data_lines)
    result["Types"] = parse_types(poketetu_poke_data_lines)
    return result
