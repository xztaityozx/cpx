// Copyright © 2019 xztaityozx
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vbauerster/mpb"
	"github.com/xztaityozx/cpx/config"
	"github.com/xztaityozx/cpx/factory"
)

var cfgFile string
var cfg config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cpx",
	Short:   "cpx: wrapper for cp command",
	Long:    `cpx is wrapper for cp command with Fuzzy Finder`,
	Version: "0.1.0",
	Args:    cobra.MinimumNArgs(2),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		src, dst := args[0], args[1]

		force, _ := cmd.Flags().GetBool("force")
		srcFf, _ := cmd.Flags().GetBool("src-finder")
		dstFf, _ := cmd.Flags().GetBool("dst-finder")
		recursive, _ := cmd.Flags().GetBool("recursive")
		parallel, _ := cmd.Flags().GetUint("parallel")
		progress, _ := cmd.Flags().GetBool("progress")

		buildPath := func(p string, b bool) []string {
			if b {
				a, e := cfg.FuzzyFinder.GetPathes(p)
				if e != nil {
					logrus.Fatal(e)
				}
				return a
			}
			return []string{p}
		}

		res, err := factory.GenerateLocalCopyers(buildPath(src, srcFf), buildPath(dst, dstFf), recursive)
		if err != nil {
			logrus.Fatal(err)
		}

		if len(res) == 0 {
			logrus.Warn("no copy task")
		}

		var wg sync.WaitGroup
		wg.Add(len(res))
		ch := make(chan struct{}, parallel)
		defer close(ch)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		act := func() func() {
			if progress {
				pb := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithContext(ctx))
				return func() {

				}
			}
		}

		for _, v := range res {
			go func() {
				ch <- struct{}{}
				act()
				<-ch
			}()
		}

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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/cpx/config.json)")

	rootCmd.Flags().Bool("force", false, "if an existing destination file cannot be opened, remove it and try again")
	rootCmd.Flags().BoolP("src-finder", "f", false, "use fuzzy-finder to select SOURCE files or directories")
	rootCmd.Flags().BoolP("dst-finder", "F", false, "use fuzzy-finder to select DESTINATION files or directories")
	rootCmd.Flags().BoolP("recursive", "r", false, "copy directories recursively")
	rootCmd.Flags().UintP("parallel", "P", 1, "number of parallel size")
	rootCmd.Flags().BoolP("progress", "p", false, "show progress bar")
	rootCmd.Flags().SortFlags = false
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

		viper.AddConfigPath(filepath.Join(home, "cpx"))
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		logrus.Fatal(err)
	}
}
