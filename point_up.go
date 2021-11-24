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

var ALL_MAX_POINT_UPS = map[int]PointUps{
	1:PointUps{MAX_POINT_UP},
	2:PointUps{MAX_POINT_UP, MAX_POINT_UP},
	3:PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
	4:PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
}
