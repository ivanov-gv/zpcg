package approxmatch

import "sort"

// skeletonMatchWeight slightly discounts skeleton-based matches vs. full-form
// matches, because a consonant skeleton discards vowel information.
const skeletonMatchWeight = 0.90

// DefaultScoreThreshold is the minimum score a candidate must reach to be
// included in Find results when no explicit threshold is provided to NewMatcher.
// It is set just above the typical score for common transliteration noise
// (e.g. "london" → 0.40 against Balkan station names), while staying below
// the score produced by single-character-deletion typos (e.g. "belgade" → 0.48).
const DefaultScoreThreshold = 0.45

// indexedWord holds all precomputed representations of one entry in the search list.
type indexedWord struct {
	original            string
	normalized          string
	skeleton            string
	normalizedStats     map[rune]RuneStat
	skeletonStats       map[rune]RuneStat
	normalizedRuneCount int
	skeletonRuneCount   int
}

// Match is a single result from Matcher.Find.
type Match struct {
	Word  string
	Score float64
}

// Matcher holds a fixed search list with precomputed statistics.
// Construct once with NewMatcher, then call Find for every user query.
type Matcher struct {
	words          []indexedWord
	scoreThreshold float64
}

// NewMatcher builds and returns a Matcher for the given fixed word list.
// All heavy preprocessing happens here so that each Find call is fast.
// threshold sets the minimum score for a result to be returned by Find;
// pass nil to use DefaultScoreThreshold.
func NewMatcher(words []string, threshold *float64) *Matcher {
	scoreThreshold := DefaultScoreThreshold
	if threshold != nil {
		scoreThreshold = *threshold
	}
	indexed := make([]indexedWord, len(words))
	for wordIndex, word := range words {
		normalized := Normalize(word)
		skeleton := ConsonantSkeleton(normalized)
		normalizedStats, normalizedRuneCount := buildRuneStats(normalized)
		skeletonStats, skeletonRuneCount := buildRuneStats(skeleton)
		indexed[wordIndex] = indexedWord{
			original:            word,
			normalized:          normalized,
			skeleton:            skeleton,
			normalizedStats:     normalizedStats,
			skeletonStats:       skeletonStats,
			normalizedRuneCount: normalizedRuneCount,
			skeletonRuneCount:   skeletonRuneCount,
		}
	}
	return &Matcher{words: indexed, scoreThreshold: scoreThreshold}
}

// Find returns all entries from the search list ranked by similarity to sample,
// best first. Entries whose score falls below the threshold configured in
// NewMatcher are omitted; pass a lower threshold to include weaker matches.
func (m *Matcher) Find(sample string) []Match {
	normalizedSample := Normalize(sample)
	skeletonSample := ConsonantSkeleton(normalizedSample)

	normalizedSampleStats, normalizedSampleRuneCount := buildRuneStats(normalizedSample)
	skeletonSampleStats, skeletonSampleRuneCount := buildRuneStats(skeletonSample)

	results := make([]Match, 0, len(m.words)/2)

	for _, entry := range m.words {
		normalizedScore := matchScore(normalizedSample, normalizedSampleStats, normalizedSampleRuneCount, entry.normalized, entry.normalizedStats, entry.normalizedRuneCount)
		skeletonScore := matchScore(skeletonSample, skeletonSampleStats, skeletonSampleRuneCount, entry.skeleton, entry.skeletonStats, entry.skeletonRuneCount) * skeletonMatchWeight

		score := normalizedScore
		if skeletonScore > score {
			score = skeletonScore
		}
		if score >= m.scoreThreshold {
			results = append(results, Match{Word: entry.original, Score: score})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}

// matchScore returns a value in [0, 1] measuring how closely word matches
// sample. Both sampleStats and wordStats must be precomputed by buildRuneStats,
// with their respective rune counts.
// 0 means no common substring was found (or either input is empty).
func matchScore(sample string, sampleStats map[rune]RuneStat, sampleRuneCount int, word string, wordStats map[rune]RuneStat, wordRuneCount int) float64 {
	if len(sample) == 0 || len(word) == 0 {
		return 0
	}

	var longestCommonSubstr int
	longestCommonSubstrIsLeading := false
	for byteOffset, char := range word {
		sampleStat, found := sampleStats[char]
		if !found {
			continue
		}
		if prefixLen := lenPrefix(word[byteOffset:], sampleStat.substrings...); prefixLen > longestCommonSubstr {
			longestCommonSubstr = prefixLen
			longestCommonSubstrIsLeading = byteOffset == 0
		}
	}

	if longestCommonSubstr == 0 {
		return 0
	}

	// Normalise LCS length against the longer of the two strings.
	longerLen := max(len(word), len(sample))
	lcsRatio := float64(longestCommonSubstr) / float64(longerLen)

	// When the sample is entirely a leading prefix of the word (both strings start
	// with the same sequence and the query is fully consumed), normalise against the
	// sample length rather than the longer length. This prevents a shorter word that
	// merely contains the query as an interior substring from outscoring a longer word
	// that starts with the full query (e.g. "beograd" should prefer "beogradcentar"
	// over "novibeograd"). The bonus only activates when word is strictly longer than
	// sample; otherwise lcsRatio is already 1.0 and leadingRatio cannot exceed it.
	if longestCommonSubstrIsLeading && longestCommonSubstr == len(sample) {
		leadingRatio := float64(longestCommonSubstr) / float64(len(sample))
		if leadingRatio > lcsRatio {
			lcsRatio = leadingRatio
		}
	}

	// Penalise by how many characters are unaccounted for (relative to total rune count).
	// Compute the abs-diff sum directly from both stats maps to avoid allocating an
	// intermediate delta map on every call.
	totalUnmatchedChars := 0
	for char, sampleStat := range sampleStats {
		delta := sampleStat.count - wordStats[char].count
		if delta > 0 {
			totalUnmatchedChars += delta
		} else if delta < 0 {
			totalUnmatchedChars -= delta
		}
	}
	for char, wordStat := range wordStats {
		if _, found := sampleStats[char]; !found {
			totalUnmatchedChars += wordStat.count
		}
	}
	unmatchedRatio := float64(totalUnmatchedChars) / float64(sampleRuneCount+wordRuneCount)

	penalty := unmatchedRatio
	if penalty > 1 {
		penalty = 1
	}
	return lcsRatio * (1.0 - penalty)
}
