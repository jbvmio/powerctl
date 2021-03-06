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
	"os"

	"github.com/jbvmio/ak8s"
	"github.com/jbvmio/powerctl/cmd/kube"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	errd error
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "powerctl",
	Short: "WiP*",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		kube.Client, errd = ak8s.NewClientFromConfig(kube.CFG)
		if errd != nil {
			fmt.Println(errd)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ROOT")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&kube.NS, "namespace", "n", "", "namespace")
	//rootCmd.PersistentFlags().BoolVarP(&exactMatches, "exact", "x", false, "return exact matches")
	//rootCmd.PersistentFlags().StringVar(&CFG, "config", "", "config file (default is $HOME/.powerctl.yaml)")
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(cmdGet)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if kube.CFG != "" {
		viper.SetConfigFile(kube.CFG)
	} else {
		// Look for existing kube/config in default location.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		switch {
		case fileExists(home + `/.powerctl.yaml`):
			viper.AddConfigPath(home)
			viper.SetConfigName(".powerctl")
		case fileExists(home + `/.kube/config`):
			err := copyKubeConfig(home+`/.kube/config`, home+`/.powerctl.yaml`)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			viper.AddConfigPath(home)
			viper.SetConfigName(".powerctl")
		default:
			fmt.Println("Error: No config file specified or located.")
			os.Exit(1)
		}
		kube.CFG = home + `/.powerctl.yaml`
	}
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
