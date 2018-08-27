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

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use:     "node",
	Aliases: []string{"nodes"},
	Short:   "WiP*",
	Run: func(cmd *cobra.Command, args []string) {
		if stdinAvailable() {
			var kind string
			in, err := ioutil.ReadAll(os.Stdin)
			h(err)
			kind, args = parseStdin(in)
			switch kind {
			case "PODNAME":
				args = filterUnique(columnReturn(in, 5)[1:])
			}
		}
		if len(args) == 0 {
			args = []string{""}
		}
		rc, err := k8s.NewRawClient(false)
		h(err)
		var results []k8s.Results
		for _, nodes := range args {
			p, err := rc.GetNodes(nodes)
			h(err)
			results = append(results, p)
		}
		var nodes []Node
		var count int
		nodeChan := make(chan Node, 100)
		for _, r := range results {
			for _, x := range r.XData {
				go makeNodes(x, nodeChan)
			}
			count = count + len(r.XData)
		}
		for i := 0; i < count; i++ {
			node := <-nodeChan
			nodes = append(nodes, node)
		}
		sortSlice(nodes)
		formatTable(nodes)
		return

	},
}

func init() {
	rootCmd.AddCommand(nodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
