package manager

import (
	"fmt"
	"strings"
)

// parseCommand parses the content using the provided pattern.
//
// content: the string to be parsed
// pattern: the pattern to look for at the beginning of the content
// (string, string, error): the parsed command, parameter, and any error encountered
func parseCommand(content, pattern string) (string, string, error) {
	if !strings.HasPrefix(content, pattern) {
		return "", "", fmt.Errorf("pattern not found")
	}

	content = strings.TrimPrefix(content, pattern)

	spaceIndex := strings.Index(content, " ")
	if spaceIndex == -1 {
		return strings.ToLower(content), "", nil
	}

	command := strings.ToLower(content[:spaceIndex])
	parameter := strings.TrimSpace(content[spaceIndex+1:])
	return command, parameter, nil
}
