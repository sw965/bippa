package bippa

// https://wiki.xn--rckteqa2e.com/wiki/%E3%83%A9%E3%83%B3%E3%82%AF%E8%A3%9C%E6%AD%A3#%E6%80%A5%E6%89%80%E3%83%A9%E3%83%B3%E3%82%AF
type CriticalRank int

// https://wiki.xn--rckteqa2e.com/wiki/%E3%83%A9%E3%83%B3%E3%82%AF%E8%A3%9C%E6%AD%A3
type Rank int

func (r Rank) Bonus() float64 {
	if r > 0 {
		return (2.0 + float64(r)) / 2.0
	} else {
		return 2.0 / (2.0 + float64(-r))
	}
}

const (
	MIN_RANK = -6
	MAX_RANK = 6
)

type RankStat struct {
	Atk Rank
	Def Rank
	SpAtk Rank
	SpDef Rank
	Speed Rank
}

func (r RankStat) Clone() RankStat {
	return r
}