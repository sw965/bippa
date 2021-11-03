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

func (rank *Rank) TotalRise() Rank_ {
	rank_ := Rank_(0)

	if rank.Atk > 0 {
		rank_ += rank.Atk
	}

	if rank.Def > 0 {
		rank_ += rank.Def
	}

	if rank.SpAtk > 0 {
		rank_ += rank.SpAtk
	}

	if rank.SpDef > 0 {
		rank_ += rank.SpDef
	}

	if rank.Speed > 0 {
		rank_ += rank.Speed
	}
	return rank_
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

type RankBonus float64

var RANK__TO_RANK_BONUS = map[Rank_]RankBonus{
	-6: 2.0 / 8.0, -5: 2.0 / 7.0, -4: 2.0 / 6.0, -3: 2.0 / 5.0, -2: 2.0 / 4.0, -1: 2.0 / 3.0,
	0: 2.0 / 2.0,
	1: 3.0 / 2.0, 2: 4.0 / 2.0, 3: 5.0 / 2.0, 4: 6.0 / 2.0, 5: 7.0 / 2.0, 6: 8.0 / 2.0,
}
