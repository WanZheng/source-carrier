package gitignore

// TODO:
//   - 支持'*'通配 // 可以用filepath.Match
//   - [bug] a/b/c 不能成功匹配 a/a/b/c

import (
	"strings"
)

type Pattern interface {
	Match(path []string) bool
}

type SimplePattern struct {
	Fields            []string
	IsDirectory       bool
	ShouldStartWithMe bool
}

func NewSimplePattern(s string) *SimplePattern {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r == '/'
	})
	pattern := SimplePattern{Fields: fields}
	pattern.ShouldStartWithMe = strings.HasPrefix(s, "/")
	pattern.IsDirectory = strings.HasSuffix(s, "/")

	return &pattern
}

func (pattern *SimplePattern) Match(path []string) bool {
	patternLen := len(pattern.Fields)
	if patternLen <= 0 {
		return false
	}

	pathLen := len(path)
	i := 0
	j := 0

	if !pattern.ShouldStartWithMe {
		found := false
		for ; i < patternLen && j < pathLen; j++ {
			if pattern.Fields[i] == path[j] {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	for i < patternLen {
		if j >= pathLen {
			return false
		}
		if pattern.Fields[i] != path[j] {
			return false
		}
		i++
		j++
	}

	if pattern.IsDirectory && j >= pathLen {
		return false
	}

	return true
}
