package bippa

import (
	"fmt"
	omwjson "github.com/sw965/omw/json"
	"github.com/sw965/omw/fn"
)

type item string

type Item int

const (
	NO_ITEM Item = iota
	LIFE_ORB //いのちのたま
	SITRUS_BERRY //オボンのみ
	BLACK_SLUDGE //くろいヘドロ
	CHOICE_SCARF //こだわりスカーフ
	CHOICE_BAND //こだわりハチマキ
	CHOICE_SPECS //こだわりメガネ
	ROCKY_HELMET //ゴツゴツメット
	WHITE_HERB //しろいハーブ
	LEFTOVERS //たべのこし
	ASSAULT_VEST //とつげきチョッキ
)

func NewItem(s string) (Item, error) {
	switch s {
		case "なし":
			return NO_ITEM, nil
		case "いのちのたま":
			return LIFE_ORB, nil
		case "オボンのみ":
			return SITRUS_BERRY, nil
		case "くろいヘドロ":
			return BLACK_SLUDGE, nil
		case "こだわりスカーフ":
			return CHOICE_SPECS, nil
		case "こだわりハチマキ":
			return CHOICE_BAND, nil
		case "こだわりメガネ":
			return CHOICE_SPECS, nil
		case "ゴツゴツメット":
			return ROCKY_HELMET, nil
		case "しろいハーブ":
			return WHITE_HERB, nil
		case "たべのこし":
			return LEFTOVERS, nil
		case "とつげきチョッキ":
			return ASSAULT_VEST, nil
		default:
			return -1, fmt.Errorf("不適なitem")
	}
}

func (item Item) IsChoice() bool {
	return item == CHOICE_BAND || item == CHOICE_SPECS || item == CHOICE_SCARF
}

type Items []Item

var ALL_ITEMS = func() Items {
	d, err := omwjson.Load[[]string](ALL_ITEMS_PATH)
	if err != nil {
		panic(err)
	}
	y, err := fn.MapError[Items](d, NewItem)
	if err != nil {
		panic(err)
	}
	return y
}()