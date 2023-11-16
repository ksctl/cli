// authors Dipankar <dipankar@dipankar-das.com>
package cmd

import (
	"fmt"
	"os"

	control_pkg "github.com/kubesimplify/ksctl/pkg/controllers"
	"github.com/kubesimplify/ksctl/pkg/utils/consts"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var credCmd = &cobra.Command{
	Use:   "cred",
	Short: "Login to your Cloud-provider Credentials",
	Long: `login to your cloud provider credentials
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		if err := control_pkg.InitializeStorageFactory(&cli.Client); err != nil {
			log.Error("Inialize Storage Driver", "Reason", err)
		}
		SetRequiredFeatureFlags(cmd)

		log.Print(`
1> AWS (EKS)
2> Azure (AKS)
3> Civo (K3s)
`)

		choice := 0

		_, err := fmt.Scanf("%d", &choice)
		if err != nil {
			panic(err.Error())
		}
		if provider, ok := cloud[choice]; ok {
			cli.Client.Metadata.Provider = consts.KsctlCloud(provider)
		} else {
			log.Error("invalid provider")
		}
		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout

		if err := controller.Credentials(&cli.Client); err != nil {
			log.Error("Failed to added the credential", "Reason", err)
			os.Exit(1)
		}
		log.Success("Credentials added successfully")
	},
}

func init() {
	rootCmd.AddCommand(credCmd)
	credCmd.Flags().BoolP("verbose", "v", true, "for verbose output")

}
