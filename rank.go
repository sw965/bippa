package bippa

import (
	omwmath "github.com/sw965/omw/math"
)

type Rank int

const (
	MIN_RANK = Rank(-6)
	MAX_RANK = Rank(6)
)

func (rank Rank) ToBonus() RankBonus {
	if rank >= 0 {
		y := (float64(rank) + 2.0) / 2.0
		return RankBonus(y)
	} else {
		abs := float64(rank) * -1
		y := 2.0 / (abs + 2.0)
		return RankBonus(y)
	}
}

type RankState struct {
	Atk   Rank
	Def   Rank
	SpAtk Rank
	SpDef Rank
	Speed Rank
	Accuracy Rank
	Evasion Rank
}

func (rs1 *RankState) Add(rs2 *RankState) RankState {
	atk := rs1.Atk + rs2.Atk
	def := rs1.Def + rs2.Def
	spAtk := rs1.SpAtk + rs2.SpAtk
	spDef := rs1.SpDef + rs2.SpDef
	speed := rs1.Speed + rs2.Speed
	return RankState{Atk: atk, Def: def, SpAtk: spAtk, SpDef: spDef, Speed: speed}
}

func (rs RankState) Regulate() RankState {
	rs.Atk = omwmath.Min(rs.Atk, MAX_RANK)
	rs.Def = omwmath.Min(rs.Def, MAX_RANK)
	rs.SpAtk = omwmath.Min(rs.SpAtk, MAX_RANK)
	rs.SpDef = omwmath.Min(rs.SpDef, MAX_RANK)
	rs.Speed = omwmath.Min(rs.Speed, MAX_RANK)
	rs.Accuracy = omwmath.Min(rs.Accuracy, MAX_RANK)

	rs.Atk = omwmath.Max(rs.Atk, MIN_RANK)
	rs.Def = omwmath.Max(rs.Def, MIN_RANK)
	rs.SpAtk = omwmath.Max(rs.SpAtk, MIN_RANK)
	rs.SpDef = omwmath.Max(rs.SpDef, MIN_RANK)
	rs.Speed = omwmath.Max(rs.Speed, MIN_RANK)
	rs.Accuracy = omwmath.Max(rs.Accuracy, MIN_RANK)
	return rs
}

func (rs *RankState) ContainsDown() bool {
	if rs.Atk < 0 {
		return true
	}

	if rs.Def < 0 {
		return true
	}

	if rs.SpAtk < 0 {
		return true
	}

	if rs.SpDef < 0 {
		return true
	}
	return rs.Speed < 0
}

func (rs RankState) ResetDown() RankState {
	rs.Atk = omwmath.Max(rs.Atk, 0)
	rs.Def = omwmath.Max(rs.Def, 0)
	rs.SpAtk = omwmath.Max(rs.SpAtk, 0)
	rs.SpDef = omwmath.Max(rs.SpDef, 0)
	rs.Speed = omwmath.Max(rs.Speed, 0)
	return rs
}

func (rs *RankState) AccuracyRankBonus() RankBonus {
	a := float64(rs.Accuracy - rs.Evasion)
	if a > 6 {
		a = 6
	} else if a < 6 {
		a = -6
	}

	var y float64
	if a > 0 {
		y = (3.0+a)/3.0
	} else {
		y = 3.0/(3.0-a)
	}
	return RankBonus(y)
}

type RankBonus float64