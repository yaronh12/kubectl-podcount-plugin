package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func main() {
	configFlags := genericclioptions.NewConfigFlags(true)

	binaryName := filepath.Base(os.Args[0])

	cmd := createCommand(configFlags, binaryName)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func createCommand(configFlags *genericclioptions.ConfigFlags, binaryName string) *cobra.Command {
	examples := `
	# Count pods in specific namespace
	%s podcount-yaron -n kube-system

	#Count pods in all namespaces
	%s podcount-yaron

	  # Multiple namespaces (comma-separated)
	%s podcount-yaron -n "default,kube-system,monitoring"
	`

	var cmdPrefix string
	if strings.HasPrefix(binaryName, "kubectl-") {
		cmdPrefix = "kubectl"
	} else if strings.HasPrefix(binaryName, "oc-") {
		cmdPrefix = "oc"
	} else {
		cmdPrefix = binaryName
	}

	formattedExamples := fmt.Sprintf(examples, cmdPrefix, cmdPrefix, cmdPrefix)

	cmd := &cobra.Command{
		Use:     "podcount-yaron",
		Short:   "Count running pods in a cluster",
		Long:    formattedExamples,
		Example: formattedExamples,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPodCount(configFlags)
		},
	}
	configFlags.AddFlags(cmd.Flags())
	return cmd
}

func runPodCount(configFlags *genericclioptions.ConfigFlags) error {
	namespaces := *configFlags.Namespace
	if namespaces == "" {
		fmt.Println("Running pod count in all namespaces")
	} else {
		fmt.Println("Running pod count in namespace:", namespaces)
	}
	factory := kcmdutil.NewFactory(configFlags)
	clientset, err := factory.KubernetesClientSet()
	if err != nil {
		return fmt.Errorf("error creating kubernetes client: %w", err)
	}
	runningCount := 0
	if namespaces == "" {
		pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("error listing pods: %w", err)
		}

		for _, pod := range pods.Items {
			if pod.Status.Phase == "Running" {
				runningCount++
			}
		}
	} else {
		namespacesList := strings.Split(namespaces, ",")
		for i, ns := range namespacesList {
			namespacesList[i] = strings.TrimSpace(ns)
		}
		for _, namespace := range namespacesList {
			pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
			if err != nil {
				return fmt.Errorf("error listing pods: %w", err)
			}
			for _, pod := range pods.Items {
				if pod.Status.Phase == "Running" {
					runningCount++
				}
			}
		}
	}

	fmt.Println("Running pods:", runningCount)
	return nil
}
