package main

import (
	"testing"

	"github.com/kubernetes/helm/pkg/repo"
)

const testDir = "testdata/"
const testFile = "testdata/local-index.yaml"

type searchTestCase struct {
	in          string
	expectedOut []string
}

var searchTestCases = []searchTestCase{
	{"foo", []string{}},
	{"alpine", []string{"alpine-1.0.0"}},
	{"sumtin", []string{"alpine-1.0.0"}},
	{"web", []string{"nginx-0.1.0"}},
}

var searchCacheTestCases = []searchTestCase{
	{"notthere", []string{}},
	{"odd", []string{"foobar/oddness-1.2.3.tgz"}},
	{"sumtin", []string{"local/alpine-1.0.0.tgz", "foobar/oddness-1.2.3.tgz"}},
	{"foobar", []string{"foobar/foobar-0.1.0.tgz"}},
	{"web", []string{"local/nginx-0.1.0.tgz"}},
}

func validateEntries(t *testing.T, in string, found []string, expected []string) {
	if len(found) != len(expected) {
		t.Errorf("Failed to search %s: Expected: %#v got: %#v", in, expected, found)
	}
	foundCount := 0
	for _, exp := range expected {
		for _, f := range found {
			if exp == f {
				foundCount = foundCount + 1
				continue
			}
		}
	}
	if foundCount != len(expected) {
		t.Errorf("Failed to find expected items for %s: Expected: %#v got: %#v", in, expected, found)
	}

}

func searchTestRunner(t *testing.T, tc searchTestCase) {
	cf, err := repo.LoadIndexFile(testFile)
	if err != nil {
		t.Errorf("Failed to load index file : %s : %s", testFile, err)
	}

	u := searchChartRefsForPattern(tc.in, cf.Entries)
	validateEntries(t, tc.in, u, tc.expectedOut)
}

func searchCacheTestRunner(t *testing.T, tc searchTestCase) {
	u, err := searchCacheForPattern(testDir, tc.in)
	if err != nil {
		t.Errorf("searchCacheForPattern failed: %#v", err)
	}
	validateEntries(t, tc.in, u, tc.expectedOut)
}

func TestSearches(t *testing.T) {
	for _, tc := range searchTestCases {
		searchTestRunner(t, tc)
	}
}

func TestCacheSearches(t *testing.T) {
	for _, tc := range searchCacheTestCases {
		searchCacheTestRunner(t, tc)
	}
}
