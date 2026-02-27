package links

import (
	"bytes"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromYouTubeURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "YouTube link with si and app parameters",
			input:    "https://youtu.be/JSKJbGi5oNA?si=b2GkFDivckm1k-Mq&app=Desktop",
			expected: "https://youtu.be/JSKJbGi5oNA",
		},
		{
			name:     "YouTube link with only si parameter",
			input:    "https://www.youtube.com/watch?v=dQw4w9WgXcQ&si=AnotherParam",
			expected: "https://youtu.be/dQw4w9WgXcQ",
		},
		{
			name:     "YouTube link with only app parameter",
			input:    "https://youtu.be/JSKJbGi5oNA?app=Desktop&other=keep",
			expected: "https://youtu.be/JSKJbGi5oNA?other=keep",
		},
		{
			name:     "YouTube link without tracking parameters",
			input:    "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			expected: "https://youtu.be/dQw4w9WgXcQ",
		},
		{
			name:     "YouTube link with mixed parameters",
			input:    "https://www.youtube.com/watch?v=dQw4w9WgXcQ&si=test123&list=PLx2ksyallYzW4WNYHD9xOFrPRYGlntAft&app=mobile",
			expected: "https://youtu.be/dQw4w9WgXcQ?list=PLx2ksyallYzW4WNYHD9xOFrPRYGlntAft",
		},
		{
			name:     "Non-YouTube link with si parameter",
			input:    "https://example.com?si=shouldnotberemoved&app=keep",
			expected: "https://example.com?si=shouldnotberemoved&app=keep",
		},
		{
			name: "Multiple YouTube links",
			input: `
Check out this video: https://youtu.be/JSKJbGi5oNA?si=b2GkFDivckm1k-Mq
And this one: https://www.youtube.com/watch?v=dQw4w9WgXcQ&si=AnotherParam&app=Desktop
`,
			expected: `
Check out this video: https://youtu.be/JSKJbGi5oNA
And this one: https://youtu.be/dQw4w9WgXcQ
`,
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this YouTube URL: https://youtu.be/JSKJbGi5oNA?si=b2GkFDivckm1k-Mq&app=Desktop

` + "```" + `
var youtubeUrl = "https://www.youtube.com/watch?v=dQw4w9WgXcQ&si=AnotherParam&app=mobile";
` + "```" + `

Another YouTube URL: https://www.youtube.com/watch?v=test123&si=remove&app=remove`,
			expected: `Check this YouTube URL: https://youtu.be/JSKJbGi5oNA

` + "```" + `
var youtubeUrl = "https://www.youtube.com/watch?v=dQw4w9WgXcQ&si=AnotherParam&app=mobile";
` + "```" + `

Another YouTube URL: https://youtu.be/test123`,
		},
		{
			name:     "YouTube short link",
			input:    "https://youtu.be/dQw4w9WgXcQ?si=test123&app=mobile&other=keep",
			expected: "https://youtu.be/dQw4w9WgXcQ?other=keep",
		},
		{
			name:     "YouTube link with feature parameter",
			input:    "https://youtu.be/yIYEQHtmUIs?feature=youtu.be",
			expected: "https://youtu.be/yIYEQHtmUIs",
		},
		{
			name:     "YouTube link with feature and other parameters",
			input:    "https://youtu.be/yIYEQHtmUIs?feature=youtu.be&si=abc123&other=keep",
			expected: "https://youtu.be/yIYEQHtmUIs?other=keep",
		},
		{
			name:     "YouTube watch URL with feature parameter",
			input:    "https://www.youtube.com/watch?v=yIYEQHtmUIs&feature=youtu.be",
			expected: "https://youtu.be/yIYEQHtmUIs",
		},
		{
			name:     "YouTube thumbnail URL with sqp and rs parameters",
			input:    "https://i.ytimg.com/vi/dDugWzbzhZw/hq720.jpg?sqp=-oaymwEhCK4FEIIDSFryq4qpAxMIARUAAAAAGAElAADIQj0AgKJD&rs=AOn4CLDZsaOqJTdkzdzhoYn_mL-Y_caRfw",
			expected: "https://i.ytimg.com/vi/dDugWzbzhZw/hq720.jpg",
		},
		{
			name:     "YouTube thumbnail URL with only sqp parameter",
			input:    "https://i.ytimg.com/vi/dDugWzbzhZw/hq720.jpg?sqp=-oaymwEhCK4FEIIDSFryq4qpAxMIARUAAAAAGAElAADIQj0AgKJD",
			expected: "https://i.ytimg.com/vi/dDugWzbzhZw/hq720.jpg",
		},
		{
			name:     "YouTube thumbnail URL without tracking parameters",
			input:    "https://i.ytimg.com/vi/dDugWzbzhZw/hq720.jpg",
			expected: "https://i.ytimg.com/vi/dDugWzbzhZw/hq720.jpg",
		},
		{
			name:     "YouTube Shorts URL",
			input:    "https://www.youtube.com/shorts/Q05BvncrHSc",
			expected: "https://youtu.be/Q05BvncrHSc",
		},
		{
			name:     "YouTube Shorts URL with tracking parameters",
			input:    "https://www.youtube.com/shorts/Q05BvncrHSc?si=abc123&app=mobile",
			expected: "https://youtu.be/Q05BvncrHSc",
		},
		{
			name:     "YouTube Shorts URL preserving non-tracking parameters",
			input:    "https://www.youtube.com/shorts/Q05BvncrHSc?si=abc123&other=keep",
			expected: "https://youtu.be/Q05BvncrHSc?other=keep",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveParamsFromYouTubeURLs(input, &output)
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

func TestIsYouTubeURL(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "YouTube.com URL",
			input:    "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			expected: true,
		},
		{
			name:     "YouTube short URL",
			input:    "https://youtu.be/dQw4w9WgXcQ",
			expected: true,
		},
		{
			name:     "Non-YouTube URL",
			input:    "https://example.com/video",
			expected: false,
		},
		{
			name:     "YouTube subdomain",
			input:    "https://music.youtube.com/watch?v=dQw4w9WgXcQ",
			expected: true,
		},
		{
			name:     "YouTube image CDN",
			input:    "https://i.ytimg.com/vi/dDugWzbzhZw/hq720.jpg",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := url.Parse(tc.input)
			if err != nil {
				t.Fatalf("Failed to parse URL: %v", err)
			}
			result := isYouTubeURL(u)
			if result != tc.expected {
				t.Errorf("isYouTubeURL(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestRemoveYouTubeCountFromLinks(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "YouTube link with count",
			input:    "[(1619) Understanding Neovim #1 - Installation - YouTube](https://www.youtube.com/watch?v=87AXw9Quy9U&list=PLx2ksyallYzW4WNYHD9xOFrPRYGlntAft \"(1619) Understanding Neovim #1 - Installation - YouTube\")",
			expected: "[Understanding Neovim #1 - Installation - YouTube](https://www.youtube.com/watch?v=87AXw9Quy9U&list=PLx2ksyallYzW4WNYHD9xOFrPRYGlntAft \"(1619) Understanding Neovim #1 - Installation - YouTube\")",
		},
		{
			name:     "YouTube short link with count",
			input:    "[(2345) Short YouTube Video](https://youtu.be/dQw4w9WgXcQ)",
			expected: "[Short YouTube Video](https://youtu.be/dQw4w9WgXcQ)",
		},
		{
			name:     "YouTube short link with count and spaces",
			input:    "[      (2345)        Short YouTube Video](https://youtu.be/dQw4w9WgXcQ)",
			expected: "[Short YouTube Video](https://youtu.be/dQw4w9WgXcQ)",
		},
		{
			name:     "YouTube link without count",
			input:    "[Understanding Neovim #1 - Installation - YouTube](https://www.youtube.com/watch?v=87AXw9Quy9U&list=PLx2ksyallYzW4WNYHD9xOFrPRYGlntAft)",
			expected: "[Understanding Neovim #1 - Installation - YouTube](https://www.youtube.com/watch?v=87AXw9Quy9U&list=PLx2ksyallYzW4WNYHD9xOFrPRYGlntAft)",
		},
		{
			name:     "Non-YouTube link",
			input:    "[Example Link](https://example.com)",
			expected: "[Example Link](https://example.com)",
		},
		{
			name: "Multiple YouTube links",
			input: `
[(1234) First Video - YouTube](https://www.youtube.com/watch?v=abc123)
[(5678) Second Video - YouTube](https://youtu.be/def456 "Some title")
[Regular Link](https://example.com)
`,
			expected: `
[First Video - YouTube](https://www.youtube.com/watch?v=abc123)
[Second Video - YouTube](https://youtu.be/def456 "Some title")
[Regular Link](https://example.com)
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveYouTubeCountFromMarkdownLinks(input, &output)
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
