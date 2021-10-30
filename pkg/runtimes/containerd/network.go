package containerd

import (
	"context"
	k3d "github.com/rancher/k3d/v5/pkg/types"
	"net"
)

func (c Containerd) GetNetwork(ctx context.Context, network *k3d.ClusterNetwork) (*k3d.ClusterNetwork, error) {
	panic("implement me")
}

func (c Containerd) GetHostIP(ctx context.Context, network string) (net.IP, error) {
	panic("implement me")

}
