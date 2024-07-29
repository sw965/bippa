package bippa

import (
	omwmath "github.com/sw965/omw/math"
)

// https://wiki.xn--rckteqa2e.com/wiki/%E3%83%A9%E3%83%B3%E3%82%AF%E8%A3%9C%E6%AD%A3#%E6%80%A5%E6%89%80%E3%83%A9%E3%83%B3%E3%82%AF
type CriticalRank int

// https://wiki.xn--rckteqa2e.com/wiki/%E3%83%A9%E3%83%B3%E3%82%AF%E8%A3%9C%E6%AD%A3
type Rank int

func (r Rank) Bonus() float64 {
	if r > 0 {
		return (2.0 + float64(r)) / 2.0
	} else {
		return 2.0 / (2.0 + float64(-r))
	}
}

const (
	MIN_RANK = -6
	MAX_RANK = 6
)

type RankStat struct {
	Atk Rank
	Def Rank
	SpAtk Rank
	SpDef Rank
	Speed Rank
}

func (r RankStat) Clone() RankStat {
	return r
}

func (r *RankStat) Fluctuation(v *RankStat) {
	if v.Atk > 0 {
		r.Atk = omwmath.Min(r.Atk+v.Atk, MAX_RANK)
	} else {
		r.Atk = omwmath.Max(r.Atk+v.Atk, MIN_RANK)
	}

	if v.Def > 0 {
		r.Def = omwmath.Min(r.Def+v.Def, MAX_RANK)
	} else {
		r.Def = omwmath.Max(r.Def+v.Def, MIN_RANK)
	}

	if v.SpAtk > 0 {
		r.SpAtk = omwmath.Min(r.SpAtk+v.SpAtk, MAX_RANK)
	} else {
		r.SpAtk = omwmath.Max(r.SpAtk+v.SpAtk, MIN_RANK)
	}

	if v.SpDef > 0 {
		r.SpDef = omwmath.Min(r.SpDef+v.SpDef, MAX_RANK)
	} else {
		r.SpDef = omwmath.Max(r.SpDef+v.SpDef, MIN_RANK)
	}

	if v.Speed > 0 {
		r.Speed = omwmath.Min(r.Speed+v.Speed, MAX_RANK)
	} else {
		r.Speed = omwmath.Max(r.Speed+v.Speed, MIN_RANK)
	}
}