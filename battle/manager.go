package battle

 import (
 	"fmt"
 	bp "github.com/sw965/bippa"
	omwrand "github.com/sw965/omw/math/rand"
	"golang.org/x/exp/slices"
)

const (
	PLAYER_NUM = 2
	LEAD_NUM = 1
	BENCH_NUM = 2
	FIGHTERS_NUM = LEAD_NUM + BENCH_NUM
	DOUBLE_BATTLE_NUM = 2
)

type RemainingTurn struct {
	Weather int
	TrickRoom int
}

type Manager struct {
	GuestHumanTitle bp.HumanTitle
	GuestHumanName bp.HumanName

	CurrentSelfLeadPokemons bp.Pokemons
	CurrentSelfBenchPokemons  bp.Pokemons
	CurrentOpponentLeadPokemons bp.Pokemons
	CurrentOpponentBenchPokemons bp.Pokemons

	CurrentSelfFollowMePokemonPointers bp.PokemonPointers
	CurrentOpponentFollowMePokemonPointers bp.PokemonPointers

	Weather Weather
	RemainingTurn RemainingTurn
	Turn int
	
	IsSingle bool
	CurrentSelfIsHost bool
	HostViewMessage string
}

func (m *Manager) GetHostLeadPokemons() bp.Pokemons {
	if m.CurrentSelfIsHost {
		return m.CurrentSelfLeadPokemons
	} else {
		return m.CurrentOpponentLeadPokemons
	}
}

func (m *Manager) Init() {
	for i := range m.CurrentSelfLeadPokemons {
		m.CurrentSelfLeadPokemons[i].IsHost = true
	}

	for i := range m.CurrentSelfBenchPokemons {
		m.CurrentSelfBenchPokemons[i].IsHost = true
	}
	m.CurrentSelfIsHost = true
}

func (m Manager) Clone() Manager {
	m.CurrentSelfLeadPokemons = m.CurrentSelfLeadPokemons.Clone()
	m.CurrentSelfBenchPokemons = m.CurrentSelfBenchPokemons.Clone()
	m.CurrentOpponentLeadPokemons = m.CurrentOpponentLeadPokemons.Clone()
	m.CurrentOpponentBenchPokemons = m.CurrentOpponentBenchPokemons.Clone()
	m.CurrentSelfFollowMePokemonPointers = slices.Clone(m.CurrentSelfFollowMePokemonPointers)
	m.CurrentOpponentFollowMePokemonPointers = slices.Clone(m.CurrentOpponentFollowMePokemonPointers)
	return m
}

func (m *Manager) IsTrickRoomState() bool {
	return m.RemainingTurn.TrickRoom > 0
}

func (m *Manager) SwapView() {
	m.CurrentSelfLeadPokemons, m.CurrentSelfBenchPokemons, m.CurrentOpponentLeadPokemons, m.CurrentOpponentBenchPokemons =
		m.CurrentOpponentLeadPokemons, m.CurrentOpponentBenchPokemons, m.CurrentSelfLeadPokemons, m.CurrentSelfBenchPokemons
	m.CurrentSelfFollowMePokemonPointers, m.CurrentOpponentFollowMePokemonPointers = m.CurrentOpponentFollowMePokemonPointers, m.CurrentSelfFollowMePokemonPointers
	m.CurrentSelfIsHost = !m.CurrentSelfIsHost
}

//攻撃する側のポケモンは瀕死ではない事が前提で呼び出す関数
func (m *Manager) TargetPokemonPointers(action *SoloAction) bp.PokemonPointers {
	moveData := bp.MOVEDEX[action.MoveName]

	single := func() bp.PokemonPointers {
		switch moveData.Target {
			case bp.NORMAL_TARGET:
				return m.CurrentOpponentLeadPokemons.ToPointers().NotFainted()
			case bp.OPPONENT_TWO_TARGET:
				return m.CurrentOpponentLeadPokemons.ToPointers().NotFainted()
			case bp.SELF_TARGET:
				return m.CurrentSelfLeadPokemons.ToPointers().NotFainted()
			case bp.OTHERS_TARGET:
				return m.CurrentOpponentLeadPokemons.ToPointers().NotFainted()
			case bp.ALL_TARGET:
				return bp.PokemonPointers{}
			case bp.OPPONENT_RANDOM_ONE_TARGET:
				return m.CurrentOpponentLeadPokemons.ToPointers().NotFainted()
		}
		return bp.PokemonPointers{}
	}

	doubleNormalTarget := func() bp.PokemonPointers {
		followMe := m.CurrentOpponentFollowMePokemonPointers.NotFainted()
		if len(followMe) != 0 {
			return bp.PokemonPointers{followMe[0]}
		}

		//味方への攻撃
		if action.IsSelfLeadTarget {
			ally := m.CurrentSelfLeadPokemons[action.TargetIndex]
			//攻撃対象の味方が瀕死ならば、ランダムに相手を攻撃する。
			if ally.IsFainted() {
				ps:= m.CurrentOpponentLeadPokemons.ToPointers().NotFainted()
				if len(ps) == 0 {
					return bp.PokemonPointers{}
				}
				return omwrand.Sample(ps, 1, GlobalContext.Rand)
			} else {
				return bp.PokemonPointers{&m.CurrentSelfLeadPokemons[action.TargetIndex]}
			}
		} else {
			target := m.CurrentOpponentLeadPokemons[action.TargetIndex]
			if target.IsFainted() {
				otherI := map[int]int{0:1, 1:0}[action.TargetIndex]
				other := m.CurrentOpponentLeadPokemons[otherI]
				if other.IsFainted() {
					return bp.PokemonPointers{}
				} else {
					return bp.PokemonPointers{&other}
				}						
			} else {
				return bp.PokemonPointers{&target}
			}
		}
	}

	doubleOpponentRandomOneTarget := func() bp.PokemonPointers {
		ps := m.CurrentOpponentLeadPokemons.ToPointers().NotFainted()
		if len(ps) == 0 {
			return bp.PokemonPointers{}
		} else {
			return omwrand.Sample(ps, 1, GlobalContext.Rand)
		}
	}

	double := func() bp.PokemonPointers {
		switch moveData.Target {
			case bp.NORMAL_TARGET:
				return doubleNormalTarget()
			case bp.OPPONENT_TWO_TARGET:
				return m.CurrentOpponentLeadPokemons.ToPointers().NotFainted()
			case bp.SELF_TARGET:
				return bp.PokemonPointers{&m.CurrentSelfLeadPokemons[action.TargetIndex]}
			case bp.OTHERS_TARGET:
				allyI := map[int]int{0:1, 1:0}[action.SrcIndex]
				ally := m.CurrentSelfLeadPokemons[allyI]
				return bp.PokemonPointers{&ally, &m.CurrentOpponentLeadPokemons[0], &m.CurrentOpponentLeadPokemons[1]}.NotFainted()
			case bp.ALL_TARGET:
				return bp.PokemonPointers{}
			case bp.OPPONENT_RANDOM_ONE_TARGET:
				return doubleOpponentRandomOneTarget()
		}
		return bp.PokemonPointers{}
	}

	if m.IsSingle {
		return single()
	} else {
		return double()
	}
}

// いかく
// https://wiki.xn--rckteqa2e.com/wiki/%E3%81%84%E3%81%8B%E3%81%8F
// https://wiki.xn--rckteqa2e.com/wiki/%E3%82%BF%E3%83%BC%E3%83%B3#%E3%82%BF%E3%83%BC%E3%83%B3%E3%81%AE%E8%A9%B3%E7%B4%B0

func (m *Manager) Switch(leadIdx, benchIdx int) error {
	beforePokeName := m.CurrentSelfLeadPokemons[benchIdx].Name
	beforePokeNameStr := beforePokeName.ToString()

	if m.CurrentSelfBenchPokemons[benchIdx].IsFainted() {
		msg := fmt.Sprintf("%d番目の %sに 交代しようとしたが、瀕死状態である為、交代出来ません。", benchIdx, beforePokeNameStr)
		return fmt.Errorf(msg)
	}

	if m.CurrentSelfIsHost {
		m.HostViewMessage = fmt.Sprintf("もどれ！ %s！", beforePokeNameStr)
	} else {
		m.HostViewMessage = fmt.Sprintf("%sの %sは %sを 引っ込めた！", m.GuestHumanTitle, m.GuestHumanName, beforePokeNameStr)
	}

	GlobalContext.Observer(m)
	m.CurrentSelfLeadPokemons[leadIdx], m.CurrentSelfBenchPokemons[benchIdx] = m.CurrentSelfBenchPokemons[benchIdx], m.CurrentSelfLeadPokemons[leadIdx]

	afterPokemon := m.CurrentSelfLeadPokemons[leadIdx]
	afterPokeName := afterPokemon.Name
	afterPokeNameStr := afterPokeName.ToString()

	if m.CurrentSelfIsHost {
		m.HostViewMessage = fmt.Sprintf("ゆけっ！ %s！", afterPokeNameStr)
	} else {
		m.HostViewMessage = fmt.Sprintf("%sの %sは %sを くりだした！", m.GuestHumanTitle, m.GuestHumanName, afterPokeNameStr)
	}
	GlobalContext.Observer(m)

	if afterPokemon.Ability == bp.INTIMIDATE {
		for i := range m.CurrentOpponentLeadPokemons {
			target := m.CurrentOpponentLeadPokemons[i]
			if target.Rank.Atk == bp.MIN_RANK {
				continue
			}

			targetPokeNameStr := target.Name.ToString()

			if target.Ability == bp.CLEAR_BODY {
				if m.CurrentSelfIsHost {
					m.HostViewMessage = fmt.Sprintf("%sの %sの クリアボディで %sの いかくは きかなかった！", m.GuestHumanName, targetPokeNameStr, afterPokeNameStr)
				} else {
					m.HostViewMessage = fmt.Sprintf("%sの クリアボディで %sの %sの いかくは きかなかった！", targetPokeNameStr, m.GuestHumanName, afterPokeNameStr)
				}
				GlobalContext.Observer(m)
				continue
			}

			if target.Rank.Atk != bp.MIN_RANK {
				m.CurrentOpponentLeadPokemons[i].Rank.Atk -= 1
				if m.CurrentSelfIsHost {
					m.HostViewMessage = fmt.Sprintf("%sの いかくで %sの %sの こうげきが さがった！", afterPokeNameStr, m.GuestHumanName, targetPokeNameStr)
				} else {
					m.HostViewMessage = fmt.Sprintf("%sの %sの いかくで %sの こうげきが さがった！", m.GuestHumanName, afterPokeNameStr, targetPokeNameStr)
				}
				GlobalContext.Observer(m)
			}
		}
	}
	return nil
}

// https://latest.pokewiki.net/%E3%83%90%E3%83%88%E3%83%AB%E4%B8%AD%E3%81%AE%E5%87%A6%E7%90%86%E3%81%AE%E9%A0%86%E7%95%AA
// https://wiki.xn--rckteqa2e.com/wiki/%E3%82%BF%E3%83%BC%E3%83%B3#1.%E3%83%9D%E3%82%B1%E3%83%A2%E3%83%B3%E3%82%92%E7%B9%B0%E3%82%8A%E5%87%BA%E3%81%99

func (m *Manager) TurnEnd() error {
	return nil
}

func (m *Manager) ToEasyRead() EasyReadManager {
	return EasyReadManager{
		SelfLeadPokemons:m.CurrentSelfLeadPokemons.ToEasyRead(),
		SelfBenchPokemons:m.CurrentSelfBenchPokemons.ToEasyRead(),

		OpponentLeadPokemons:m.CurrentOpponentLeadPokemons.ToEasyRead(),
		OpponentBenchPokemons:m.CurrentOpponentBenchPokemons.ToEasyRead(),

		Turn:m.Turn,
		Weather:m.Weather.ToString(),
		RemainingTurn:m.RemainingTurn,
	}
}

type EasyReadManager struct {
	SelfLeadPokemons bp.EasyReadPokemons
	SelfBenchPokemons bp.EasyReadPokemons

	OpponentLeadPokemons bp.EasyReadPokemons
	OpponentBenchPokemons bp.EasyReadPokemons

	Turn int
	IsHostView bool

	Weather string
	RemainingTurn RemainingTurn
}