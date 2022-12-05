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

type ItemWithTier map[Item]Tier

func (itemWithTier ItemWithTier) KeysAndValues() (Items, Tiers) {
	length := len(itemWithTier)
	keys := make(Items, 0, length)
	values := make(Tiers, 0, length)

	for k, v := range itemWithTier {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}
