package main

import (
	"context"
	"fmt"
	"github.com/containernetworking/cni/pkg/skel"
	cnitypes "github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/current"
	cniversion "github.com/containernetworking/cni/pkg/version"
	"github.com/k8snetworkplumbingwg/whereabouts/pkg/allocate"
	"github.com/k8snetworkplumbingwg/whereabouts/pkg/config"
	"github.com/k8snetworkplumbingwg/whereabouts/pkg/logging"
	"github.com/k8snetworkplumbingwg/whereabouts/pkg/storage/kubernetes"
	"github.com/k8snetworkplumbingwg/whereabouts/pkg/types"
	"github.com/k8snetworkplumbingwg/whereabouts/pkg/version"
	"net"
)

func main() {
	skel.PluginMain(func(args *skel.CmdArgs) error {
		ipamConf, confVersion, err := config.LoadIPAMConfig(args.StdinData, args.Args)
		if err != nil {
			logging.Errorf("IPAM configuration load failed: %s", err)
			return err
		}
		logging.Debugf("ADD - IPAM configuration successfully read: %+v", *ipamConf)
		ipam, err := kubernetes.NewKubernetesIPAM(args.ContainerID, args.IfName, *ipamConf)
		if err != nil {
			return logging.Errorf("failed to create Kubernetes IPAM manager: %v", err)
		}
		defer func() { safeCloseKubernetesBackendConnection(ipam) }()
		return cmdAdd(ipam, confVersion)
	},
		cmdCheck,
		func(args *skel.CmdArgs) error {
			ipamConf, _, err := config.LoadIPAMConfig(args.StdinData, args.Args)
			if err != nil {
				logging.Errorf("IPAM configuration load failed: %s", err)
				return err
			}
			logging.Debugf("DEL - IPAM configuration successfully read: %+v", *ipamConf)

			ipam, err := kubernetes.NewKubernetesIPAM(args.ContainerID, args.IfName, *ipamConf)
			if err != nil {
				return logging.Errorf("IPAM client initialization error: %v", err)
			}
			defer func() { safeCloseKubernetesBackendConnection(ipam) }()
			return cmdDel(ipam)
		},
		cniversion.All,
		fmt.Sprintf("whereabouts %s", version.GetFullVersionWithRuntimeInfo()),
	)
}

func safeCloseKubernetesBackendConnection(ipam *kubernetes.KubernetesIPAM) {
	if err := ipam.Close(); err != nil {
		_ = logging.Errorf("failed to close the connection to the K8s backend: %v", err)
	}
}

func cmdCheck(args *skel.CmdArgs) error {
	// TODO
	return fmt.Errorf("CNI CHECK method is not implemented")
}

func cmdAdd(client *kubernetes.KubernetesIPAM, cniVersion string) error {
	// Initialize our result, and assign DNS & routing.
	result := &current.Result{}
	result.DNS = client.Config.DNS
	result.Routes = client.Config.Routes

	var newips []net.IPNet

	ctx, cancel := context.WithTimeout(context.Background(), types.AddTimeLimit)
	defer cancel()

	newips, err := kubernetes.IPManagement(ctx, types.Allocate, client.Config, client)
	if err != nil {
		logging.Errorf("Error at storage engine: %s", err)
		return fmt.Errorf("error at storage engine: %w", err)
	}

	var useVersion string
	for _, newip := range newips {
		// Determine if v4 or v6.
		if allocate.IsIPv4(newip.IP) {
			useVersion = "4"
		} else {
			useVersion = "6"
		}

		result.IPs = append(result.IPs, &current.IPConfig{
			Version: useVersion,
			Address: newip,
			Gateway: client.Config.Gateway})
	}

	// Assign all the static IP elements.
	for _, v := range client.Config.Addresses {
		result.IPs = append(result.IPs, &current.IPConfig{
			Version: v.Version,
			Address: v.Address,
			Gateway: v.Gateway})
	}

	return cnitypes.PrintResult(result, cniVersion)
}

func cmdDel(client *kubernetes.KubernetesIPAM) error {
	ctx, cancel := context.WithTimeout(context.Background(), types.DelTimeLimit)
	defer cancel()

	_, _ = kubernetes.IPManagement(ctx, types.Deallocate, client.Config, client)

	return nil
}
