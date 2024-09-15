package bippa

import (
	omwjson "github.com/sw965/omw/json"
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

func (i Item) ToString() string {
	return ITEM_TO_STRING[i]
}

type Items []Item

var ALL_ITEMS = func() Items {
	ss, err := omwjson.Load[[]string](ALL_ITEMS_PATH)
	if err != nil {
		panic(err)
	}

	is, err := StringsToItems(ss)
	if err != nil {
		panic(err)
	}
	return is
}()

func (is Items) ToStrings() []string {
	ss := make([]string, len(is))
	for i, item := range is {
		ss[i] = item.ToString()
	}
	return ss
}