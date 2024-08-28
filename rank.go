package bippa

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

func (r RankStat) Add(v *RankStat) RankStat {
	r.Atk += v.Atk
	r.Def += v.Def
	r.SpAtk += v.SpAtk
	r.SpDef += v.SpDef
	r.Speed += v.Speed
	return r
}

func (r RankStat) DownToZero() RankStat {
	if r.Atk < 0 {
		r.Atk = 0
	}

	if r.Def < 0 {
		r.Def = 0
	}

	if r.SpAtk < 0 {
		r.SpAtk = 0
	}

	if r.SpDef < 0 {
		r.SpDef = 0
	}

	if r.Speed < 0 {
		r.Speed = 0
	}

	return r
}

func (r RankStat) AdjustFluctuation(v *RankStat) RankStat {
	a := v.Add(&r)
	if a.Atk > MAX_RANK {
		r.Atk -= a.Atk - MAX_RANK
	} else if a.Atk < MIN_RANK {
		r.Atk -= a.Atk - MIN_RANK
	}

	if a.Def > MAX_RANK {
		r.Def -= a.Def - MAX_RANK
	} else if a.Def < MIN_RANK {
		r.Def -= a.Def - MIN_RANK
	}

	if a.SpAtk > MAX_RANK {
		r.SpAtk -= a.SpAtk - MAX_RANK
	} else if a.SpAtk < MIN_RANK {
		r.SpAtk -= a.SpAtk - MIN_RANK
	}

	if a.SpDef > MAX_RANK {
		r.SpDef -= a.SpDef - MAX_RANK
	} else if a.SpDef < MIN_RANK {
		r.SpDef -= a.SpDef - MIN_RANK
	}

	if a.Speed > MAX_RANK {
		r.Speed -= a.Speed - MAX_RANK
	} else if a.Speed < MIN_RANK {
		r.Speed -= a.Speed - MIN_RANK
	}
	return r
}