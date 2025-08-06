//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	injectBeginMarker = "// matcha:inject-begin "
	injectEndMarker   = "// matcha:inject-end "
	injectedHeader    = "// Code injected by Matcha generator. DO NOT EDIT."
)

func main() {
	snippets := make(map[string][]string)

	// First pass: collect all export blocks
	_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		data, _ := os.ReadFile(path)
		lines := strings.Split(string(data), "\n")

		var collecting bool
		var name string
		var block []string

		for _, line := range lines {
			if strings.HasPrefix(line, "// matcha:export ") {
				collecting = true
				name = strings.TrimSpace(strings.TrimPrefix(line, "// matcha:export "))
				block = nil
				continue
			}
			if strings.HasPrefix(line, "// matcha:end") && collecting {
				snippets[name] = append([]string(nil), block...)
				collecting = false
				continue
			}
			if collecting {
				block = append(block, line)
			}
		}
		return nil
	})

	// Second pass: inject snippets
	_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		originalData, _ := os.ReadFile(path)
		lines := strings.Split(string(originalData), "\n")
		var output []string
		changed := false

		for i := 0; i < len(lines); i++ {
			line := lines[i]

			if strings.HasPrefix(line, "// matcha:import ") {
				name := strings.TrimSpace(strings.TrimPrefix(line, "// matcha:import "))
				code, ok := snippets[name]
				if !ok {
					fmt.Fprintf(os.Stderr, "warning: snippet %q not found for %s\n", name, path)
					output = append(output, line)
					continue
				}

				// Skip old injected block if it exists
				j := i + 1
				if j < len(lines) && lines[j] == injectedHeader {
					j++
					for j < len(lines) {
						if strings.TrimSpace(lines[j]) == injectEndMarker+name {
							j++
							break
						}
						j++
					}
					i = j - 1 // set i to end of old block
					changed = true
				}

				// Write new block
				output = append(output, line)
				output = append(output, injectedHeader)
				output = append(output, injectBeginMarker+name)
				output = append(output, code...)
				output = append(output, injectEndMarker+name)
				continue
			}

			output = append(output, line)
		}

		if changed || strings.Join(output, "\n") != string(originalData) {
			return os.WriteFile(path, []byte(strings.Join(output, "\n")), 0644)
		}

		return nil
	})
}
