package bippa

type MoveCategory string

const (
	PHYSICS = MoveCategory("物理")
	SPECIAL = MoveCategory("特殊")
	STATUS  = MoveCategory("変化")
)
