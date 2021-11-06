package bippa

type PointUp int

const (
	MIN_POINT_UP = PointUp(0)
	MAX_POINT_UP = PointUp(3)
)

func (pointUp PointUp) IsValid() bool {
	return pointUp >= MIN_POINT_UP && pointUp <= MAX_POINT_UP
}

type PointUps []PointUp

func NewAllMaxPointUps(length int) PointUps {
	result := make(PointUps, length)
	for i := 0; i < length; i++ {
		result[i] = MAX_POINT_UP
	}
	return result
}
