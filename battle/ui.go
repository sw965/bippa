package battle

 import (
  	"fmt"
  	bp "github.com/sw965/bippa"
)

type DisplayUI struct {
	P1LeadPokeNames []string
	P1LeadLevels []bp.Level
	P1LeadMaxHPs []int
	P1LeadCurrentHPs []int

	P2LeadPokeNames []string
	P2LeadLevels []bp.Level
	P2LeadMaxHPs []int
	P2LeadCurrentHPs []int

	Message string
}

type ObserverUI struct {
	LastP1ViewManager Manager
	LastP2ViewManager Manager
}

func (ui *ObserverUI) LastManager(isPlayer1View bool) Manager {
	if isPlayer1View {
		return ui.LastP1ViewManager
	} else {
		return ui.LastP2ViewManager
	}
}

func (ui *ObserverUI) Observer(m *Manager) {
	selfIsHost := m.CurrentSelfIsHost
	if !selfIsHost {
		m.SwapView()
	}

	fmt.Println(
		"self",
		m.CurrentSelfLeadPokemons.Names().ToStrings(), m.CurrentSelfLeadPokemons.CurrentHPs(),
		m.CurrentSelfBenchPokemons.Names().ToStrings(), m.CurrentSelfBenchPokemons.CurrentHPs(),
	)

	fmt.Println(
		"opponent",
		m.CurrentOpponentLeadPokemons.Names().ToStrings(), m.CurrentOpponentLeadPokemons.CurrentHPs(),
		m.CurrentOpponentBenchPokemons.Names().ToStrings(), m.CurrentOpponentBenchPokemons.CurrentHPs(),
	)

	if !selfIsHost {
		m.SwapView()
	}
	fmt.Println(m.HostViewMessage)
}