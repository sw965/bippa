package bippa

import (
	omwmaps "github.com/sw965/omw/maps"
)

type Ability int

const (
    TORRENT Ability = iota // げきりゅう
    IMMUNITY               // めんえき
    THICK_FAT              // あついしぼう
    INTIMIDATE             // いかく
    LEVITATE               // ふゆう
	HEATPROOF              // たいねつ
    OWN_TEMPO              // マイペース
    TECHNICIAN             // テクニシャン
    ANTICIPATION           // きけんよち
    DRY_SKIN               // かんそうはだ
    CLEAR_BODY             // クリアボディ
)

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

var ABILITY_TO_STRING = omwmaps.Invert[map[Ability]string](STRING_TO_ABILITY)

func (a Ability) ToString() string {
	return ABILITY_TO_STRING[a]
}

type Abilities []Ability

func (as Abilities) ToStrings() []string {
	ret := make([]string, len(as))
	for i, a := range as {
		ret[i] = a.ToString()
	}
	return ret
}