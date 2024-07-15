package single

func Thunderbolt(battle *Battle, selfLeadIdx, targetId int) {
	targetPokemon := battle.GetTargetPokemon(targetId)
	dmg := battle.CalcDamage(selfLeadIdx, targetId, false)
	dmg = omwmath.Min(dmg, 1)
	targetPokemon.CurrentHP -= dmg
}

func RockSlide(battle *Battle, selfLeadIdx, targetId int) {
	
}