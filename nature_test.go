package bippa_test

import (
	"testing"
	"fmt"
	bp "github.com/sw965/bippa"
)

func TestNATUREDEX(t *testing.T) {
	for k, v := range bp.NATUREDEX {
		fmt.Println(k.ToString(), v)
	}
}