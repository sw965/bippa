package bippa

type EasyReadNaturedex map[string]NatureData

func (e EasyReadNaturedex) From() (Naturedex, error) {
	d := Naturedex{}
	for k, v := range e {
		n, err := StringToNature(k)
		if err != nil {
			return Naturedex{}, err
		}
		d[n] = &v
	}
	return d, nil
}

type EasyReadMoveset map[string]PowerPoint

func (e EasyReadMoveset) From() (Moveset, error) {
	m := Moveset{}
	for k, v := range e {
		n, err := StringToMoveName(k)
		if err != nil {
			return Moveset{}, err
		}
		pp := PowerPoint{Max:v.Max, Current:v.Current}
		m[n] = &pp
	}
	return m, nil
}

type EasyReadMoveData struct {
    Type         string
    Category     string
    Power        int
    Accuracy     int
    BasePP       int
	IsContact    bool
	PriorityRank int
	CriticalRank CriticalRank
	Target       string
}

func (m *EasyReadMoveData) From() (MoveData, error) {
	t, err := StringToType(m.Type)
	if err != nil {
		return MoveData{}, err
	}

	category, err := StringToMoveCategory(m.Category)
	if err != nil {
		return MoveData{}, err
	}

	target, err := StringToTargetRange(m.Target)
	if err != nil {
		return MoveData{}, err
	}

	return MoveData{
		Type:t,
		Category:category,
		Power:m.Power,
		Accuracy:m.Accuracy,
		BasePP:m.BasePP,
		IsContact:m.IsContact,
		PriorityRank:m.PriorityRank,
		CriticalRank:m.CriticalRank,
		Target:target,
	}, nil
}

type EasyReadMovedex map[string]EasyReadMoveData
