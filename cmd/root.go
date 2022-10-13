// Copyright © 2022.  Douglas Chimento <dchimento@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var cfgFile string

type Config struct {
	domain   string
	record   string
	resolver string
	verbose  bool
	profile  string
}

var config Config = Config{
	resolver: "resolver1.opendns.com",
}

var logger, _ = zap.NewProduction()

var rootCmd = &cobra.Command{
	Use:   "aws-dyndns",
	Short: "aws-dyndns",
	Long:  "blah",
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aws-dyndns.yaml)")
	rootCmd.PersistentFlags().StringVarP(&config.domain, "domain", "d", "", "domain")
	rootCmd.PersistentFlags().StringVar(&config.profile, "profile", "aws-dyndns", "set aws profile or use AWS_PROFILE environment variable ")
	rootCmd.PersistentFlags().BoolVar(&config.verbose, "verbose", false, "verbose")
	os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
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

		// Search config in home directory with name ".aws-dyndns" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".aws-dyndns")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()
	if config.verbose {
		logger, _ = zap.NewDevelopment()
	}
	p := os.Getenv("AWS_PROFILE")
	if p == "" {
		os.Setenv("AWS_PROFILE", config.profile)
	}
	logger.Debug(fmt.Sprintf("Profile set to %s", os.Getenv("AWS_PROFILE")))
	if config.domain == "" {
		config.domain = viper.GetString("domain")
	}
	if config.profile == "" {
		config.profile = viper.GetString("profile")
	}
	err := setupAws()
	if err != nil {
		panic(err)
	}
}
