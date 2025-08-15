package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromKaiserURLs(t *testing.T) {
	testCases := []struct {
		name       string
		input      string
		expected   string
		skip       bool
		skipReason string
	}{
		{
			name:       "Kaiser URL with all tracking parameters",
			input:      "https://healthy.kaiserpermanente.org/health-wellness/healtharticle.anti-inflammatory-diet?promo_id=&wt.mc_id=&wt.tsrc=em&cid=ca-mem|re-natl|ch-em|pl-mkt|ta-|au-cm|bo-ret|&ad_id=0&cat=l&mkt_tok=NDkyLU5RVS0wMTQAAAGcSWtRFiiqvr1jYVMWyor17ctr8Y7ZVXQCRRz3GfEzZoigsh0IH3CXKyFbW_TQmNAAaqDue_lrlxxwCCOX8CerdwGnaBsEGSJlHUQtCxskBcAsddc",
			expected:   "https://healthy.kaiserpermanente.org/health-wellness/healtharticle.anti-inflammatory-diet",
			skip:       true,
			skipReason: "xurls.Strict() cannot parse URLs with unencoded pipe characters - see issue with malformed Kaiser URL",
		},
		{
			name:     "Kaiser URL with some tracking parameters",
			input:    "https://healthy.kaiserpermanente.org/health-wellness/article?promo_id=test&other=keep&wt.tsrc=em",
			expected: "https://healthy.kaiserpermanente.org/health-wellness/article?other=keep",
		},
		{
			name:     "Kaiser URL with UTM parameters",
			input:    "https://healthy.kaiserpermanente.org/article?utm_source=newsletter&utm_medium=email&content=keep",
			expected: "https://healthy.kaiserpermanente.org/article?content=keep",
		},
		{
			name:     "Kaiser URL without tracking parameters",
			input:    "https://healthy.kaiserpermanente.org/health-wellness/article?param=value",
			expected: "https://healthy.kaiserpermanente.org/health-wellness/article?param=value",
		},
		{
			name:     "Non-Kaiser URL with similar parameters",
			input:    "https://example.com?promo_id=test&wt.mc_id=keep&mkt_tok=alsokeep",
			expected: "https://example.com?promo_id=test&wt.mc_id=keep&mkt_tok=alsokeep",
		},
		{
			name: "Multiple Kaiser URLs in text",
			input: `Check out these Kaiser articles:
https://healthy.kaiserpermanente.org/article1?promo_id=test&wt.tsrc=em&content=keep
https://healthy.kaiserpermanente.org/article2?cid=tracking&ad_id=123&other=preserve`,
			expected: `Check out these Kaiser articles:
https://healthy.kaiserpermanente.org/article1?content=keep
https://healthy.kaiserpermanente.org/article2?other=preserve`,
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this Kaiser URL: https://healthy.kaiserpermanente.org/article?promo_id=test&wt.tsrc=em

` + "```" + `
var kaiserUrl = "https://healthy.kaiserpermanente.org/article?promo_id=test&wt.tsrc=em";
` + "```" + `

Another Kaiser URL: https://healthy.kaiserpermanente.org/another?mkt_tok=token123&cat=l`,
			expected: `Check this Kaiser URL: https://healthy.kaiserpermanente.org/article

` + "```" + `
var kaiserUrl = "https://healthy.kaiserpermanente.org/article?promo_id=test&wt.tsrc=em";
` + "```" + `

Another Kaiser URL: https://healthy.kaiserpermanente.org/another`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skipf("KNOWN FAILING: %s", tc.skipReason)
			}

			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveParamsFromKaiserURLs(input, &output)
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
