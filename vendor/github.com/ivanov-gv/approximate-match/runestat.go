package approxmatch

import "unicode/utf8"

// RuneStat holds the occurrence count of a rune in a string and all substrings
// of that string starting at each occurrence position.
type RuneStat struct {
	count      int
	substrings []string
}

// lenPrefix returns the byte length of the longest common prefix between sample
// and any of the given candidate strings.
// Each candidate tracks its own byte offset so multi-byte UTF-8 characters
// (e.g. Cyrillic, 2 bytes each) are compared correctly by rune, not by byte.
func lenPrefix(sample string, candidates ...string) int {
	candidateByteOffsets := make([]int, len(candidates))
	candidateEnded := make([]bool, len(candidates))

	for sampleByteOffset, sampleRune := range sample {
		matched := false
		for candidateIndex, candidate := range candidates {
			if candidateEnded[candidateIndex] {
				continue
			}
			if candidateByteOffsets[candidateIndex] >= len(candidate) {
				candidateEnded[candidateIndex] = true
				continue
			}
			candidateRune, runeSize := utf8.DecodeRuneInString(candidate[candidateByteOffsets[candidateIndex]:])
			if sampleRune == candidateRune {
				matched = true
				candidateByteOffsets[candidateIndex] += runeSize
			} else {
				candidateEnded[candidateIndex] = true
			}
		}
		if !matched {
			return sampleByteOffset
		}
	}
	return len(sample)
}

// buildRuneStats returns, for every character in input, its count and all
// substrings of input starting at that character's position, plus the total
// rune count of input. Used by the substring-prefix matching in matchScore.
func buildRuneStats(input string) (map[rune]RuneStat, int) {
	stats := make(map[rune]RuneStat, len(input))
	runeCount := 0
	for byteOffset, char := range input {
		stat := stats[char]
		stat.count++
		stat.substrings = append(stat.substrings, input[byteOffset:])
		stats[char] = stat
		runeCount++
	}
	return stats, runeCount
}
