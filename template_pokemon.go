package bippa

//ギャラドス
// https://matsu-1129.hatenadiary.org/entry/20090308/1236586122
func NewRomanStan2009Gyarados() Pokemon {
	pokemon, err := NewPokemon(
		GYARADOS, STANDARD_LEVEL, ADAMANT, INTIMIDATE, WACAN_BERRY,
		MoveNames{WATERFALL, STONE_EDGE, THUNDER_WAVE, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT,
		&EffortStat{HP:164, Atk:92, Speed:MAX_EFFORT},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//カビゴン
func NewMoruhu2007Snorlax() Pokemon {
	ivStat := MAX_INDIVIDUAL_STAT.Clone()
	ivStat.Speed = MAX_INDIVIDUAL
	pokemon, err := NewPokemon(
		SNORLAX, STANDARD_LEVEL, RELAXED, THICK_FAT, SITRUS_BERRY,
		MoveNames{RETURN, FIRE_PUNCH, BELLY_DRUM, SUBSTITUTE},
		MAX_POINT_UPS,
		&ivStat, &HB252_D4,
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//カビゴン
func NewMoruhu2008Snorlax() Pokemon {
	ivStat := NewMaxIndividualStat()
	ivStat.Speed = MIN_INDIVIDUAL
	pokemon, err := NewPokemon(
		SNORLAX, STANDARD_LEVEL, RELAXED, THICK_FAT, SITRUS_BERRY,
		MoveNames{RETURN, FIRE_PUNCH, BELLY_DRUM, PROTECT},
		MAX_POINT_UPS,
		&ivStat, &HB252_D4,
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}


//カビゴン
// https://matsu-1129.hatenadiary.org/entry/20090308/1236586122
func NewRomanStan2009Snorlax() Pokemon {
	ivStat := NewMaxIndividualStat()
	ivStat.Speed = MIN_INDIVIDUAL
	pokemon, err := NewPokemon(
		SNORLAX, STANDARD_LEVEL, BRAVE, THICK_FAT, SITRUS_BERRY,
		MoveNames{RETURN, CRUNCH, SELF_DESTRUCT, PROTECT},
		MAX_POINT_UPS,
		&ivStat, &EffortStat{HP:204, Atk:52, Def:156, SpDef:96},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//カビゴン
// https://detail.chiebukuro.yahoo.co.jp/qa/question_detail/q1267938327
func NewKusanagi2009Snorlax() Pokemon {
	pokemon, err := NewPokemon(
		SNORLAX, STANDARD_LEVEL, BRAVE, THICK_FAT, SITRUS_BERRY,
		MoveNames{RETURN, CRUNCH, SELF_DESTRUCT, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT,
		&EffortStat{HP:204, Atk:52, Def:156, SpDef:60, Speed:36},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//ドーブル
func NewMoruhu2007Smeargle() Pokemon {
	pokemon, err := NewPokemon(
		SMEARGLE, MIN_LEVEL, BRAVE, OWN_TEMPO, FOCUS_SASH,
		MoveNames{FAKE_OUT, FOLLOW_ME, DARK_VOID, ENDEAVOR},
		MAX_POINT_UPS,
		&MIN_INDIVIDUAL_STAT, &EffortStat{},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//ドーブル
func NewMoruhu2008Smeargle() Pokemon {
	return NewMoruhu2007Smeargle()
}

//ボーマンダ
// https://detail.chiebukuro.yahoo.co.jp/qa/question_detail/q1267938327
func NewKusanagiSalamence2009() Pokemon {
	pokemon, err := NewPokemon(
		SALAMENCE, STANDARD_LEVEL, MODEST, INTIMIDATE, SITRUS_BERRY,
		MoveNames{DRACO_METEOR, HEAT_WAVE, RAIN_DANCE, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT, &EffortStat{HP:20, SpAtk:236, Speed:252},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//メタグロス
func NewMoruhu2007Metagross() Pokemon {
	ivStat := MAX_INDIVIDUAL_STAT.Clone()
	ivStat.Speed = MIN_INDIVIDUAL
	pokemon, err := NewPokemon(
		METAGROSS, STANDARD_LEVEL, RELAXED, CLEAR_BODY, LUM_BERRY,
		MoveNames{EARTHQUAKE, BULLET_PUNCH, ROCK_SLIDE, RECOVER},
		MAX_POINT_UPS,
		&ivStat, &EffortStat{HP:MAX_EFFORT, Def:128, SpDef:128},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//メタグロス
func NewMoruhu2008Metagross() Pokemon {
	ivStat := MAX_INDIVIDUAL_STAT.Clone()
	ivStat.Speed = MIN_INDIVIDUAL
	pokemon, err := NewPokemon(
		METAGROSS, STANDARD_LEVEL, RELAXED, CLEAR_BODY, LUM_BERRY,
		MoveNames{HAMMER_ARM, BULLET_PUNCH, ROCK_SLIDE, RECOVER},
		MAX_POINT_UPS,
		&ivStat, &EffortStat{HP:MAX_EFFORT, Def:128, SpDef:128},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//メタグロス
// https://matsu-1129.hatenadiary.org/entry/20090308/1236586122
func NewRomanStan2009Metagross() Pokemon {
	pokemon, err := NewPokemon(
		METAGROSS, STANDARD_LEVEL, ADAMANT, CLEAR_BODY, LUM_BERRY,
		MoveNames{COMET_PUNCH, BULLET_PUNCH, EARTHQUAKE, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT, &EffortStat{HP:236, Atk:36, Def:4, SpDef:172, Speed:60},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//ラティオス
// https://matsu-1129.hatenadiary.org/entry/20090308/1236586122
func NewRomanStan2009Latios() Pokemon {
	pokemon, err := NewPokemon(
		LATIOS, STANDARD_LEVEL, TIMID, LEVITATE, FOCUS_SASH,
		MoveNames{DRACO_METEOR, THUNDERBOLT, RAIN_DANCE, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT, &CS252_H4,
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//エンペルト
// https://detail.chiebukuro.yahoo.co.jp/qa/question_detail/q1267938327
func NewKusanagi2009Empoleon() Pokemon {
	pokemon, err := NewPokemon(
		EMPOLEON, STANDARD_LEVEL, MODEST, TORRENT, WACAN_BERRY,
		MoveNames{HYDRO_PUMP, SURF, ICY_WIND, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT, &EffortStat{HP:68, Def:12, SpAtk:252, SpDef:4, Speed:172},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//ドータクン
func NewMoruhu2007Bronzong() Pokemon {
	ivStat := NewMaxIndividualStat()
	ivStat.Speed = MIN_INDIVIDUAL
	pokemon, err := NewPokemon(
		BRONZONG, STANDARD_LEVEL, SASSY, HEATPROOF, CHESTO_BERRY,
		MoveNames{GYRO_BALL, EXPLOSION, TRICK_ROOM, HYPNOSIS},
		MAX_POINT_UPS,
		&ivStat, &HD252_B4,
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//ドータクン
func NewMoruhu2008Bronzong() Pokemon {
	ivStat := NewMaxIndividualStat()
	ivStat.Speed = MIN_INDIVIDUAL
	pokemon, err := NewPokemon(
		BRONZONG, STANDARD_LEVEL, SASSY, HEATPROOF, CHESTO_BERRY,
		MoveNames{PSYCHIC, EXPLOSION, TRICK_ROOM, HYPNOSIS},
		MAX_POINT_UPS,
		&ivStat, &HD252_B4,
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//ドクロッグ
// https://detail.chiebukuro.yahoo.co.jp/qa/question_detail/q1267938327
func NewKusanagi2009Toxicroak() Pokemon {
	pokemon, err := NewPokemon(
		TOXICROAK, STANDARD_LEVEL, ADAMANT, DRY_SKIN, FOCUS_SASH,
		MoveNames{CROSS_CHOP, SUCKER_PUNCH, FAKE_OUT, TAUNT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT, &AS252_B4,
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}