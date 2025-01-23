package util

import (
	"regexp"
	"strings"
)

func NormalizeAsaVersion(version string) string {
	re := regexp.MustCompile(`(\d+)\.(\d+)\(([\d|.]+)\)\(?(\d*)\)?`)
	matches := re.FindStringSubmatch(version)
	if len(matches) == 0 {
		return version
	}
	parts := []string{matches[1], matches[2]}
	if matches[3] != "" {
		parts = append(parts, matches[3])
	}
	if matches[4] != "" {
		parts = append(parts, matches[4])
	}

	return strings.Join(parts, ".")
}

func DoNormalisedVersionsMatch(a string, b string) bool {
	normalisedA := NormalizeAsaVersion(a)
	normalisedB := NormalizeAsaVersion(b)

	return normalisedA == normalisedB
}
