package bippa

import (
	"fmt"
	"github.com/sw965/omw/fn"
)

type Effort int

const (
	EMPTY_EFFORT = -1
	MIN_EFFORT = 0
	MAX_EFFORT = 252
)

const (
    MIN_SUM_EFFORT Effort = 0
    MAX_SUM_EFFORT Effort = 510
)

type Efforts []Effort

var ALL_EFFORTS = func() Efforts {
    ret := make(Efforts, MAX_EFFORT+1)
    for i := 0; i < int(MAX_EFFORT); i++ {
        ret[i] = Effort(i)
    }
    return ret
}()

var EFFECTIVE_EFFORTS = fn.Filter[Efforts](ALL_EFFORTS, func(e Effort) bool { return e%4 == 0 })

type EffortStat struct {
	HP Effort
	Atk Effort
	Def Effort
	SpAtk Effort
	SpDef Effort
	Speed Effort
}

var (
    EMPTY_EFFORT_STAT = EffortStat{HP:EMPTY_EFFORT, Atk:EMPTY_EFFORT, Def:EMPTY_EFFORT, SpAtk:EMPTY_EFFORT, SpDef:EMPTY_EFFORT, Speed:EMPTY_EFFORT}

    HA252_B4 = EffortStat{HP: MAX_EFFORT, Atk: MAX_EFFORT, Def: 4}
    HA252_C4 = EffortStat{HP: MAX_EFFORT, Atk: MAX_EFFORT, SpAtk: 4}
    HA252_D4 = EffortStat{HP: MAX_EFFORT, Atk: MAX_EFFORT, SpDef: 4}
    HA252_S4 = EffortStat{HP: MAX_EFFORT, Atk: MAX_EFFORT, Speed: 4}

    HB252_A4 = EffortStat{HP: MAX_EFFORT, Def: MAX_EFFORT, Atk: 4}
    HB252_C4 = EffortStat{HP: MAX_EFFORT, Def: MAX_EFFORT, SpAtk: 4}
    HB252_D4 = EffortStat{HP: MAX_EFFORT, Def: MAX_EFFORT, SpDef: 4}
    HB252_S4 = EffortStat{HP: MAX_EFFORT, Def: MAX_EFFORT, Speed: 4}

    HC252_A4 = EffortStat{HP: MAX_EFFORT, SpAtk: MAX_EFFORT, Atk: 4}
    HC252_B4 = EffortStat{HP: MAX_EFFORT, SpAtk: MAX_EFFORT, Def: 4}
    HC252_D4 = EffortStat{HP: MAX_EFFORT, SpAtk: MAX_EFFORT, SpDef: 4}
    HC252_S4 = EffortStat{HP: MAX_EFFORT, SpAtk: MAX_EFFORT, Speed: 4}

    HD252_A4 = EffortStat{HP: MAX_EFFORT, SpDef: MAX_EFFORT, Atk: 4}
    HD252_B4 = EffortStat{HP: MAX_EFFORT, SpDef: MAX_EFFORT, Def: 4}
    HD252_C4 = EffortStat{HP: MAX_EFFORT, SpDef: MAX_EFFORT, SpAtk: 4}
    HD252_S4 = EffortStat{HP: MAX_EFFORT, SpDef: MAX_EFFORT, Speed: 4}

    HS252_A4 = EffortStat{HP: MAX_EFFORT, Speed: MAX_EFFORT, Atk: 4}
    HS252_B4 = EffortStat{HP: MAX_EFFORT, Speed: MAX_EFFORT, Def: 4}
    HS252_C4 = EffortStat{HP: MAX_EFFORT, Speed: MAX_EFFORT, SpAtk: 4}
    HS252_D4 = EffortStat{HP: MAX_EFFORT, Speed: MAX_EFFORT, SpDef: 4}

    AB252_H4 = EffortStat{Atk: MAX_EFFORT, Def: MAX_EFFORT, HP: 4}
    AB252_C4 = EffortStat{Atk: MAX_EFFORT, Def: MAX_EFFORT, SpAtk: 4}
    AB252_D4 = EffortStat{Atk: MAX_EFFORT, Def: MAX_EFFORT, SpDef: 4}
    AB252_S4 = EffortStat{Atk: MAX_EFFORT, Def: MAX_EFFORT, Speed: 4}

    AC252_H4 = EffortStat{Atk: MAX_EFFORT, SpAtk: MAX_EFFORT, HP: 4}
    AC252_B4 = EffortStat{Atk: MAX_EFFORT, SpAtk: MAX_EFFORT, Def: 4}
    AC252_D4 = EffortStat{Atk: MAX_EFFORT, SpAtk: MAX_EFFORT, SpDef: 4}
    AC252_S4 = EffortStat{Atk: MAX_EFFORT, SpAtk: MAX_EFFORT, Speed: 4}

    AD252_H4 = EffortStat{Atk: MAX_EFFORT, SpDef: MAX_EFFORT, HP: 4}
    AD252_B4 = EffortStat{Atk: MAX_EFFORT, SpDef: MAX_EFFORT, Def: 4}
    AD252_C4 = EffortStat{Atk: MAX_EFFORT, SpDef: MAX_EFFORT, SpAtk: 4}
    AD252_S4 = EffortStat{Atk: MAX_EFFORT, SpDef: MAX_EFFORT, Speed: 4}

    AS252_H4 = EffortStat{Atk: MAX_EFFORT, Speed: MAX_EFFORT, HP: 4}
    AS252_B4 = EffortStat{Atk: MAX_EFFORT, Speed: MAX_EFFORT, Def: 4}
    AS252_C4 = EffortStat{Atk: MAX_EFFORT, Speed: MAX_EFFORT, SpAtk: 4}
    AS252_D4 = EffortStat{Atk: MAX_EFFORT, Speed: MAX_EFFORT, SpDef: 4}

    BC252_H4 = EffortStat{Def: MAX_EFFORT, SpAtk: MAX_EFFORT, HP: 4}
    BC252_A4 = EffortStat{Def: MAX_EFFORT, SpAtk: MAX_EFFORT, Atk: 4}
    BC252_D4 = EffortStat{Def: MAX_EFFORT, SpAtk: MAX_EFFORT, SpDef: 4}
    BC252_S4 = EffortStat{Def: MAX_EFFORT, SpAtk: MAX_EFFORT, Speed: 4}

    BD252_H4 = EffortStat{Def: MAX_EFFORT, SpDef: MAX_EFFORT, HP: 4}
    BD252_A4 = EffortStat{Def: MAX_EFFORT, SpDef: MAX_EFFORT, Atk: 4}
    BD252_C4 = EffortStat{Def: MAX_EFFORT, SpDef: MAX_EFFORT, SpAtk: 4}
    BD252_S4 = EffortStat{Def: MAX_EFFORT, SpDef: MAX_EFFORT, Speed: 4}

    BS252_H4 = EffortStat{Def: MAX_EFFORT, Speed: MAX_EFFORT, HP: 4}
    BS252_A4 = EffortStat{Def: MAX_EFFORT, Speed: MAX_EFFORT, Atk: 4}
    BS252_C4 = EffortStat{Def: MAX_EFFORT, Speed: MAX_EFFORT, SpAtk: 4}
    BS252_D4 = EffortStat{Def: MAX_EFFORT, Speed: MAX_EFFORT, SpDef: 4}

    CD252_H4 = EffortStat{SpAtk: MAX_EFFORT, SpDef: MAX_EFFORT, HP: 4}
    CD252_A4 = EffortStat{SpAtk: MAX_EFFORT, SpDef: MAX_EFFORT, Atk: 4}
    CD252_B4 = EffortStat{SpAtk: MAX_EFFORT, SpDef: MAX_EFFORT, Def: 4}
    CD252_S4 = EffortStat{SpAtk: MAX_EFFORT, SpDef: MAX_EFFORT, Speed: 4}

    CS252_H4 = EffortStat{SpAtk: MAX_EFFORT, Speed: MAX_EFFORT, HP: 4}
    CS252_A4 = EffortStat{SpAtk: MAX_EFFORT, Speed: MAX_EFFORT, Atk: 4}
    CS252_B4 = EffortStat{SpAtk: MAX_EFFORT, Speed: MAX_EFFORT, Def: 4}
    CS252_D4 = EffortStat{SpAtk: MAX_EFFORT, Speed: MAX_EFFORT, SpDef: 4}

    DS252_H4 = EffortStat{SpDef: MAX_EFFORT, Speed: MAX_EFFORT, HP: 4}
    DS252_A4 = EffortStat{SpDef: MAX_EFFORT, Speed: MAX_EFFORT, Atk: 4}
    DS252_B4 = EffortStat{SpDef: MAX_EFFORT, Speed: MAX_EFFORT, Def: 4}
    DS252_C4 = EffortStat{SpDef: MAX_EFFORT, Speed: MAX_EFFORT, SpAtk: 4}
)

func (ev *EffortStat) Sum() Effort {
    return ev.HP + ev.Atk + ev.Def + ev.SpAtk + ev.SpDef + ev.Speed
}

func (ev *EffortStat) SumError() error {
	if ev.Sum() > MAX_SUM_EFFORT {
		m := fmt.Sprintf("努力値の合計が%dを超えている", MAX_SUM_EFFORT)
		return fmt.Errorf(m)
	}
	return nil
}