package cmd

import (
	"io"

	"github.com/gkwa/littlewill/core/links"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TransformDefinition defines a URL transformation with its metadata
type TransformDefinition struct {
	Name           string
	ConfigKey      string
	FlagName       string
	Description    string
	Function       func(io.Reader, io.Writer) error
	DefaultEnabled bool
}

// AllTransforms is the single source of truth for all available transformations
var AllTransforms = []TransformDefinition{
	{
		Name:           "generic-tracking",
		ConfigKey:      "transforms.generic_tracking",
		FlagName:       "enable-generic-tracking",
		Description:    "Enable generic tracking parameter removal",
		Function:       links.RemoveGenericTrackingParams,
		DefaultEnabled: true,
	},
	{
		Name:           "google",
		ConfigKey:      "transforms.google",
		FlagName:       "enable-google",
		Description:    "Enable Google URL parameter removal",
		Function:       links.RemoveParamsFromGoogleURLs,
		DefaultEnabled: true,
	},
	{
		Name:           "youtube",
		ConfigKey:      "transforms.youtube",
		FlagName:       "enable-youtube",
		Description:    "Enable YouTube URL parameter removal",
		Function:       links.RemoveParamsFromYouTubeURLs,
		DefaultEnabled: true,
	},
	{
		Name:           "substack",
		ConfigKey:      "transforms.substack",
		FlagName:       "enable-substack",
		Description:    "Enable Substack URL parameter removal",
		Function:       links.RemoveParamsFromSubstackURLs,
		DefaultEnabled: true,
	},
	{
		Name:           "thesweekly",
		ConfigKey:      "transforms.thesweekly",
		FlagName:       "enable-thesweekly",
		Description:    "Enable TheSweekly URL parameter removal",
		Function:       links.RemoveParamsFromTheSweeklyURLs,
		DefaultEnabled: true,
	},
	{
		Name:           "techcrunch",
		ConfigKey:      "transforms.techcrunch",
		FlagName:       "enable-techcrunch",
		Description:    "Enable TechCrunch URL parameter removal",
		Function:       links.RemoveParamsFromTechCrunchURLs,
		DefaultEnabled: true,
	},
	{
		Name:           "linkedin",
		ConfigKey:      "transforms.linkedin",
		FlagName:       "enable-linkedin",
		Description:    "Enable LinkedIn URL parameter removal",
		Function:       links.RemoveParamsFromLinkedInURLs,
		DefaultEnabled: true,
	},
	{
		Name:           "mailchimp",
		ConfigKey:      "transforms.mailchimp",
		FlagName:       "enable-mailchimp",
		Description:    "Enable Mailchimp URL parameter removal",
		Function:       links.RemoveParamsFromMailchimpURLs,
		DefaultEnabled: true,
	},
	{
		Name:           "wsj",
		ConfigKey:      "transforms.wsj",
		FlagName:       "enable-wsj",
		Description:    "Enable WSJ URL parameter removal",
		Function:       links.RemoveParamsFromWSJURLs,
		DefaultEnabled: true,
	},
	{
		Name:           "reddit",
		ConfigKey:      "transforms.reddit",
		FlagName:       "enable-reddit",
		Description:    "Enable Reddit URL parameter removal",
		Function:       links.RemoveParamsFromRedditURLs,
		DefaultEnabled: true,
	},
	{
		Name:           "shopify",
		ConfigKey:      "transforms.shopify",
		FlagName:       "enable-shopify",
		Description:    "Enable Shopify URL parameter removal",
		Function:       links.RemoveParamsFromShopifyURLs,
		DefaultEnabled: true,
	},
	{
		Name:           "amazon",
		ConfigKey:      "transforms.amazon",
		FlagName:       "enable-amazon",
		Description:    "Enable Amazon URL parameter removal",
		Function:       links.RemoveParamsFromAmazonURLs,
		DefaultEnabled: true,
	},
	{
		Name:           "conditional",
		ConfigKey:      "transforms.conditional",
		FlagName:       "enable-conditional",
		Description:    "Enable conditional parameter removal",
		Function:       links.RemoveConditionalParams,
		DefaultEnabled: true,
	},
	{
		Name:           "text-fragments",
		ConfigKey:      "transforms.text_fragments",
		FlagName:       "enable-text-fragments",
		Description:    "Enable text fragment removal",
		Function:       links.RemoveTextFragmentsFromURLs,
		DefaultEnabled: true,
	},
	{
		Name:           "youtube-count",
		ConfigKey:      "transforms.youtube_count",
		FlagName:       "enable-youtube-count",
		Description:    "Enable YouTube count removal from markdown links",
		Function:       links.RemoveYouTubeCountFromMarkdownLinks,
		DefaultEnabled: true,
	},
}

// buildLinkTransforms creates the list of enabled transformations based on configuration
func buildLinkTransforms() []func(io.Reader, io.Writer) error {
	var transforms []func(io.Reader, io.Writer) error

	for _, transform := range AllTransforms {
		if viper.GetBool(transform.ConfigKey) {
			transforms = append(transforms, transform.Function)
		}
	}

	return transforms
}

// setupTransformFlags adds flags and config bindings for all transforms
func setupTransformFlags(cmd *cobra.Command) {
	for _, transform := range AllTransforms {
		cmd.PersistentFlags().Bool(transform.FlagName, transform.DefaultEnabled, transform.Description)
		viper.BindPFlag(transform.ConfigKey, cmd.PersistentFlags().Lookup(transform.FlagName))
	}
}

// setTransformDefaults sets default values for all transforms in viper
func setTransformDefaults() {
	for _, transform := range AllTransforms {
		viper.SetDefault(transform.ConfigKey, transform.DefaultEnabled)
	}
}
