package gitignore

import (
	"testing"
)

func testNewSimplePattern(t *testing.T, s string, arr []string, prefix, suffix bool) {
	p := NewSimplePattern(s)
	if p.ShouldStartWithMe != prefix || p.IsDirectory != suffix {
		t.Error("failed: s=", s, ", prefix=", prefix, ", suffix=", suffix)
	}

	for i := 0; i < len(arr); i++ {
		if i >= len(p.Fields) || p.Fields[i] != arr[i] {
			t.Error("failed: s=", s, ", pattern=", p)
		}
	}
}

func TestNewSimplePattern(t *testing.T) {
	testNewSimplePattern(t, "a/b/c", []string{"a", "b", "c"}, false, false)
	testNewSimplePattern(t, "/a/b/c", []string{"a", "b", "c"}, true, false)
	testNewSimplePattern(t, "a/b/c/", []string{"a", "b", "c"}, false, true)
	testNewSimplePattern(t, "aa/b/c/", []string{"aa", "b", "c"}, false, true)

	testNewSimplePattern(t, "", []string{}, false, false)
}

func testMatching(t *testing.T, s string, path []string, result bool) {
	pattern := NewSimplePattern(s)
	if pattern.Match(path) != result {
		t.Error("failed: pattern=", s, ", path=", path, ", expect=", result)
	}
}

func TestMatching(t *testing.T) {
	testMatching(t, "", []string{"a"}, false)
	testMatching(t, "", []string{}, false)

	testMatching(t, "a/b/c", []string{"a", "b", "c"}, true)
	testMatching(t, "/a/b/c", []string{"a", "b", "c"}, true)
	testMatching(t, "a/b/c/", []string{"a", "b", "c"}, false)

	testMatching(t, "a/b/c", []string{"a", "b", "c", "d"}, true)
	testMatching(t, "/a/b/c", []string{"a", "b", "c", "d"}, true)
	testMatching(t, "/a/b/c/", []string{"a", "b", "c", "d"}, true)

	testMatching(t, "a/b/c", []string{"a", "b", "d"}, false)

	testMatching(t, "a/b/c", []string{"0", "a", "b", "c"}, true)
	testMatching(t, "a/b/c", []string{"0", "a", "b", "d"}, false)
	testMatching(t, "/a/b/c", []string{"0", "a", "b", "c"}, false)
	testMatching(t, "a/b/c/", []string{"0", "a", "b", "c"}, false)
	testMatching(t, "a/b/c/", []string{"0", "a", "b", "c", "d"}, true)

	testMatching(t, "0/a/b/c", []string{"a", "b", "c"}, false)

	testMatching(t, "long/b/c", []string{"long", "b", "c"}, true)
	testMatching(t, "long/b/c", []string{"longlong", "b", "c"}, false)

	testMatching(t, "a", []string{"a"}, true)
	testMatching(t, "/a", []string{"a"}, true)
	testMatching(t, "a/", []string{"a"}, false)

	testMatching(t, "a", []string{"0", "a"}, true)
	testMatching(t, "/a", []string{"0", "a"}, false)
	testMatching(t, "a/", []string{"0", "a"}, false)

	testMatching(t, "a", []string{"a", "b"}, true)
	testMatching(t, "/a", []string{"a", "b"}, true)
	testMatching(t, "a/", []string{"a", "b"}, true)
}
