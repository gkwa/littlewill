package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromFacebookURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Facebook marketplace URL with referral tracking",
			input:    "https://www.facebook.com/marketplace/item/809426435533151/?referral_code=null&referral_story_type=post&tracking=browse_serp%3A5e3660a7-7186-4915-a72f-d570475178e5",
			expected: "https://www.facebook.com/marketplace/item/809426435533151/",
		},
		{
			name:     "Facebook URL without tracking parameters",
			input:    "https://www.facebook.com/marketplace/item/809426435533151",
			expected: "https://www.facebook.com/marketplace/item/809426435533151",
		},
		{
			name:     "Facebook reel URL with referral_source and surface params",
			input:    "https://www.facebook.com/reel/762313063202501/?referral_source=video_home&surface_type=tab&in_reels_tab_context=TRUE",
			expected: "https://www.facebook.com/reel/762313063202501/",
		},
		{
			name:     "Non-Facebook URL with referral_code is unchanged",
			input:    "https://example.com/signup?referral_code=FRIEND10",
			expected: "https://example.com/signup?referral_code=FRIEND10",
		},
		{
			name:     "Non-Facebook URL with tracking is unchanged",
			input:    "https://fedex.com/track?tracking=123456789",
			expected: "https://fedex.com/track?tracking=123456789",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveParamsFromFacebookURLs(input, &output)
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
