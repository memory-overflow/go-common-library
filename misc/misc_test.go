package misc_test

import (
	"fmt"
	"testing"

	"github.com/memory-overflow/go-common-library/misc"
)

func goFunc() {
	fmt.Println("start goFunc")
	panic("test panic")
}

func TestGoroutineHelp(t *testing.T) {
	help := misc.GoroutineHelp{}
	help.SafeGo(goFunc, true)

	for {

	}
}

func TestFlodLog(t *testing.T) {
	log := `{"RetrieveInputData":"其他","RetrieveInputType":0,"FilterSet":[],"PageNumber":1,"PageSize":6,"TIBusinessID":1,"TIProjectID":1}`
	fmt.Println(misc.FoldLog([]byte(log)))
}
