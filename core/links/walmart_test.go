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
		{
			name:     "Walmart ad campaign URL strips all ad-tracking params and keeps functional params",
			input:    "https://www.walmart.com/ip/Raindrops-Rotating-Shelf-Wall-Mounted-Swivel-Towel-Rail-Bathroom-Clothes-Rack-Coat/16114150174?adid=22222222297D492DB7E57C8399FB4537265C3FFA506_0000000000_21407473164&cn=FY25-ENTP-PMAX_cnv_dps_dsn_dis_ad_entp_e_n&conditionGroupCode=1&gclsrc=aw.ds&selectedOfferId=D492DB7E57C8399FB4537265C3FFA506&selectedSellerId=102510489&veh=sem&wl0=&wl1=g&wl10=5398217977&wl11=online&wl12=D492DB7E57C8399FB4537265C3FFA506&wl2=c&wl3=&wl4=&wl5=9015483&wl6=&wl7=&wl8=&wl9=pla&wmlspartner=wlpa&wmlspartner=wlpa",
			expected: "https://www.walmart.com/ip/Raindrops-Rotating-Shelf-Wall-Mounted-Swivel-Towel-Rail-Bathroom-Clothes-Rack-Coat/16114150174?conditionGroupCode=1&selectedOfferId=D492DB7E57C8399FB4537265C3FFA506&selectedSellerId=102510489",
		},
		{
			name:     "Walmart URL with wl-prefixed ad label params are removed",
			input:    "https://www.walmart.com/ip/some-product/12345?wl0=&wl1=g&wl2=c&wl9=pla&wl10=5398217977&wl11=online&wl12=ABCD",
			expected: "https://www.walmart.com/ip/some-product/12345",
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
