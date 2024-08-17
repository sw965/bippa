package bippa

import (
	"fmt"
	"github.com/sw965/omw/fn"
)

var STRING_TO_POKE_NAME = map[string]PokeName{
    "":          EMPTY_POKE_NAME,
    "ギャラドス": GYARADOS,
    "カビゴン":   SNORLAX,
    "ドーブル":   SMEARGLE,
    "ボーマンダ": SALAMENCE,
    "メタグロス": METAGROSS,
    "ラティオス": LATIOS,
    "エンペルト": EMPOLEON,
    "ドータクン": BRONZONG,
    "ドクロッグ": TOXICROAK,
}

func StringToPokeName(s string) (PokeName, error) {
	if n, ok := STRING_TO_POKE_NAME[s]; !ok {
		m := fmt.Sprintf("%s は PokeNameに変換出来ません", s)
		return EMPTY_POKE_NAME, fmt.Errorf(m)
	} else {
		return n, nil
	}
}

func StringsToPokeNames(ss []string) (PokeNames, error) {
	return fn.MapWithError[PokeNames](ss, StringToPokeName)
}

var STRING_TO_TYPE = map[string]Type{
	"ノーマル":NORMAL,
	"ほのお":FIRE,
	"みず":WATER,
	"くさ":GRASS,
	"でんき":ELECTRIC,
	"こおり":ICE,
	"かくとう":FIGHTING,
	"どく":POISON,
	"じめん":GROUND,
	"ひこう":FLYING,
	"エスパー":PSYCHIC_TYPE,
	"むし":BUG,
	"いわ":ROCK,
	"ゴースト":GHOST,
	"ドラゴン":DRAGON,
	"あく":DARK,
	"はがね":STEEL,
	"フェアリー":FAIRY,
}

func StringToType(s string) (Type, error) {
	if t, ok := STRING_TO_TYPE[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_TYPE に含まれていない為、Typeに変換出来ません。", s)
		return t, fmt.Errorf(msg)
	} else {
		return t, nil
	}
}

func StringsToTypes(ss []string) (Types, error) {
	return fn.MapWithError[Types](ss, StringToType)
}

var STRING_TO_NATURE = map[string]Nature{
    ""         : EMPTY_NATURE,
    "さみしがり": LONELY,
    "いじっぱり": ADAMANT,
    "やんちゃ":   NAUGHTY,
    "ゆうかん":   BRAVE,

    "ずぶとい":   BOLD,
    "わんぱく":   IMPISH,
    "のうてんき": LAX,
    "のんき":     RELAXED,

    "ひかえめ":   MODEST,
    "おっとり":   MILD,
    "うっかりや": RASH,
    "れいせい":   QUIET,

    "おだやか":   CALM,
    "おとなしい": GENTLE,
    "しんちょう": CAREFUL,
    "なまいき":   SASSY,

    "おくびょう": TIMID,
    "せっかち":   HASTY,
    "ようき":     JOLLY,
    "むじゃき":   NAIVE,

    "てれや":     BASHFUL,
    "がんばりや": HARDY,
    "すなお":     DOCILE,
    "きまぐれ":   QUIRKY,
    "まじめ":     SERIOUS,
}

func StringToNature(s string) (Nature, error) {
	if nature, ok := STRING_TO_NATURE[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_NATURE に含まれていない為、Natureに変換出来ません。", s)
		return nature, fmt.Errorf(msg)
	} else {
		return nature, nil
	}
}

func StringsToNatures(ss []string) (Natures, error) {
	return fn.MapWithError[Natures](ss, StringToNature)
}

var STRING_TO_ABILITY = map[string]Ability{
    "げきりゅう":TORRENT,
    "めんえき": IMMUNITY,
    "あついしぼう": THICK_FAT,
    "いかく": INTIMIDATE,
    "ふゆう": LEVITATE,
	"たいねつ":HEATPROOF,
    "マイペース": OWN_TEMPO,
    "テクニシャン": TECHNICIAN,
    "きけんよち": ANTICIPATION,
    "かんそうはだ": DRY_SKIN,
    "クリアボディ": CLEAR_BODY,
}

func StringToAbility(s string) (Ability, error) {
	if ability, ok := STRING_TO_ABILITY[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_ABILITY に含まれていない為、Abilityに変換出来ません。", s)
		return ability, fmt.Errorf(msg)
	} else {
		return ability, nil
	}
}

func StringsToAbilities(ss []string) (Abilities, error) {
	return fn.MapWithError[Abilities](ss, StringToAbility)
}

var STRING_TO_ITEM = map[string]Item{
	"イバンのみ":IAPAPA_BERRY,
	"オボンのみ":SITRUS_BERRY,
	"カゴのみ":CHESTO_BERRY,
	"きあいのタスキ":FOCUS_SASH,
	"ソクノのみ":WACAN_BERRY,
	"ラムのみ":LUM_BERRY,
}

func StringToItem(s string) (Item, error) {
	if item, ok := STRING_TO_ITEM[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_ITEM に含まれていない為、Itemに変換出来ません。", s)
		return item, fmt.Errorf(msg)
	} else {
		return item, nil
	}
}

func StringsToItems(ss []string) (Items, error) {
	return fn.MapWithError[Items](ss, StringToItem)
}

var STRING_TO_MOVE_NAME = map[string]MoveName{
    "":              EMPTY_MOVE_NAME,
    "10まんボルト":   THUNDERBOLT,
    "アームハンマー": HAMMER_ARM,
    "ストーンエッジ": STONE_EDGE,
    "なみのり":      SURF,
    "れいとうビーム":  ICE_BEAM,
    "わるあがき":    STRUGGLE,
    "あまごい":      RAIN_DANCE,
    "いわなだれ":    ROCK_SLIDE,
    "おんがえし":    RETURN,
    "かみくだく":    CRUNCH,
    "がむしゃら":    ENDEAVOR,
    "こごえるかぜ":  ICY_WIND,
    "このゆびとまれ": FOLLOW_ME,
    "さいみんじゅつ": HYPNOSIS,
    "じこあんじ":    RECOVER,
    "じしん":        EARTHQUAKE,
    "じばく":        SELF_DESTRUCT,
    "たきのぼり":    WATERFALL,
    "だいばくはつ":  EXPLOSION,
    "ちょうはつ":    TAUNT,
    "でんじは":      THUNDER_WAVE,
    "ねこだまし":    FAKE_OUT,
    "ねっぷう":      HEAT_WAVE,
    "はらだいこ":    BELLY_DRUM,
    "ふいうち":      SUCKER_PUNCH,
    "ほのおのパンチ": FIRE_PUNCH,
    "まもる":        PROTECT,
    "みがわり":      SUBSTITUTE,
    "りゅうせいぐん": DRACO_METEOR,
    "クロスチョップ": CROSS_CHOP,
    "コメットパンチ": COMET_PUNCH,
    "サイコキネシス": PSYCHIC,
    "ジャイロボール": GYRO_BALL,
    "ダークホール":  DARK_VOID,
    "トリックルーム": TRICK_ROOM,
    "ハイドロポンプ": HYDRO_PUMP,
    "バレットパンチ": BULLET_PUNCH,
}

func StringToMoveName(s string) (MoveName, error) {
	if moveName, ok := STRING_TO_MOVE_NAME[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_MOVE_NAME に含まれていない為、MoveNameに変換出来ません。", s)
		return moveName, fmt.Errorf(msg)
	} else {
		return moveName, nil
	}
}

func StringsToMoveNames(ss []string) (MoveNames, error) {
	return fn.MapWithError[MoveNames](ss, StringToMoveName)
}

var STRING_TO_MOVE_CATEGORY = map[string]MoveCategory{
	"物理":PHYSICS,
	"特殊":SPECIAL,
	"変化":STATUS,
}

func StringToMoveCategory(s string) (MoveCategory, error) {
	if category, ok := STRING_TO_MOVE_CATEGORY[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_MOVE_CATEGORY に含まれていない為、MoveCategoryに変換出来ません。", s)
		return category, fmt.Errorf(msg)
	} else {
		return category, nil
	}
}

var STRING_TO_MOVE_TARGET = map[string]MoveTarget{
	"通常":NORMAL_TARGET,
	"相手2体":OPPONENT_TWO_TARGET,
	"自分":SELF_TARGET,
	"自分以外":OTHERS_TARGET,
	"全体":ALL_TARGET,
	"相手ランダム1体":OPPONENT_RANDOM_ONE_TARGET,
}

func StringToMoveTarget(s string) (MoveTarget, error) {
	if target, ok := STRING_TO_MOVE_TARGET[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_MOVE_TARGETに含まれていない為、MoveTargetに変換出来ません。", s)
		return target, fmt.Errorf(msg)
	} else {
		return target, nil
	}
}