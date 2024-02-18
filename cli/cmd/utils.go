package cmd

// authors Dipankar <dipankar@dipankar-das.com>
import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ksctl/ksctl/pkg/helpers/consts"
)

func createManaged(approval bool) {
	cli.Client.Metadata.ManagedNodeType = nodeSizeMP
	cli.Client.Metadata.NoMP = noMP

	cli.Client.Metadata.ClusterName = clusterName
	cli.Client.Metadata.K8sDistro = consts.KsctlKubernetes(distro)
	cli.Client.Metadata.K8sVersion = k8sVer
	cli.Client.Metadata.Region = region

	cli.Client.Metadata.CNIPlugin = cni
	cli.Client.Metadata.Applications = apps
	cli.Client.Metadata.LogWritter = os.Stdout
	if err := createApproval(approval); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	err := controller.CreateManagedCluster(&cli.Client)
	if err != nil {
		log.Error("Failed to create managed cluster", "Reason", err)
		os.Exit(1)
	}
	log.Success("Created the managed cluster successfully")
}

func createHA(approval bool) {
	cli.Client.Metadata.IsHA = true

	cli.Client.Metadata.ClusterName = clusterName
	cli.Client.Metadata.K8sDistro = consts.KsctlKubernetes(distro)
	cli.Client.Metadata.K8sVersion = k8sVer
	cli.Client.Metadata.Region = region

	cli.Client.Metadata.NoCP = noCP
	cli.Client.Metadata.NoWP = noWP
	cli.Client.Metadata.NoDS = noDS

	cli.Client.Metadata.LoadBalancerNodeType = nodeSizeLB
	cli.Client.Metadata.ControlPlaneNodeType = nodeSizeCP
	cli.Client.Metadata.WorkerPlaneNodeType = nodeSizeWP
	cli.Client.Metadata.DataStoreNodeType = nodeSizeDS

	cli.Client.Metadata.CNIPlugin = cni
	cli.Client.Metadata.Applications = apps

	cli.Client.Metadata.LogWritter = os.Stdout
	if err := createApproval(approval); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	err := controller.CreateHACluster(&cli.Client)
	if err != nil {
		log.Error("Failed to create self-managed HA cluster", "Reason", err)
		os.Exit(1)
	}
	log.Success("Created the self-managed HA cluster successfully")
}

func deleteManaged(approval bool) {

	cli.Client.Metadata.ClusterName = clusterName
	cli.Client.Metadata.K8sDistro = consts.KsctlKubernetes(distro)
	cli.Client.Metadata.Region = region

	cli.Client.Metadata.LogWritter = os.Stdout
	if err := deleteApproval(approval); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	err := controller.DeleteManagedCluster(&cli.Client)
	if err != nil {
		log.Error("Failed to delete managed cluster", "Reason", err)
		os.Exit(1)
	}
	log.Success("Deleted the managed cluster successfully")
}

func deleteHA(approval bool) {

	cli.Client.Metadata.IsHA = true

	cli.Client.Metadata.ClusterName = clusterName
	cli.Client.Metadata.K8sDistro = consts.KsctlKubernetes(distro)
	cli.Client.Metadata.Region = region

	cli.Client.Metadata.LogWritter = os.Stdout
	if err := deleteApproval(approval); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	err := controller.DeleteHACluster(&cli.Client)
	if err != nil {
		log.Error("Failed to delete self-managed HA cluster", "Reason", err)
		os.Exit(1)
	}
	log.Success("Deleted the self-managed HA cluster successfully")
}

func getRequestPayload() ([]byte, error) {
	a, err := json.MarshalIndent(cli.Client.Metadata, "", " ")
	if err != nil {
		return nil, err
	}
	return a, nil
}

func deleteApproval(showMsg bool) error {

	a, err := getRequestPayload()
	if err != nil {
		return err
	}
	log.Box("Input in Json", string(a))

	if !showMsg {
		log.Box("Warning 🚨", "THIS IS A DESTRUCTIVE STEP MAKE SURE IF YOU WANT TO DELETE THE CLUSTER")

		log.Print("Enter your choice to continue..[y/N]")
		choice := "n"
		unsafe := false
		fmt.Scanf("%s", &choice)
		if strings.Compare("y", choice) == 0 ||
			strings.Compare("yes", choice) == 0 ||
			strings.Compare("Y", choice) == 0 {
			unsafe = true
		}

		if !unsafe {
			return log.NewError("approval cancelled")
		}
	}
	return nil
}

func createApproval(showMsg bool) error {

	a, err := getRequestPayload()
	if err != nil {
		return err
	}
	log.Box("Input in Json", string(a))

	if !showMsg {
		log.Box("Warning 🚨", "THIS IS A CREATION STEP MAKE SURE IF YOU WANT TO CREATE THE CLUSTER")

		log.Print("Enter your choice to continue..[y/N]")
		choice := "n"
		unsafe := false
		fmt.Scanf("%s", &choice)
		if strings.Compare("y", choice) == 0 ||
			strings.Compare("yes", choice) == 0 ||
			strings.Compare("Y", choice) == 0 {
			unsafe = true
		}

		if !unsafe {
			return log.NewError("approval cancelled")
		}
	}
	return nil
}

func SetDefaults(provider consts.KsctlCloud, clusterType consts.KsctlClusterType) {
	switch string(provider) + string(clusterType) {
	case string(consts.CloudLocal) + string(consts.ClusterTypeMang):
		if noMP == -1 {
			noMP = 2
		}
		if len(k8sVer) == 0 {
			k8sVer = "1.27.1"
		}

	case string(consts.CloudAzure) + string(consts.ClusterTypeMang):
		if len(nodeSizeMP) == 0 {
			nodeSizeMP = "Standard_DS2_v2"
		}
		if noMP == -1 {
			noMP = 1
		}
		if len(region) == 0 {
			region = "eastus"
		}
		if len(k8sVer) == 0 {
			k8sVer = "1.27"
		}

	case string(consts.CloudCivo) + string(consts.ClusterTypeMang):
		if len(nodeSizeMP) == 0 {
			nodeSizeMP = "g4s.kube.small"
		}
		if noMP == -1 {
			noMP = 1
		}
		if len(region) == 0 {
			region = "LON1"
		}
		if len(k8sVer) == 0 {
			k8sVer = "1.27.1"
		}

	case string(consts.CloudAzure) + string(consts.ClusterTypeHa):
		if len(nodeSizeCP) == 0 {
			nodeSizeCP = "Standard_F2s"
		}
		if len(nodeSizeWP) == 0 {
			nodeSizeWP = "Standard_F2s"
		}
		if len(nodeSizeDS) == 0 {
			nodeSizeDS = "Standard_F2s"
		}
		if len(nodeSizeLB) == 0 {
			nodeSizeLB = "Standard_F2s"
		}
		if len(region) == 0 {
			region = "eastus"
		}
		if noWP == -1 {
			noWP = 1
		}
		if noCP == -1 {
			noCP = 3
		}
		if noDS == -1 {
			noDS = 1
		}
		if len(k8sVer) == 0 {
			k8sVer = "1.27.1"
		}
		if len(distro) == 0 {
			distro = string(consts.K8sK3s)
		}

	case string(consts.CloudCivo) + string(consts.ClusterTypeHa):
		if len(nodeSizeCP) == 0 {
			nodeSizeCP = "g3.small"
		}
		if len(nodeSizeWP) == 0 {
			nodeSizeWP = "g3.large"
		}
		if len(nodeSizeDS) == 0 {
			nodeSizeDS = "g3.small"
		}
		if len(nodeSizeLB) == 0 {
			nodeSizeLB = "g3.small"
		}
		if len(region) == 0 {
			region = "LON1s"
		}
		if noWP == -1 {
			noWP = 1
		}
		if noCP == -1 {
			noCP = 3
		}
		if noDS == -1 {
			noDS = 1
		}
		if len(k8sVer) == 0 {
			k8sVer = "1.27.1"
		}
		if len(distro) == 0 {
			distro = string(consts.K8sK3s)
		}
	}
}

func argsFlags() {
	// Managed Azure
	clusterNameFlag(createClusterAzure)
	nodeSizeManagedFlag(createClusterAzure)
	regionFlag(createClusterAzure)
	noOfMPFlag(createClusterAzure)
	k8sVerFlag(createClusterAzure)
	distroFlag(createClusterAzure)
	appsFlag(createClusterAzure)
	cniFlag(createClusterAzure)

	// Managed Civo
	clusterNameFlag(createClusterCivo)
	nodeSizeManagedFlag(createClusterCivo)
	regionFlag(createClusterCivo)
	appsFlag(createClusterCivo)
	cniFlag(createClusterCivo)
	noOfMPFlag(createClusterCivo)
	distroFlag(createClusterCivo)
	k8sVerFlag(createClusterCivo)

	// Managed Local
	clusterNameFlag(createClusterLocal)
	appsFlag(createClusterLocal)
	cniFlag(createClusterLocal)
	noOfMPFlag(createClusterLocal)
	distroFlag(createClusterLocal)
	k8sVerFlag(createClusterLocal)

	// HA Civo
	clusterNameFlag(createClusterHACivo)
	nodeSizeCPFlag(createClusterHACivo)
	nodeSizeDSFlag(createClusterHACivo)
	nodeSizeWPFlag(createClusterHACivo)
	nodeSizeLBFlag(createClusterHACivo)
	regionFlag(createClusterHACivo)
	appsFlag(createClusterHACivo)
	cniFlag(createClusterHACivo)
	noOfWPFlag(createClusterHACivo)
	noOfCPFlag(createClusterHACivo)
	noOfDSFlag(createClusterHACivo)
	distroFlag(createClusterHACivo)
	k8sVerFlag(createClusterHACivo)

	// HA Azure
	clusterNameFlag(createClusterHAAzure)
	nodeSizeCPFlag(createClusterHAAzure)
	nodeSizeDSFlag(createClusterHAAzure)
	nodeSizeWPFlag(createClusterHAAzure)
	nodeSizeLBFlag(createClusterHAAzure)
	regionFlag(createClusterHAAzure)
	appsFlag(createClusterHAAzure)
	cniFlag(createClusterHAAzure)
	noOfWPFlag(createClusterHAAzure)
	noOfCPFlag(createClusterHAAzure)
	noOfDSFlag(createClusterHAAzure)
	distroFlag(createClusterHAAzure)
	k8sVerFlag(createClusterHAAzure)

	// Delete commands
	// Managed Local
	clusterNameFlag(deleteClusterLocal)

	// managed Azure
	clusterNameFlag(deleteClusterAzure)
	regionFlag(deleteClusterAzure)

	// Managed Civo
	clusterNameFlag(deleteClusterCivo)
	regionFlag(deleteClusterCivo)

	// HA Civo
	clusterNameFlag(deleteClusterHACivo)
	regionFlag(deleteClusterHACivo)

	// HA Azure
	clusterNameFlag(deleteClusterHAAzure)
	regionFlag(deleteClusterHAAzure)

	AllFeatures()
}

func AllFeatures() {

	featureFlag(createClusterAzure)
	featureFlag(createClusterHAAzure)
	featureFlag(createClusterCivo)
	featureFlag(createClusterHACivo)
	featureFlag(createClusterLocal)

	featureFlag(deleteClusterAzure)
	featureFlag(deleteClusterHAAzure)
	featureFlag(deleteClusterCivo)
	featureFlag(deleteClusterHACivo)
	featureFlag(deleteClusterLocal)

	featureFlag(addMoreWorkerNodesHACivo)
	featureFlag(addMoreWorkerNodesHAAzure)

	featureFlag(deleteNodesHAAzure)
	featureFlag(deleteNodesHACivo)

	featureFlag(getClusterCmd)
	featureFlag(switchCluster)
}
