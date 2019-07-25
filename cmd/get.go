package cmd

import (
	"github.com/jbvmio/powerctl/cmd/pods"
	"github.com/spf13/cobra"
)

var cmdGet = &cobra.Command{
	Use:     "get",
	Example: "  powerctl get pod",
	Short:   "Get K8s Info",
	Run: func(cmd *cobra.Command, args []string) {
		switch true {
		case len(args) > 0:
			//out.Failf("No such resource: %v", args[0])
		default:
			cmd.Help()
		}
	},
}

func init() {
	//CmdGet.PersistentFlags().StringVarP(&outFlags.Format, "out", "o", "", "Change Output Format - yaml|json.")

	cmdGet.AddCommand(pods.CmdGetPods)

}
