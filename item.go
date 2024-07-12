package bippa

import (
	omwmaps "github.com/sw965/omw/maps"
)

type Item int

const (
	EMPTY_ITEM Item = iota
    IAPAPA_BERRY // イバンのみ
    SITRUS_BERRY // オボンのみ
    CHESTO_BERRY // カゴのみ
    FOCUS_SASH   // きあいのタスキ
    WACAN_BERRY  // ソクノのみ
    LUM_BERRY    // ラムのみ
)

var STRING_TO_ITEM = map[string]Item{
	"イバンのみ":IAPAPA_BERRY,
	"オボンのみ":SITRUS_BERRY,
	"カゴのみ":CHESTO_BERRY,
	"きあいのタスキ":FOCUS_SASH,
	"ソクノのみ":WACAN_BERRY,
	"ラムのみ":LUM_BERRY,
}

var ITEM_TO_STRING = omwmaps.Invert[map[Item]string](STRING_TO_ITEM)

type Items []Item