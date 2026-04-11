# Development Guidelines

This project follows guidelines from https://github.com/ivanov-gv/claude-code-setup/tree/main/.claude/shared/guidelines 

See the files for Project structure convention, Naming conventions and others. 
Start with https://github.com/ivanov-gv/claude-code-setup/tree/main/.claude/CLAUDE.md


## Codebase Overview

This is a fuzzy string matching library for station/location names, designed to handle linguistic variation: diacritics, phonetic equivalences, dialect variants (Serbian ekavica/ijekavica), Cyrillic transliterations, and minor typos. The public API is:

```go
matcher := approximatematch.NewMatcher(wordList, nil)  // nil = use DefaultScoreThreshold
results := matcher.Find(query)  // []Match, sorted by Score descending
```

### Matching Algorithm

`NewMatcher` preprocesses each word in the list:
1. Runs `Normalize()` — Unicode NFD, diacritic stripping, lowercase, phonetic substitutions (see below)
2. Runs `ConsonantSkeleton()` — strips vowels from the normalized form
3. Calls `buildRuneStats()` — maps each rune to its frequency and all substrings starting at that position, and returns the total rune count; both are stored in `indexedWord` for reuse across every `Find` call

`Find` scores every candidate against the query using `matchScore()`, which:
- Computes the **longest common substring (LCS)** byte length between normalized forms
- Computes the absolute character-frequency difference directly from both precomputed stats maps (no intermediate allocation)
- Combines them: `score = lcsRatio * (1 - unmatchedRatio)`, where `lcsRatio = lcs/longerByteLen` and `unmatchedRatio = absDiff/totalRuneCount`
- Runs the same computation on the **consonant skeletons** (weighted by `skeletonMatchWeight = 0.90`)
- Takes the max of the two scores

**Unit note:** `lcsRatio` uses byte lengths (consistent with `lenPrefix` which returns byte offsets); `unmatchedRatio` uses rune counts (consistent with the per-rune frequency stats). Both are valid [0, 1] proportions — do not "unify" them without careful measurement.

Results below `scoreThreshold` are filtered out and the rest are returned sorted descending.

### Normalization Pipeline (`normalize.go`)

`Normalize()` applies these steps in order:
1. Unicode NFD decomposition + remove category-`Mn` nonspacing marks (handles most diacritics)
2. Remove spaces, lowercase everything
3. Multi-char phonetic substitutions (applied in order — longer patterns first):
   - Slavic: `ije → e`, `lj → l`, `nj → n`, `dj → d`, `đ → d` (no NFD decomposition)
   - Germanic: `w → v`
   - Foreign clusters: `sch → s`, `sh → s`, `zh → z`, `ch → c`, `ph → f`, `th → t`, `ck → k`
   - Vowel collapses: `ee → i`, `oo → u`, `ou → u`
   - Double consonant reduction: `bb → b`, `cc → c`, … `zz → z`
   - Cyrillic: `ль → л`, `нь → н`, `ь/ъ → ∅`, `ю → у`, `ы → и`, `љ → л`, `њ → н`, `ћ → ч`, `ђ → д`, `ј → и`

`ConsonantSkeleton()` takes an **already-normalized** string and strips all vowels (`a e i o u` and their Cyrillic equivalents). It does **not** call `Normalize()` — callers are responsible for running `Normalize()` first.

`transform.String` can fail on invalid UTF-8. `Normalize()` handles this by falling back to the original input rather than propagating the error or using a partial result. The public API stays `func Normalize(input string) string`.

## Dependencies

- `golang.org/x/text` — Unicode normalization (`norm`, `runes`, `transform`).
- `github.com/stretchr/testify` — test assertions (`assert`, `require`).
- `github.com/samber/lo` — generic slice/map utilities (`lo.Keys`, `lo.Filter`, `lo.Without`, `lo.Uniq`, …).
