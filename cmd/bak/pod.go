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
	"fmt"
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
			case "NoResultsFound":
				fmt.Println("NoResultsFound")
				os.Exit(1)
			case "NODENAME":
				makePrintPods(nodesToPods(args))
				return
			case "REPLICASET":
				makePrintPods(rsToPods(args))
				return
			case "DEPLOYMENT":
				var replicasetNames []string
				rs := deploysToRS(args)
				for _, r := range rs {
					replicasetNames = append(replicasetNames, r.Name)
				}
				makePrintPods(rsToPods(filterUnique(replicasetNames)))
				return
			default:
				fmt.Println("NoResultsFound")
				os.Exit(1)
			}
		}
		if len(args) == 0 {
			args = []string{""}
		}
		rc, err := k8s.NewRawClient(false)
		h(err)
		rc.SetNS(targetNamespace)
		rc.ExactMatches(exactMatches)
		results, err := rc.GetPods(args[:]...)
		h(err)
		validateResults(results)
		makePrintPods(results.XData)
		return
	},
}

func init() {
	rootCmd.AddCommand(podCmd)

	// podCmd.PersistentFlags().String("foo", "", "A help for foo")
	// podCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}