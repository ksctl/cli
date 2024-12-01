package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/ksctl/ksctl/pkg/controllers"
	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/ksctl/ksctl/pkg/types"
)

func newLogo() string {
	x := strings.Split(logoKsctl, "\n")

	y := []string{}

	colorCode := map[int]func(str string) string{
		0: func(str string) string { return color.New(color.BgHiMagenta).Add(color.FgHiBlack).SprintFunc()(str) },
		1: func(str string) string { return color.New(color.BgHiBlue).Add(color.FgHiBlack).SprintFunc()(str) },
		2: func(str string) string { return color.New(color.BgHiCyan).Add(color.FgHiBlack).SprintFunc()(str) },
		3: func(str string) string { return color.New(color.BgHiGreen).Add(color.FgHiBlack).SprintFunc()(str) },
		4: func(str string) string { return color.New(color.BgHiYellow).Add(color.FgHiBlack).SprintFunc()(str) },
		5: func(str string) string { return color.New(color.BgHiRed).Add(color.FgHiBlack).SprintFunc()(str) },
	}

	for i, _x := range x {
		if _y, ok := colorCode[i]; ok {
			y = append(y, _y(_x))
		}
	}
	return strings.Join(y, "\n")
}

func LongMessage(message string) string {
	if !DisableHeaderBanner {
		return "Ksctl ascii [logo]"
	}
	return fmt.Sprintf("%s\n\n%s", newLogo(), color.New(color.BgHiYellow).Add(color.FgBlack).SprintFunc()(message))
}

func createManaged(ctx context.Context, log types.LoggerFactory, approval bool) {
	cli.Client.Metadata.ManagedNodeType = nodeSizeMP
	cli.Client.Metadata.NoMP = noMP

	cli.Client.Metadata.ClusterName = clusterName
	cli.Client.Metadata.K8sDistro = consts.KsctlKubernetes(distro)
	cli.Client.Metadata.K8sVersion = k8sVer
	cli.Client.Metadata.Region = region
	cli.Client.Metadata.StateLocation = consts.KsctlStore(storage)
	if len(cni) > 0 {
		_cni := strings.Split(cni, "@")
		cli.Client.Metadata.CNIPlugin = types.KsctlApp{
			StackName: _cni[0],
			Overrides: map[string]map[string]any{
				_cni[0]: {
					"version": func() string {
						if len(_cni) != 2 {
							return "latest"
						}
						return _cni[1]
					}(),
				},
			},
		}
	}
	if err := createApproval(ctx, log, approval); err != nil {
		log.Error("createApproval", "Reason", err)
		os.Exit(1)
	}
	m, err := controllers.NewManagerClusterManaged(
		ctx,
		log,
		&cli.Client,
	)
	if err != nil {
		log.Error("Failed to create self-managed HA cluster", "Reason", err)
		os.Exit(1)
	}

	err = m.CreateCluster()
	if err != nil {
		log.Error("Failed to create managed cluster", "Reason", err)
		os.Exit(1)
	}
	log.Success(ctx, "Created the managed cluster successfully")
}

func createHA(ctx context.Context, log types.LoggerFactory, approval bool) {
	cli.Client.Metadata.IsHA = true

	cli.Client.Metadata.ClusterName = clusterName
	cli.Client.Metadata.K8sDistro = consts.KsctlKubernetes(distro)
	cli.Client.Metadata.K8sVersion = k8sVer
	cli.Client.Metadata.Region = region

	cli.Client.Metadata.NoCP = noCP
	cli.Client.Metadata.NoWP = noWP
	cli.Client.Metadata.NoDS = noDS
	cli.Client.Metadata.StateLocation = consts.KsctlStore(storage)

	cli.Client.Metadata.LoadBalancerNodeType = nodeSizeLB
	cli.Client.Metadata.ControlPlaneNodeType = nodeSizeCP
	cli.Client.Metadata.WorkerPlaneNodeType = nodeSizeWP
	cli.Client.Metadata.DataStoreNodeType = nodeSizeDS

	if len(cni) > 0 {
		_cni := strings.Split(cni, "@")
		cli.Client.Metadata.CNIPlugin = types.KsctlApp{
			StackName: _cni[0],
			Overrides: map[string]map[string]any{
				_cni[0]: {
					"version": func() string {
						if len(_cni) != 2 {
							return "latest"
						}
						return _cni[1]
					}(),
				},
			},
		}
	}

	if err := createApproval(ctx, log, approval); err != nil {
		log.Error("createApproval", "Reason", err)
		os.Exit(1)
	}
	m, err := controllers.NewManagerClusterSelfManaged(
		ctx,
		log,
		&cli.Client,
	)
	if err != nil {
		log.Error("Failed to create self-managed HA cluster", "Reason", err)
		os.Exit(1)
	}

	err = m.CreateCluster()
	if err != nil {
		log.Error("Failed to create self-managed HA cluster", "Reason", err)
		os.Exit(1)
	}
	log.Success(ctx, "Created the self-managed HA cluster successfully")
}

func deleteManaged(ctx context.Context, log types.LoggerFactory, approval bool) {

	cli.Client.Metadata.ClusterName = clusterName
	cli.Client.Metadata.K8sDistro = consts.KsctlKubernetes(distro)
	cli.Client.Metadata.Region = region
	cli.Client.Metadata.StateLocation = consts.KsctlStore(storage)

	if err := deleteApproval(ctx, log, approval); err != nil {
		log.Error("deleteApproval", "Reason", err)
		os.Exit(1)
	}

	m, err := controllers.NewManagerClusterManaged(
		ctx,
		log,
		&cli.Client,
	)
	if err != nil {
		log.Error("Failed to create self-managed HA cluster", "Reason", err)
		os.Exit(1)
	}
	err = m.DeleteCluster()
	if err != nil {
		log.Error("Failed to delete managed cluster", "Reason", err)
		os.Exit(1)
	}
	log.Success(ctx, "Deleted the managed cluster successfully")
}

func deleteHA(ctx context.Context, log types.LoggerFactory, approval bool) {

	cli.Client.Metadata.IsHA = true

	cli.Client.Metadata.ClusterName = clusterName
	cli.Client.Metadata.K8sDistro = consts.KsctlKubernetes(distro)
	cli.Client.Metadata.Region = region
	cli.Client.Metadata.StateLocation = consts.KsctlStore(storage)

	if err := deleteApproval(ctx, log, approval); err != nil {
		log.Error("deleteApproval", "Reason", err)
		os.Exit(1)
	}
	m, err := controllers.NewManagerClusterSelfManaged(
		ctx,
		log,
		&cli.Client,
	)
	if err != nil {
		log.Error("Failed to create self-managed HA cluster", "Reason", err)
		os.Exit(1)
	}

	err = m.DeleteCluster()
	if err != nil {
		log.Error("Failed to delete self-managed HA cluster", "Reason", err)
		os.Exit(1)
	}
	log.Success(ctx, "Deleted the self-managed HA cluster successfully")
}

func getRequestPayload() ([]byte, error) {
	a, err := json.MarshalIndent(cli.Client.Metadata, "", " ")
	if err != nil {
		return nil, err
	}
	return a, nil
}

func deleteApproval(ctx context.Context, log types.LoggerFactory, showMsg bool) error {

	a, err := getRequestPayload()
	if err != nil {
		return err
	}
	log.Box(ctx, "Input in Json", string(a))

	if !showMsg {
		log.Box(ctx, "Warning 🚨", "THIS IS A DESTRUCTIVE STEP MAKE SURE IF YOU WANT TO DELETE THE CLUSTER")

		log.Print(ctx, "Enter your choice to continue..[y/N]")
		choice := "n"
		unsafe := false
		fmt.Scanf("%s", &choice)
		if strings.Compare("y", choice) == 0 ||
			strings.Compare("yes", choice) == 0 ||
			strings.Compare("Y", choice) == 0 {
			unsafe = true
		}

		if !unsafe {
			return log.NewError(ctx, "approval cancelled by user")
		}
	}
	return nil
}

func createApproval(ctx context.Context, log types.LoggerFactory, showMsg bool) error {

	a, err := getRequestPayload()
	if err != nil {
		return err
	}
	log.Box(ctx, "Input in Json", string(a))

	if !showMsg {
		log.Box(ctx, "Warning 🚨", "THIS IS A CREATION STEP MAKE SURE IF YOU WANT TO CREATE THE CLUSTER")

		log.Print(ctx, "Enter your choice to continue..[y/N]")
		choice := "n"
		unsafe := false
		fmt.Scanf("%s", &choice)
		if strings.Compare("y", choice) == 0 ||
			strings.Compare("yes", choice) == 0 ||
			strings.Compare("Y", choice) == 0 {
			unsafe = true
		}

		if !unsafe {
			return log.NewError(ctx, "approval cancelled by user")
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
			k8sVer = "1.30.0"
		}
		if len(storage) == 0 {
			storage = string(consts.StoreLocal)
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
		if len(storage) == 0 {
			storage = string(consts.StoreLocal)
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
			k8sVer = "1.28.7"
		}
		if len(storage) == 0 {
			storage = string(consts.StoreLocal)
		}

	case string(consts.CloudAzure) + string(consts.ClusterTypeHa):
		if len(nodeSizeCP) == 0 {
			if distro == string(consts.K8sKubeadm) {
				nodeSizeCP = "Standard_F2s"
			} else {
				nodeSizeCP = "Standard_F2s"
			}
		}
		if len(nodeSizeWP) == 0 {
			if distro == string(consts.K8sKubeadm) {
				nodeSizeWP = "Standard_F4s"
			} else {
				nodeSizeWP = "Standard_F2s"
			}
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
			if distro == string(consts.K8sKubeadm) {
				noWP = 2
			} else {
				noWP = 1
			}
		}
		if noCP == -1 {
			noCP = 3
		}
		if noDS == -1 {
			noDS = 3
		}
		if len(distro) == 0 {
			distro = string(consts.K8sK3s)
		}
		if len(storage) == 0 {
			storage = string(consts.StoreLocal)
		}

	case string(consts.CloudAws) + string(consts.ClusterTypeMang):
		if len(nodeSizeMP) == 0 {
			nodeSizeMP = "t2.micro"
		}
		if noMP == -1 {
			noMP = 2
		}
		if len(region) == 0 {
			region = "ap-south-1"
		}
		if len(k8sVer) == 0 {
			k8sVer = "1.30"
		}
		if len(storage) == 0 {
			storage = string(consts.StoreLocal)
		}

	case string(consts.CloudAws) + string(consts.ClusterTypeHa):
		if len(nodeSizeCP) == 0 {
			if distro == string(consts.K8sKubeadm) {
				nodeSizeCP = "t2.medium"
			} else {
				nodeSizeCP = "t2.micro"
			}
		}
		if len(nodeSizeWP) == 0 {
			if distro == string(consts.K8sKubeadm) {
				nodeSizeWP = "t2.medium"
			} else {
				nodeSizeWP = "t2.micro"
			}
		}
		if len(nodeSizeDS) == 0 {
			nodeSizeDS = "t2.micro"
		}
		if len(nodeSizeLB) == 0 {
			nodeSizeLB = "t2.micro"
		}
		if len(region) == 0 {
			region = "us-east-1"
		}
		if noWP == -1 {
			if distro == string(consts.K8sKubeadm) {
				noWP = 2
			} else {
				noWP = 1
			}
		}
		if noCP == -1 {
			noCP = 3
		}
		if noDS == -1 {
			noDS = 3
		}
		if len(distro) == 0 {
			distro = string(consts.K8sK3s)
		}
		if len(storage) == 0 {
			storage = string(consts.StoreLocal)
		}

	case string(consts.CloudCivo) + string(consts.ClusterTypeHa):
		if len(nodeSizeCP) == 0 {
			if distro == string(consts.K8sKubeadm) {
				nodeSizeCP = "g3.large"
			} else {
				nodeSizeCP = "g3.medium"
			}
		}
		if len(nodeSizeWP) == 0 {
			if distro == string(consts.K8sKubeadm) {
				nodeSizeWP = "g3.large"
			} else {
				nodeSizeWP = "g3.medium"
			}
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
			if distro == string(consts.K8sKubeadm) {
				noWP = 2
			} else {
				noWP = 1
			}
		}
		if noCP == -1 {
			noCP = 3
		}
		if noDS == -1 {
			noDS = 3
		}
		if len(storage) == 0 {
			storage = string(consts.StoreLocal)
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
	cniFlag(createClusterAzure)
	storageFlag(createClusterAzure)

	// Managed Aws
	clusterNameFlag(createClusterAws)
	nodeSizeManagedFlag(createClusterAws)
	regionFlag(createClusterAws)
	noOfMPFlag(createClusterAws)
	k8sVerFlag(createClusterAws)
	distroFlag(createClusterAws)
	cniFlag(createClusterAws)
	storageFlag(createClusterAws)

	// Managed Civo
	clusterNameFlag(createClusterCivo)
	nodeSizeManagedFlag(createClusterCivo)
	regionFlag(createClusterCivo)
	cniFlag(createClusterCivo)
	noOfMPFlag(createClusterCivo)
	distroFlag(createClusterCivo)
	k8sVerFlag(createClusterCivo)
	storageFlag(createClusterCivo)

	// Managed Local
	clusterNameFlag(createClusterLocal)
	cniFlag(createClusterLocal)
	noOfMPFlag(createClusterLocal)
	distroFlag(createClusterLocal)
	k8sVerFlag(createClusterLocal)
	storageFlag(createClusterLocal)

	// HA Civo
	clusterNameFlag(createClusterHACivo)
	nodeSizeCPFlag(createClusterHACivo)
	nodeSizeDSFlag(createClusterHACivo)
	nodeSizeWPFlag(createClusterHACivo)
	nodeSizeLBFlag(createClusterHACivo)
	regionFlag(createClusterHACivo)
	cniFlag(createClusterHACivo)
	noOfWPFlag(createClusterHACivo)
	noOfCPFlag(createClusterHACivo)
	noOfDSFlag(createClusterHACivo)
	distroFlag(createClusterHACivo)
	k8sVerFlag(createClusterHACivo)
	storageFlag(createClusterHACivo)

	// HA Aws
	clusterNameFlag(createClusterHAAws)
	nodeSizeCPFlag(createClusterHAAws)
	nodeSizeDSFlag(createClusterHAAws)
	nodeSizeWPFlag(createClusterHAAws)
	nodeSizeLBFlag(createClusterHAAws)
	regionFlag(createClusterHAAws)
	cniFlag(createClusterHAAws)
	noOfWPFlag(createClusterHAAws)
	noOfCPFlag(createClusterHAAws)
	noOfDSFlag(createClusterHAAws)
	distroFlag(createClusterHAAws)
	k8sVerFlag(createClusterHAAws)
	storageFlag(createClusterHAAws)

	// HA Azure
	clusterNameFlag(createClusterHAAzure)
	nodeSizeCPFlag(createClusterHAAzure)
	nodeSizeDSFlag(createClusterHAAzure)
	nodeSizeWPFlag(createClusterHAAzure)
	nodeSizeLBFlag(createClusterHAAzure)
	regionFlag(createClusterHAAzure)
	cniFlag(createClusterHAAzure)
	noOfWPFlag(createClusterHAAzure)
	noOfCPFlag(createClusterHAAzure)
	noOfDSFlag(createClusterHAAzure)
	distroFlag(createClusterHAAzure)
	k8sVerFlag(createClusterHAAzure)
	storageFlag(createClusterHAAzure)

	// Delete commands
	// Managed Local
	clusterNameFlag(deleteClusterLocal)
	storageFlag(deleteClusterLocal)

	// managed Azure
	clusterNameFlag(deleteClusterAzure)
	regionFlag(deleteClusterAzure)
	storageFlag(deleteClusterAzure)

	// managed Aws
	clusterNameFlag(deleteClusterAws)
	regionFlag(deleteClusterAws)
	storageFlag(deleteClusterAws)

	// Managed Civo
	clusterNameFlag(deleteClusterCivo)
	regionFlag(deleteClusterCivo)
	storageFlag(deleteClusterCivo)

	// HA Civo
	clusterNameFlag(deleteClusterHAAws)
	regionFlag(deleteClusterHAAws)
	storageFlag(deleteClusterHAAws)

	// HA Aws
	clusterNameFlag(deleteClusterHACivo)
	regionFlag(deleteClusterHACivo)
	storageFlag(deleteClusterHACivo)

	// HA Azure
	clusterNameFlag(deleteClusterHAAzure)
	regionFlag(deleteClusterHAAzure)
	storageFlag(deleteClusterHAAzure)

	AllFeatures()
}

func AllFeatures() {

	featureFlag(createClusterAzure)
	featureFlag(createClusterHAAzure)
	featureFlag(createClusterCivo)
	featureFlag(createClusterHACivo)
	featureFlag(createClusterLocal)
	featureFlag(createClusterHAAws)
	featureFlag(createClusterAws)

	featureFlag(deleteClusterAzure)
	featureFlag(deleteClusterHAAzure)
	featureFlag(deleteClusterCivo)
	featureFlag(deleteClusterHACivo)
	featureFlag(deleteClusterLocal)
	featureFlag(deleteClusterHAAws)

	featureFlag(addMoreWorkerNodesHACivo)
	featureFlag(addMoreWorkerNodesHAAzure)

	featureFlag(deleteNodesHAAzure)
	featureFlag(deleteNodesHACivo)

	featureFlag(getClusterCmd)
	featureFlag(switchCluster)
	featureFlag(infoClusterCmd)
}
