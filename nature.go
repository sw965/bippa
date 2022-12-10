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

type NatureBonus float64

const (
	NO_NATURE_BONUS   = NatureBonus(1.0)
	UP_NATURE_BONUS   = NatureBonus(1.1)
	DOWN_NATURE_BONUS = NatureBonus(0.9)
)
