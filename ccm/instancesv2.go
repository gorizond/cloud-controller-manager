package ccm

import (
	"context"
	"fmt"
	"k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
	"strings"
)

type InstancesV2 struct{}

func (i InstancesV2) InstanceExists(ctx context.Context, node *v1.Node) (bool, error) {
	return true, nil
}

func (i InstancesV2) InstanceShutdown(ctx context.Context, node *v1.Node) (bool, error) {
	return false, nil
}

func (i InstancesV2) InstanceMetadata(ctx context.Context, node *v1.Node) (*cloudprovider.InstanceMetadata, error) {
	var batmanIp string
	if node.Spec.ProviderID != "" {
		var err error
		batmanIp, err = getProviderID(node)
		if err != nil {
			return nil, err
		}
	} else {
		// TODO Check is valid ip
		batmanIp = node.Labels["batman/ip"]
		if batmanIp == "" {
			return nil, fmt.Errorf("batman/ip is empty for %s", node.Name)
		}
		klog.Infof("Set ip %s for node %s ", batmanIp, node.Name)
	}

	return &cloudprovider.InstanceMetadata{
		ProviderID:   fmt.Sprintf("%s://%s", providerName, batmanIp),
		InstanceType: "",
		NodeAddresses: []v1.NodeAddress{
			{
				Type:    v1.NodeInternalIP,
				Address: batmanIp,
			},
		},
		Zone:   "",
		Region: "",
	}, nil
}

func getProviderID(node *v1.Node) (string, error) {
	providerID, found := strings.CutPrefix(node.Spec.ProviderID, fmt.Sprintf("%s://", providerName))
	if !found {
		return "", fmt.Errorf("ProviderID does not follow expected format: %s", node.Spec.ProviderID)
	}
	id := providerID
	return id, nil
}
