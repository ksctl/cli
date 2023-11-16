package cmd

import (
	"os"
	"strings"

	"github.com/kubesimplify/ksctl/pkg/utils/consts"
	"github.com/spf13/cobra"
)

const (
	ksctl_feature_auto_scale   consts.KsctlSpecialFlags = "autoscale"
	ksctl_feature_applications consts.KsctlSpecialFlags = "application"
)

func featureFlag(f *cobra.Command) {
	f.Flags().StringP("feature-flags", "", "", `Experimental Features: Supported values with comma seperated: [autoscale,application]`)
	// f.Flags().StringArrayP("feature-flags", "", nil, `Supported values: [autoscale]`)
}

func SetRequiredFeatureFlags(cmd *cobra.Command) {
	rawFeatures, err := cmd.Flags().GetString("feature-flags")
	if err != nil {
		log.Error(err.Error())
		return
	}
	features := strings.Split(rawFeatures, ",")

	for _, feature := range features {

		switch consts.KsctlSpecialFlags(feature) {
		case ksctl_feature_auto_scale:
			if err := os.Setenv(string(consts.KsctlFeatureFlagHaAutoscale), "true"); err != nil {
				log.Error("Unable to set the ha autoscale feature")
			}

		case ksctl_feature_applications:
			if err := os.Setenv(string(consts.KsctlFeatureFlagApplications), "true"); err != nil {
				log.Error("Unable to set applications feature")
			}
		}
	}
}
