package helper

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ParseCookieFile parses a Netscape cookie file and returns cookies as name=value pairs.
func ParseCookiesFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open cookie file: %w", err)
	}
	defer file.Close()

	var cookies []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Handle HttpOnly cookies.
		if after, ok := strings.CutPrefix(line, "#HttpOnly_"); ok {
			line = after
		} else if strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Split(line, "\t")
		if len(fields) < 7 {
			continue
		}

		name := fields[5]
		value := fields[6]

		if name == "" {
			continue
		}

		cookies = append(cookies, name+"="+value)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse cookie file: %w", err)
	}

	return cookies, nil
}
