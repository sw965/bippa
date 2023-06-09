package bippa

type Nature string
type Natures []Nature
type NatureBonus float64

const (
	NO_NATURE_BONUS   = NatureBonus(1.0)
	UP_NATURE_BONUS   = NatureBonus(1.1)
	DOWN_NATURE_BONUS = NatureBonus(0.9)
)

type NatureBonuses []NatureBonus

var ALL_NATURE_BONUSES = NatureBonuses{NO_NATURE_BONUS, UP_NATURE_BONUS, DOWN_NATURE_BONUS}