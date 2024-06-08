package bippa

type IV int

const (
	EMPTY_IV IV = -1
	MIN_IV IV = 0
	MAX_IV IV = 31
)

type IVs []IV

var ALL_IVS = func() IVs {
	ret := make(IVs, int(MAX_IV + 1))
	for i := range ret {
		ret[i] = IV(i)
	}
	return ret
}()

type IVStat struct {
	HP IV
	Atk IV
	Def IV
	SpAtk IV
	SpDef IV
	Speed IV
}

var EMPTY_IV_STAT = IVStat{HP:EMPTY_IV, Atk:EMPTY_IV, Def:EMPTY_IV, SpAtk:EMPTY_IV, SpDef:EMPTY_IV, Speed:EMPTY_IV}
var MIN_IV_STAT = IVStat{}
var MAX_IV_STAT = IVStat{HP:MAX_IV, Atk:MAX_IV, Def:MAX_IV, SpAtk:MAX_IV, SpDef:MAX_IV, Speed:MAX_IV}


func (iv *IVStat) IsAnyEmpty() bool {
	if iv.HP == EMPTY_IV {
		return true
	}

	if iv.Atk == EMPTY_IV {
		return true
	}

	if iv.Def == EMPTY_IV {
		return true
	}

	if iv.SpAtk == EMPTY_IV {
		return true
	}

	if iv.SpDef == EMPTY_IV {
		return true
	}

	if iv.Speed == EMPTY_IV {
		return true
	}
	return false
}