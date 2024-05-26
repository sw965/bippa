package bippa_test

import (
	"fmt"
	"testing"
	"golang.org/x/exp/slices"
	bp "github.com/sw965/bippa"
)

func TestSort(t *testing.T) {
	types := bp.Types{bp.GRASS, bp.NORMAL, bp.POISON, bp.ELECTRIC, bp.FAIRY, bp.STEEL, bp.DARK, bp.ROCK, bp.DRAGON}
	result := types.Sort()
	expected := bp.Types{bp.NORMAL, bp.GRASS, bp.ELECTRIC, bp.POISON, bp.ROCK, bp.DRAGON, bp.DARK, bp.STEEL, bp.FAIRY}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}