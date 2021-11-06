package bippa

type Effort_ int

const (
	MIN_EFFORT_     = Effort_(0)
	MAX_EFFORT_     = Effort_(252)
	MAX_SUM_EFFORT_ = Effort_(510)
)

func (effort_ Effort_) IsValid() bool {
	return effort_ >= MIN_EFFORT_ && effort_ <= MAX_EFFORT_
}

type Effort struct {
	HP    Effort_
	Atk   Effort_
	Def   Effort_
	SpAtk Effort_
	SpDef Effort_
	Speed Effort_
}

var (
	HA252_B4 = Effort{HP: 252, Atk: 252, Def: 4, SpAtk: 0, SpDef: 0, Speed: 0}
	HA252_C4 = Effort{HP: 252, Atk: 252, Def: 0, SpAtk: 4, SpDef: 0, Speed: 0}
	HA252_D4 = Effort{HP: 252, Atk: 252, Def: 0, SpAtk: 0, SpDef: 4, Speed: 0}
	HA252_S4 = Effort{HP: 252, Atk: 252, Def: 0, SpAtk: 0, SpDef: 0, Speed: 4}

	HB252_A4 = Effort{HP: 252, Atk: 4, Def: 252, SpAtk: 0, SpDef: 0, Speed: 0}
	HB252_C4 = Effort{HP: 252, Atk: 0, Def: 252, SpAtk: 4, SpDef: 0, Speed :0}
	HB252_D4 = Effort{HP: 252, Atk: 0, Def: 252, SpAtk: 0, SpDef: 4, Speed: 0}
	HB252_S4 = Effort{HP: 252, Atk: 0, Def: 252, SpAtk: 0, SpDef: 0, Speed: 4}

	HC252_S4 = Effort{HP: 252, Atk: 0, Def: 0, SpAtk: 252, SpDef: 0, Speed: 4}

	HS252_C4 = Effort{HP: 252, Atk: 0, Def: 0, SpAtk: 4, SpDef: 0, Speed: 252}

	HD252_B4 = Effort{HP: 252, Atk: 0, Def: 4, SpAtk: 0, SpDef: 252, Speed: 0}
	HD252_S4 = Effort{HP: 252, Atk: 0, Def: 0, SpAtk: 0, SpDef: 252, Speed: 4}

	AD252_H4 = Effort{HP: 4, Atk: 252, Def: 0, SpAtk: 0, SpDef: 252, Speed: 0}

	AS252_H4 = Effort{HP: 4, Atk: 252, Def: 0, SpAtk: 0, SpDef: 0, Speed: 252}

	BC252_H4 = Effort{HP: 4, Atk: 0, Def: 252, SpAtk: 252, SpDef: 0, Speed: 0}
	BC252_S4 = Effort{HP: 0, Atk: 0, Def: 252, SpAtk: 252, SpDef: 0, Speed: 4}

	CD252_H4 = Effort{HP: 4, Atk: 0, Def: 0, SpAtk: 252, SpDef: 252, Speed: 0}
	CS252_H4 = Effort{HP: 4, Atk: 0, Def: 0, SpAtk: 252, SpDef: 0, Speed: 252}
	CS252_B4 = Effort{HP: 0, Atk: 0, Def: 4, SpAtk: 252, SpDef: 0, Speed: 252}
)

func (effort *Effort) Sum() Effort_ {
	return effort.HP + effort.Atk + effort.Def + effort.SpAtk + effort.SpDef + effort.Speed
}

func (effort *Effort) IsAllValid() bool {
	return effort.HP.IsValid() && effort.Atk.IsValid() && effort.Def.IsValid() &&
		effort.SpAtk.IsValid() && effort.SpDef.IsValid() && effort.Speed.IsValid()
}

func (effort *Effort) IsValidSum() bool {
	return effort.Sum() <= MAX_SUM_EFFORT_
}
