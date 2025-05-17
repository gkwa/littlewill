package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromTechCrunchURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "TechCrunch email URL with tracking parameters",
			input:    "https://email.techcrunch.com/coinbase-gets-hacked?ecid=ACsprvvVVSUJyBE8uiAbYcNXHHDbdsTM6Pk3U-p8roL3rVUXUUL8jIa3qPVobOn0rf1I1SahrunG&utm_campaign=Week%20in%20Review&utm_medium=email&_hsenc=p2ANqtz-88FH_esTJ4bG0z6k5JuPuqumhKhlBLcqE3liRMxfkQaWnxnBmmhR7Gys85l4BP0hjaNEWynl2OE9E6PMMYv3kjGv-GMvK3Ai9FCOOFgbh1eph51QA&_hsmi=361875404&utm_content=361875404&utm_source=hs_email",
			expected: "https://email.techcrunch.com/coinbase-gets-hacked",
		},
		{
			name:     "TechCrunch email URL with some tracking parameters",
			input:    "https://email.techcrunch.com/some-article?utm_source=newsletter&param=value",
			expected: "https://email.techcrunch.com/some-article?param=value",
		},
		{
			name:     "TechCrunch URL without tracking parameters",
			input:    "https://email.techcrunch.com/article?param=value",
			expected: "https://email.techcrunch.com/article?param=value",
		},
		{
			name:     "Non-TechCrunch URL",
			input:    "https://example.com?utm_source=newsletter",
			expected: "https://example.com?utm_source=newsletter",
		},
		{
			name: "Multiple URLs in text",
			input: `Check out these links:
https://email.techcrunch.com/coinbase-gets-hacked?ecid=ABC123&utm_source=newsletter
https://example.com/article?id=123&utm_source=newsletter`,
			expected: `Check out these links:
https://email.techcrunch.com/coinbase-gets-hacked
https://example.com/article?id=123&utm_source=newsletter`,
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this URL: https://email.techcrunch.com/article?utm_source=newsletter

` + "```" + `
var url = "https://email.techcrunch.com/article?utm_source=newsletter";
` + "```" + `

Another URL: https://email.techcrunch.com/another?utm_source=twitter`,
			expected: `Check this URL: https://email.techcrunch.com/article

` + "```" + `
var url = "https://email.techcrunch.com/article?utm_source=newsletter";
` + "```" + `

Another URL: https://email.techcrunch.com/another`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveParamsFromTechCrunchURLs(input, &output)
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
