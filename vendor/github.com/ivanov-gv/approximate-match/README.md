# approximate-match

Let's assume we have a list of the words, and we want to find the closest match to another one - sample.
For example:

aaaaa - our sample

List:
aaaab
aaaac
abcde
bbbbb


Let's define the closest to sample word as:
1. It has the longest common substring, but with length of this substring has to be more than 0
2. It has the same set of letters or, at least, this set of letters has minimal difference with sample's set

As in our example before:
aaaab  - has substring "aaaa", 1 less of letters a, 1 more of letters b. So 4 - is a lenght of a common substring, 2 - set difference
aaaac  - the same as before - 4 and 2
abcde  - length - 1 , 8 - difference (-4 'a', + 'b', + 'c', + 'd', + 'e' )
bbbbb  - 0 common letters, difference is 10

According to our definition "aaaab", "aaaac" and "abcde" are our candidates, but "bbbbb" is not.

Let's use it for searching another closest match.
Our next sample - Ella
List: Adele Elaine Elizabeth Harriet Ingrid Michelle Ella

Output:
Ella      - exact match, no differences
Michelle  - has the same substring 'ell', but has 5+ more letters and doesn't have 1 letter 'a'
Adele
Elaine
Elizabeth
Harriet

and no word "Ingrid" in output list, because it doesn't match on any letter.

## Usage

```go
m := NewMatcher([]string{"Adele", "Elaine", "Elizabeth", "Harriet", "Ingrid", "Michelle", "Ella"}, nil)
for _, match := range m.Find("Ella") {
    fmt.Printf("%s\t%.3f\n", match.Word, match.Score)
}
```

`NewMatcher` preprocesses the word list once. Call `Find` for each query — it returns matches sorted by score (1.0 = identical, 0 = no commonality). Pass a custom threshold as the second argument to control sensitivity; nil uses the default (0.45).

---

## Real-world use case: railway station matching in Montenegro and the Balkans

The library was extended from the original interview exercise to serve as the station-name resolution layer for a Montenegrin railway application. Users search for train stations by typing station names in various forms — different scripts, dialects, and transliteration styles — and the matcher must find the right station regardless.

See: https://github.com/ivanov-gv/zpcg

### The problem

The Montenegrin rail network connects stations across Montenegro, Serbia, Bosnia, and Albania. A single station may be known by several legitimate spellings:

| Station | Latin (official) | Cyrillic | English | Common misspelling |
|---------|-----------------|----------|---------|-------------------|
| Bijelo Polje | Bijelo Polje | Бијело Поље | Bijelo Polje | belo pole, belo polje |
| Nikšić | Nikšić | Никшић | Niksic | nickshicsh |
| Herceg Novi | Herceg Novi | Херцег Нови | Herceg Novi | hertzeg novi |
| Podgorica | Podgorica | Подгорица | Podgorica | padgareeka, podgoritsa |

Users may type in:
- **Serbian ekavica or ijekavica** — "belo pole" and "bijelo polje" are the same place
- **Cyrillic script** — Serbian, or Russian transliteration conventions (soft signs, ю/ы variants)
- **Latin diacritics or their ASCII approximations** — Nikšić vs Niksic
- **German/English phonetics** — "sh", "ch", "sch", "ph", "th" clusters
- **Minor typos** — transpositions, missing letters, doubled vowels

### How the matcher handles this

`Normalize()` collapses all equivalent forms to a single canonical representation before indexing or querying:

1. **Unicode NFD** — strips combining diacritical marks automatically (š→s, č→c, ž→z, ć→c, and hundreds of others)
2. **Explicit rules** — `đ→d` (no NFD decomposition), `w→v` (Germanic)
3. **Dialect normalisation** — `ije→e` collapses ijekavica into ekavica
4. **Digraph folding** — `lj→l`, `nj→n`, `dj→d`, `sh→s`, `ch→c`, `sch→s`, `ph→f`, `th→t`
5. **Cyrillic normalisation** — Serbian ligatures (љ→л, њ→н), Russian soft-sign sequences (ль→л, нь→н), ю→у, ы→и
6. **Double-consonant reduction** — `bb→b`, `ss→s`, etc.
7. **Space removal** — "novi sad" and "novisad" are treated identically

`ConsonantSkeleton()` then strips all vowels from the normalised form. When the full-form match score is low, the matcher falls back to comparing consonant skeletons — this handles severe vowel confusion like "padgareeka" → `pdgrk` ≈ `pdgrc` ← "podgorica".

### Scoring

Each candidate is scored against the query as:

```
score = lcsRatio × (1 − unmatchedRatio)
```

where `lcsRatio` is the longest common substring length divided by the longer string's length, and `unmatchedRatio` is the sum of absolute character-frequency differences divided by total character count. The final score is the maximum of the full-form score and `0.9 × skeleton score`.

Results below the score threshold (default 0.45) are filtered out, which eliminates unrelated city names ("london", "berlin") while retaining weak but legitimate phonetic matches.
