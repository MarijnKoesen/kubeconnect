/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"kubeconnect/k8s"
	"kubeconnect/lib"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "kubeconnect",
	Short: "Connect to any running pod in k8s with ease",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		// Get context
		contexts, err := k8s.GetContexts()
		var listItems []lib.ListItem
		for index, item := range contexts {
			listItems = append(listItems, lib.ListItem{Number: index + 1, Label: item.Name})
		}
		index, err := lib.SelectFromList("\033[38;5;3mWhat context do you want to connect to?\033[0m", listItems)

		if err != nil {
			fmt.Print(err.Error())
			return
		}
		context := contexts[index]


		// Get namespace
		namespaces, err := k8s.GetNamespaces(context)
		if err != nil {
			fmt.Print(err.Error())
			return
		}

		var listItems2 []lib.ListItem
		for index, item := range namespaces {
			listItems2 = append(listItems2, lib.ListItem{Number: index + 1, Label: item.Name})
		}
		index2, err2 := lib.SelectFromList("\033[38;5;3mWhat namespace do you want to connect to?\033[0m", listItems2)

		if err2 != nil {
			fmt.Print(err2.Error())
			return
		}
		namespace := namespaces[index2]


		// Get pods
		pods, err := k8s.GetPods(context, namespace)
		if err != nil {
			fmt.Print(err.Error())
			return
		}

		var listItems3 []lib.ListItem
		for index, item := range pods {
			listItems3 = append(listItems3, lib.ListItem{Number: index + 1, Label: item.Name})
		}
		index3, err3 := lib.SelectFromList("\033[38;5;3mWhat pod do you want to connect to?\033[0m", listItems3)

		if err3 != nil {
			fmt.Print(err3.Error())
			return
		}
		pod := pods[index3]


		// Connect
		// kubectl exec -it -n instapro-master rabbitmq-556d8c78df-vwppc /bin/bash

		// Get the current working directory.
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		// Transfer stdin, stdout, and stderr to the new process
		// and also set target directory for the shell to start in.
		pa := os.ProcAttr {
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
			Dir: cwd,
		}

		proc, err := os.StartProcess(
			"/usr/local/bin/kubectl",
			[]string{"kubectl", "exec", "-it", "--namespace", namespace.Name, "--context", context.Name, pod.Name, "/bin/bash"}, &pa)

		if err != nil {
			panic(err)
		}

		// Wait until user exits the shell
		_, err = proc.Wait()
		if err != nil {
			panic(err)
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubeconnect.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".kubeconnect" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".kubeconnect")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
