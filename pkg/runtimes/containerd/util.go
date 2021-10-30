package containerd

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/opencontainers/runtime-spec/specs-go"
	k3d "github.com/rancher/k3d/v5/pkg/types"
)

func containerToK3DNode(c containerd.Container, ctx context.Context) (*k3d.Node, error) {
	info, err := c.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get container information: %w", err)
	}

	spec, err := c.Spec(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get container information: %w", err)
	}
	mounts := mountsToString(spec.Mounts)
	var node = &k3d.Node{
		Name:    info.ID,
		Image:   info.Image,
		Env:     spec.Process.Env,
		Args:    spec.Process.Args,
		Volumes: mounts,
		//Cmd:     spec.Process.CommandLine,
		//Restart
		//Created
		//RuntimeLabels
		//K3sNodeLabels
		//Networks
		//ExtraHosts
		//ServerOpts
		//AgentOpts
		//GPURequest
		//Memory
		//State
		//IP
		//HookActions,
		//Role:  nil,
	}
	return node, nil
}

func mountsToString(mounts []specs.Mount) []string {
	var mount []string
	for _, mnt := range mounts {
		mount = append(mount, fmt.Sprintf("%s %s %s %s", mnt.Type, mnt.Source, mnt.Destination, mnt.Options))
	}
	return mount
}
