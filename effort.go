package bippa

type Effort int

const (
	MIN_EFFORT     = Effort(0)
	MAX_EFFORT     = Effort(252)
	MAX_SUM_EFFORT = Effort(510)
)

func (effort Effort) IsValid() bool {
	return effort >= MIN_EFFORT && effort <= MAX_EFFORT
}

type EffortState struct {
	HP    Effort
	Atk   Effort
	Def   Effort
	SpAtk Effort
	SpDef Effort
	Speed Effort
}

var (
	HA252_B4 = EffortState{HP: 252, Atk: 252, Def: 4, SpAtk: 0, SpDef: 0, Speed: 0}
	HA252_C4 = EffortState{HP: 252, Atk: 252, Def: 0, SpAtk: 4, SpDef: 0, Speed: 0}
	HA252_D4 = EffortState{HP: 252, Atk: 252, Def: 0, SpAtk: 0, SpDef: 4, Speed: 0}
	HA252_S4 = EffortState{HP: 252, Atk: 252, Def: 0, SpAtk: 0, SpDef: 0, Speed: 4}

	HB252_A4 = EffortState{HP: 252, Atk: 4, Def: 252, SpAtk: 0, SpDef: 0, Speed: 0}
	HB252_D4 = EffortState{HP: 252, Atk: 0, Def: 252, SpAtk: 0, SpDef: 4, Speed: 0}
	HB252_S4 = EffortState{HP: 252, Atk: 0, Def: 252, SpAtk: 0, SpDef: 0, Speed: 4}

	HC252_S4 = EffortState{HP: 252, Atk: 0, Def: 0, SpAtk: 252, SpDef: 0, Speed: 4}

	HS252_C4 = EffortState{HP: 252, Atk: 0, Def: 0, SpAtk: 4, SpDef: 0, Speed: 252}

	HD252_B4 = EffortState{HP: 252, Atk: 0, Def: 4, SpAtk: 0, SpDef: 252, Speed: 0}
	HD252_S4 = EffortState{HP: 252, Atk: 0, Def: 0, SpAtk: 0, SpDef: 252, Speed: 4}

	AD252_H4 = EffortState{HP: 4, Atk: 252, Def: 0, SpAtk: 0, SpDef: 252, Speed: 0}

	AS252_H4 = EffortState{HP: 4, Atk: 252, Def: 0, SpAtk: 0, SpDef: 0, Speed: 252}

	BC252_H4 = EffortState{HP: 4, Atk: 0, Def: 252, SpAtk: 252, SpDef: 0, Speed: 0}
	BC252_S4 = EffortState{HP: 0, Atk: 0, Def: 252, SpAtk: 252, SpDef: 0, Speed: 4}

	CD252_H4 = EffortState{HP: 4, Atk: 0, Def: 0, SpAtk: 252, SpDef: 252, Speed: 0}
	CS252_H4 = EffortState{HP: 4, Atk: 0, Def: 0, SpAtk: 252, SpDef: 0, Speed: 252}
	CS252_B4 = EffortState{HP: 0, Atk: 0, Def: 4, SpAtk: 252, SpDef: 0, Speed: 252}
)

func (effortState *EffortState) Sum() Effort {
	return effortState.HP + effortState.Atk + effortState.Def + effortState.SpAtk + effortState.SpDef + effortState.Speed
}

func (effortState *EffortState) IsAllValid() bool {
	return effortState.HP.IsValid() && effortState.Atk.IsValid() && effortState.Def.IsValid() &&
		effortState.SpAtk.IsValid() && effortState.SpDef.IsValid() && effortState.Speed.IsValid()
}

func (effortState *EffortState) IsValidSum() bool {
	return effortState.Sum() <= MAX_SUM_EFFORT
}
