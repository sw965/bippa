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

func (r RankStat) DownToZero() RankStat {
	r.Atk = omwmath.Max(r.Atk, 0)
	r.Def = omwmath.Max(r.Def, 0)
	r.SpAtk = omwmath.Max(r.SpAtk, 0)
	r.SpDef = omwmath.Max(r.SpDef, 0)
	r.Speed = omwmath.Max(r.Speed, 0)
	return r
}

func (r RankStat) ContainsNotZero() bool {
	if r.Atk != 0 {
		return true
	}

	if r.Def != 0 {
		return true
	}

	if r.SpAtk != 0 {
		return true
	}

	if r.SpDef != 0 {
		return true
	}

	if r.Speed != 0 {
		return true
	}

	return false
}

func (r *RankStat) Fluctuation(v *RankStat) RankStatFluctuationResult {
	var atk Rank
	if v.Atk > 0 {
		atk = omwmath.Min(v.Atk, MAX_RANK - r.Atk)
	} else {
		atk = omwmath.Max(v.Atk, MIN_RANK - r.Atk)
	}

	var def Rank
	if v.Def > 0 {
		def = omwmath.Min(v.Atk, MAX_RANK - r.Atk)
	} else {
		def = omwmath.Max(v.Atk, MIN_RANK - r.Atk)
	}

	var spAtk Rank
	if v.SpAtk > 0 {
		spAtk = omwmath.Min(v.Atk, MAX_RANK - r.Atk)
	} else {
		spAtk = omwmath.Max(v.Atk, MIN_RANK - r.Atk)
	}

	var spDef Rank
	if v.SpDef > 0 {
		spDef = omwmath.Min(v.Atk, MAX_RANK - r.Atk)
	} else {
		spDef = omwmath.Max(v.Atk, MIN_RANK - r.Atk)
	}

	var speed Rank
	if v.Speed > 0 {
		speed = omwmath.Min(v.Atk, MAX_RANK - r.Atk)
	} else {
		speed = omwmath.Max(v.Atk, MIN_RANK - r.Atk)
	}

	r.Atk += atk
	r.Def += def
	r.SpAtk += spAtk
	r.SpDef += spDef
	r.Speed += speed

	return RankStatFluctuationResult{
		Atk:atk,
		Def:def,
		SpAtk:spAtk,
		SpDef:spDef,
		Speed:speed,
	}
}

func (r RankStat) Sub(v *RankStat) RankStat {
	return r
}

type RankStatFluctuationResult struct {
	Atk Rank
	Def Rank
	SpAtk Rank
	SpDef Rank
	Speed Rank
	IsClearBody bool
}