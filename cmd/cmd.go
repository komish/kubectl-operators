package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	// "k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/tools/clientcmd"
	"github.com/spf13/pflag"
)

func init() {
		pflag.StringP("kubeconfig", "k", "", "Path to kubeconfig file.")
}

var (
	kubeconfig string
	availableKubeconfigs []string
)

// Run starts the command line interface for kubectl-operators and returns
// an exit code which is the exit code returned by the cli.
func Run() int {
	pflag.Parse()
	populateAvailableKubeconfigs()

	if len(availableKubeconfigs) == 0 {
		printKubeConfigHelpOutput()
		return 2
	}

	// DEBUG
	fmt.Println(availableKubeconfigs)
	return 0
}

// homeDir gets the home directory of the user.
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// kubeconfigExistsAndIsFile does a quick check to make sure that the file path provided
// is a file type and can be opened.
func kubeconfigExistsAndIsFile(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	// if we couldn't open the file, error out
	if err != nil {
		return false, err
	}
	// if we have a directory instead of a file, error out
	if info.IsDir() {
		return false, errors.New("The provided path was a directory. Expected a file.")
	}
	return true, err
}

func populateAvailableKubeconfigs() {
	// Check for the kubeconfig flag
	if kubeconfig != "" {
		if res, _ := kubeconfigExistsAndIsFile(kubeconfig); res {
			availableKubeconfigs = append(availableKubeconfigs, kubeconfig)
		}
	}

	// Checking KUBECONFIG environment variable
	if kubeconfigEnvFilePath := os.Getenv("KUBECONFIG"); kubeconfigEnvFilePath != "" {
		if res, _ := kubeconfigExistsAndIsFile(kubeconfigEnvFilePath); res {
			availableKubeconfigs = append(availableKubeconfigs, kubeconfigEnvFilePath)
		} 
	}

	// Checking ~/.kube/config
	if home := homeDir(); home != "" {
		fullFilePath := filepath.Join(home, ".kube", "config")
		if res, _ := kubeconfigExistsAndIsFile(fullFilePath); res {
			availableKubeconfigs = append(availableKubeconfigs, fullFilePath)
		}	
	}
}

// printKubeConfigHelpOutput responds with information on KUBECONFIG precedence
// that this plugin uses.
func printKubeConfigHelpOutput() {
	fmt.Println("No authentication context was found. This plugin looks for a configuration in one of the following places:")
	fmt.Printf("\t(1) The --kubeconfig flag passed to this plugin.\n")
	fmt.Printf("\t(2) The KUBECONFIG environment variable.\n")
	fmt.Printf("\t(3) The .kube directory in your home directory.\n")
	fmt.Println("Configuration file precedence is as listed above.")
}