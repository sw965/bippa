package bippa

type Rank int

const (
	MIN_RANK = Rank(-6)
	MAX_RANK = Rank(6)
)

type RankBonus float64

var RANK_TO_RANK_BONUS = map[Rank]RankBonus{
	-6: 2.0 / 8.0, -5: 2.0 / 7.0, -4: 2.0 / 6.0, -3: 2.0 / 5.0, -2: 2.0 / 4.0, -1: 2.0 / 3.0,
	0: 2.0 / 2.0,
	1: 3.0 / 2.0, 2: 4.0 / 2.0, 3: 5.0 / 2.0, 4: 6.0 / 2.0, 5: 7.0 / 2.0, 6: 8.0 / 2.0,
}

type RankState struct {
	Atk   Rank
	Def   Rank
	SpAtk Rank
	SpDef Rank
	Speed Rank
}

func (rankState RankState) Reset() RankState {
	return RankState{}
}

func (rankState *RankState) TotalRise() Rank {
	rank := Rank(0)

	if rankState.Atk > 0 {
		rank += rankState.Atk
	}

	if rankState.Def > 0 {
		rank += rankState.Def
	}

	if rankState.SpAtk > 0 {
		rank += rankState.SpAtk
	}

	if rankState.SpDef > 0 {
		rank += rankState.SpDef
	}

	if rankState.Speed > 0 {
		rank += rankState.Speed
	}
	return rank
}

func (rankState RankState) Regulate() RankState {
	if rankState.Atk > MAX_RANK {
		rankState.Atk = MAX_RANK
	}

	if rankState.Def > MAX_RANK {
		rankState.Def = MAX_RANK
	}

	if rankState.SpAtk > MAX_RANK {
		rankState.SpAtk = MAX_RANK
	}

	if rankState.SpDef > MAX_RANK {
		rankState.SpDef = MAX_RANK
	}

	if rankState.Speed > MAX_RANK {
		rankState.Speed = MAX_RANK
	}

	if rankState.Atk < MIN_RANK {
		rankState.Atk = MIN_RANK
	}

	if rankState.Def < MIN_RANK {
		rankState.Def = MIN_RANK
	}

	if rankState.SpAtk < MIN_RANK {
		rankState.SpAtk = MIN_RANK
	}

	if rankState.SpDef < MIN_RANK {
		rankState.SpDef = MIN_RANK
	}

	if rankState.Speed < MIN_RANK {
		rankState.Speed = MIN_RANK
	}
	return rankState
}
