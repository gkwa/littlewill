package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromWalmartURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Walmart URL with tracking parameters",
			input:    "https://www.walmart.com/ip/Great-Value-Mango-Chunks-Frozen-16-oz/31712522?classType=REGULAR&athbdg=L1102&from=/search",
			expected: "https://www.walmart.com/ip/Great-Value-Mango-Chunks-Frozen-16-oz/31712522",
		},
		{
			name:     "Walmart URL without parameters is unchanged",
			input:    "https://www.walmart.com/ip/Great-Value-Mango-Chunks-Frozen-16-oz/31712522",
			expected: "https://www.walmart.com/ip/Great-Value-Mango-Chunks-Frozen-16-oz/31712522",
		},
		{
			name:     "Non-Walmart URL with same parameters is unchanged",
			input:    "https://www.example.com/product/123?classType=REGULAR&from=/search",
			expected: "https://www.example.com/product/123?classType=REGULAR&from=/search",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer

			err := RemoveParamsFromWalmartURLs(input, &output)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			result := output.String()
			if diff := cmp.Diff(tc.expected, result); diff != "" {
				t.Errorf("Unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}
