package bippa

import (
	omwjson "github.com/sw965/omw/json"
)

type NatureData struct {
	AtkBonus NatureBonus
	DefBonus NatureBonus
	SpAtkBonus NatureBonus
	SpDefBonus NatureBonus
	SpeedBonus NatureBonus
}

type Naturedex map[Nature]*NatureData

var NATUREDEX = func() Naturedex {
	e, err := omwjson.Load[EasyReadNaturedex](NATUREDEX_PATH)
	if err != nil {
		panic(err)
	}
	d, err := e.From()
	if err != nil {
		panic(err)
	}
	return d
}()

func (n Naturedex) ToEasyRead() EasyReadNaturedex {
	e := EasyReadNaturedex{}
	for k, v := range n {
		e[k.ToString()] = *v
	}
	return e
}

type Nature int

const (
    EMPTY_NATURE Nature = iota
    LONELY               // さみしがり
    ADAMANT              // いじっぱり
    NAUGHTY              // やんちゃ
    BRAVE                // ゆうかん

    BOLD                 // ずぶとい
    IMPISH               // わんぱく
    LAX                  // のうてんき
    RELAXED              // のんき

    MODEST               // ひかえめ
    MILD                 // おっとり
    RASH                 // うっかりや
    QUIET                // れいせい

    CALM                 // おだやか
    GENTLE               // おとなしい
    CAREFUL              // しんちょう
    SASSY                // なまいき

    TIMID                // おくびょう
    HASTY                // せっかち
    JOLLY                // ようき
    NAIVE                // むじゃき

    BASHFUL              // てれや
    HARDY                // がんばりや
    DOCILE               // すなお
    QUIRKY               // きまぐれ
    SERIOUS              // まじめ
)

func (n Nature) ToString() string {
    return NATURE_TO_STRING[n]
}

type Natures []Nature

var ALL_NATURES = func() Natures {
	ss, err := omwjson.Load[[]string](ALL_NATURES_PATH)
	if err != nil {
		panic(err)
	}

	ns, err := StringsToNatures(ss)
	if err != nil {
		panic(err)
	}
	return ns
}()

func (ns Natures) ToStrings() []string {
    ret := make([]string, len(ns))
    for i, n := range ns {
        ret[i] = n.ToString()
    } 
    return ret
}

type NatureBonus float64

const (
	GOOD_NATURE_BONUS = 1.1
	NEUTRAL_NATURE_BONUS = 1.0
	BAD_NATURE_BONUS = 0.9
)