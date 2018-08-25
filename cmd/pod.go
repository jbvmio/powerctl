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
				stdInDetected = true
				return
			}
		}
		if len(args) == 0 {
			args = []string{""}
		}
		rc, err := k8s.NewRawClient(false)
		h(err)
		rc.SetNS(targetNamespace)

		var results []k8s.Results
		for _, pods := range args {
			p, err := rc.GetPods(pods)
			h(err)
			results = append(results, p)
		}
		var pods []Pod
		var count int
		podChan := make(chan Pod, 100)
		for _, r := range results {
			for _, x := range r.XData {
				go makePods(x, podChan)
			}
			count = count + len(r.XData)
		}
		for i := 0; i < count; i++ {
			pod := <-podChan
			pods = append(pods, pod)
		}
		sortSlice(pods)
		formatTable(pods)
		return
	},
}

func init() {
	rootCmd.AddCommand(podCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// podCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// podCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
