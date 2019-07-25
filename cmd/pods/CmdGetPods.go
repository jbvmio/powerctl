package pods

import (
	"fmt"
	"os"

	"github.com/jbvmio/powerctl/cmd/kube"
	"github.com/spf13/cobra"
)

//var outFlags out.OutFlags

// CmdGetPods gets pods.
var CmdGetPods = &cobra.Command{
	Use:     "pods",
	Example: "  powerctl get pod",
	Short:   "Get Pod Info",
	Run: func(cmd *cobra.Command, args []string) {
		pods, err := kube.Client.GetAllPods()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		switch {
		case len(args) > 0:
			pods = pods.Search(args...)
		}
		for _, p := range pods.Items {
			fmt.Println(p.Name)
		}
		/*
			switch true {
			case len(args) > 0:
				out.Failf("No such resource: %v", args[0])
			default:
				cmd.Help()
			}
		*/
	},
}

func init() {
	//CmdGet.PersistentFlags().StringVarP(&outFlags.Format, "out", "o", "", "Change Output Format - yaml|json.")

	//CmdGet.AddCommand(pods.CmdGetBroker)

}
