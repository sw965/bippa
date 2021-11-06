package bippa

import (
  "fmt"
)

type PhysicsAttackBonus int

const (
  INIT_PHYSICS_ATTACK_BONUS = PhysicsAttackBonus(4096)
)

func NewPhysicsAttackBonus(spovb *SelfPointOfViewBattle) PhysicsAttackBonus {
  result := INIT_PHYSICS_ATTACK_BONUS
	if spovb.SelfFighters[0].Item == "こだわりハチマキ" {
    result = result.MulChoiceBand()
	}
	return result
}

func (physicsAttackBonus PhysicsAttackBonus) MulChoiceBand() PhysicsAttackBonus {
  result := RoundingZeroPointFiveOrMore(float64(physicsAttackBonus) * 6144.0 / 4096.0)
  return PhysicsAttackBonus(result)
}

type SpecialAttackBonus int

const (
  INIT_SPECIAL_ATTACK_BONUS = SpecialAttackBonus(4096)
)

func NewSpecialAttackBonus(spovb *SelfPointOfViewBattle) SpecialAttackBonus {
	result := INIT_SPECIAL_ATTACK_BONUS
	if spovb.SelfFighters[0].Item == "こだわりメガネ" {
    result = result.MulChoiceSpecs()
	}
	return result
}

func (specialAttackBonus SpecialAttackBonus) MulChoiceSpecs() SpecialAttackBonus {
  result := RoundingZeroPointFiveOrMore(float64(specialAttackBonus) * 6144.0 / 4096.0)
  return SpecialAttackBonus(result)
}

type AttackBonus int

const (
  INIT_ATTACK_BONUS = AttackBonus(4096)
)

func NewAttackBonus(spovb *SelfPointOfViewBattle, moveName MoveName) (AttackBonus, error) {
	moveData := MOVEDEX[moveName]

	if moveData.Category == PHYSICS {
    physicsAttackBonus := NewPhysicsAttackBonus(spovb)
		return AttackBonus(physicsAttackBonus), nil
	}

	if moveData.Category == SPECIAL {
    specialAttackBonus := NewSpecialAttackBonus(spovb)
		return AttackBonus(specialAttackBonus), nil
	}

	return 0, fmt.Errorf("変化技以外でなければならない")
}

type FinalAttack int

func NewFinalAttack(spovb *SelfPointOfViewBattle, moveName MoveName, isCritical bool) (FinalAttack, error) {
	moveData := MOVEDEX[moveName]

	var attackState_ State_
	var rank_ Rank_

	switch moveData.Category {
		case PHYSICS:
			attackState_ = spovb.SelfFighters[0].State.Atk
			rank_ = spovb.SelfFighters[0].Rank.Atk
		case SPECIAL:
			attackState_ = spovb.SelfFighters[0].State.SpAtk
			rank_ = spovb.SelfFighters[0].Rank.SpAtk
	}

	//変化技の場合、ここでエラーが起きるので、上のswitch文ではチェック不要
	attackBonus, err := NewAttackBonus(spovb, moveName)

	if err != nil {
		return 0, err
	}

	if rank_ < 0 && isCritical {
		rank_ = 0
	}

  rankBonus := RANK__TO_RANK_BONUS[rank_]

	result := int(float64(attackState_) * float64(rankBonus))
	result = RoundingZeroPointFiveOver(float64(result) * float64(attackBonus) / 4096.0)
	if result < 1 {
		return 1, nil
	}
	return FinalAttack(result), nil
}
