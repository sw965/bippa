package bippa

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