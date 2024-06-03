package bippa

import (
    omwmaps "github.com/sw965/omw/maps"
)

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

var NATURE_TO_STRING = omwmaps.Invert[map[Nature]string](STRING_TO_NATURE)

type NatureBonus float64

const (
	GOOD_NATURE_BONUS = 1.1
	NEUTRAL_NATURE_BONUS = 1.0
	BAD_NATURE_BONUS = 0.9
)