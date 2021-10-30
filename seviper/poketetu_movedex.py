import re
import boa
import base_data

def read_poketetu_move_data_lines(path):
    return [re.sub("\t", " ", line) for line in boa.readlines_txt(path, True)]

def parse_type(poketetu_move_data_lines):
    data = poketetu_move_data_lines[3].split()
    assert data[0] == "タイプ"
    result = data[1]
    assert result in base_data.TYPEDEX
    return result

def parse_category(poketetu_move_data_lines):
    data = poketetu_move_data_lines[3].split()
    assert data[2] == "分類"
    category = data[3]
    assert category in ["物理", "特殊", "変化"]
    return category

def parse_power(poketetu_move_data_lines):
    data = poketetu_move_data_lines[4].split()
    assert data[0] == "威力"
    power = data[1]
    if power == "-":
        return -1
    else:
        return int(power)

def parse_accuracy(poketetu_move_data_lines):
    data = poketetu_move_data_lines[4].split()
    assert data[2] == "命中率"
    accuracy = data[3]
    if accuracy == "-":
        return -1
    else:
        return int(accuracy)

def parse_base_pp(poketetu_move_data_lines):
    data = poketetu_move_data_lines[5].split()
    assert data[0] == "PP"
    return int(data[1])

def parse_target(poketetu_move_data_lines):
    data = poketetu_move_data_lines[5].split()
    assert data[2] == "対象"
    return data[3]

def parse_contact(poketetu_move_data_lines):
    data = poketetu_move_data_lines[6].split()
    assert data[0] == "直接攻撃"
    return data[1]

def parse_protect(poketetu_move_data_lines):
    data = poketetu_move_data_lines[6].split()
    assert data[2] == "まもる"
    return data[3]

def parse_magic_coat(poketetu_move_data_lines):
    data = poketetu_move_data_lines[7].split()
    assert data[0] == "マジックコート"
    return data[1]

def parse_mirror_move(poketetu_move_data_lines):
    data = poketetu_move_data_lines[8].split()
    assert data[0] == "オウムがえし"
    return data[1]

def parse_substitute(poketetu_move_data_lines):
    data = poketetu_move_data_lines[8].split()
    assert data[2] == "みがわり"
    return data[3]

def parse_gigantamax_move(poketetu_move_data_lines):
    data = poketetu_move_data_lines[9].split()
    assert data[0] == "ダイマックス技"

    if "威" in data[1]:
        index = data[1].index("威")
        return data[1][:index - 1]
    else:
        assert parse_category(poketetu_move_data_lines) == "変化"
        return data[1]

def parse_gigantamax_power(poketetu_move_data_lines):
    data = poketetu_move_data_lines[9].split()
    assert data[0] == "ダイマックス技"

    if ":" in data[1]:
        index = data[1].index(":")
        return int(data[1][index + 1:-1])
    else:
        assert parse_category(poketetu_move_data_lines) == "変化"
        return -1

def parse_priority_rank(poketetu_move_data_lines):
    data = poketetu_move_data_lines[10].split()
    assert len(data) == 2
    assert data[0] == "効果"

    effect_description = data[1]

    if "優先度" in effect_description:
        index = effect_description.index(":")
        return int(effect_description[index + 1:index + 3])
    else:
        return 0

def parse_critical_rank(poketetu_move_data_lines):
    data = poketetu_move_data_lines[10].split()
    assert len(data) == 2
    assert data[0] == "効果"

    effect_description = data[1]

    if "急所に当たりやすい(急所ランク:" in data[1]:
        index = effect_description.index(":")
        return int(effect_description[index + 1:index + 3])
    else:
        return 0

def parse_move_data(path):
    poketetu_move_data_lines = read_poketetu_move_data_lines(path)
    result = {}
    result["Type"] = parse_type(poketetu_move_data_lines)
    result["Category"] = parse_category(poketetu_move_data_lines)
    result["Power"] = parse_power(poketetu_move_data_lines)
    result["Accuracy"] = parse_accuracy(poketetu_move_data_lines)
    result["BasePP"] = parse_base_pp(poketetu_move_data_lines)
    result["Target"] = parse_target(poketetu_move_data_lines)
    result["Contact"] = parse_contact(poketetu_move_data_lines)
    result["Protect"] = parse_protect(poketetu_move_data_lines)
    result["MagicCoat"] = parse_magic_coat(poketetu_move_data_lines)
    result["MirrorMove"] = parse_mirror_move(poketetu_move_data_lines)
    result["Substitute"] = parse_substitute(poketetu_move_data_lines)
    result["GigantamaxMove"] = parse_gigantamax_move(poketetu_move_data_lines)
    result["GigantamaxPower"] = parse_gigantamax_power(poketetu_move_data_lines)
    result["PriorityRank"] = parse_priority_rank(poketetu_move_data_lines)
    result["CriticalRank"] = parse_critical_rank(poketetu_move_data_lines)
    return result
