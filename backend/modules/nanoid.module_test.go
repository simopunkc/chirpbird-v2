package modules

import (
	"testing"
)

func almostEqual(v1, v2 int) bool {
	return v1 == v2
}

type testCase struct {
	prefix   string
	expected int
}

func TestGenerateUniqueID(t *testing.T) {
	testCases := []testCase{
		{"a", 11},
		{"aa", 12},
		{"aaa", 13},
		{"aaaa", 14},
	}

	for _, tc := range testCases {
		t.Run("parameter "+tc.prefix, func(t *testing.T) {
			out := GenerateUniqueID(tc.prefix)
			if !almostEqual(len(out), tc.expected) {
				t.Fatalf("%d != %d", len(out), tc.expected)
			}
		})
	}
}
