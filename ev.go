package bippa

import (
    "github.com/sw965/omw/fn"
)

type EV int

const (
    EMPTY_EV EV = -1
	MIN_EV EV = 0
	MAX_EV EV = 252
    MAX_SUM_EV = 510
)

func IsEffectiveEV(ev EV) bool {
    return ev%4 == 0
}

type EVs []EV

var ALL_EVS = func() EVs {
    evs := make(EVs, MAX_EV+1)
    for i := 0; i < int(MAX_EV); i++ {
        evs[i] = EV(i)
    }
    return evs
}()

var EFFECTIVE_EVS = fn.Filter[EVs](ALL_EVS, IsEffectiveEV)

type EVStat struct {
	HP EV
	Atk EV
	Def EV
	SpAtk EV
	SpDef EV
	Speed EV
}

var (
    HA252_B4 = EVStat{HP: MAX_EV, Atk: MAX_EV, Def: 4}
    HA252_C4 = EVStat{HP: MAX_EV, Atk: MAX_EV, SpAtk: 4}
    HA252_D4 = EVStat{HP: MAX_EV, Atk: MAX_EV, SpDef: 4}
    HA252_S4 = EVStat{HP: MAX_EV, Atk: MAX_EV, Speed: 4}

    HB252_A4 = EVStat{HP: MAX_EV, Def: MAX_EV, Atk: 4}
    HB252_C4 = EVStat{HP: MAX_EV, Def: MAX_EV, SpAtk: 4}
    HB252_D4 = EVStat{HP: MAX_EV, Def: MAX_EV, SpDef: 4}
    HB252_S4 = EVStat{HP: MAX_EV, Def: MAX_EV, Speed: 4}

    HC252_A4 = EVStat{HP: MAX_EV, SpAtk: MAX_EV, Atk: 4}
    HC252_B4 = EVStat{HP: MAX_EV, SpAtk: MAX_EV, Def: 4}
    HC252_D4 = EVStat{HP: MAX_EV, SpAtk: MAX_EV, SpDef: 4}
    HC252_S4 = EVStat{HP: MAX_EV, SpAtk: MAX_EV, Speed: 4}

    HD252_A4 = EVStat{HP: MAX_EV, SpDef: MAX_EV, Atk: 4}
    HD252_B4 = EVStat{HP: MAX_EV, SpDef: MAX_EV, Def: 4}
    HD252_C4 = EVStat{HP: MAX_EV, SpDef: MAX_EV, SpAtk: 4}
    HD252_S4 = EVStat{HP: MAX_EV, SpDef: MAX_EV, Speed: 4}

    HS252_A4 = EVStat{HP: MAX_EV, Speed: MAX_EV, Atk: 4}
    HS252_B4 = EVStat{HP: MAX_EV, Speed: MAX_EV, Def: 4}
    HS252_C4 = EVStat{HP: MAX_EV, Speed: MAX_EV, SpAtk: 4}
    HS252_D4 = EVStat{HP: MAX_EV, Speed: MAX_EV, SpDef: 4}

    AB252_H4 = EVStat{Atk: MAX_EV, Def: MAX_EV, HP: 4}
    AB252_C4 = EVStat{Atk: MAX_EV, Def: MAX_EV, SpAtk: 4}
    AB252_D4 = EVStat{Atk: MAX_EV, Def: MAX_EV, SpDef: 4}
    AB252_S4 = EVStat{Atk: MAX_EV, Def: MAX_EV, Speed: 4}

    AC252_H4 = EVStat{Atk: MAX_EV, SpAtk: MAX_EV, HP: 4}
    AC252_B4 = EVStat{Atk: MAX_EV, SpAtk: MAX_EV, Def: 4}
    AC252_D4 = EVStat{Atk: MAX_EV, SpAtk: MAX_EV, SpDef: 4}
    AC252_S4 = EVStat{Atk: MAX_EV, SpAtk: MAX_EV, Speed: 4}

    AD252_H4 = EVStat{Atk: MAX_EV, SpDef: MAX_EV, HP: 4}
    AD252_B4 = EVStat{Atk: MAX_EV, SpDef: MAX_EV, Def: 4}
    AD252_C4 = EVStat{Atk: MAX_EV, SpDef: MAX_EV, SpAtk: 4}
    AD252_S4 = EVStat{Atk: MAX_EV, SpDef: MAX_EV, Speed: 4}

    AS252_H4 = EVStat{Atk: MAX_EV, Speed: MAX_EV, HP: 4}
    AS252_B4 = EVStat{Atk: MAX_EV, Speed: MAX_EV, Def: 4}
    AS252_C4 = EVStat{Atk: MAX_EV, Speed: MAX_EV, SpAtk: 4}
    AS252_D4 = EVStat{Atk: MAX_EV, Speed: MAX_EV, SpDef: 4}

    BC252_H4 = EVStat{Def: MAX_EV, SpAtk: MAX_EV, HP: 4}
    BC252_A4 = EVStat{Def: MAX_EV, SpAtk: MAX_EV, Atk: 4}
    BC252_D4 = EVStat{Def: MAX_EV, SpAtk: MAX_EV, SpDef: 4}
    BC252_S4 = EVStat{Def: MAX_EV, SpAtk: MAX_EV, Speed: 4}

    BD252_H4 = EVStat{Def: MAX_EV, SpDef: MAX_EV, HP: 4}
    BD252_A4 = EVStat{Def: MAX_EV, SpDef: MAX_EV, Atk: 4}
    BD252_C4 = EVStat{Def: MAX_EV, SpDef: MAX_EV, SpAtk: 4}
    BD252_S4 = EVStat{Def: MAX_EV, SpDef: MAX_EV, Speed: 4}

    BS252_H4 = EVStat{Def: MAX_EV, Speed: MAX_EV, HP: 4}
    BS252_A4 = EVStat{Def: MAX_EV, Speed: MAX_EV, Atk: 4}
    BS252_C4 = EVStat{Def: MAX_EV, Speed: MAX_EV, SpAtk: 4}
    BS252_D4 = EVStat{Def: MAX_EV, Speed: MAX_EV, SpDef: 4}

    CD252_H4 = EVStat{SpAtk: MAX_EV, SpDef: MAX_EV, HP: 4}
    CD252_A4 = EVStat{SpAtk: MAX_EV, SpDef: MAX_EV, Atk: 4}
    CD252_B4 = EVStat{SpAtk: MAX_EV, SpDef: MAX_EV, Def: 4}
    CD252_S4 = EVStat{SpAtk: MAX_EV, SpDef: MAX_EV, Speed: 4}

    CS252_H4 = EVStat{SpAtk: MAX_EV, Speed: MAX_EV, HP: 4}
    CS252_A4 = EVStat{SpAtk: MAX_EV, Speed: MAX_EV, Atk: 4}
    CS252_B4 = EVStat{SpAtk: MAX_EV, Speed: MAX_EV, Def: 4}
    CS252_D4 = EVStat{SpAtk: MAX_EV, Speed: MAX_EV, SpDef: 4}

    DS252_H4 = EVStat{SpDef: MAX_EV, Speed: MAX_EV, HP: 4}
    DS252_A4 = EVStat{SpDef: MAX_EV, Speed: MAX_EV, Atk: 4}
    DS252_B4 = EVStat{SpDef: MAX_EV, Speed: MAX_EV, Def: 4}
    DS252_C4 = EVStat{SpDef: MAX_EV, Speed: MAX_EV, SpAtk: 4}
)

func (ev *EVStat) Sum() EV {
    return ev.HP + ev.Atk + ev.Def + ev.SpAtk + ev.SpDef + ev.Speed
}

func (ev *EVStat) IsAnyEmpty() bool {
	if ev.HP == EMPTY_EV {
		return true
	}

	if ev.Atk == EMPTY_EV {
		return true
	}

	if ev.Def == EMPTY_EV {
		return true
	}

	if ev.SpAtk == EMPTY_EV {
		return true
	}

	if ev.SpDef == EMPTY_EV {
		return true
	}

	if ev.Speed == EMPTY_EV {
		return true
	}
	return false
}