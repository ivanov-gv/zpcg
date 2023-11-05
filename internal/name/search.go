package name

import (
	"slices"
)

// RuneStat contains num of occurrences and all substrings starting with this letter
type RuneStat struct {
	num        int
	substrings []string
}

// Word contains fields for ranking word matches
type Word struct {
	word         string
	longestMatch int // len of longest matching substring
	diff         int // difference between word and the given sample
}

// lenPrefix returns number of the first matching runes in all given strings, len of the longest prefix
func lenPrefix(sample string, strings ...string) int {
	prefixEnds := make(map[int]bool, len(strings))

	for i, letter := range sample {
		match := false
		for j, str := range strings {
			if !prefixEnds[j] && i < len(str) && letter == rune(str[i]) {
				match = true
			} else {
				prefixEnds[j] = true
			}
		}

		if !match {
			return i
		}
	}

	return len(sample)
}

// diff returns summary difference on wordStats
func diff(wordStat map[rune]int) int {
	difference := 0
	for _, numRunes := range wordStat {
		if numRunes > 0 {
			difference += numRunes
		} else {
			difference -= numRunes
		}
	}
	return difference
}

// FindBestMatch returns words, first of them are more likely to be the same as sample
func FindBestMatch(sample string, searchList []string) string {
	sampleRunes := map[rune]RuneStat{}
	// define number of every letter in the sample, save every substring starting with that letter
	for i, char := range sample {
		stat, _ := sampleRunes[char]
		stat.num += 1
		stat.substrings = append(stat.substrings, sample[i:])
		sampleRunes[char] = stat
	}

	result := make([]Word, 0)
	for _, word := range searchList {
		// wordDiff shows difference between amount of every letter in sample and a given word
		wordDiff := map[rune]int{}

		for k, v := range sampleRunes {
			wordDiff[k] = v.num
		}

		var maxSubstrLen int
		for i, char := range word {
			wordDiff[char] -= 1

			if _, found := sampleRunes[char]; !found {
				continue
			}

			substrLen := lenPrefix(word[i:], sampleRunes[char].substrings...)
			if maxSubstrLen < substrLen {
				maxSubstrLen = substrLen
			}
		}

		// if no matching substring found - exclude the word from results
		if maxSubstrLen == 0 {
			continue
		}

		result = append(result, Word{
			word:         word,
			longestMatch: maxSubstrLen,
			diff:         diff(wordDiff),
		})
	}

	if len(result) == 0 { // no matches at all
		return ""
	}

	// find min of the result - the min word must have the longest matching substring and the lowest difference
	match := slices.MaxFunc(result, func(a, b Word) int {
		// longest match
		if a.longestMatch != b.longestMatch {
			if a.longestMatch > b.longestMatch {
				return 1
			} else {
				return -1
			}
		}
		// differences
		if a.diff != b.diff {
			if a.diff < b.diff {
				return 1
			} else {
				return -1
			}
		}
		return 0
	})
	return match.word
}
