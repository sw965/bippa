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

var STRING_TO_ITEM = map[string]Item{
	"イバンのみ":IAPAPA_BERRY,
	"オボンのみ":SITRUS_BERRY,
	"カゴのみ":CHESTO_BERRY,
	"きあいのタスキ":FOCUS_SASH,
	"ソクノのみ":WACAN_BERRY,
	"ラムのみ":LUM_BERRY,
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