package gitignore

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	GITIGNORE = ".gitignore"
)

type Gitignore struct {
	m map[string][]Pattern
}

func parseGitignore(r io.Reader) ([]Pattern, error) {
	buf := bufio.NewReader(r)

	arr := make([]Pattern, 0)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return arr, err
		}

		line = strings.TrimSpace(line)
		if len(line) <= 0 || strings.HasPrefix(line, "#") {
			continue
		}

		m := NewSimplePattern(line)
		arr = append(arr, m)
	}

	return arr, nil
}

func ScanGitignore(root string) (*Gitignore, error) {
	m := make(map[string][]Pattern, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() != GITIGNORE {
			return nil
		}

		file, err := os.OpenFile(path, os.O_RDONLY, 0)
		if err != nil {
			return err
		}
		defer file.Close()

		patterns, err := parseGitignore(file)
		if err != nil {
			return err
		}

		s := len(root)
		e := len(path) - len(GITIGNORE)
		path = path[s:e]
		if path[0] != '/' {
			path = fmt.Sprintf("/%v", path)
		}
		m[path] = patterns

		return nil
	})

	/*
		for k, v := range m {
			log.Printf("%s -> %v", k, v)
			for _, p := range v {
				log.Print("  ", p)
			}
		}
	*/

	return &Gitignore{m}, err
}

func (g *Gitignore) IsIgnored(path string) bool {
	// XXX: never ignore the root path
	if len(path) <= 0 {
		return false
	}

	fields := strings.FieldsFunc(path, func(r rune) bool {
		return r == '/'
	})
	if len(fields) <= 0 {
		return true
	}

	prefix := "/"
	for i := 0; i < len(fields); i++ {
		// XXX Always ignore .git directory
		if fields[i] == ".git" {
			return true
		}

		if i >= 1 {
			prefix = prefix + fields[i-1] + "/"
		}

		patterns := g.m[prefix]
		if patterns == nil {
			continue
		}

		tail := fields[i:]
		for _, p := range patterns {
			if p.Match(tail) {
				return true
			}
		}
	}

	return false
}
