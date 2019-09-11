package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kubeconnect/k8s"
	"kubeconnect/lib"
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
		context, err := getContext()
		if err != nil {
			fmt.Print(err.Error())
			return
		}

		namespace, err := getNamespace(context)
		if err != nil {
			fmt.Print(err.Error())
			return
		}

		pod, err := getPod(context, namespace)
		if err != nil {
			fmt.Print(err.Error())
			return
		}

		// Connect
		// kubectl exec -it --context my-conext --namespace my-namespace my-pod /bin/sh

		// Get the current working directory.
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		// Transfer stdin, stdout, and stderr to the new process
		// and also set target directory for the shell to start in.
		pa := os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
			Dir:   cwd,
		}

		proc, err := os.StartProcess(
			"/usr/local/bin/kubectl",
			[]string{"kubectl", "exec", "-it", "--namespace", namespace.Name, "--context", context.Name, pod.Name, "/bin/sh"}, &pa)

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

func getContext() (k8s.Context, error) {
	contexts, err := k8s.GetContexts()
	if err != nil {
		return k8s.Context{}, err
	}

	index, err := lib.SelectFromList(
		"\033[38;5;3mWhat context do you want to connect to?\033[0m",
		"Context",
		k8s.ContextListItems(contexts))

	if err != nil {
		return k8s.Context{}, err
	}

	return contexts[index], nil
}

func getNamespace(context k8s.Context) (k8s.Namespace, error) {
	namespaces, err := k8s.GetNamespaces(context)
	if err != nil {
		return k8s.Namespace{}, err
	}

	index, err := lib.SelectFromList(
		"\033[38;5;3mWhat namespace do you want to connect to?\033[0m",
		"Namespace",
		k8s.NamespaceListItems(namespaces))
	if err != nil {
		return k8s.Namespace{}, err
	}

	return namespaces[index], nil
}

func getPod(context k8s.Context, namespace k8s.Namespace) (k8s.Pod, error) {
	pods, err := k8s.GetPods(context, namespace)
	if err != nil {
		return k8s.Pod{}, err
	}

	index, err := lib.SelectFromList(
		"\033[38;5;3mWhat pod do you want to connect to?\033[0m",
		"Pod",
		k8s.PodListItems(pods))
	if err != nil {
		return k8s.Pod{}, err
	}

	return pods[index], nil
}
