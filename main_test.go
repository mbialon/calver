package main

import (
	"strings"
	"testing"

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
