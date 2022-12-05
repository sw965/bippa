package bippa

type Nature string

type Natures []Nature

func (natures Natures) In(nature Nature) bool {
	for _, iNature := range natures {
		if iNature == nature {
			return true
		}
	}
	return false
}

func (natures Natures) InAll(nature ...Nature) bool {
	for _, iNature := range nature {
		if !natures.In(iNature) {
			return false
		}
	}
	return true
}

type NatureWithTier map[Nature]Tier

func (natureWithTier NatureWithTier) KeysAndValues() (Natures, Tiers) {
	length := len(natureWithTier)
	keys := make(Natures, 0, length)
	values := make(Tiers, 0, length)

	for k, v := range natureWithTier {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

type NatureBonus float64

const (
	NO_NATURE_BONUS   = NatureBonus(1.0)
	UP_NATURE_BONUS   = NatureBonus(1.1)
	DOWN_NATURE_BONUS = NatureBonus(0.9)
)
