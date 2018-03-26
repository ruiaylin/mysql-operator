package backup

import (
	crv1 "github.com/grtl/mysql-operator/pkg/apis/cr/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/grtl/mysql-operator/cli/cmd/config"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/resource"

	backupop "github.com/grtl/mysql-operator/operator/backup"
	"github.com/grtl/mysql-operator/util"
)

const backupTemplate = "artifacts/backup-cr.yaml"

var clusterName string

type backupSchedule struct {
	Name        string
	ClusterName string
	Time        string
	Storage     resource.Quantity
}

var backupCreateCmd = &cobra.Command{
	Use:   "create [backup name]",
	Short: "Schedules MySQL backups",
	Long: `Creates a recurring job for MySQL backup creation:
msp backup create "my-backup" --cluster "my-cluster --storage 1Gi"`,
	Run: func(cmd *cobra.Command, args []string) {
		namespace, err := cmd.InheritedFlags().GetString("namespace")
		if err != nil {
			panic(err)
		}

		backupName := clusterName + "-backup"
		if len(args) >= 1 {
			backupName = args[0]
		}

		time := args[1]

		cluster, err := config.GetConfig().Clientset().CrV1().MySQLClusters(namespace).Get(clusterName, metav1.GetOptions{})
		if err != nil {
			panic(err)
		}

		storage := cluster.Spec.Storage
		if len(args) >= 3 {
			storage, err = resource.ParseQuantity(args[2])
			if err != nil {
				panic(err)
			}
		}

		backupData := &backupSchedule{
			Name:        backupName,
			ClusterName: clusterName,
			Time:        time,
			Storage:     storage,
		}

		createBackup(backupData, namespace)
	},
}

func init() {
	Cmd.AddCommand(backupCreateCmd)

	backupCreateCmd.Flags().StringVar(&clusterName, "cluster", "",
		"name of cluster for which the backup is made")
	backupCreateCmd.MarkFlagRequired("cluster")
}

func createBackup(backupData interface{}, namespace string) {
	backup := new(crv1.MySQLBackup)
	err := util.ObjectFromTemplate(backupData, backup, backupTemplate, backupop.FuncMap)

	if err != nil {
		panic(err)
	}

	backupInterface := config.GetConfig().Clientset().CrV1().MySQLBackups(namespace)
	backup, err = backupInterface.Create(backup)

}
