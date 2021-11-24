package bippa

import (
  "testing"
)

func TestVenusaur(t *testing.T) {
  pokemon := TEST_POKEMONS["フシギバナ"]()
  if pokemon.State.MaxHP != 187 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.CurrentHP != 187 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.Atk != 102 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.Def != 103 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.SpAtk != 108 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.SpDef !=167 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.Speed != 101 {
    t.Errorf("テスト失敗")
  }
}

func TestCharizard(t *testing.T) {
  pokemon := TEST_POKEMONS["リザードン"]()
  if pokemon.State.MaxHP != 154 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.CurrentHP != 154 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.Atk != 93 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.Def != 98 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.SpAtk != 161 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.SpDef !=105 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.Speed != 167 {
    t.Errorf("テスト失敗")
  }
}

func TestBlastoise(t *testing.T) {
  pokemon := TEST_POKEMONS["カメックス"]()
  if pokemon.State.MaxHP != 155 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.CurrentHP != 155 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.Atk != 92 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.Def != 120 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.SpAtk != 150 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.SpDef !=125 {
    t.Errorf("テスト失敗")
  }

  if pokemon.State.Speed != 130 {
    t.Errorf("テスト失敗")
  }
}
