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
		results, err := rc.GetNodes(args[:]...)
		h(err)
		var nodes []Node
		nodeChan := make(chan Node, 100)
		for _, x := range results.XData {
			go makeNodes(x, nodeChan)
		}
		for i := 0; i < len(results.XData); i++ {
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

	// nodeCmd.PersistentFlags().String("foo", "", "A help for foo")
	// nodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
