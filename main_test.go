package main

import (
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestScanner(t *testing.T) {
	scanner := NewScanner(strings.NewReader(`v2020.01.1
v2020.01.2
v2020.02.1`))
	var got []Version
	for scanner.Scan() {
		got = append(got, scanner.Version())
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}

	want := []Version{
		{2020, 1, 1},
		{2020, 1, 2},
		{2020, 2, 1},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestVersionBump(t *testing.T) {
	cases := []struct {
		Now     time.Time
		Version Version
		Want    Version
	}{
		{
			Now: time.Date(2020, 3, 4, 0, 0, 0, 0, time.Local),
			Version: Version{
				Major(2020),
				Minor(3),
				Patch(2),
			},
			Want: Version{
				Major(2020),
				Minor(3),
				Patch(3),
			},
		},
		{
			Now: time.Date(2020, 3, 4, 0, 0, 0, 0, time.Local),
			Version: Version{
				Major(2020),
				Minor(2),
				Patch(3),
			},
			Want: Version{
				Major(2020),
				Minor(3),
				Patch(1),
			},
		},
	}
	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			got := tc.Version.Bump(tc.Now)
			if diff := cmp.Diff(tc.Want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
