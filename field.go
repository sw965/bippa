package bippa

type FieldType int

const (
	GRASS_FIELD FieldType = iota
)

type Field struct {
	Type FieldType
	RemainingTurn int
}