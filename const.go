package bippa

const (
	PHYSICS = "物理"
	SPECIAL = "特殊"
	STATUS  = "変化"
)

const (
	MAX_SPIKES = 3
	MAX_TOXIC_SPIKES = 2
)

var SPIKES_COUNT_TO_DAMAGE_PERCENT = map[int]float64{
	1:1.0 / 8.0, 2:1.0 / 6.0, 3:1.0 / 4.0,
}

var TOXIC_SPIKES_COUNT_TO_STATUS_AILMENT_ = map[int]StatusAilment_{
	1:NORMAL_POISON, 2:BAD_POISON,
}
