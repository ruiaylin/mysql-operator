package cluster

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/grtl/mysql-operator/cli/cmd/config"
)

var removePVC bool

// Command
var clusterDeleteCmd = &cobra.Command{
	Use:   "delete [cluster names]",
	Short: "Deletes MySQL clusters",
	Long: `Deletes MySQL clusters and resources associated with them.
Unless explicitly specified does not remove PersistentVolumeClaims.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace, err := cmd.InheritedFlags().GetString("namespace")
		if err != nil {
			panic(err)
		}

		for _, arg := range args {
			err := deleteCluster(arg, namespace)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to remove %s: %v", arg, err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	clusterDeleteCmd.PersistentFlags().BoolVarP(
		&removePVC,
		"remove-pvc",
		"r",
		false,
		"remove PersistentVolumeClaims along with the cluster",
	)

	Cmd.AddCommand(clusterDeleteCmd)
}

func deleteCluster(clusterName string, namespace string) error {
	clustersInterface := config.GetConfig().Clientset().CrV1().MySQLClusters(namespace)
	err := clustersInterface.Delete(clusterName, &v1.DeleteOptions{})
	if err != nil {
		return err
	}

	if removePVC {
		return deletePVC(clusterName, namespace)
	}

	return err
}

func deletePVC(clusterName string, namespace string) error {
	pvcInterface := config.GetConfig().KubeClientset().CoreV1().PersistentVolumeClaims(namespace)
	return pvcInterface.DeleteCollection(&v1.DeleteOptions{}, v1.ListOptions{
		LabelSelector: labels.Set{"app": clusterName}.AsSelector().String(),
	})
}
