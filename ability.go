package bippa

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

func (a Ability) ToString() string {
	return ABILITY_TO_STRING[a]
}

type Abilities []Ability

func (as Abilities) ToStrings() []string {
	ss := make([]string, len(as))
	for i, a := range as {
		ss[i] = a.ToString()
	}
	return ss
}