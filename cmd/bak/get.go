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
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var (
	getTarget string
	getArgs   []string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "WiP*",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			log.Fatal("Get Args?\n")
		}
		getTarget = args[0]
		getArgs = args[1:]
		if strings.Contains(getTarget, "pod") {
			podCmd.Run(cmd, getArgs)
			return
		}
		if strings.Contains(getTarget, "node") {
			nodeCmd.Run(cmd, getArgs)
			return
		}
		if strings.Contains(getTarget, "rs") || strings.Contains(getTarget, "replica") {
			rsCmd.Run(cmd, getArgs)
			return
		}
		if strings.Contains(getTarget, "dep") {
			deployCmd.Run(cmd, getArgs)
			return
		}
		// Leave to ensure return statements exist
		fmt.Println("GET DONE")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
