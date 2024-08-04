package battle

import (
	"math/rand"
	omwrand "github.com/sw965/omw/math/rand"
	"github.com/sw965/bippa/battle/dmgtools"
)

type Observer func(*Manager, EventType)

func EmptyObserver(_ *Manager, _ EventType) {}

type Context struct {
	DamageRandBonuses []float64
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

func (c *Context) DamageRandBonus() float64 {
	return omwrand.Choice(c.DamageRandBonuses, c.Rand)
}