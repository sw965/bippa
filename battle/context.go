package battle

import (
	"math/rand"
	omwrand "github.com/sw965/omw/math/rand"
)

type Observer func(*Manager)

func EmptyObserver(_ *Manager) {}

type Context struct {
	Rand *rand.Rand
	DamageRandBonuses []float64
	Observer Observer
}

var GlobalContext = Context{
	Rand:omwrand.NewMt19937(),
	DamageRandBonuses:DAMAGE_RAND_BONUSES,
	Observer:EmptyObserver,
}

func (c *Context) GetDamageRandBonus() float64 {
	return omwrand.Choice(c.DamageRandBonuses, c.Rand)
}