package bippa

type IV int

const (
	MIN_IV IV = 0
	MAX_IV IV = 31
)

type IVStat struct {
	HP IV
	Atk IV
	Def IV
	SpAtk IV
	SpDef IV
	Speed IV
}

var MIN_IV_STAT = IVStat{}
var MAX_IV_STAT = IVStat{HP:MAX_IV, Atk:MAX_IV, Def:MAX_IV, SpAtk:MAX_IV, SpDef:MAX_IV, Speed:MAX_IV}
