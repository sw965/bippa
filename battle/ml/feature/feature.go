package feature

func FirePowerIndex(pokemon *Pokemon) tensor.D1 {
	ret := make(tensor.D1, len(ALL_MOVE_NAMES))
	for i, moveName := range ALL_MOVE_NAMES {
		if _, ok := pokemon.Moveset[moveName]; ok {
			moveData := MOVEDEX[moveName]
			if moveData.Category == PHYSICS {
				ret[i] = float64(moveData.Power * pokemon.Atk) / 100.0
			} else if moveData.Category == SPECIAL {
				ret[i] = float64(moveData.Power * pokemon.SpAtk) / 100.0
			} else {
				ret[i] = 1.0
			}
		}
	}
	return ret
}

func DefenseIndex(pokemon *Pokemon) tensor.D1 {
	defIndexFeature := make(tensor.D1, len(bp.ALL_TYPES))
	spDefIndexFeature := make(tensor.D1, len(bp.ALL_TYPES))
	for i, t := range ALL_TYPES {
		typeData := TYPEDEX[t]
		effect := dmgtools.Effectiveness(t, pokeTypes)
		defIndexFeature[i] = float64(effect * pokemon.Def) / 100.0
		spDefIndexFeature[i] = float64(effect * pokemon.SpDef) / 100.0
	}
	return omwslices.Concat(defIndexFeature, spDefIndexFeature)
}

func ExpectedDamageRatioToCurrentHP(p1Pokemon, p2Pokemon *Pokemon) tensor.D1 {
		
}

func KORatioToCurrentHP(p1Pokemon, p2Pokemon *Pokemon) tensor.D1 {

}