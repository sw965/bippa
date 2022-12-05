package bippa

import (
	"github.com/sw965/omw"
)

type EffortVal int

var (
	MIN_EFFORT_VAL     = EffortVal(0)
	MAX_EFFORT_VAL     = EffortVal(252)
	MAX_SUM_EFFORT_VAL = EffortVal(510)
)

type EffortVals []EffortVal

var ALL_EFFORT_VALS = func() EffortVals {
	length := int(MAX_EFFORT_VAL + 1)
	result := make(EffortVals, length)
	for i := 0; i < length; i++ {
		result[i] = EffortVal(i)
	}
	return result
}()

type Effort struct {
	HP    EffortVal
	Atk   EffortVal
	Def   EffortVal
	SpAtk EffortVal
	SpDef EffortVal
	Speed EffortVal
}

func (effort *Effort) Sum() EffortVal {
	return effort.HP + effort.Atk + effort.Def + effort.SpAtk + effort.SpDef + effort.Speed
}

func (effort *Effort) Remaining() EffortVal {
	return MAX_SUM_EFFORT_VAL - effort.Sum()
}

func (effort Effort) SetUpperLimitHP() Effort {
	remaining := effort.Remaining()
	hp := omw.MinInt(int(MAX_EFFORT_VAL-effort.HP), int(remaining))
	effort.HP = EffortVal(hp)
	return effort
}

func (effort Effort) SetUpperLimitAtk() Effort {
	remaining := effort.Remaining()
	atk := omw.MinInt(int(MAX_EFFORT_VAL-effort.Atk), int(remaining))
	effort.Atk = EffortVal(atk)
	return effort
}

func (effort Effort) SetUpperLimitDef() Effort {
	remaining := effort.Remaining()
	def := omw.MinInt(int(MAX_EFFORT_VAL-effort.Def), int(remaining))
	effort.Def = EffortVal(def)
	return effort
}

func (effort Effort) SetUpperLimitSpAtk() Effort {
	remaining := effort.Remaining()
	spAtk := omw.MinInt(int(MAX_EFFORT_VAL-effort.SpAtk), int(remaining))
	effort.SpAtk = EffortVal(spAtk)
	return effort
}

func (effort Effort) SetUpperLimitSpDef() Effort {
	remaining := effort.Remaining()
	spDef := omw.MinInt(int(MAX_EFFORT_VAL-effort.SpDef), int(remaining))
	effort.SpDef = EffortVal(spDef)
	return effort
}

func (effort Effort) SetUpperLimitSpeed() Effort {
	remaining := effort.Remaining()
	speed := omw.MinInt(int(MAX_EFFORT_VAL-effort.Speed), int(remaining))
	effort.Speed = EffortVal(speed)
	return effort
}

type Efforts []Effort

type EffortWithTier map[Effort]Tier

func (effortWithTier EffortWithTier) KeysAndValues() (Efforts, Tiers) {
	length := len(effortWithTier)
	keys := make(Efforts, 0, length)
	values := make(Tiers, 0, length)
	for k, v := range effortWithTier {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}
