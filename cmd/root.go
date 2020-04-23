package cmd

import (
	"fmt"
	"kubeconnect/k8s"
	"kubeconnect/lib"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, shell string

// Version compiled
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "kubeconnect",
	Version: Version,
	Short:   "Connect to any running pod in k8s with ease",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		context, err := getContext()
		if err != nil {
			return
		}

		namespace, err := getNamespace(context)
		if err != nil {
			return
		}

		pod, err := getPod(namespace)
		if err != nil {
			return
		}

		container, err := getContainer(pod)
		if err != nil {
			return
		}

		err = connect(pod, container)
		return
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubeconnect.yaml)")
	rootCmd.PersistentFlags().StringVar(&shell, "shell", "zsh ksh bash sh", "Shell to be used")
	if err := viper.BindPFlag("shell", rootCmd.PersistentFlags().Lookup("shell")); err != nil {
		panic(err)
	}
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

	viper.SetEnvPrefix("KUBECONNECT")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getContext() (context k8s.Context, err error) {
	contexts, err := k8s.GetContexts()
	if err != nil {
		return
	}

	index, err := lib.SelectFromList(
		"\033[38;5;3mWhat context do you want to connect to?\033[0m",
		"Context",
		k8s.ContextListItems(contexts))

	if err != nil {
		return
	}

	return contexts[index], nil
}

func getNamespace(context k8s.Context) (namespace k8s.Namespace, err error) {
	namespaces, err := context.GetNamespaces()
	if err != nil {
		return
	}

	index, err := lib.SelectFromList(
		"\033[38;5;3mWhat namespace do you want to connect to?\033[0m",
		"Namespace",
		k8s.NamespaceListItems(namespaces))
	if err != nil {
		return
	}

	return namespaces[index], nil
}

func getPod(namespace k8s.Namespace) (pod k8s.Pod, err error) {
	pods, err := namespace.GetPods()
	if err != nil {
		return
	}

	index, err := lib.SelectFromList(
		"\033[38;5;3mWhat pod do you want to connect to?\033[0m",
		"Pod",
		k8s.PodListItems(pods))
	if err != nil {
		return
	}

	return pods[index], nil
}

func getContainer(pod k8s.Pod) (container string, err error) {
	if len(pod.Containers) == 1 {
		return pod.Containers[0], nil
	}

	index, err := lib.SelectFromList(
		"\033[38;5;3mWhat container do you want to connect to?\033[0m",
		"Container",
		k8s.PodContainerListItems(pod))
	if err != nil {
		return
	}

	return pod.Containers[index], nil
}

func connect(pod k8s.Pod, container string) (err error) {
	// Get the current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	// Transfer stdin, stdout, and stderr to the new process
	// and also set target directory for the shell to start in.
	pa := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir:   cwd,
	}

	kubectlPath, err := exec.LookPath("kubectl")
	if err != nil {
		return
	}

	shell, err := getShell(pod, container)
	if err != nil {
		return
	}

	cmd := []string{
		"kubectl",
		"exec",
		"-it",
		"--namespace", pod.Namespace,
		"--context", pod.Context,
		"--container", container,
		pod.Name,
		shell}

	fmt.Println("Connecting using: " + strings.Join(cmd, " "))
	proc, err := os.StartProcess(kubectlPath, cmd, &pa)

	if err != nil {
		return
	}

	// Wait until user exits the shell
	_, err = proc.Wait()

	return err
}

func getShell(pod k8s.Pod, container string) (shell string, err error) {
	shellList := viper.GetString("shell")

	if regexp.MustCompile(`^/[^ ]+[a-z]+`).MatchString(shellList) {
		shell = shellList
		return
	}

	// Make sure we always have a valid fallback to sh
	fmt.Println("Looking up the shell to use from: " + shellList)
	shellList += " sh"
	cmd := "command -v " + strings.Join(strings.Fields(shellList), " || command -v ")

	shell, err = k8s.RunCmd(""+
		"exec",
		"-it",
		"--context", pod.Context,
		"--namespace", pod.Namespace,
		"--container", container,
		pod.Name,
		"--",
		"sh",
		"-c",
		cmd)

	return
}
