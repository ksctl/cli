package cmd

import (
	"context"
	"strings"

	"github.com/ksctl/ksctl/v2/pkg/helpers/consts"
	"github.com/ksctl/ksctl/v2/pkg/types"
	"github.com/spf13/cobra"
)

const (
	ksctl_feature_auto_scale consts.KsctlSpecialFlags = "autoscale"
)

func featureFlag(f *cobra.Command) {
	f.Flags().StringP("feature-flags", "", "", `Experimental Features: Supported values with comma seperated: [autoscale]`)
}

func SetRequiredFeatureFlags(ctx context.Context, log types.LoggerFactory, cmd *cobra.Command) {
	rawFeatures, err := cmd.Flags().GetString("feature-flags")
	if err != nil {
		log.Error("Error in setting feature flags", "Reason", err)
		return
	}
	features := strings.Split(rawFeatures, ",")

	for _, feature := range features {

		switch consts.KsctlSpecialFlags(feature) {
		// case ksctl_feature_auto_scale:
		// 	if err := os.Setenv(string(consts.KsctlFeatureFlagHaAutoscale), "true"); err != nil {
		// 		log.Error("Unable to set the ha autoscale feature")
		// 	}
		default:
			return
		}
	}
}
