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
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, partnerUserID, partnerUserSecret string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "justexpenseit",
	Short: "Expensify CLI client",
	Long:  `A smart Expensify client that does `,
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.justexpenseit.yaml)")
	rootCmd.PersistentFlags().StringVar(&partnerUserID, "expensify-user-id", "", "User ID to log into Expensify API")
	rootCmd.PersistentFlags().StringVar(&partnerUserSecret, "expensify-user-secret", "", "User secret to log into Expensify API")

	rootCmd.AddCommand(NewTotalCmd())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cmd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".justexpenseit")
	}

	// read in environment variables that match
	viper.SetEnvPrefix("justexpenseit")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.BindPFlag("expensify-user-id", rootCmd.PersistentFlags().Lookup("expensify-user-id"))
	viper.BindPFlag("expensify-user-secret", rootCmd.PersistentFlags().Lookup("expensify-user-secret"))
	partnerUserID = viper.GetString("expensify-user-id")
	partnerUserSecret = viper.GetString("expensify-user-secret")
}
