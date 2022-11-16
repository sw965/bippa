package bippa

type Type string

const (
	NORMAL   = Type("ノーマル")
	FIRE     = Type("ほのお")
	WATER    = Type("みず")
	GRASS    = Type("くさ")
	ELECTRIC = Type("でんき")
	ICE      = Type("こおり")
	FIGHTING = Type("かくとう")
	POISON   = Type("どく")
	GROUND   = Type("じめん")
	FLYING   = Type("ひこう")
	PSYCHIC  = Type("エスパー")
	BUG      = Type("むし")
	ROCK     = Type("いわ")
	GHOST    = Type("ゴースト")
	DRAGON   = Type("ドラゴン")
	DARK     = Type("あく")
	STEEL    = Type("はがね")
	FAIRY    = Type("フェアリー")
)

type Types []Type

func (types Types) In(type_ Type) bool {
	for _, iType := range types {
		if iType == type_ {
			return true
		}
	}
	return false
}