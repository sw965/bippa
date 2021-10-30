package bippa

type PhysicsBonus int

const (
  INIT_PHYSICS_BONUS = PhysicsBonus(4096)
)

func NewPhysicsBonus(spovb *SelfPointOfViewBattle) PhysicsBonus {
  physicsBonus := INIT_PHYSICS_BONUS
	if spovb.SelfFighters[0].Item == "こだわりハチマキ" {
    physicsBonus = physicsBonus.MulChoiceBand()
	}
	return physicsBonus
}

func (physicsBonus PhysicsBonus) MulChoiceBand() PhysicsBonus {
  result := RoundingZeroPointFiveOrMore(float64(physicsBonus) * 6144.0 / 4096.0)
  return PhysicsBonus(result)
}

type SpecialBonus int

const (
  INIT_SPECIAL_BONUS = SpecialBonus(4096)
)

func NewSpecialBonus(spovb *SelfPointOfViewBattle) SpecialBonus {
	specialBonus := INIT_SPECIAL_BONUS
	if spovb.SelfFighters[0].Item == "こだわりメガネ" {
    specialBonus = specialBonus.MulChoiceSpecs()
	}
	return specialBonus
}

func (specialBonus SpecialBonus) MulChoiceSpecs() SpecialBonus {
  result := RoundingZeroPointFiveOrMore(float64(specialBonus) * 6144.0 / 4096.0)
  return SpecialBonus(result)
}

type AttackBonus int

const (
  INIT_ATTACK_BONUS = AttackBonus(4096)
)

func NewAttackBonus(spovb *SelfPointOfViewBattle, moveName MoveName) (AttackBonus, error) {
	moveData := MOVEDEX[moveName]

	if moveData.Category == PHYSICS {
    physicsBonus := NewPhysicsBonus(spovb)
		return AttackBonus(physicsBonus), nil
	}

	if moveData.Category == SPECIAL {
    specialBonus := NewSpecialBonus(spovb)
		return SpecialBonus(pecialBonus), nil
	}

	return 0, fmt.Errorf("変化技以外でなければならない")
}

type FinalAttack int

func FinalAttackCalc(spovb *SelfPointOfViewBattle, moveName MoveName, isCritical bool) (FinalAttack, error) {
	moveData := MOVEDEX[moveName]

	var atkValue int
	var rank int

	switch moveData.Category {
		case PHYSICS:
			atkValue = spovb.SelfFighters[0].State.Atk
			rank = spovb.SelfFighters[0].RankState.Atk
		case SPECIAL:
			atkValue = spovb.SelfFighters[0].State.SpAtk
			rank = spovb.SelfFighters[0].RankState.SpAtk
	}

	//変化技の場合、ここでエラーが起きるので、上のswitch文ではチェック不要
	attackBonus, err := NewAttackBonus(spovb, moveName)

	if err != nil {
		return 0, err
	}

	if rank < 0 && isCritical {
		rank = 0
	}

  rankBonus := RANK_TO_RANK_BONUS[rank]

	finalAttack := int(float64(atkValue) * rankBonus)
	finalAttack = RoundingZeroPointFiveOver(float64(finalAttack) * float64(attackBonus) / 4096.0)
	if finalAttack < 1 {
		return 1, nil
	}
	return finalAttack, nil
}
