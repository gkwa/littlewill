package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromShopifyURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Shopify product URL with all tracking parameters",
			input:    "https://shanassourdough.myshopify.com/products/duralex-lys-bowl?pr_prod_strat=jac&pr_rec_id=c79d8a927&pr_rec_pid=6901841690677&pr_ref_pid=6771797491765&pr_seq=uniform",
			expected: "https://shanassourdough.myshopify.com/products/duralex-lys-bowl",
		},
		{
			name:     "Shopify URL with some tracking parameters",
			input:    "https://example.myshopify.com/products/test-product?pr_prod_strat=jac&pr_rec_id=123456&other=keep",
			expected: "https://example.myshopify.com/products/test-product?other=keep",
		},
		{
			name:     "Shopify URL with UTM parameters",
			input:    "https://store.myshopify.com/products/item?utm_source=newsletter&utm_medium=email&pr_seq=uniform",
			expected: "https://store.myshopify.com/products/item",
		},
		{
			name:     "Shopify URL with mixed parameters",
			input:    "https://demo.myshopify.com/products/demo?id=123&pr_prod_strat=test&variant=456&pr_rec_id=abc&color=red",
			expected: "https://demo.myshopify.com/products/demo?color=red&id=123&variant=456",
		},
		{
			name:     "Shopify URL without tracking parameters",
			input:    "https://store.myshopify.com/products/item?variant=123&color=blue",
			expected: "https://store.myshopify.com/products/item?color=blue&variant=123",
		},
		{
			name:     "Non-Shopify URL with pr_ parameters",
			input:    "https://example.com?pr_prod_strat=shouldnotberemoved&pr_rec_id=keep",
			expected: "https://example.com?pr_prod_strat=shouldnotberemoved&pr_rec_id=keep",
		},
		{
			name:     "Shopify main domain URL",
			input:    "https://shopify.com/partners?pr_prod_strat=test&utm_source=ads",
			expected: "https://shopify.com/partners",
		},
		{
			name: "Multiple Shopify URLs in text",
			input: `Check out these products:
https://shanassourdough.myshopify.com/products/duralex-lys-bowl?pr_prod_strat=jac&pr_rec_id=c79d8a927&pr_rec_pid=6901841690677&pr_ref_pid=6771797491765&pr_seq=uniform
https://store.myshopify.com/products/another?pr_rec_id=xyz123&utm_campaign=summer`,
			expected: `Check out these products:
https://shanassourdough.myshopify.com/products/duralex-lys-bowl
https://store.myshopify.com/products/another`,
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this Shopify URL: https://shanassourdough.myshopify.com/products/duralex-lys-bowl?pr_prod_strat=jac&pr_rec_id=c79d8a927
` + "```" + `
var shopifyUrl = "https://shanassourdough.myshopify.com/products/duralex-lys-bowl?pr_prod_strat=jac&pr_rec_id=c79d8a927";
` + "```" + `
Another Shopify URL: https://demo.myshopify.com/products/test?pr_seq=uniform&pr_ref_pid=123456`,
			expected: `Check this Shopify URL: https://shanassourdough.myshopify.com/products/duralex-lys-bowl
` + "```" + `
var shopifyUrl = "https://shanassourdough.myshopify.com/products/duralex-lys-bowl?pr_prod_strat=jac&pr_rec_id=c79d8a927";
` + "```" + `
Another Shopify URL: https://demo.myshopify.com/products/test`,
		},
		{
			name:     "Shopify URL with individual tracking parameters",
			input:    "https://test.myshopify.com/products/item1?pr_prod_strat=test&keep=this&pr_rec_pid=456",
			expected: "https://test.myshopify.com/products/item1?keep=this",
		},
		{
			name:     "Shopify URL with pr_ref_pid only",
			input:    "https://example.myshopify.com/collections/all?pr_ref_pid=6771797491765&sort_by=price",
			expected: "https://example.myshopify.com/collections/all?sort_by=price",
		},
		{
			name:     "Shopify URL with pr_seq parameter",
			input:    "https://store.myshopify.com/cart?pr_seq=uniform&discount=SAVE10",
			expected: "https://store.myshopify.com/cart?discount=SAVE10",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer

			err := RemoveParamsFromShopifyURLs(input, &output)
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
