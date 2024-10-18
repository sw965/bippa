package bippa_test

import (
	"fmt"
	bp "github.com/sw965/bippa"
	"testing"
)

func TestNATUREDEX(t *testing.T) {
	for k, v := range bp.NATUREDEX {
		fmt.Println(k.ToString(), v)
	}
}
