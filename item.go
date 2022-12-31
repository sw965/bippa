package bippa

import (
	"math/rand"
)

type Item string

const (
	EMPTY_ITEM = Item("なし")
)

func (item Item) IsValid() bool {
	for _, iItem := range ALL_ITEMS {
		if iItem == item {
			return true
		}
	}
	return item == EMPTY_ITEM
}

func (item Item) IsChoice() bool {
	return item == "こだわりハチマキ" || item == "こだわりメガネ" || item == "こだわりスカーフ"
}

type Items []Item

func (items Items) In(item Item) bool {
	for _, iItem := range items {
		if iItem == item {
			return true
		}
	}
	return false
}

func (items Items) InAll(item ...Item) bool {
	for _, iItem := range item {
		if !items.In(iItem) {
			return false
		}
	}
	return true
}

func (items Items) RandomChoice(random *rand.Rand) Item {
	index := random.Intn(len(items))
	return items[index]
}

type ItemWithFloat64 map[Item]float64

func (itemWithFloat64 ItemWithFloat64) KeysAndValues() (Items, []float64) {
	length := len(itemWithFloat64)
	keys := make(Items, 0, length)
	values := make([]float64, 0, length)
	for k, v := range itemWithFloat64 {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}