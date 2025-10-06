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
			input:    "https://www.amazon.com/dp/B08N5WRWNW?ref=ppx_yo2ov_dt_b_fed_asin_title&th=1",
			expected: "https://www.amazon.com/dp/B08N5WRWNW",
		},
		{
			name:     "Amazon URL with multiple tracking parameters",
			input:    "https://www.amazon.com/dp/B08N5WRWNW?keywords=laptop&qid=1234567890&sr=8-1&ref=sr_1_1",
			expected: "https://www.amazon.com/dp/B08N5WRWNW",
		},
		{
			name:     "Amazon URL with UTM parameters",
			input:    "https://www.amazon.com/dp/B08N5WRWNW?utm_source=newsletter&utm_medium=email&utm_campaign=summer&variant=blue",
			expected: "https://www.amazon.com/dp/B08N5WRWNW?variant=blue",
		},
		{
			name:     "Amazon URL with search parameters",
			input:    "https://www.amazon.com/s?k=laptop&crid=ABC123&sprefix=lap%2Caps%2C123&ref=nb_sb_noss_1",
			expected: "https://www.amazon.com/s?k=laptop",
		},
		{
			name:     "Amazon URL with pd_ parameters",
			input:    "https://www.amazon.com/dp/B08N5WRWNW?pd_rd_i=B08N5WRWNW&pd_rd_r=abc-123&pd_rd_w=xyz&pd_rd_wg=def-456",
			expected: "https://www.amazon.com/dp/B08N5WRWNW",
		},
		{
			name:     "Amazon URL with pf_ parameters",
			input:    "https://www.amazon.com/dp/B08N5WRWNW?pf_rd_p=abc123&pf_rd_r=xyz789&other=keep",
			expected: "https://www.amazon.com/dp/B08N5WRWNW?other=keep",
		},
		{
			name:     "Amazon URL with encoding and content-id",
			input:    "https://www.amazon.com/dp/B08N5WRWNW?_encoding=UTF8&content-id=amzn1.sym.abc&th=1",
			expected: "https://www.amazon.com/dp/B08N5WRWNW",
		},
		{
			name:     "Amazon URL with sbo and sp_csd parameters",
			input:    "https://www.amazon.com/dp/B08N5WRWNW?sbo=RZvfv%2F%2FHxDF%2BO5021pAnSA%3D%3D&sp_csd=d2lkZ2V0TmFtZT1zcF9kZXRhaWw",
			expected: "https://www.amazon.com/dp/B08N5WRWNW",
		},
		{
			name:     "Amazon URL with cv_ct_cx parameter",
			input:    "https://www.amazon.com/dp/B08N5WRWNW?cv_ct_cx=laptop&keywords=laptop&pd_rd_i=B08N5WRWNW",
			expected: "https://www.amazon.com/dp/B08N5WRWNW",
		},
		{
			name:     "Amazon URL with ref_ prefix parameter",
			input:    "https://www.amazon.com/dp/B08N5WRWNW?ref_=nav_logo&th=1",
			expected: "https://www.amazon.com/dp/B08N5WRWNW",
		},
		{
			name:     "Amazon URL with mixed tracking and non-tracking parameters",
			input:    "https://www.amazon.com/dp/B08N5WRWNW?variant=123&ref=ppx_yo2ov&color=blue&th=1&size=large",
			expected: "https://www.amazon.com/dp/B08N5WRWNW?color=blue&size=large&variant=123",
		},
		{
			name:     "Amazon URL with mixed UTM and Amazon-specific parameters",
			input:    "https://www.amazon.com/dp/B08N5WRWNW?ref=ppx_yo2ov&utm_source=twitter&th=1&variant=blue&utm_campaign=promo",
			expected: "https://www.amazon.com/dp/B08N5WRWNW?variant=blue",
		},
		{
			name:     "Amazon URL without tracking parameters",
			input:    "https://www.amazon.com/dp/B08N5WRWNW?variant=123&color=blue",
			expected: "https://www.amazon.com/dp/B08N5WRWNW?color=blue&variant=123",
		},
		{
			name:     "Non-Amazon URL with similar parameters",
			input:    "https://example.com?ref=something&th=1&keywords=test",
			expected: "https://example.com?ref=something&th=1&keywords=test",
		},
		{
			name: "Multiple Amazon URLs in text",
			input: `Check out these products:
https://www.amazon.com/dp/B08N5WRWNW?ref=ppx_yo2ov_dt_b_fed_asin_title&th=1
https://www.amazon.com/s?k=laptop&crid=ABC123&sprefix=lap%2Caps%2C123&ref=nb_sb_noss_1`,
			expected: `Check out these products:
https://www.amazon.com/dp/B08N5WRWNW
https://www.amazon.com/s?k=laptop`,
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this Amazon URL: https://www.amazon.com/dp/B08N5WRWNW?ref=ppx_yo2ov&th=1

` + "```" + `
var amazonUrl = "https://www.amazon.com/dp/B08N5WRWNW?ref=ppx_yo2ov_dt_b_fed_asin_title&th=1";
` + "```" + `

Another Amazon URL: https://www.amazon.com/s?k=laptop&keywords=laptop&qid=1234567890&sr=8-1`,
			expected: `Check this Amazon URL: https://www.amazon.com/dp/B08N5WRWNW

` + "```" + `
var amazonUrl = "https://www.amazon.com/dp/B08N5WRWNW?ref=ppx_yo2ov_dt_b_fed_asin_title&th=1";
` + "```" + `

Another Amazon URL: https://www.amazon.com/s?k=laptop`,
		},
		{
			name:     "Amazon URL with all tracking parameter types",
			input:    "https://www.amazon.com/dp/B08N5WRWNW?_encoding=UTF8&content-id=amzn1&crid=ABC&cv_ct_cx=laptop&keywords=laptop&pd_rd_i=B08N5WRWNW&pd_rd_r=abc&pd_rd_w=xyz&pd_rd_wg=def&pf_rd_p=123&pf_rd_r=456&qid=1234567890&ref=sr_1_1&ref_=nav&sbo=test&sp_csd=test&sprefix=lap&sr=8-1&th=1&keep=this",
			expected: "https://www.amazon.com/dp/B08N5WRWNW?keep=this",
		},
		{
			name:     "Amazon subdomain URL",
			input:    "https://smile.amazon.com/dp/B08N5WRWNW?ref=ppx_yo2ov&th=1&variant=blue",
			expected: "https://smile.amazon.com/dp/B08N5WRWNW?variant=blue",
		},
		{
			name:     "Amazon international domain",
			input:    "https://www.amazon.co.uk/dp/B08N5WRWNW?ref=ppx_yo2ov&th=1&color=red",
			expected: "https://www.amazon.co.uk/dp/B08N5WRWNW?color=red",
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
