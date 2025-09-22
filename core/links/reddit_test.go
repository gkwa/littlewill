package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromRedditURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Reddit post with share_id parameter",
			input:    "https://www.reddit.com/r/Sourdough/comments/1lb7mga/crumb_read_customer_complains_my_sourdough_rises/?share_id=mHY19WpSs_UfZsKQhbBrx",
			expected: "https://www.reddit.com/r/Sourdough/comments/1lb7mga/crumb_read_customer_complains_my_sourdough_rises/",
		},
		{
			name:     "Reddit search with cId and iId parameters",
			input:    "https://www.reddit.com/search/?q=sourdough+discard+scallion+pancakes&cId=f99a9b72-a2aa-4582-81db-bb5b366b6895&iId=8fc793d9-7555-4dab-8b72-cfcdb061bf1d",
			expected: "https://www.reddit.com/search/?q=sourdough+discard+scallion+pancakes",
		},
		{
			name:     "Reddit short URL with share_id",
			input:    "https://redd.it/1lb7mga?share_id=mHY19WpSs_UfZsKQhbBrx",
			expected: "https://redd.it/1lb7mga",
		},
		{
			name:     "Reddit URL with UTM parameters using shared logic",
			input:    "https://www.reddit.com/r/test/comments/789?utm_source=newsletter&utm_medium=email&sort=new",
			expected: "https://www.reddit.com/r/test/comments/789?sort=new",
		},
		{
			name:     "Reddit URL with share_id and UTM parameters",
			input:    "https://www.reddit.com/r/programming/comments/123456?share_id=abc123&utm_source=newsletter&utm_medium=email",
			expected: "https://www.reddit.com/r/programming/comments/123456",
		},
		{
			name:     "Reddit URL with cId, iId and other parameters",
			input:    "https://www.reddit.com/search/?q=test&cId=abc-123&iId=def-456&sort=new&other=keep",
			expected: "https://www.reddit.com/search/?other=keep&q=test&sort=new",
		},
		{
			name:     "Reddit URL with share_id and other parameters",
			input:    "https://www.reddit.com/r/test/comments/789?share_id=xyz&sort=new&other=keep",
			expected: "https://www.reddit.com/r/test/comments/789?other=keep&sort=new",
		},
		{
			name:     "Reddit URL without tracking parameters",
			input:    "https://www.reddit.com/r/golang/comments/456789/some_post/",
			expected: "https://www.reddit.com/r/golang/comments/456789/some_post/",
		},
		{
			name:     "Non-Reddit URL with share_id",
			input:    "https://example.com?share_id=shouldnotberemoved",
			expected: "https://example.com?share_id=shouldnotberemoved",
		},
		{
			name:     "Non-Reddit URL with cId and iId",
			input:    "https://example.com?cId=shouldnotberemoved&iId=alsokeep",
			expected: "https://example.com?cId=shouldnotberemoved&iId=alsokeep",
		},
		{
			name:     "Reddit URL with Branch.io parameters",
			input:    "https://www.reddit.com/r/esp32/comments/1nl6fs9/built_my_own_esp32_rc_cars_firmware_opensource/?%243p=e_as&%24deep_link=true&_branch_match_id=1076535883026722802&_branch_referrer=H4sIAAAAAAAAA22Q207DMAyGn6a76zaaMlGkCSEQr2GlidtZ5CTHVeGGZ8fdJq6QEuXP9%2FuUXERKfT4cGL0n2dtS9oHS58GUl6brTTkj2LpTmZlmSjbAwuF82bIa89p0H7rWdd3f812OClg31mI6PZVETFJVPqRwmuqgalwoCMRvyGuCaySwA2e5wkQcV8sIuWCqeWGHWxejjR673iMW2AZszLvwotbJZWYMVignIK%2FcOttPw4DtaMantn%2Fq%2BnbAo2sHP6De%2FLG3TvNKrgLTEkKyEbdyBv4GvJmUPH6pc1TAOKnCaCmApxmr3KAOHYulOf3v3h9w9xQuEsHlJfolSq9thCTg7kfDkZnSDCPntSKf3y6cI%2F4CSgQps6IBAAA%3D&correlation_id=aca4f99e-b3b8-4824-9e0c-9d9e482d04ac&post_fullname=t3_1nl6fs9&post_index=0&ref_campaign=email_digest&ref_source=email",
			expected: "https://www.reddit.com/r/esp32/comments/1nl6fs9/built_my_own_esp32_rc_cars_firmware_opensource/",
		},
		{
			name:     "Reddit URL with individual Branch.io parameters",
			input:    "https://www.reddit.com/r/test/comments/123?%243p=e_as&correlation_id=aca4f99e-b3b8-4824-9e0c-9d9e482d04ac&other=keep",
			expected: "https://www.reddit.com/r/test/comments/123?other=keep",
		},
		{
			name:     "Reddit URL with ref_campaign and ref_source",
			input:    "https://www.reddit.com/r/programming/comments/456?ref_campaign=email_digest&ref_source=email&sort=top",
			expected: "https://www.reddit.com/r/programming/comments/456?sort=top",
		},
		{
			name: "Multiple Reddit URLs in text",
			input: `Check out these Reddit posts:
https://www.reddit.com/r/Sourdough/comments/1lb7mga/crumb_read_customer_complains_my_sourdough_rises/?share_id=mHY19WpSs_UfZsKQhbBrx
https://www.reddit.com/search/?q=test&cId=abc-123&iId=def-456
https://redd.it/456789?share_id=another123&utm_source=app
https://www.reddit.com/r/esp32/comments/1nl6fs9/?%243p=e_as&%24deep_link=true&correlation_id=aca4f99e`,
			expected: `Check out these Reddit posts:
https://www.reddit.com/r/Sourdough/comments/1lb7mga/crumb_read_customer_complains_my_sourdough_rises/
https://www.reddit.com/search/?q=test
https://redd.it/456789
https://www.reddit.com/r/esp32/comments/1nl6fs9/`,
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this Reddit URL: https://www.reddit.com/r/test/comments/123?share_id=abc123&cId=def456&%243p=e_as
` + "```" + `
var redditUrl = "https://www.reddit.com/search/?q=test&cId=abc-123&iId=def-456&%24deep_link=true";
` + "```" + `
Another Reddit URL: https://redd.it/456?share_id=def456&utm_source=test&iId=ghi789&correlation_id=test123`,
			expected: `Check this Reddit URL: https://www.reddit.com/r/test/comments/123
` + "```" + `
var redditUrl = "https://www.reddit.com/search/?q=test&cId=abc-123&iId=def-456&%24deep_link=true";
` + "```" + `
Another Reddit URL: https://redd.it/456`,
		},
		{
			name:     "Reddit subdomain URL",
			input:    "https://old.reddit.com/r/programming/comments/123?share_id=test123&cId=uuid123&iId=uuid456&post_fullname=t3_123456&post_index=1",
			expected: "https://old.reddit.com/r/programming/comments/123",
		},
		{
			name:     "Reddit URL with custom UTM parameter",
			input:    "https://www.reddit.com/r/golang/?utm_custom_param=value&sort=new&cId=test&_branch_match_id=12345",
			expected: "https://www.reddit.com/r/golang/?sort=new",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer

			err := RemoveParamsFromRedditURLs(input, &output)
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
