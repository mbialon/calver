package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	bump := flag.Bool("bump", false, "bump patch")
	flag.Parse()

	vers, err := Scan()
	if err != nil {
		log.Fatal(err)
	}
	if len(vers) == 0 {
		return
	}
	sort.Sort(sort.Reverse(versionSlice(vers)))
	ver := vers[0]
	if *bump {
		ver = ver.Bump(time.Now())
	}
	fmt.Println(ver.String())
}

type versionSlice []Version

func (s versionSlice) Len() int {
	return len(s)
}

func (s versionSlice) Less(i, j int) bool {
	a, b := s[i], s[j]
	if a.Major < b.Major {
		return true
	}
	if a.Minor < b.Minor {
		return true
	}
	return a.Patch < b.Patch
}

func (s versionSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type (
	Major int
	Minor int
	Patch int
)

type Version struct {
	Major Major
	Minor Minor
	Patch Patch
}

func (v Version) String() string {
	return fmt.Sprintf("v%d.%02d.%d", v.Major, v.Minor, v.Patch)
}

func (v Version) Bump(now time.Time) Version {
	major := Major(now.Year())
	minor := Minor(now.Month())
	patch := v.Patch
	if minor > v.Minor {
		patch = 1
	} else {
		patch++
	}
	return Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

func Scan() ([]Version, error) {
	var versions []Version
	b, err := exec.Command("git", "tag", "--list", "v*").Output()
	if err != nil {
		return nil, err
	}
	scanner := NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		versions = append(versions, scanner.Version())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return versions, nil
}

type Scanner struct {
	*bufio.Scanner
	err error
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{Scanner: bufio.NewScanner(r)}
}

func (s *Scanner) Scan() bool {
	if s.err != nil {
		return false
	}
	return s.Scanner.Scan()
}

func (s *Scanner) Version() Version {
	text := strings.TrimPrefix(s.Text(), "v")
	parts := strings.SplitN(text, ".", 3)
	if len(parts) != 3 {
		s.err = errors.New("invalid format")
		return Version{}
	}
	segment := make([]int, len(parts))
	for i, part := range parts {
		v, err := strconv.Atoi(part)
		if err != nil {
			s.err = err
			return Version{}
		}
		segment[i] = v
	}
	return Version{
		Major: Major(segment[0]),
		Minor: Minor(segment[1]),
		Patch: Patch(segment[2]),
	}
}

func (s *Scanner) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.Scanner.Err()
}
