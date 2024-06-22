package single

import (
	"math/rand"
	omwrand "github.com/sw965/omw/math/rand"
	"github.com/sw965/bippa/battle/dmgtools"
)

type Observer func(*Battle, Step)

func EmptyObserver(_ *Battle, _ Step) {}

type Context struct {
	DamageRandBonuses dmgtools.RandBonuses
	Rand *rand.Rand
	Observer Observer
}

func (c *Context) DamageRandBonus() dmgtools.RandBonus {
	return omwrand.Choice(c.DamageRandBonuses, c.Rand)
}