package timetable

import (
	"testing"
	"time"
)

func date(year int) time.Time {
	return time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
}

func season(start, end time.Time) Season {
	return Season{Start: start, End: end}
}

var (
	inf   = time.Time{} // zero == ±∞ per convention
	d2020 = date(2020)
	d2021 = date(2021)
	d2022 = date(2022)
	d2023 = date(2023)
)

func TestIntersects(t *testing.T) {
	tests := []struct {
		name  string
		s     Season
		other Season
		want  bool
	}{
		// ── infinite seasons ────────────────────────────────────────────
		{
			name:  "both fully infinite",
			s:     season(inf, inf),
			other: season(inf, inf),
			want:  true,
		},
		{
			name:  "s infinite, other finite",
			s:     season(inf, inf),
			other: season(d2020, d2021),
			want:  true,
		},
		{
			name:  "s finite, other infinite",
			s:     season(d2020, d2021),
			other: season(inf, inf),
			want:  true,
		},

		// ── half-infinite, disjoint ──────────────────────────────────────
		{
			name:  "s=(-inf,2021) other=(2022,+inf): gap between",
			s:     season(inf, d2021),
			other: season(d2022, inf),
			want:  false,
		},
		{
			name:  "s=(2022,+inf) other=(-inf,2021): gap between",
			s:     season(d2022, inf),
			other: season(inf, d2021),
			want:  false,
		},

		// ── half-infinite, overlapping ───────────────────────────────────
		{
			name:  "s=(-inf,2022) other=(2021,+inf): overlap",
			s:     season(inf, d2022),
			other: season(d2021, inf),
			want:  true,
		},
		{
			name:  "s=(2021,+inf) other=(-inf,2022): overlap",
			s:     season(d2021, inf),
			other: season(inf, d2022),
			want:  true,
		},

		// ── touching boundaries (inclusive) ─────────────────────────────
		{
			name:  "s ends exactly where other starts",
			s:     season(d2020, d2021),
			other: season(d2021, d2022),
			want:  true,
		},
		{
			name:  "other ends exactly where s starts",
			s:     season(d2021, d2022),
			other: season(d2020, d2021),
			want:  true,
		},

		// ── strictly disjoint ───────────────────────────────────────────
		{
			name:  "s strictly before other",
			s:     season(d2020, d2021),
			other: season(d2022, d2023),
			want:  false,
		},
		{
			name:  "other strictly before s",
			s:     season(d2022, d2023),
			other: season(d2020, d2021),
			want:  false,
		},

		// ── partial overlap ──────────────────────────────────────────────
		{
			name:  "s starts before other, ends inside",
			s:     season(d2020, d2022),
			other: season(d2021, d2023),
			want:  true,
		},
		{
			name:  "other starts before s, ends inside",
			s:     season(d2021, d2023),
			other: season(d2020, d2022),
			want:  true,
		},

		// ── containment ─────────────────────────────────────────────────
		{
			name:  "s contains other",
			s:     season(d2020, d2023),
			other: season(d2021, d2022),
			want:  true,
		},
		{
			name:  "other contains s",
			s:     season(d2021, d2022),
			other: season(d2020, d2023),
			want:  true,
		},
		{
			name:  "identical seasons",
			s:     season(d2020, d2021),
			other: season(d2020, d2021),
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Intersects(tt.other); got != tt.want {
				t.Errorf("Intersects() = %v, want %v", got, tt.want)
			}
		})

		// Intersects must be symmetric: A∩B ⟺ B∩A
		t.Run(tt.name+"/symmetric", func(t *testing.T) {
			if tt.s.Intersects(tt.other) != tt.other.Intersects(tt.s) {
				t.Errorf("symmetry violated: s.Intersects(other) != other.Intersects(s)")
			}
		})
	}
}

func TestIsInSeason(t *testing.T) {
	tests := []struct {
		name string
		s    Season
		t    time.Time
		want bool
	}{
		// ── fully infinite ───────────────────────────────────────────────
		{
			name: "infinite season contains any time",
			s:    season(inf, inf),
			t:    d2021,
			want: true,
		},

		// ── inclusive boundaries ─────────────────────────────────────────
		{
			name: "exactly on start",
			s:    season(d2020, d2022),
			t:    d2020,
			want: true,
		},
		{
			name: "exactly on end",
			s:    season(d2020, d2022),
			t:    d2022,
			want: true,
		},

		// ── inside / outside ─────────────────────────────────────────────
		{
			name: "strictly inside",
			s:    season(d2020, d2022),
			t:    d2021,
			want: true,
		},
		{
			name: "strictly before start",
			s:    season(d2021, d2023),
			t:    d2020,
			want: false,
		},
		{
			name: "strictly after end",
			s:    season(d2020, d2021),
			t:    d2022,
			want: false,
		},

		// ── half-infinite ────────────────────────────────────────────────
		{
			name: "no start bound: time before end",
			s:    season(inf, d2022),
			t:    d2020,
			want: true,
		},
		{
			name: "no start bound: time after end",
			s:    season(inf, d2022),
			t:    d2023,
			want: false,
		},
		{
			name: "no end bound: time after start",
			s:    season(d2020, inf),
			t:    d2022,
			want: true,
		},
		{
			name: "no end bound: time before start",
			s:    season(d2020, inf),
			t:    d2021, // 2021 > 2020, still in
			want: true,
		},
		{
			name: "no end bound: time strictly before start",
			s:    season(d2022, inf),
			t:    d2020,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsInSeason(tt.t); got != tt.want {
				t.Errorf("IsInSeason(%v) = %v, want %v", tt.t, got, tt.want)
			}
		})
	}
}

func TestBefore(t *testing.T) {
	tests := []struct {
		name  string
		s     Season
		other Season
		want  bool
	}{
		// ── infinite end never precedes anything ─────────────────────────
		{
			name:  "s has no end: never before",
			s:     season(d2020, inf),
			other: season(d2023, inf),
			want:  false,
		},
		{
			name:  "both infinite: not before",
			s:     season(inf, inf),
			other: season(inf, inf),
			want:  false,
		},

		// ── strict ordering ──────────────────────────────────────────────
		{
			name:  "s ends before other starts",
			s:     season(d2020, d2021),
			other: season(d2022, d2023),
			want:  true,
		},
		{
			name:  "s ends after other starts",
			s:     season(d2020, d2023),
			other: season(d2021, d2022),
			want:  false,
		},

		// ── touching boundary is not "before" ────────────────────────────
		{
			name:  "s ends exactly when other starts",
			s:     season(d2020, d2021),
			other: season(d2021, d2022),
			want:  false,
		},

		// ── other has no start (−∞) ───────────────────────────────────────
		{
			name:  "other starts at -inf: s.End cannot precede it",
			s:     season(d2020, d2021),
			other: season(inf, d2022),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Before(tt.other); got != tt.want {
				t.Errorf("Before() = %v, want %v", got, tt.want)
			}
		})
	}
}
