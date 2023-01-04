package text_test

import (
	"testing"

	"github.com/memory-overflow/go-common-library/text"
)

func TestActrie(t *testing.T) {
	ac := text.BuildAcTrie([]string{"哈哈哈", "234", "dfg"})
	list, index := ac.Search("哈哈哈哈23434dfgdd")
	for i, l := range list {
		t.Log(l, index[i])
	}
}

func TestTextSim(t *testing.T) {
	sim := text.TextSim("坚持干部一线值守常态化防控不放松让市民群众度过一个健康放心的假期", "这个口现在要开始封死")
	if sim != 0 {
		t.Error("Failed")
	}
}
