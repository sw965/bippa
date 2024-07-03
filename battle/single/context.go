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

func NewContext(r *rand.Rand) Context {
	return Context{
		DamageRandBonuses:dmgtools.RAND_BONUSES,
		Rand:r,
		Observer:EmptyObserver,
	}
}

func (c *Context) SetDamageRandBonuses(dmgRandBonuses ...dmgtools.RandBonus) {
	c.DamageRandBonuses = dmgRandBonuses
}

func (c *Context) DamageRandBonus() dmgtools.RandBonus {
	return omwrand.Choice(c.DamageRandBonuses, c.Rand)
}