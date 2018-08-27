// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"io/ioutil"
	"os"

	"github.com/jbvmio/k8s"
	"github.com/spf13/cobra"
)

// podCmd represents the pod command
var podCmd = &cobra.Command{
	Use:     "pod",
	Aliases: []string{"pods"},
	Short:   "WiP*",
	Run: func(cmd *cobra.Command, args []string) {
		if stdinAvailable() {
			var kind string
			in, err := ioutil.ReadAll(os.Stdin)
			h(err)
			kind, args = parseStdin(in)
			switch kind {
			case "NODENAME":
				var xdata []k8s.XD
				allPods := getAllPods()
				for _, pod := range allPods.XData {
					for _, node := range args {
						if pod.NodeName == node {
							xdata = append(xdata, pod)
						}
					}
				}
				makePrintPods(xdata)
				return
			}
		}
		if len(args) == 0 {
			args = []string{""}
		}
		rc, err := k8s.NewRawClient(false)
		h(err)
		rc.SetNS(targetNamespace)
		results, err := rc.GetPods(args[:]...)
		h(err)
		makePrintPods(results.XData)
		return
	},
}

func init() {
	rootCmd.AddCommand(podCmd)

	// podCmd.PersistentFlags().String("foo", "", "A help for foo")
	// podCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
