package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromAmazonURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Amazon product URL with ref and th parameters",
			input:    "https://www.amazon.com/dp/B07F3VRTJ8?ref=ppx_yo2ov_dt_b_product_details&th=1",
			expected: "https://www.amazon.com/dp/B07F3VRTJ8",
		},
		{
			name:     "Amazon URL with only ref parameter",
			input:    "https://www.amazon.com/Seagate-BarraCuda-Internal-Drive-3-5-Inch/dp/B07H289S7C?ref=sr_1_3",
			expected: "https://www.amazon.com/Seagate-BarraCuda-Internal-Drive-3-5-Inch/dp/B07H289S7C",
		},
		{
			name:     "Amazon URL with only th parameter",
			input:    "https://www.amazon.com/dp/B07H289S7C?th=1",
			expected: "https://www.amazon.com/dp/B07H289S7C",
		},
		{
			name:     "Amazon URL with ref, th and other parameters",
			input:    "https://www.amazon.com/dp/B07H289S7C?crid=32O8AIOG7T1E6&ref=sr_1_3&th=1&keywords=sata+drive",
			expected: "https://www.amazon.com/dp/B07H289S7C?crid=32O8AIOG7T1E6&keywords=sata+drive",
		},
		{
			name:     "Amazon URL with UTM parameters",
			input:    "https://www.amazon.com/dp/B07H289S7C?utm_source=newsletter&utm_medium=email&ref=sr_1_3&th=1",
			expected: "https://www.amazon.com/dp/B07H289S7C",
		},
		{
			name:     "Amazon URL without tracking parameters",
			input:    "https://www.amazon.com/dp/B07H289S7C?keywords=sata+drive",
			expected: "https://www.amazon.com/dp/B07H289S7C?keywords=sata+drive",
		},
		{
			name:     "Amazon short URL with parameters",
			input:    "https://amzn.to/3abc123?ref=sr_1_3&th=1",
			expected: "https://amzn.to/3abc123",
		},
		{
			name:     "Amazon country-specific domain",
			input:    "https://www.amazon.co.uk/dp/B07H289S7C?ref=sr_1_3&th=1",
			expected: "https://www.amazon.co.uk/dp/B07H289S7C",
		},
		{
			name:     "Non-Amazon URL with ref and th parameters",
			input:    "https://example.com?ref=shouldnotberemoved&th=alsokeep",
			expected: "https://example.com?ref=shouldnotberemoved&th=alsokeep",
		},
		{
			name: "Multiple Amazon URLs in text",
			input: `Check out these products:
https://www.amazon.com/dp/B07F3VRTJ8?ref=ppx_yo2ov_dt_b_product_details&th=1
https://www.amazon.com/Seagate-BarraCuda-Internal-Drive-3-5-Inch/dp/B07H289S7C?ref=sr_1_3&crid=32O8AIOG7T1E6&th=1
https://amzn.to/3abc123?ref=sr_1_3&th=1`,
			expected: `Check out these products:
https://www.amazon.com/dp/B07F3VRTJ8
https://www.amazon.com/Seagate-BarraCuda-Internal-Drive-3-5-Inch/dp/B07H289S7C?crid=32O8AIOG7T1E6
https://amzn.to/3abc123`,
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this Amazon URL: https://www.amazon.com/dp/B07F3VRTJ8?ref=ppx_yo2ov_dt_b_product_details&th=1

` + "```" + `
var amazonUrl = "https://www.amazon.com/dp/B07F3VRTJ8?ref=ppx_yo2ov_dt_b_product_details&th=1";
` + "```" + `

Another Amazon URL: https://www.amazon.com/dp/B07H289S7C?ref=sr_1_3&th=1`,
			expected: `Check this Amazon URL: https://www.amazon.com/dp/B07F3VRTJ8

` + "```" + `
var amazonUrl = "https://www.amazon.com/dp/B07F3VRTJ8?ref=ppx_yo2ov_dt_b_product_details&th=1";
` + "```" + `

Another Amazon URL: https://www.amazon.com/dp/B07H289S7C`,
		},
		{
			name:     "Amazon URL with complex ref parameter",
			input:    "https://www.amazon.com/dp/B07H289S7C?ref=ppx_yo2ov_dt_b_product_details&crid=32O8AIOG7T1E6&dib=eyJ2IjoiMSJ9.test&th=1",
			expected: "https://www.amazon.com/dp/B07H289S7C?crid=32O8AIOG7T1E6&dib=eyJ2IjoiMSJ9.test",
		},
		{
			name:     "Amazon search URL with ref and th",
			input:    "https://www.amazon.com/s?k=laptop&ref=nb_sb_noss&th=1",
			expected: "https://www.amazon.com/s?k=laptop",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveParamsFromAmazonURLs(input, &output)
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
