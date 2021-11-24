package bippa

type Rank_ int

const (
	MIN_RANK_ = Rank_(-6)
	MAX_RANK_ = Rank_(6)
)

type Rank struct {
	Atk   Rank_
	Def   Rank_
	SpAtk Rank_
	SpDef Rank_
	Speed Rank_
}

func (rank1 *Rank) Add(rank2 *Rank) Rank {
	atk := rank1.Atk + rank2.Atk
	def := rank1.Def + rank2.Def
	spAtk := rank1.SpAtk + rank2.SpAtk
	spDef := rank1.SpDef + rank2.SpDef
	speed := rank1.Speed + rank2.Speed
	return Rank{Atk:atk, Def:def, SpAtk:spAtk, SpDef:spDef, Speed:speed}
}

func (rank Rank) Regulate() Rank {
	if rank.Atk > MAX_RANK_ {
		rank.Atk = MAX_RANK_
	}

	if rank.Def > MAX_RANK_ {
		rank.Def = MAX_RANK_
	}

	if rank.SpAtk > MAX_RANK_ {
		rank.SpAtk = MAX_RANK_
	}

	if rank.SpDef > MAX_RANK_ {
		rank.SpDef = MAX_RANK_
	}

	if rank.Speed > MAX_RANK_ {
		rank.Speed = MAX_RANK_
	}

	if rank.Atk < MIN_RANK_ {
		rank.Atk = MIN_RANK_
	}

	if rank.Def < MIN_RANK_ {
		rank.Def = MIN_RANK_
	}

	if rank.SpAtk < MIN_RANK_ {
		rank.SpAtk = MIN_RANK_
	}

	if rank.SpDef < MIN_RANK_ {
		rank.SpDef = MIN_RANK_
	}

	if rank.Speed < MIN_RANK_ {
		rank.Speed = MIN_RANK_
	}
	return rank
}

//アシストパワー・つけあがる用
func (rank *Rank) TotalRise() Rank_ {
	result := Rank_(0)

	if rank.Atk > 0 {
		result += rank.Atk
	}

	if rank.Def > 0 {
		result += rank.Def
	}

	if rank.SpAtk > 0 {
		result += rank.SpAtk
	}

	if rank.SpDef > 0 {
		result += rank.SpDef
	}

	if rank.Speed > 0 {
		result += rank.Speed
	}
	return result
}

//しろいハーブ用
func (rank *Rank) InDown() bool {
	if rank.Atk < 0 {
		return true
	}

	if rank.Def < 0 {
		return true
	}

	if rank.SpAtk < 0 {
		return true
	}

	if rank.SpDef < 0 {
		return true
	}

	if rank.Speed < 0 {
		return true
	}
	return false
}

func (rank *Rank) ResetDown() Rank {
	result := Rank{}
	if rank.Atk < 0 {
		result.Atk = 0
	} else {
		result.Atk = rank.Atk
	}

	if rank.Def < 0 {
		result.Def = 0
	} else {
		result.Def = rank.Def
	}

	if rank.SpAtk < 0 {
		result.SpAtk = 0
	} else {
		result.SpAtk = rank.SpAtk
	}

	if rank.SpDef < 0 {
		result.SpDef = 0
	} else {
		result.SpDef = rank.SpDef
	}

	if rank.Speed < 0 {
		result.Speed = 0
	} else {
		result.Speed = rank.Speed
	}
	
	return result
}

type RankBonus float64

var RANK__TO_RANK_BONUS = map[Rank_]RankBonus{
	-6: 2.0 / 8.0, -5: 2.0 / 7.0, -4: 2.0 / 6.0, -3: 2.0 / 5.0, -2: 2.0 / 4.0, -1: 2.0 / 3.0,
	0: 2.0 / 2.0,
	1: 3.0 / 2.0, 2: 4.0 / 2.0, 3: 5.0 / 2.0, 4: 6.0 / 2.0, 5: 7.0 / 2.0, 6: 8.0 / 2.0,
}
