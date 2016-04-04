package gitignore

import (
	"strings"
	"testing"
)

func TestParseGitignore(t *testing.T) {
	str := "#comment\na/b/c\n  c/d/e  \n\n"
	p, _ := parseGitignore(strings.NewReader(str))
	t.Logf("str=\n%v", str)
	for i, v := range p {
		t.Logf("[%v] %v", i, v)
	}
}

func TestTravel(t *testing.T) {
	g, _ := ScanGitignore("/Users/cos/wdj/git/phoenix2/RB_4.2_em")
	for k, arr := range g.m {
		t.Logf("%v =>", k)
		for i, v := range arr {
			t.Logf("  [%v] %v", i, v)
		}
	}
}
